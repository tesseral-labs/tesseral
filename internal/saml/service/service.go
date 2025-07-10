package service

import (
	"fmt"
	"html/template"
	"net/http"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/tesseral-labs/tesseral/internal/common/accesstoken"
	"github.com/tesseral-labs/tesseral/internal/cookies"
	"github.com/tesseral-labs/tesseral/internal/emailaddr"
	"github.com/tesseral-labs/tesseral/internal/saml/authn"
	"github.com/tesseral-labs/tesseral/internal/saml/internal/saml"
	"github.com/tesseral-labs/tesseral/internal/saml/store"
)

type Service struct {
	AccessTokenIssuer *accesstoken.Issuer
	Store             *store.Store
	Cookier           *cookies.Cookier
}

func (s *Service) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("GET /api/saml/v1/{samlConnectionID}/init", withErr(s.init))
	mux.Handle("POST /api/saml/v1/{samlConnectionID}/acs", withErr(s.acs))
	mux.Handle("POST /api/saml/v1/{samlConnectionID}/finish", withErr(s.finish))

	return mux
}

type initTemplateData struct {
	SignOnURL   string
	SAMLRequest string
}

var initTemplate = template.Must(template.New("init").Parse(`
<html>
	<body>
		<form method="POST" action="{{ .SignOnURL }}">
			<input type="hidden" name="SAMLRequest" value="{{ .SAMLRequest }}"></input>
		</form>
		<script>
			document.forms[0].submit();
		</script>
	</body>
</html>
`))

func (s *Service) init(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	samlConnectionID := r.PathValue("samlConnectionID")

	samlConnectionInitData, err := s.Store.GetSAMLConnectionInitData(ctx, samlConnectionID)
	if err != nil {
		return err
	}

	initRes := saml.Init(&saml.InitRequest{
		RequestID:  uuid.NewString(),
		SPEntityID: samlConnectionInitData.SPEntityID,
		Now:        time.Now(),
	})

	w.Header().Set("Content-Type", "text/html")
	if err := initTemplate.Execute(w, initTemplateData{
		SignOnURL:   samlConnectionInitData.IDPRedirectURL,
		SAMLRequest: initRes.SAMLRequest,
	}); err != nil {
		return fmt.Errorf("execute template: %w", err)
	}

	return nil
}

type acsTemplateData struct {
	FinishURL    string
	SAMLResponse string
}

var acsTemplate = template.Must(template.New("acs").Parse(`
<html>
	<body>
		<form method="POST" action="{{ .FinishURL }}">
			<input type="hidden" name="SAMLResponse" value="{{ .SAMLResponse }}"></input>
		</form>
		<script>
			document.forms[0].submit();
		</script>
	</body>
</html>
`))

func (s *Service) acs(w http.ResponseWriter, r *http.Request) error {
	samlConnectionID := r.PathValue("samlConnectionID")

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}

	w.Header().Set("Content-Type", "text/html")
	if err := acsTemplate.Execute(w, acsTemplateData{
		FinishURL:    fmt.Sprintf("/api/saml/v1/%s/finish", samlConnectionID),
		SAMLResponse: r.Form.Get("SAMLResponse"),
	}); err != nil {
		return fmt.Errorf("execute template: %w", err)
	}

	return nil
}

func (s *Service) finish(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()
	samlConnectionID := r.PathValue("samlConnectionID")

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}

	samlConnectionACSData, err := s.Store.GetSAMLConnectionACSData(ctx, samlConnectionID)
	if err != nil {
		// todo handle specifically a saml connection not yet fully configured
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}

	validateRes, err := saml.Validate(&saml.ValidateRequest{
		SAMLResponse:   r.Form.Get("SAMLResponse"),
		IDPCertificate: samlConnectionACSData.IDPX509Certificate,
		IDPEntityID:    samlConnectionACSData.IDPEntityID,
		SPEntityID:     samlConnectionACSData.SPEntityID,
		Now:            time.Now(),
	})
	if err != nil {
		return err
	}

	email := validateRes.SubjectID
	domain, err := emailaddr.Parse(email)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return nil
	}

	if !slices.Contains(samlConnectionACSData.OrganizationDomains, domain) {
		http.Error(w, "bad domain", http.StatusBadRequest)
		return nil
	}

	// For IdP-initiated SAML, the intermediate session is not present and must be created.
	intermediateSessionID := authn.IntermediateSessionID(ctx)
	if intermediateSessionID == nil {
		intermediateSession, err := s.Store.CreateIntermediateSession(ctx)
		if err != nil {
			return fmt.Errorf("create intermediate session: %w", err)
		}
		intermediateAccessToken, err := s.Cookier.NewIntermediateAccessToken(ctx, authn.ProjectID(ctx), intermediateSession.SecretToken)
		if err != nil {
			return fmt.Errorf("create intermediate access token cookie: %w", err)
		}
		w.Header().Set("Set-Cookie", intermediateAccessToken)
		ctx = authn.NewContext(ctx, intermediateSession.IntermediateSession, authn.ProjectID(ctx))
	}

	redirectURL, err := s.Store.FinishLogin(ctx, store.FinishLoginRequest{
		Email:                    email,
		VerifiedSAMLConnectionID: samlConnectionID,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return nil
	}

	w.Header().Add("Location", redirectURL)
	w.WriteHeader(http.StatusFound)
	return nil
}

func withErr(f func(w http.ResponseWriter, r *http.Request) error) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			http.Error(w, "", http.StatusInternalServerError)
			panic(err)
		}
	})
}

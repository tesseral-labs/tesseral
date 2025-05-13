package service

import (
	"fmt"
	"html/template"
	"net/http"
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

func (s *Service) acs(w http.ResponseWriter, r *http.Request) error {
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

	var domainOk bool
	for _, orgDomain := range samlConnectionACSData.OrganizationDomains {
		if orgDomain == domain {
			domainOk = true
			break
		}
	}

	if !domainOk {
		http.Error(w, "bad domain", http.StatusBadRequest)
	}

	createSessionRes, err := s.Store.CreateSession(ctx, &store.CreateSessionRequest{
		SAMLConnectionID: samlConnectionID,
		Email:            email,
	})
	if err != nil {
		return fmt.Errorf("create session: %w", err)
	}

	accessToken, err := s.AccessTokenIssuer.NewAccessToken(ctx, createSessionRes.RefreshToken)
	if err != nil {
		return fmt.Errorf("issue access token: %w", err)
	}

	refreshTokenCookie, err := s.Cookier.NewRefreshToken(ctx, authn.ProjectID(ctx), createSessionRes.RefreshToken)
	if err != nil {
		return fmt.Errorf("issue refresh token cookie: %w", err)
	}

	accessTokenCookie, err := s.Cookier.NewAccessToken(ctx, authn.ProjectID(ctx), accessToken)
	if err != nil {
		return fmt.Errorf("issue access token cookie: %w", err)
	}

	w.Header().Add("Set-Cookie", refreshTokenCookie)
	w.Header().Add("Set-Cookie", accessTokenCookie)

	w.Header().Add("Location", createSessionRes.RedirectURI)
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

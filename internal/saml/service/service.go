package service

import (
	"fmt"
	"html/template"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/tesseral-labs/tesseral/internal/emailaddr"
	"github.com/tesseral-labs/tesseral/internal/saml/internal/saml"
	"github.com/tesseral-labs/tesseral/internal/saml/store"
)

type Service struct {
	Store *store.Store
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

	validateRes, validateProblems, err := saml.Validate(&saml.ValidateRequest{
		SAMLResponse:   r.Form.Get("SAMLResponse"),
		IDPCertificate: samlConnectionACSData.IDPX509Certificate,
		IDPEntityID:    samlConnectionACSData.IDPEntityID,
		SPEntityID:     samlConnectionACSData.SPEntityID,
		Now:            time.Now(),
	})
	if err != nil {
		return err
	}

	// todo visual treatment here
	if validateProblems != nil {
		http.Error(w, fmt.Sprintf("%#v", validateProblems), http.StatusBadRequest)
		return nil
	}

	domain, err := emailaddr.Parse(validateRes.SubjectID)
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
		// todo visual treatment
		http.Error(w, "bad domain", http.StatusBadRequest)
	}

	// todo issue session
	// todo redirect

	// just to prove the concept
	if _, err := fmt.Fprintf(w, "hi %s in organization %s!", validateRes.SubjectID, samlConnectionACSData.OrganizationID); err != nil {
		return err
	}

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

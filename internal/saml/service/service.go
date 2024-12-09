package service

import (
	"fmt"
	"net/http"
	"time"

	"github.com/openauth/openauth/internal/saml/internal/emailaddr"
	"github.com/openauth/openauth/internal/saml/internal/saml"
	"github.com/openauth/openauth/internal/saml/store"
)

type Service struct {
	Store *store.Store
}

func (s *Service) Handler() http.Handler {
	mux := http.NewServeMux()

	mux.Handle("POST /saml/v1/{samlConnectionID}/acs", withErr(s.acs))

	return mux
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
	if _, err := w.Write([]byte(fmt.Sprintf("hi %s in organization %s!", validateRes.SubjectID, samlConnectionACSData.OrganizationID))); err != nil {
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

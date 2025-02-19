package saml_test

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tesseral-labs/tesseral/internal/saml/internal/saml"
)

func TestValidate(t *testing.T) {
	entries, err := os.ReadDir("testdata/assertions")
	assert.NoError(t, err)

	for _, entry := range entries {
		t.Run(entry.Name(), func(t *testing.T) {
			assertion, err := os.ReadFile(fmt.Sprintf("testdata/assertions/%s/assertion.xml", entry.Name()))
			require.NoError(t, err)

			metadata, err := os.ReadFile(fmt.Sprintf("testdata/assertions/%s/metadata.xml", entry.Name()))
			require.NoError(t, err)

			params, err := os.ReadFile(fmt.Sprintf("testdata/assertions/%s/params.json", entry.Name()))
			require.NoError(t, err)

			parseMetadataRes, err := saml.ParseMetadata(metadata)
			require.NoError(t, err)

			var paramData struct {
				SPEntityID string    `json:"sp_entity_id"`
				Now        time.Time `json:"now"`
			}
			err = json.Unmarshal(params, &paramData)
			require.NoError(t, err)

			_, problems, err := saml.Validate(&saml.ValidateRequest{
				SAMLResponse:   base64.StdEncoding.EncodeToString(assertion),
				IDPCertificate: parseMetadataRes.IDPCertificate,
				IDPEntityID:    parseMetadataRes.IDPEntityID,
				SPEntityID:     paramData.SPEntityID,
				Now:            paramData.Now,
			})
			assert.NoError(t, err)
			assert.Nil(t, problems)
		})
	}
}

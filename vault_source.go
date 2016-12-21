package configo

import (
	"crypto/tls"
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
)

type Vault struct {
	token         string
	address       string
	document_name string
}

func newVault(token, address, document_name string) *Vault {
	return &Vault{
		token:         token,
		address:       address,
		document_name: document_name,
	}
}

func FromVaultDocument(token, address, document_name string) *JSONSource {
	vault := newVault(token, address, document_name)
	vault.token = token

	document, err := vault.getDocument()
	if err != nil {
		log.Panic("[ERROR] Vault document read error:", err)
	}

	return FromJSONObject(document)
}

/////////////////////////////////////////

func (this *Vault) getDocument() (map[string]interface{}, error) {
	response, err := this.requestDocument()
	if err != nil {
		return nil, err
	}

	document, err := parseDocument(response.Body)
	if err != nil {
		return nil, err
	}

	return document.Data, nil
}

func parseDocument(responseBody io.Reader) (*SecretDocument, error) {
	var document SecretDocument
	decoder := json.NewDecoder(responseBody)
	if err := decoder.Decode(&document); err != nil {
		return nil, err
	}
	return &document, nil
}

/////////////////////////////////////////

func (this *Vault) dialTLS(network, address string) (net.Conn, error) {
	return tls.Dial(network, address, &tls.Config{ServerName: this.address})
}

func (this *Vault) requestDocument() (*http.Response, error) {
	httpClient := &http.Client{
		Transport: &http.Transport{DialTLS: this.dialTLS},
	}
	request, err := http.NewRequest("GET", "https://"+this.address+":8200/v1/"+this.document_name, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("X-Vault-Token", this.token)
	response, err := httpClient.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

/////////////////////////////////////////

// From -> https://github.com/hashicorp/vault/blob/master/api/secret.go
type SecretDocument struct {
	LeaseID       string `json:"lease_id"`
	LeaseDuration int    `json:"lease_duration"`
	Renewable     bool   `json:"renewable"`

	// Data is the actual contents of the secret. The format of the data
	// is arbitrary and up to the secret backend.
	Data map[string]interface{} `json:"data"`

	// Warnings contains any warnings related to the operation. These
	// are not issues that caused the command to fail, but that the
	// client should be aware of.
	Warnings []string `json:"warnings"`

	// Auth, if non-nil, means that there was authentication information
	// attached to this response.
	// Auth *SecretAuth `json:"auth,omitempty"`
	Auth SecretAuth `json:"auth,omitempty"`
}
type SecretAuth struct {
	ClientToken string            `json:"client_token"`
	Accessor    string            `json:"accessor"`
	Policies    []string          `json:"policies"`
	Metadata    map[string]string `json:"metadata"`

	LeaseDuration int  `json:"lease_duration"`
	Renewable     bool `json:"renewable"`
}

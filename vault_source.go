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
	token        string
	address      string
	documentName string
}

func newVault(token, address, documentName string) *Vault {
	return &Vault{
		token:        token,
		address:      address,
		documentName: documentName,
	}
}

func FromVaultDocument(token, address, documentName string) *JSONSource {
	vault := newVault(token, address, documentName)
	vault.token = token

	document, err := vault.getDocument()
	if err != nil {
		log.Panic("[ERROR] Vault document read error:", err) // TODO: only panic after N retries failed
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

func parseDocument(responseBody io.Reader) (*vaultDocument, error) {
	var document vaultDocument
	decoder := json.NewDecoder(responseBody)
	if err := decoder.Decode(&document); err != nil {
		return nil, err
	}
	return &document, nil // TODO: log any warnings that came back
}

/////////////////////////////////////////

func (this *Vault) dialTLS(network, address string) (net.Conn, error) {
	return tls.Dial(network, address, &tls.Config{ServerName: this.address})
}

func (this *Vault) requestDocument() (*http.Response, error) {
	httpClient := &http.Client{
		Transport: &http.Transport{DialTLS: this.dialTLS}, // TODO: timeouts
	}
	request, err := http.NewRequest("GET", "https://"+this.address+":8200/v1/"+this.documentName, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("X-Vault-Token", this.token)
	response, err := httpClient.Do(request) // TODO: retry
	if err != nil {
		return nil, err
	}

	return response, nil
}

/////////////////////////////////////////

// From https://github.com/hashicorp/vault/blob/master/api/secret.go
type vaultDocument struct {
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
	// Auth *vaultAuthentication `json:"auth,omitempty"`
	Auth vaultAuthentication `json:"auth,omitempty"`
}
type vaultAuthentication struct {
	ClientToken string            `json:"client_token"`
	Accessor    string            `json:"accessor"`
	Policies    []string          `json:"policies"`
	Metadata    map[string]string `json:"metadata"`

	LeaseDuration int  `json:"lease_duration"`
	Renewable     bool `json:"renewable"`
}

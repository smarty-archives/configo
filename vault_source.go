package configo

import (
	"crypto/tls"
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
	"time"
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

	for _, ip := range vault.getIPList() {
		document, err := vault.getDocument(ip)
		if err != nil {
			log.Println("[WARN] Vault document read error:", err)
			continue
		}
		return FromJSONObject(document)
	}
	log.Panic("Unable to get document from Vault")
	return nil
}

/////////////////////////////////////////

func (this *Vault) getDocument(ip string) (map[string]interface{}, error) {
	response, err := this.requestDocument(ip)
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
		log.Println("[WARN] vault source document decode error:", err)
		return nil, err
	}
	return &document, nil
}

/////////////////////////////////////////

func (this *Vault) dialTLS(network, address string) (net.Conn, error) {
	return tls.Dial(network, address, &tls.Config{ServerName: this.address})
}

func (this *Vault) requestDocument(ip string) (*http.Response, error) {
	httpClient := &http.Client{
		Transport: &http.Transport{DialTLS: this.dialTLS},
		Timeout:   time.Duration(5 * time.Second),
	}
	retryClient := NewRetryClient(httpClient, maxRetries, requestTimeout)

	request, err := http.NewRequest("GET", "https://"+ip+":8200/v1/"+this.documentName, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("X-Vault-Token", this.token)
	response, err := retryClient.Do(request)
	if err != nil {
		return nil, err
	}

	this.checkResponse(response)
	return response, nil
}

func (this *Vault) getIPList() []string {
	ips, err := net.LookupHost(this.address)
	if err != nil {
		log.Fatalf("[ERROR] %s", err)
	}
	return ips
}

func (this *Vault) checkResponse(response *http.Response) {
	if response != nil {
		switch response.StatusCode {
		case 200:
			log.Println("[INFO] Success, data returned")
		case 204:
			log.Println("[INFO] Success, no data returned")
		case 400:
			log.Println("[INFO] Invalid request, missing or invalid data")
		case 403:
			log.Println("[INFO] Forbidden. Credentials are wrong or you do not have permission")
		case 404:
			log.Println("[INFO] Invalid path. Path may be invalid or you do not have permission to view the path")
		case 429:
			log.Println("[INFO] Rate limite exceeded")
		case 500:
			log.Println("[INFO] Internal server error")
		case 503:
			log.Println("[INFO] Vault is down for maintenance or sealed")
		}
	}
}

/////////////////////////////////////////

type Client interface {
	Do(*http.Request) (*http.Response, error)
}

type RetryClient struct {
	inner   Client
	retries int
	timeout int
}

func NewRetryClient(inner Client, retries, timeout int) *RetryClient {
	return &RetryClient{
		inner:   inner,
		retries: retries,
		timeout: timeout,
	}
}

func (this *RetryClient) Do(request *http.Request) (response *http.Response, err error) {
	for current := 0; current <= this.retries; current++ {
		response, err = this.inner.Do(request)
		if err == nil {
			break
		}
	}
	return response, err
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

/////////////////////////////////////////

const (
	maxRetries     = 2
	requestTimeout = 2
)

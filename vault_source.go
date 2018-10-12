package configo

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

type VaultSource struct {
	token        string
	address      url.URL
	documentName string
}

func newVault(token, address, documentName string) *VaultSource {
	ret := &VaultSource{
		documentName: documentName,
	}
	ret.token = ret.checkToken(token)
	ret.address = *ret.checkAddress(address)
	return ret
}

func FromVaultDocument(token, address, documentName string) *JSONSource {
	vault := newVault(token, address, documentName)

	for _, ip := range vault.getIPList() {
		addrCopy := vault.address
		setURLHostPort(&addrCopy, ip, addrCopy.Port())
		document, err := vault.getDocument(addrCopy)
		if err != nil {
			log.Println("[WARN] Vault document read error:", err)
			continue
		}
		return FromJSONObject(document)
	}
	log.Panic("Unable to get document from Vault")
	return nil
}

//func (this *VaultSource) Initialize() {}

// Work around url.URL's broken ports with ip6
func setURLHostPort(u *url.URL, host, port string) {
	if strings.Contains(host, ":") && host[0] != '[' { //ip6
		host = "[" + host + "]"
	}

	u.Host = host + ":" + port
}

/////////////////////////////////////////

func (this *VaultSource) checkAddress(address string) (ret *url.URL) {
	if address == "" {
		address = os.Getenv("VAULT_ADDR")
	}

	address = strings.TrimSpace(address)

	if address == "" {
		log.Panic("No Vault address provided nor in environment VAULT_ADDR")
	}

	if ! strings.Contains(address, "//") {
		address = "//" + address
	}

	var err error
	if ret, err = url.Parse(address); err == nil {
		if ret.Scheme != "http" {
			ret.Scheme = "https"
		}

		if ret.Port() == "" {
			setURLHostPort(ret, ret.Hostname(), "8200")
		}

		ret.User = nil
	} else {
		log.Panic("Unable to parse Vault URL " + address + ": " + err.Error())
	}

	return
}

func (this *VaultSource) checkToken(token string) string {
	if token == "" {
		token = os.Getenv("VAULT_TOKEN")
	}

	token = strings.TrimSpace(token)

	if token == "" {
		log.Println("[WARN] No Vault token provided nor in environment VAULT_TOKEN")
	}

	return token
}

func (this *VaultSource) getDocument(addr url.URL) (map[string]interface{}, error) {
	response, err := this.requestDocument(addr)
	if err != nil {
		return nil, err
	}

	document, err := parseDocument(response.Body)
	if err != nil {
		return nil, err
	}

	if response.StatusCode >= 400 {
		err = errors.New(response.Status)
	}

	return document.Data, err
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

func (this *VaultSource) dialTLS(network, address string) (net.Conn, error) {
	return tls.Dial(network, address, &tls.Config{
		ServerName: this.address.Hostname(),
		InsecureSkipVerify: "" != os.Getenv("VAULT_SKIP_VERIFY"),
	})
}

func (this *VaultSource) requestDocument(addr url.URL) (*http.Response, error) {
	httpClient := &http.Client{
		Timeout:   time.Duration(5 * time.Second),
	}

	if addr.Scheme == "https" {
		httpClient.Transport = &http.Transport{DialTLS: this.dialTLS}
	}

	retryClient := NewRetryClient(httpClient, maxRetries, requestTimeout)

	addr.Path = "/v1/" + strings.Trim(this.documentName, "/")
	request, err := http.NewRequest("GET", addr.String(), nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("X-Vault-Token", this.token)
	response, err := retryClient.Do(request)
	if err != nil {
		return nil, err
	}

	this.checkResponse(response, request)
	return response, nil
}

func (this *VaultSource) getIPList() (ips []string) {
	var err error
	for i := 0; i < 3; i++ {
		if ips, err = net.LookupHost(this.address.Hostname()); err == nil {
			return ips
		}
		time.Sleep(100 * time.Millisecond)
		log.Println("[WARN] DNS lookup error")
	}

	log.Fatalf("[ERROR] %s", err)
	return
}

func (this *VaultSource) checkResponse(response *http.Response, request *http.Request) {
	if response != nil {
		switch response.StatusCode {
		case 200:
			log.Println("[INFO] Success, data returned")
		case 204:
			log.Println("[INFO] Success, no data returned")
		case 400:
			log.Println("[INFO] Invalid request, missing or invalid data:", request.URL.Path)
		case 403:
			log.Println("[INFO] Forbidden. Credentials are wrong or you do not have permission:", request.URL.Path)
		case 404:
			log.Println("[INFO] Invalid path. Path may be invalid or you do not have permission to view the path:", request.URL.Path)
		case 429:
			log.Println("[INFO] Rate limit exceeded")
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

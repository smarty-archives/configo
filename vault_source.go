package configo

import (
	"crypto/tls"
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"
)

type HostTools struct {
	winnerChannel  chan string
	errorChannel   chan bool
	failureChannel chan bool

	hostsAddress  string
	numberOfHosts int
	token         string
	documentName  string
}

func FromVaultDocument(vault_token, hosts_address, document_name string) *JSONSource {
	hostTools := HostTools{
		winnerChannel:  make(chan string, 1),
		errorChannel:   make(chan bool),
		failureChannel: make(chan bool),

		hostsAddress: hosts_address,
		token:        vault_token,
		documentName: document_name,
	}

	maxRetry := 2
	for i := 0; i < maxRetry; i++ {
		log.Println("[INFO] Finding best host...")
		hostTools.findBestHost()
		select {
		case winner := <-hostTools.winnerChannel:
			log.Println("[INFO] Winner:", winner)
			document, err := hostTools.getVaultDocument(winner)
			if err != nil {
				log.Println("[ERROR] Vault document read error:", err)
			}
			if document != nil {
				return FromJSONObject(document)
			}
		case <-hostTools.failureChannel:
			log.Println("[INFO] All host addresses have been checked, no winner found")
		}
	}
	panic("Unable to procure vault document")
}

/////////////////////////////////////////

func (this *HostTools) getVaultDocument(address string) (map[string]interface{}, error) {
	response, err := this.vaultClient(address, this.documentName)
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

func (this *HostTools) findBestHost() {
	go this.watchForProblems()
	for _, ip := range this.getIPList() {
		go this.getVaultNode(ip)
	}
}

func (this *HostTools) dialTLS(network, address string) (net.Conn, error) {
	return tls.Dial(network, address, &tls.Config{ServerName: this.hostsAddress})
}
func (this *HostTools) vaultClient(address, path string) (*http.Response, error) {
	client := &http.Client{
		Transport: &http.Transport{DialTLS: this.dialTLS},
	}
	request, err := http.NewRequest("GET", "https://"+address+":8200/v1/"+path, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("X-Vault-Token", this.token)
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (this *HostTools) getVaultNode(host string) {
	response, err := this.vaultClient(host, "sys/health")
	if err != nil {
		log.Printf("[ERROR] Vault server is not accessible (%s): %s", host, err)
	}
	if response != nil {
		switch response.StatusCode {
		case 200:
			this.winnerChannel <- host
			return
		case 429:
			log.Printf("[ERROR] Vault server is unsealed but in standby mode (%s)", host)
		case 500:
			log.Printf("[ERROR] Vault server is sealed or not initialized (%s)", host)
		}
	}
	this.errorChannel <- true
}

func (this *HostTools) getIPList() []string {
	ips, err := net.LookupHost(this.hostsAddress)
	if err != nil {
		log.Fatalf("[ERROR] %s", err)
	}
	this.numberOfHosts = len(ips)
	return ips
}

func (this *HostTools) watchForProblems() {
	for i := 1; ; i++ {
		<-this.errorChannel
		if i >= this.numberOfHosts {
			break
		}
	}
	this.failureChannel <- true
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

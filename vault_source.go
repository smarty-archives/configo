package configo

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"

	// vaultapi "github.com/hashicorp/vault/api"
)

type HostTools struct {
	winnerChannel  chan string
	errorChannel   chan bool
	failureChannel chan bool

	hostsAddress  string
	numberOfHosts int
}

var (
	token string
)

func FromVaultDocument(vault_token, hosts_address, document_name string) *JSONSource {
	token = vault_token
	hostTools := HostTools{
		winnerChannel:  make(chan string, 1),
		errorChannel:   make(chan bool),
		failureChannel: make(chan bool),

		hostsAddress: hosts_address,
	}

	maxRetry := 2
	for i := 0; i < maxRetry; i++ {
		log.Println("[INFO] Finding best host...")
		hostTools.findBestHost()
		select {
		case winner := <-hostTools.winnerChannel:
			log.Println("[INFO] Winner:", winner)
			document, err := getVaultDocument(winner)
			if err != nil {
				log.Println("[ERROR] Vault document read error:", err)
			}
			if document != nil {
				return &JSONSource{values: document}
			} else {
				// else retry
				log.Printf("[INFO] Document from %s was empty", winner)
			}
		case <-hostTools.failureChannel:
			log.Println("[INFO] All host addresses have been checked, no winner found")
		}
	}
	panic("AIEEEEEEEE!!!!!")
}

/////////////////////////////////////////

func getVaultDocument(address string) (map[string]interface{}, error) {
	response, err := vaultClient(address, "secret/smartystreets")
	if err != nil {
		log.Println("vaultclient error:", err)
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
	// decoder.UseNumber()
	if err := decoder.Decode(&document); err != nil {
		return nil, err
	}
	return &document, nil
}

/////////////////////////////////////////

func (this *HostTools) findBestHost() {
	go this.watchForProblems()
	for _, host := range this.getHostList() {
		go this.getVaultNode(host)
	}
}

func vaultClient(address, path string) (*http.Response, error) {
	client := &http.Client{}
	request, err := http.NewRequest("GET", "http://"+address+":8200/v1/"+path, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("X-Vault-Token", token)
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (this *HostTools) getVaultNode(host string) {
	response, err := vaultClient(host, "sys/health")
	if err != nil {
		log.Printf("[ERROR] Vault server is not accessible (%s): %s", host, err)
	}
	// fmt.Printf("\nseal-status raw: %#v (%s)\n", response, host)
	// fmt.Printf("seal-status: %v (%s)\n", response, host)
	if response != nil {
		// fmt.Printf("status code: %#v (%s)\n", response.StatusCode, host)
		switch response.StatusCode {
		case 200:
			// vault is ready to rock and roll
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

func (this *HostTools) getHostList() []string {
	hosts, err := net.LookupHost(this.hostsAddress)
	if err != nil {
		log.Fatalf("[ERROR] %s", err)
	}
	this.numberOfHosts = len(hosts)
	return hosts
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

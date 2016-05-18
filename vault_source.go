package configo

// TODO: this is a future configo source

import (
	"encoding/json"
	"io"
	"log"
	"net"
	"net/http"

	vaultapi "github.com/hashicorp/vault/api"
	// "github.com/smartystreets/configo"
)

type HostTools struct {
	winnerChannel  chan *VaultNode
	errorChannel   chan bool
	failureChannel chan bool

	hostsAddress  string
	numberOfHosts int
}

var (
	token string
)

// FromVaultDocument(vault_token, vault.ss.net, document_name) *jsonsource
// type JSONSource struct {
// 	values map[string]interface{}
// }
func FromVaultDocument(vault_token, hosts_address, document_name string) *JSONSource {
	token = vault_token
	hostTools := HostTools{
		winnerChannel:  make(chan *VaultNode, 1),
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
			log.Println("[INFO] Winner:", winner.address)
			// we have a winner, where shall we send it?
			// document, err := winner.client.Logical().Read("secret/smartystreets")
			document, err := getVaultDocument(winner.address)
			if err != nil {
				log.Println("[ERROR] Vault read error:", err)
			}
			if document != nil {
				return &JSONSource{values: document}
			} else {
				// else retry
				log.Printf("[INFO] Document from %s was empty", winner.address)
			}
		case <-hostTools.failureChannel:
			log.Println("[INFO] All host addresses have been checked, no winner found")
		}
	}
	panic("AIEEEEEEEE!!!!!")
}

/////////////////////////////////////////

func getVaultDocument(address string) (map[string]interface{}, error) {
	client := &http.Client{}
	// GET http://127.0.0.1:8200/v1/secret/smartystreets
	request, err := http.NewRequest("GET", "http://"+address+":8200/v1/secret/smartystreets2", nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("X-Vault-Token", token)
	response, err := client.Do(request)
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

type VaultNode struct {
	address string
	client  *vaultapi.Client
}

func NewVaultNode(address string, client *vaultapi.Client) *VaultNode {
	return &VaultNode{
		address: address,
		client:  client,
	}
}

func (this *HostTools) getVaultNode(host string) {
	client := newVaultClient(host)
	status, err := client.Sys().SealStatus()
	if err != nil {
		log.Println("[ERROR] SealStatus check:", err)
		this.errorChannel <- true
		return
	}
	if status.Sealed {
		log.Printf("[WARN] Node at %s is still sealed\n", host)
		this.errorChannel <- true
		return
	}

	this.winnerChannel <- NewVaultNode(host, client)
}

func newVaultClient(host string) *vaultapi.Client {
	config := vaultapi.DefaultConfig()

	config.Address = "http://" + host + ":8200"
	client, err := vaultapi.NewClient(config)
	if err != nil {
		log.Println("[ERROR] Initializing vault client:", err)
	}

	client.SetToken(token)

	return client
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

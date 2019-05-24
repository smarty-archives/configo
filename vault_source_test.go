package configo

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/smartystreets/assertions/should"
	"github.com/smartystreets/gunit"
)

func TestVaultSourceFixture(t *testing.T) {
	if testing.Short() {
		return
	}
	gunit.Run(new(VaultTestFixture), t)
}

type VaultTestFixture struct {
	*gunit.Fixture
}

func (this *VaultTestFixture) Setup() {}

////////////////////////////////////////////////////////////////

func (this *VaultTestFixture) TestEnvironment() {
	path := "/my/path/here"
	tokenCmd := "22-25-12"
	tokenEnv := "88-55-44"
	os.Setenv("VAULT_TOKEN", tokenEnv)
	os.Setenv("VAULT_ADDR", "[::1]")

	vault := newVault("", "", path)
	this.So(vault.address.String(), should.Equal, "https://[::1]:8200")
	this.So(vault.token, should.Equal, tokenEnv)
	this.So(vault.documentName, should.Equal, path)

	addrCmd := "http://169.254.0.1:1111"
	vault = newVault(tokenCmd, addrCmd, path)
	this.So(vault.address.String(), should.Equal, addrCmd)
	this.So(vault.token, should.Equal, tokenCmd)
	this.So(vault.documentName, should.Equal, path)
}

func (this *VaultTestFixture) TestBadEndpointPanics() {
	this.So(func() { FromVaultDocument("xyz", "localhost:29999", "") }, should.Panic)
}

func (this *VaultTestFixture) TestValidHTTPEndpoint() {
	svr := dummyHTTP(false, nil)
	defer svr.Close()
	src := FromVaultDocument("xyz", svr.URL, "")
	this.So(src, should.NotBeEmpty)

	data, err := src.Strings("string")
	this.So(err, should.BeEmpty)
	this.So(data, should.NotBeEmpty)

	data, err = src.Strings("bool")
	this.So(err, should.BeEmpty)
	this.So(data, should.NotBeEmpty)
}

func (this *VaultTestFixture) TestValidHTTPSEndpoint() {
	svr := dummyHTTP(true, nil)
	defer svr.Close()

	os.Setenv("VAULT_SKIP_VERIFY", "")
	this.So(func() { FromVaultDocument("xyz", svr.URL, "") }, should.Panic)

	os.Setenv("VAULT_SKIP_VERIFY", "1")
	src := FromVaultDocument("9-10-11", svr.URL, "/secret/my/doc")
	this.So(src, should.NotBeEmpty)

	data, err := src.Strings("string")
	this.So(err, should.BeEmpty)
	this.So(data, should.NotBeEmpty)

	data, err = src.Strings("bool")
	this.So(err, should.BeEmpty)
	this.So(data, should.NotBeEmpty)
}

func (this *VaultTestFixture) TestValidUnauthenticatedEndpoint() {
	svr := dummyHTTP(false, httpHandleUnauthenticated)
	defer svr.Close()

	this.So(func() { FromVaultDocument("xyz", svr.URL, "") }, should.Panic)
}

////////////////////////////////////////////////////////////////

// Construct a dummy http(s) server
func dummyHTTP(isTLS bool, handler func(http.ResponseWriter, *http.Request)) (server *httptest.Server) {
	if handler == nil {
		handler = httpHandleData
	}

	if isTLS {
		server = httptest.NewTLSServer(http.HandlerFunc(handler))
	} else {
		server = httptest.NewServer(http.HandlerFunc(handler))
	}

	return
}

// Valid response with simplistic data
func httpHandleData(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	response.WriteHeader(http.StatusOK)
	fmt.Fprintln(response, `{"lease_id": "1", "data": {"bool": true, "string": "String"}}`)
}

// Forbidden response
func httpHandleUnauthenticated(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("content-type", "application/json")
	response.WriteHeader(http.StatusForbidden)
	fmt.Fprintln(response, `{"errors":["permission denied"]}`)
}

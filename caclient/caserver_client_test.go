package caclient

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/useurmind/certbird/caserver"
)

func TestCAClientRequestCertWorks(t *testing.T) {
	ctx := caserver.NewTestContext(t)
	ctx.StartServer()

	defer func() {
		ctx.CleanupTest()
	}()

	dnsName := "mydummyhost.com"
	client := CAServerClient{
		RequestedDNSNames: []string { dnsName },
		CAServers: []string { fmt.Sprintf("http://%s:%d", ctx.ServerConfig.Address, ctx.ServerConfig.Port) },
	}

	cert, err := client.RequestCert()
	assert.Nil(t, err)

	assert.NotEqual(t, "", cert)
	ctx.ValidateCertificatePEM(cert, dnsName)
}
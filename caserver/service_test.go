package caserver

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCACertReturnsCACert(t *testing.T) {
	ctx := NewTestContext(t)
	serverConfig := ctx.PrepareTest()
	defer ctx.CleanupTest()

	service := NewService(serverConfig)

	caCertPEM, err := service.GetCACertificate()
	assert.Nil(t, err)

	assert.NotEqual(t, "", caCertPEM)
}

func TestGetCACertReturns404IfCACertMissing(t *testing.T) {
	ctx := NewTestContext(t)
	serverConfig := ctx.PrepareTest()
	defer ctx.CleanupTest()

	os.Remove(serverConfig.CACertFilePath)

	service := NewService(serverConfig)

	_, err := service.GetCACertificate()
	assert.NotNil(t, err)
	assert.Equal(t, 404, err.HTTPStatusCode)
}

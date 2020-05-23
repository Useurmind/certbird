package caserver

import (
	"os"
	"time"
	"testing"
	"github.com/stretchr/testify/assert"
)

// TestContext holds common information required for most tests.
type TestContext struct {
	// storageFolder is the folder holding all test data
	// there should be nothing in there except the data for the current test
	// putting all files into a deletable folder that makes cleanup easier
	storageFolder string

	t *testing.T
}

func NewTestContext(t *testing.T) *TestContext {
	return &TestContext{
		storageFolder: "test_data",
		t: t,
	}
}

func (c *TestContext) PrepareTest() *ServerConfig {
	serverConfig := TestServerConfig(c.storageFolder)
	
	validDuration, _ := time.ParseDuration("1h")

	err := os.MkdirAll(c.storageFolder, 666)
	assert.Nil(c.t, err)

	certConfig := &CertConfig{
		IsCA: true,
		ValidDuration: validDuration,
	}
	err = EnsureCACertificate(certConfig, serverConfig)
	assert.Nil(c.t, err)

	return &serverConfig
}

func (c *TestContext) CleanupTest() {	
	os.RemoveAll(c.storageFolder)
}
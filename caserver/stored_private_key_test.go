package caserver

import (
	"os"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/useurmind/certbird/utils"
)

func TestStoredPrivateKeyEnsure(t*  testing.T) {
	keyFilePath := "./ca_test_key.key"
	defer os.Remove(keyFilePath)

	spk := StoredPrivateKey {
		filePath: keyFilePath,
	}

	err := spk.ensure()
	assert.Nil(t, err)

	assert.Equal(t, true, utils.DoesFileExist(keyFilePath))
	assert.NotNil(t, spk.privKey)
}
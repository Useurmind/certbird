package caserver

import (
	"testing"
)

func TestGetCACertReturnsCACert(t *testing.T) {
	// certFilePath := "./ca_test_cert.pem"
	// keyFilePath := "./ca_test_cert.key"
	// defer os.Remove(certFilePath)
	// defer os.Remove(keyFilePath)

	// serverConfig := ServerConfig {
	// 	CACertFilePath: certFilePath,
	// 	CAKeyFilePath: keyFilePath,
	// }

	// certConfig := &CertConfig{
	// 	IsCA: true,
	// }
	// certConfig.ParseIPAddresses([]string { "121.12.34.12", "233.23.54.87" })

	// err := EnsureCACertificate(certConfig, serverConfig)
	// assert.Nil(t, err)

	// assert.Equal(t, true, utils.DoesFileExist(certFilePath))
	// assert.Equal(t, true, utils.DoesFileExist(keyFilePath))
}

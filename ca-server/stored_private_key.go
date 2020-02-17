package caserver

import (
	"fmt"
	"log"
	"bytes"
	"crypto/rand"
	"crypto/x509"
	"crypto/rsa"
	"io/ioutil"
	"encoding/pem"
	"github.com/useurmind/certbird/utils"
)

// StoredPrivateKey encapsulates fields that are commonly used for handling private keys.
type StoredPrivateKey struct {
	privKey *rsa.PrivateKey
	filePath string
}

func (spk *StoredPrivateKey) ensure() error {
	err := spk.loadFromPEM()
	if err != nil {
		return err
	}

	if spk.privKey == nil {		
		log.Println("Private key is missing, creating private key file", spk.filePath)
		privKey, err := rsa.GenerateKey(rand.Reader, 4096)
		if err != nil {
			return err
		}
		spk.privKey = privKey

		err = spk.writeToPEM()
		if err != nil {
			return err
		}
	}

	return nil
}

func (spk *StoredPrivateKey) loadFromPEM() error {
	filePath := spk.filePath

	log.Println("Trying to load private key from PEM file", filePath)
	if !utils.DoesFileExist(filePath) {
		log.Println("Key file", filePath, "does not exist yet")
		return nil
	}

	keyBytes, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	keyBlock, _ := pem.Decode(keyBytes)
	if keyBlock == nil {
		return fmt.Errorf("Could not read key block from PEM file %s", filePath)
	}

	if keyBlock.Type != utils.PEM_TYPE_PRIVATE_KEY {
		return fmt.Errorf("The PEM file %s does not contain a private key, found type %s", filePath, keyBlock.Type)
	}

	privKey, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)
	if err != nil {
		return err
	}

	spk.privKey = privKey

	return nil
}

func (spk *StoredPrivateKey) writeToPEM() error {
	filePath := spk.filePath
	privKey := spk.privKey

	log.Println("Writing private key as PEM to", filePath)

	privKeyPEM := new(bytes.Buffer)
	err := pem.Encode(privKeyPEM, &pem.Block{
		Type:  utils.PEM_TYPE_PRIVATE_KEY,
		Bytes: x509.MarshalPKCS1PrivateKey(privKey),
	})
	if err != nil {
		return err
	}
	
	err = ioutil.WriteFile(filePath, privKeyPEM.Bytes(), 0777)
	if err != nil {
		return err
	}

	return nil
}

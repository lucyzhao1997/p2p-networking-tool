package security

import (
    "crypto/rand"
    "crypto/rsa"
    "crypto/x509"
    "encoding/pem"
    "fmt"
)

func GenerateKeyPair() (*rsa.PrivateKey, *rsa.PublicKey, error) {
    //use crypto rsa library to generate key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
    if err != nil {
        return nil, nil, err
    }
    return privateKey, &privateKey.PublicKey, nil
package date

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"fmt"
)

func RSA_OAEP_Encrypt(secretMessage []byte, key rsa.PublicKey) []byte {
	label := []byte("OAEP Encrypted")
	rng := rand.Reader
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rng, &key,
		secretMessage, label)

	if err != nil {
		fmt.Println(err.Error())
	}

	return ciphertext
}

func RSA_OAEP_Decrypt(cipherText []byte, privKey rsa.PrivateKey) []byte {
	label := []byte("OAEP Encrypted")
	rng := rand.Reader
	plaintext, err := rsa.DecryptOAEP(sha256.New(), rng, &privKey, cipherText, label)

	if err != nil {
		fmt.Println(err.Error())
	}

	return plaintext
}

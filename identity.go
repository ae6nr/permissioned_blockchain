package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path"
	"reflect"
)

type identity_t struct {
	label  string
	prvKey *ecdsa.PrivateKey
	pubKey *ecdsa.PublicKey
}

func (id *identity_t) GetPubBytes() []byte {
	pubBytes, err := x509.MarshalPKIXPublicKey(id.pubKey)
	if err != nil {
		panic(err)
	}
	return pubBytes
}

func encode(privateKey *ecdsa.PrivateKey, publicKey *ecdsa.PublicKey) (string, string) {
	x509Encoded, _ := x509.MarshalECPrivateKey(privateKey)
	pemEncoded := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: x509Encoded})

	x509EncodedPub, _ := x509.MarshalPKIXPublicKey(publicKey)
	pemEncodedPub := pem.EncodeToMemory(&pem.Block{Type: "PUBLIC KEY", Bytes: x509EncodedPub})

	return string(pemEncoded), string(pemEncodedPub)
}

func decode(pemEncoded string, pemEncodedPub string) (*ecdsa.PrivateKey, *ecdsa.PublicKey) {
	block, _ := pem.Decode([]byte(pemEncoded))
	x509Encoded := block.Bytes
	privateKey, _ := x509.ParseECPrivateKey(x509Encoded)

	blockPub, _ := pem.Decode([]byte(pemEncodedPub))
	x509EncodedPub := blockPub.Bytes
	genericPublicKey, _ := x509.ParsePKIXPublicKey(x509EncodedPub)
	publicKey := genericPublicKey.(*ecdsa.PublicKey)

	return privateKey, publicKey
}

func GenerateKeys(id string) {

	fmt.Printf("Generating keys for identity %s.\r\n", id)

	// generate keys
	privateKey, _ := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	publicKey := &privateKey.PublicKey

	// encode keys to save to file
	encPriv, encPub := encode(privateKey, publicKey)

	// decode keys
	priv2, pub2 := decode(encPriv, encPub)

	// check equality of decoded keys to originals
	if !reflect.DeepEqual(privateKey, priv2) {
		panic("private keys do not match")
	}
	if !reflect.DeepEqual(publicKey, pub2) {
		panic("public keys do not match")
	}

	// create keys directory if it doesn't exist
	if _, err := os.Stat("keys"); os.IsNotExist(err) {
		err := os.Mkdir("keys", os.ModeDir)
		if err != nil {
			panic(err)
		}
	}

	// check if keys already exist
	fname_pub := path.Join("keys", id+"_pub.pem")
	fname_prv := path.Join("keys", id+"_prv.pem")
	if _, err := os.Stat(fname_pub); errors.Is(err, os.ErrNotExist) {
		// good
	} else {
		fmt.Println("A public key already exists for that identity")
		return
	}
	if _, err := os.Stat(fname_prv); errors.Is(err, os.ErrNotExist) {
		// good
	} else {
		fmt.Println("A private key already exist for that identity")
		return
	}

	// save keys to files
	err := os.WriteFile(fname_pub, []byte(encPub), fs.ModeAppend) // save public key
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(fname_prv, []byte(encPriv), fs.ModeAppend) // save private key
	if err != nil {
		panic(err)
	}

}

func LoadIdentity(label string) (id identity_t) {
	var encPub, encPriv []byte
	var err error

	// get filenames to keys
	fname_pub := path.Join("keys", label+"_pub.pem")
	fname_prv := path.Join("keys", label+"_prv.pem")

	// read files
	encPub, err = os.ReadFile(fname_pub)
	if err != nil {
		GenerateKeys(label) // make new keys if identity doesn't already exist
		encPub, err = os.ReadFile(fname_pub)
		if err != nil {
			panic(err)
		}
	}
	encPriv, err = os.ReadFile(fname_prv)
	if err != nil {
		panic(err)
	}

	// decode
	privateKey, publicKey := decode(string(encPriv), string(encPub))

	// ensure keys match
	challenge := make([]byte, 16)
	rand.Reader.Read(challenge)
	hash := sha256.Sum256(challenge)
	sig, err := ecdsa.SignASN1(rand.Reader, privateKey, hash[:])
	if err != nil {
		panic(err)
	}
	valid := ecdsa.VerifyASN1(publicKey, hash[:], sig)
	if !valid {
		panic("signature verification failed")
	}

	id.prvKey = privateKey
	id.pubKey = publicKey
	id.label = label
	return id
}

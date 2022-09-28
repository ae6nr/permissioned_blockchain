package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"testing"
)

// This test is really just an example of how to sign hashes using an ECDSA algorithm
func TestSignature(t *testing.T) {
	privateKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		t.Error(err)
	}

	msg := "Hello, world!"
	hash := sha256.Sum256([]byte(msg))
	sig, err := ecdsa.SignASN1(rand.Reader, privateKey, hash[:])
	if err != nil {
		t.Error(err)
	}

	valid := ecdsa.VerifyASN1(&privateKey.PublicKey, hash[:], sig)
	if !valid {
		t.Errorf("error verifying signature")
	}
}

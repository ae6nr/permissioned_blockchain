package main

// https://gist.github.com/LukaGiorgadze/85b9e09d2008a03adfdfd5eea5964f93

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
)

type EllipticCurve struct {
	pubKeyCurve elliptic.Curve
	privateKey  *ecdsa.PrivateKey
	publicKey   *ecdsa.PublicKey
}

func (ec *EllipticCurve) GenerateKeys() (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	var err error
	privKey, err := ecdsa.GenerateKey(ec.pubKeyCurve, rand.Reader)

	if err != nil {
		ec.privateKey = privKey
		ec.publicKey = &privKey.PublicKey
	}

	return ec.privateKey, ec.publicKey, err
}

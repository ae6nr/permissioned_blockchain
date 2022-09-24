package main

import (
	"crypto"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path"

	"crypto/ecdsa"
	"crypto/rand"
	_ "crypto/sha256"
	"crypto/x509"
)

type block_t struct {
	timestamp   []byte
	prev_hash   []byte // hash of previous block
	hash        []byte // hash of entire block, including signature
	signed_hash []byte // hash corresponding to the signature
	data        []byte // data included in block
	validator   []byte // validator public key
	signature   []byte // signature corresponding to validator's public key
}

func (block *block_t) ComputeSignedHash() []byte {
	hash := crypto.SHA256.New()
	hash.Write(block.timestamp)
	hash.Write(block.prev_hash)
	hash.Write(block.validator)
	hash.Write(block.data)
	block.signed_hash = hash.Sum(nil)
	return hash.Sum(nil)
}

func (block *block_t) ComputeBlockHash() []byte {
	hash := crypto.SHA256.New()
	// hash.Write(block.timestamp)
	// hash.Write(block.prev_hash)
	// hash.Write(block.validator)
	// hash.Write(block.signature)
	// hash.Write(block.data)
	hash.Write(block.Marshal())
	block.hash = hash.Sum(nil)
	return block.hash
}

func (block *block_t) VerifySignature() (bool, error) {
	genericPublicKey, _ := x509.ParsePKIXPublicKey(block.validator)
	publicKey := genericPublicKey.(*ecdsa.PublicKey)
	if !ecdsa.VerifyASN1(publicKey, block.ComputeSignedHash(), block.signature) {
		return false, errors.New("invalid signature")
	} else {
		return true, nil
	}
}

func (block *block_t) VerifyData() (bool, error) {
	if len(block.data) >= 1024 {
		return false, errors.New("block data too long")
	}

	return true, nil
}

func (block *block_t) Verify() (bool, error) {
	_, err := block.VerifyData()
	if err != nil {
		return false, err
	}

	_, err = block.VerifySignature()
	if err != nil {
		return false, err
	}

	return true, nil
}

func NewBlock(timestamp int64, prev_hash []byte, data []byte, id identity_t) (block block_t) {
	block.timestamp = make([]byte, 8)
	binary.BigEndian.PutUint64(block.timestamp, uint64(timestamp))
	block.prev_hash = prev_hash
	block.data = data
	block.validator = id.GetPubBytes()

	sig, err := ecdsa.SignASN1(rand.Reader, id.prvKey, block.ComputeSignedHash())
	if err != nil {
		panic(err)
	}
	block.signature = sig
	block.ComputeBlockHash()

	_, err = block.Verify()
	if err != nil {
		panic(err)
	}

	return block
}

// produce the genesis block
func Genesis() block_t {
	var block block_t

	block.timestamp = make([]byte, 8)
	block.prev_hash = make([]byte, 32)

	hash := crypto.SHA256.New()
	hash.Write([]byte("redd"))
	block.data = hash.Sum(nil)

	validator_bytes, err := hex.DecodeString("3076301006072a8648ce3d020106052b8104002203620004d985ce1893c962f0dfe389b6193e4149a54eca746f9c1ba1f56b1ed898009a4669520de0b5e53e91336115c668e304b6d6a9b1e98bee50c0b0f1cf80b13e0f554c9df3a51bbee2ab1f7c37f12d563d6fb174bd7315cfbac97c09ad47e852afc9")
	if err != nil {
		panic("err")
	}
	signature_bytes, err := hex.DecodeString("3066023100ba5b934c6a39d563eb3f61d6db09bed020e8e29e39519013ccdb67445cf51aba148a8f755705db649cee33c2efa134260231009cbe7264fbce3c794f46b0bb45b72eb32543f5acfb43f044422c1e4bbcfecd069b668b52b9818a1ff6afee99783e5a2b")
	if err != nil {
		panic("err")
	}
	block.validator = validator_bytes
	block.signature = signature_bytes
	block.ComputeBlockHash()
	return block
}

func (block *block_t) Print() {
	fmt.Printf("Block %x\r\n", block.hash)
	fmt.Printf("  timestamp: %x\r\n", block.timestamp)
	fmt.Printf("  prev_hash:  %x\r\n", block.prev_hash)
	fmt.Printf("  validator: %x\r\n", block.validator)
	fmt.Printf("  signature: %x\r\n", block.signature)
	fmt.Printf("  data:      %x\r\n", block.data)
}

func (block *block_t) Marshal() []byte {
	var d []byte
	d = append(d, block.timestamp...)
	d = append(d, block.prev_hash...)
	d = append(d, block.validator...)
	d = append(d, block.signature...)
	d = append(d, block.data...)
	return d
}

// save block to a file
func (block *block_t) Save() (bool, error) {
	// create blocks directory if it doesn't exist
	if _, err := os.Stat("blocks"); os.IsNotExist(err) {
		err := os.Mkdir("blocks", os.ModeDir)
		if err != nil {
			return false, err
		}
	}

	_, err := block.Verify()
	if err != nil {
		return false, err
	}

	hashhex := hex.EncodeToString(block.hash)
	fname := path.Join("blocks", hashhex+".dat")
	err = os.WriteFile(fname, block.Marshal(), 0777)
	if err != nil {
		return false, err
	}

	return true, nil
}

func BlockTest() {
	id := LoadIdentity("main")

	block0 := Genesis()
	block0.Print()
	_, err := block0.Verify()
	if err != nil {
		panic(err)
	}

	block1 := NewBlock(GetCurrentTimestamp(), block0.hash, []byte("I added a block!"), id)
	block1.Print()

	block2 := NewBlock(GetCurrentTimestamp(), block1.hash, []byte("Yet another block."), id)
	block2.Print()

}

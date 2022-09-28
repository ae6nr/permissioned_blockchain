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

// the block struct
type block_t struct {
	timestamp        [TIMESTAMP_SIZE]byte
	prev_hash        [HASH_SIZE]byte   // hash of previous block
	hash             [HASH_SIZE]byte   // hash of entire block, including signature
	signed_hash      [HASH_SIZE]byte   // hash corresponding to the signature
	validator        [PUBKEY_SIZE]byte // validator public key
	signature_length uint8             // length of signature (variable)
	signature        []byte            // signature corresponding to validator's public key
	tx               transaction_t     // data included in block
}

// create a new block
func NewBlock(timestamp int64, prev_hash []byte, tx transaction_t, id identity_t) (block block_t, err error) {
	binary.BigEndian.PutUint64(block.timestamp[:], uint64(timestamp))
	n := copy(block.prev_hash[:], prev_hash)
	if n != len(prev_hash) {
		s := fmt.Sprintf("length of prev_hash (%d) was too long for block.validator (%d)", len(prev_hash), n)
		return block, errors.New(s)
	}

	block.tx = tx
	pubbytes := id.GetPubBytes()
	n = copy(block.validator[:], pubbytes)
	if n != len(pubbytes) {
		s := fmt.Sprintf("length of pubbytes (%d) was too long for block.validator (%d)", len(pubbytes), n)
		return block, errors.New(s)
	}

	sig, err := ecdsa.SignASN1(rand.Reader, id.prvKey, block.ComputeSignedHash())
	if err != nil {
		return block, err
	}
	block.signature = sig
	block.signature_length = byte(len(block.signature))
	block.ComputeBlockHash()

	_, err = block.Verify()
	if err != nil {
		return block, err
	}

	return block, nil
}

// compute the hash that will be signed
func (block *block_t) ComputeSignedHash() []byte {
	hash := crypto.SHA256.New()
	hash.Write(block.timestamp[:])
	hash.Write(block.prev_hash[:])
	hash.Write(block.validator[:])
	hash.Write(block.tx.Marshal())
	copy(block.signed_hash[:], hash.Sum(nil))
	return block.signed_hash[:]
}

// compute the block hash
func (block *block_t) ComputeBlockHash() []byte {
	hash := crypto.SHA256.New()
	hash.Write(block.Marshal())
	copy(block.hash[:], hash.Sum(nil))
	return block.hash[:]
}

// verify the validity of a block
// transaction portion must be less than TX_MAX_SIZE
// also checks for a valid signature
func (block *block_t) Verify() (bool, error) {

	// verify data
	if len(block.tx.Marshal()) >= int(TX_MAX_SIZE) {
		return false, errors.New("block data too long")
	}

	// verify signature
	genericPublicKey, err := x509.ParsePKIXPublicKey(block.validator[:])
	if err != nil {
		return false, err
	}
	publicKey := genericPublicKey.(*ecdsa.PublicKey)
	if !ecdsa.VerifyASN1(publicKey, block.ComputeSignedHash(), block.signature) {
		return false, errors.New("invalid signature")
	}

	// check if valid transaction type
	if block.tx.txtype < Entry || block.tx.txtype > Permission {
		return false, errors.New("invalid transaction type")
	}

	return true, nil
}

// returns the hash of the block
func (block *block_t) GetHash() []byte {
	return block.hash[:]
}

// returns the validator public key as a hex-encoded string
func (block *block_t) GetValidatorString() string {
	return hex.EncodeToString(block.validator[:])
}

// prints information about the block
func (block *block_t) Print() {
	fmt.Printf("Block %x\r\n", block.hash)
	fmt.Printf("  timestamp:  %x\r\n", block.timestamp)
	fmt.Printf("  prev_hash:  %x\r\n", block.prev_hash)
	fmt.Printf("  validator:  %x\r\n", block.validator)
	fmt.Printf("  sig_length: %x\r\n", block.signature_length)
	fmt.Printf("  signature:  %x\r\n", block.signature)
	fmt.Printf("  data:       %x\r\n", block.tx.Marshal())
}

// creates a binary representation of the block's data
func (block *block_t) Marshal() []byte {
	var d []byte
	d = append(d, block.timestamp[:]...)
	d = append(d, block.prev_hash[:]...)
	d = append(d, block.validator[:]...)
	d = append(d, block.signature_length)
	d = append(d, block.signature[:]...)
	d = append(d, block.tx.Marshal()...)
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

	hashhex := hex.EncodeToString(block.hash[:])
	fname := path.Join("blocks", hashhex+".dat")
	err = os.WriteFile(fname, block.Marshal(), 0777)
	if err != nil {
		return false, err
	}

	return true, nil
}

func getBounds(offset int, length int) (int, int) {
	return offset, offset + length
}

// loads the block from a file
func LoadBlock(hash []byte) (block block_t, err error) {

	hashhex := hex.EncodeToString(hash)
	fname := path.Join("blocks", hashhex+".dat")
	data, err := os.ReadFile(fname)
	if err != nil {
		return block, err
	}

	var i, j int
	i, j = getBounds(0, int(TIMESTAMP_SIZE))
	copy(block.timestamp[:], data[i:j])
	i, j = getBounds(j, int(HASH_SIZE))
	copy(block.prev_hash[:], data[i:j])
	i, j = getBounds(j, int(PUBKEY_SIZE))
	copy(block.validator[:], data[i:j])
	i, j = getBounds(j, 1)
	block.signature_length = data[i]
	i, j = getBounds(j, int(block.signature_length))
	block.signature = data[i:j]
	err = block.tx.Unmarshal(data[j:])
	if err != nil {
		return block, err
	}

	block.ComputeSignedHash()
	block.ComputeBlockHash()
	_, err = block.Verify()
	return block, err
}

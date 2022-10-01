package main

import (
	"crypto"
	"encoding/hex"
	"fmt"
)

// create a genesis block signed by the specified id
func GenesisBootstrap(id identity_t) {
	// var block block_t
	hash := crypto.SHA256.New()
	hash.Write([]byte("redd"))
	tx := NewTx_Entry(hash.Sum(nil))
	block, err := NewBlock(0, make([]byte, HASH_SIZE), tx, id)
	if err != nil {
		panic(err)
	}

	block.Print()
	block.Save()

	loaded_block, err := LoadBlock(block.hash[:])
	if err != nil {
		fmt.Println("error loading block")
		fmt.Println("please ensure the data loaded correctly")
		loaded_block.Print()
		panic(err)
	}
}

// produce the genesis block
func Genesis() block_t {
	var block block_t

	hash := crypto.SHA256.New()
	hash.Write([]byte("redd"))
	block.tx = NewTx_Entry(hash.Sum(nil))

	validator_bytes, err := hex.DecodeString(GENESIS_VALIDATOR)
	if err != nil {
		panic("err")
	}
	signature_bytes, err := hex.DecodeString(GENESIS_SIGNATURE)
	if err != nil {
		panic("err")
	}
	n := copy(block.validator[:], validator_bytes)
	if n != len(validator_bytes) {
		s := fmt.Sprintf("length of validator_bytes (%d) was too long for block.validator (%d)", len(validator_bytes), n)
		panic(s)
	}

	block.signature_length = GENESIS_SIGNATURE_LENGTH
	block.signature = signature_bytes
	if n != len(validator_bytes) {
		s := fmt.Sprintf("length of signature_bytes (%d) was too long for block.validator (%d)", len(signature_bytes), n)
		panic(s)
	}

	block.ComputeBlockHash()
	return block
}

func NewChainFromGenesis() {
	blockchain.Init("main")
	blockchain.Save()
}

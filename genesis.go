package main

import (
	"crypto"
	"encoding/hex"
	"fmt"
)

// create a genesis block using the specified id
func GenesisBootstrap(id identity_t) {
	// var block block_t
	hash := crypto.SHA256.New()
	hash.Write([]byte("redd"))
	tx := NewTx_Entry(hash.Sum(nil))
	block, err := NewBlock(0, make([]byte, 32), tx, id)
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

	validator_bytes, err := hex.DecodeString("3076301006072a8648ce3d020106052b8104002203620004d985ce1893c962f0dfe389b6193e4149a54eca746f9c1ba1f56b1ed898009a4669520de0b5e53e91336115c668e304b6d6a9b1e98bee50c0b0f1cf80b13e0f554c9df3a51bbee2ab1f7c37f12d563d6fb174bd7315cfbac97c09ad47e852afc9")
	if err != nil {
		panic("err")
	}
	signature_bytes, err := hex.DecodeString("3066023100fe4adf500ae5f67c0274283315c1430e11eec469c7a7a5b68135615b5d80a540cd49d22eb593e5ebae16f257614a2559023100d8a760339f133c73db4eda3ef300959a29fa453271c6b7c4b166d4a783f8c97c84f166e8cb5621d1a9190403f08f7d1a")
	if err != nil {
		panic("err")
	}
	n := copy(block.validator[:], validator_bytes)
	if n != len(validator_bytes) {
		s := fmt.Sprintf("length of validator_bytes (%d) was too long for block.validator (%d)", len(validator_bytes), n)
		panic(s)
	}

	block.signature_length = 0x68
	block.signature = signature_bytes
	if n != len(validator_bytes) {
		s := fmt.Sprintf("length of signature_bytes (%d) was too long for block.validator (%d)", len(signature_bytes), n)
		panic(s)
	}

	block.ComputeBlockHash()
	return block
}

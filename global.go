package main

import "time"

// the blockchain in memory
var blockchain blockchain_t

func GetCurrentTimestamp() int64 {
	return time.Now().UTC().Unix()
}

const (
	MAIN_CHAIN_NAME string = "main"
)

// constants related to block formatting
const (
	TIMESTAMP_SIZE uint8  = 8
	HASH_SIZE      uint8  = 32
	PUBKEY_SIZE    uint8  = 120
	TX_MAX_SIZE    uint32 = 1048577 // 1 MB limit
)

// constants related to file structure
const (
	BLOCKCHAIN_DIR string = "data/blockchains"
	BLOCKS_DIR     string = "data/blocks"
	KEYS_DIR       string = "data/keys"
)

// constants related to the genesis block
const (
	GENESIS_VALIDATOR        = "3076301006072a8648ce3d020106052b8104002203620004d985ce1893c962f0dfe389b6193e4149a54eca746f9c1ba1f56b1ed898009a4669520de0b5e53e91336115c668e304b6d6a9b1e98bee50c0b0f1cf80b13e0f554c9df3a51bbee2ab1f7c37f12d563d6fb174bd7315cfbac97c09ad47e852afc9"
	GENESIS_SIGNATURE        = "3066023100fe4adf500ae5f67c0274283315c1430e11eec469c7a7a5b68135615b5d80a540cd49d22eb593e5ebae16f257614a2559023100d8a760339f133c73db4eda3ef300959a29fa453271c6b7c4b166d4a783f8c97c84f166e8cb5621d1a9190403f08f7d1a"
	GENESIS_SIGNATURE_LENGTH = 0x68
)

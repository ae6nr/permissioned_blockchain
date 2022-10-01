package main

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path"
)

type blockchain_t struct {
	label      string
	blocks     []block_t
	validators map[string]uint32 // validators and how many blocks they are allowed to mint
}

// initialize the blockchain with the genesis block
func (bc *blockchain_t) Init(label string) {
	bc.label = label
	bc.blocks = append(bc.blocks, Genesis())
	bc.validators = make(map[string]uint32)
	bc.validators[bc.blocks[0].GetValidatorString()] = 4294967295
}

// appends a block to the chain if the block is valid
func (bc *blockchain_t) AppendBlock(block block_t) (bool, error) {
	_, err := block.Verify()
	if err != nil {
		return false, err
	}

	// check that block hashes actually form a chain
	if !bytes.Equal(block.prev_hash[:], bc.blocks[len(bc.blocks)-1].ComputeBlockHash()) {
		s := fmt.Sprintf("candidate block doesn't contain hash of previous block\r\n  candidate prev_hash %x\r\n  hash of prev        %x\r\n", block.prev_hash, bc.blocks[len(bc.blocks)-1].ComputeBlockHash())
		return false, errors.New(s)
	}

	// check if validator has authorization to publish to the chain
	if val, ok := bc.validators[block.GetValidatorString()]; !ok || val <= 0 {
		return false, errors.New("validator is not authorized")
	}

	// check for various transaction types
	if block.tx.txtype == Entry {
		// no need to do anything
	} else if block.tx.txtype == Permission {
		n, validator, err := block.tx.ParseTx_Permission()
		if err != nil {
			return false, err
		}

		// check that validator has enough blocks to delegate
		// take them away if so
		if val, ok := bc.validators[block.GetValidatorString()]; !ok || val <= n {
			return false, errors.New("validator is not authorized to delegate that many blocks")
		}
		bc.validators[block.GetValidatorString()] -= n

		// give blocks to other validator
		nval := hex.EncodeToString(validator[:])
		if _, ok := bc.validators[nval]; !ok {
			bc.validators[nval] = n
		} else {
			bc.validators[nval] += n
		}

	}

	// add the block to the blockchain
	bc.validators[block.GetValidatorString()] -= 1
	bc.blocks = append(bc.blocks, block)
	return true, nil
}

// verifies the validitity of the chain
// basically rebuilds a new chain and checks if it hits any errors
func (bc *blockchain_t) Verify() (bool, error) {
	var nc blockchain_t
	nc.Init("verification_chain")
	for i, blk := range bc.blocks {
		if i > 0 {
			_, err := nc.AppendBlock(blk)
			if err != nil {
				return false, err
			}
		}
	}
	return true, nil
}

// saves all of the blocks in the blockchain to files corresponding to their hashes
func (bc *blockchain_t) SaveBlocks() error {
	for _, blk := range bc.blocks {
		_, err := blk.Save()
		if err != nil {
			return err
		}
	}
	return nil
}

// returns the block at the tip of the blockchain
func (bc *blockchain_t) GetTip() block_t {
	return bc.blocks[len(bc.blocks)-1]
}

// gets the hash of the block at the tip of the blockchain
func (bc *blockchain_t) GetTipHash() []byte {
	return bc.blocks[len(bc.blocks)-1].GetHash()
}

// returns the genesis block
func (bc *blockchain_t) GetGenesisBlock() block_t {
	return bc.blocks[0]
}

// returns the hash of the genesis block
func (bc *blockchain_t) GetGenesisHash() []byte {
	return bc.blocks[0].GetHash()
}

// saves the entire blockchain
// first verifies the chain
// if verified, then all of the blocks are saved to files
// also saves a file that contains the hashes of the tip and genesis blocks
func (bc *blockchain_t) Save() error {

	_, err := bc.Verify() // only save the chain if it's valid
	if err != nil {
		panic(err)
	}

	bc.SaveBlocks() // to reconstruct chain later

	if _, err := os.Stat(BLOCKCHAIN_DIR); os.IsNotExist(err) {
		err := os.Mkdir(BLOCKCHAIN_DIR, os.ModeDir)
		if err != nil {
			return err
		}
	}

	d := make(map[string][]byte)
	d["genesis"] = bc.GetGenesisHash()
	d["tip"] = bc.GetTipHash()
	data, err := json.MarshalIndent(d, "", "\t")
	if err != nil {
		return err
	}

	fname := path.Join(BLOCKCHAIN_DIR, bc.label+".json")
	return os.WriteFile(fname, data, 0777)

}

// prints every block in the blockchain, along with validator information
func (bc *blockchain_t) Print() {
	fmt.Println("Blocks")
	for i, blk := range bc.blocks {
		fmt.Printf("%d ", i)
		blk.Print()
		fmt.Println("")
	}

	fmt.Println("Validators")
	for v, c := range bc.validators {
		fmt.Printf("%d %s\r\n", c, v)
	}
}

// loads the a chain from a file
func LoadChain(label string) (bc blockchain_t, err error) {

	fname := path.Join(BLOCKCHAIN_DIR, label+".json")
	data, err := os.ReadFile(fname)
	if err != nil {
		return bc, err
	}

	d := make(map[string][]byte)
	err = json.Unmarshal(data, &d)
	if err != nil {
		return bc, err
	}

	var current_hash, genesis_hash []byte
	current_hash = d["tip"]
	genesis_hash = d["genesis"]
	// load tip
	block, err := LoadBlock(current_hash)
	if err != nil {
		return bc, err
	}
	bc.blocks = append(bc.blocks, block)

	for !bytes.Equal(bc.GetTipHash(), genesis_hash) {
		// set current_hash to prev_hash and load next block
		last_block := bc.GetTip()
		block, err := LoadBlock(last_block.prev_hash[:])
		if err != nil {
			return bc, err
		}
		bc.blocks = append(bc.blocks, block)
	}

	i := 0
	j := len(bc.blocks) - 1
	for i < j {
		bc.blocks[i], bc.blocks[j] = bc.blocks[j], bc.blocks[i]
		i += 1
		j -= 1
	}

	bc.label = label
	bc.validators = make(map[string]uint32)
	bc.validators[bc.blocks[0].GetValidatorString()] = 4294967295

	_, err = bc.Verify()
	return bc, err
}

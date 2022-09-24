package main

import (
	"bytes"
	"errors"
	"fmt"
)

type blockchain_t struct {
	blocks []block_t
}

func (bc *blockchain_t) Init() {
	bc.blocks = append(bc.blocks, Genesis())
}

// appends a block to the chain if the block is valid
func (bc *blockchain_t) AppendBlock(block block_t) (bool, error) {
	_, err := block.Verify()
	if err != nil {
		return false, err
	}

	if !bytes.Equal(block.prev_hash, bc.blocks[len(bc.blocks)-1].ComputeBlockHash()) {
		return false, errors.New("candidate block doesn't contain hash of previous block")
	}

	bc.blocks = append(bc.blocks, block)
	return true, nil
}

// verifies hash of previous block matches the hash in the next block
func (bc *blockchain_t) VerifyHashes() error {
	for i := range bc.blocks {
		if i == 0 {
			continue
		}

		if !bytes.Equal(bc.blocks[i].prev_hash, bc.blocks[i-1].ComputeBlockHash()) {
			return errors.New("block hashes don't form a valid chain")
		}
	}

	return nil
}

// verifies that each block is valid
func (bc *blockchain_t) VerifyBlocks() error {
	for _, blk := range bc.blocks {
		_, err := blk.Verify()
		if err != nil {
			return err
		}
	}
	return nil
}

// verifies the validitity of the chain
func (bc *blockchain_t) Verify() (bool, error) {
	err := bc.VerifyHashes()
	if err != nil {
		return false, err
	}

	err = bc.VerifyBlocks()
	if err != nil {
		return false, err
	}

	return true, nil
}

func (bc *blockchain_t) SaveBlocks() error {
	for _, blk := range bc.blocks {
		_, err := blk.Save()
		if err != nil {
			return err
		}
	}
	return nil
}

// verifies the validitity of the chain
func (bc *blockchain_t) GetLastBlock() block_t {
	return bc.blocks[len(bc.blocks)-1]
}

// func ReplaceChain(oldchain *blockchain_t, newchain *blockchain_t) *blockchain_t {
// 	if len(newchain.blocks) <= len(oldchain.blocks) {
// 		fmt.Println("Length of new chain is less or equal to the length of the old chain.")
// 		return oldchain
// 	} else if !newchain.IsValidChain() {
// 		fmt.Println("New chain is invalid.")
// 		return oldchain
// 	} else {
// 		return newchain
// 	}
// }

func (bc *blockchain_t) Print() {
	for i, blk := range bc.blocks {
		fmt.Printf("%d ", i)
		blk.Print()
		fmt.Println("")
	}
}

func BlockchainTest() {

	id := LoadIdentity("main")

	var bc1 blockchain_t
	bc1.Init()

	blk1 := NewBlock(GetCurrentTimestamp(), bc1.GetLastBlock().hash, []byte("hello!"), id)
	_, err := bc1.AppendBlock(blk1)
	if err != nil {
		panic(err)
	}

	_, err = bc1.AppendBlock(blk1) // this should throw an error because prev_hash won't match
	if err == nil {
		panic(errors.New("this should be invalid"))
	}

	blk2 := NewBlock(GetCurrentTimestamp(), bc1.GetLastBlock().hash, []byte("hey!"), id)
	_, err = bc1.AppendBlock(blk2)
	if err != nil {
		panic(err)
	}

	bc1.Print()

	_, err = bc1.Verify()
	if err != nil {
		panic(err)
	}

	fmt.Println("Chain is valid!")
	bc1.SaveBlocks()

	// var bc2 blockchain_t
	// bc2.Init()
	// bc2.AddBlock([]byte("Sup?"))
	// bc = *ReplaceChain(&bc, &bc2)

	// bc2.AddBlock([]byte(""))
	// bc2.AddBlock([]byte("aaaa"))
	// bc2.AddBlock([]byte("aaaa"))
	// bc = *ReplaceChain(&bc, &bc2)

	// if bc.IsValidChain() {
	// 	bc.Print()
	// } else {
	// 	fmt.Println("Warning! Invalid Chain.")
	// 	bc.Print()
	// }
}

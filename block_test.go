package main

import (
	"testing"
)

func TestBlockCreation(t *testing.T) {
	id := LoadIdentity("main")

	block0 := Genesis()

	_, err := block0.Verify()
	if err != nil {
		t.Errorf("genesis block verification failed (%s)", err)
	}

	tx1 := NewTx_Entry([]byte("I added a block!"))

	block1, err := NewBlock(GetCurrentTimestamp(), block0.hash[:], tx1, id)
	if err != nil {
		t.Errorf("verification failed (%s)", err)
	}

	tx2 := NewTx_Entry([]byte("Yet another block."))
	for i := 0; i < 2; i++ {
		_, err := NewBlock(GetCurrentTimestamp(), block1.hash[:], tx2, id)
		if err != nil {
			t.Errorf("verification failed (%s)", err)
		}
	}
}

func TestBlockSaving(t *testing.T) {
	id := LoadIdentity("test")

	block0 := Genesis()

	tx1 := NewTx_Entry([]byte("Saved."))
	block1, err := NewBlock(GetCurrentTimestamp(), block0.hash[:], tx1, id)
	if err != nil {
		t.Errorf("error creating block (%s)", err)
	}
	_, err = block1.Verify()
	if err != nil {
		t.Errorf("error verifying block (%s)", err)
	}
	_, err = block1.Save()
	if err != nil {
		t.Errorf("error saving block (%s)", err)
	}
	block_loaded, err := LoadBlock(block1.GetHash())
	if err != nil {
		block1.Print()
		block_loaded.Print()
		t.Errorf("error loading block %x (%s)", block1.GetHash(), err)
	}
	_, err = block_loaded.Verify()
	if err != nil {
		t.Errorf("verification failed after loading block %x (%s)", block1.GetHash(), err)
	}

}

package main

import (
	"testing"
)

func TestBlockchain_NewChain(t *testing.T) {
	id := LoadIdentity("main")

	var bc1 blockchain_t
	bc1.Init("test")

	tx1 := NewTx_Entry([]byte("hello!"))
	blk1, err := NewBlock(GetCurrentTimestamp(), bc1.GetTipHash(), tx1, id)
	if err != nil {
		bc1.Print()
		t.Errorf("error creating first block (%s)", err)
	}

	_, err = bc1.AppendBlock(blk1)
	if err != nil {
		bc1.Print()
		t.Errorf("error appending first block (%s)", err)
	}

	_, err = bc1.AppendBlock(blk1) // this should throw an error because prev_hash won't match
	if err == nil {
		bc1.Print()
		t.Errorf("error appending this block (should be invalid) (%s)", err)
	}

	tx2 := NewTx_Entry([]byte("hey."))
	blk2, err := NewBlock(GetCurrentTimestamp(), bc1.GetTipHash(), tx2, id)
	if err != nil {
		bc1.Print()
		t.Errorf("error creating second block (%s)", err)
	}

	_, err = bc1.AppendBlock(blk2)
	if err != nil {
		bc1.Print()
		t.Errorf("error appending second block (%s)", err)
	}

	_, err = bc1.Verify()
	if err != nil {
		bc1.Print()
		t.Errorf("error verifying blockchain (%s)", err)
	}

	err = bc1.Save()
	if err != nil {
		bc1.Print()
		t.Errorf("error saving blockchain (%s)", err)
	}
}

func TestBlockchainDelegation(t *testing.T) {
	id_main := LoadIdentity("main")
	id_bar := LoadIdentity("bar")
	id_foo := LoadIdentity("foo")

	var bc1 blockchain_t
	bc1.Init("bar")

	tx1 := NewTx_Permission(100, id_bar.GetPubBytes())
	blk1, err := NewBlock(GetCurrentTimestamp(), bc1.GetTipHash(), tx1, id_main)
	if err != nil {
		bc1.Print()
		t.Errorf("error creating first block (%s)", err)
	}
	_, err = bc1.AppendBlock(blk1)
	if err != nil {
		bc1.Print()
		t.Errorf("error appending first block to blockchain (%s)", err)
	}

	tx2 := NewTx_Permission(200, id_foo.GetPubBytes())
	blk2i, err := NewBlock(GetCurrentTimestamp(), bc1.GetTipHash(), tx2, id_bar)
	if err != nil {
		bc1.Print()
		t.Errorf("error creating this invalid block (%s)", err)
	}
	_, err = bc1.AppendBlock(blk2i)
	if err == nil {
		bc1.Print()
		t.Errorf("this should result in an error because bar can't delegate 200 blocks, but you allowed it to happen")
	}

	blk2, err := NewBlock(GetCurrentTimestamp(), bc1.GetTipHash(), tx2, id_main)
	if err != nil {
		bc1.Print()
		t.Errorf("error creating second block (%s)", err)
	}
	_, err = bc1.AppendBlock(blk2)
	if err != nil {
		bc1.Print()
		t.Errorf("error appending second block to blockchain (%s)", err)
	}

	tx3 := NewTx_Permission(50, id_foo.GetPubBytes())
	blk3, err := NewBlock(GetCurrentTimestamp(), bc1.GetTipHash(), tx3, id_bar)
	if err != nil {
		bc1.Print()
		t.Errorf("error creating third block (%s)", err)
	}
	_, err = bc1.AppendBlock(blk3)
	if err != nil {
		bc1.Print()
		t.Errorf("error appending second block to blockchain (%s)", err)
	}

}

func TestBlockchainLoad(t *testing.T) {
	_, err := LoadChain("test")
	if err != nil {
		t.Errorf("error loading chain (%s)", err)
	}
}

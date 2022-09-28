package main

import (
	"encoding/binary"
	"errors"
	"fmt"
)

type txtype_t byte

const (
	Entry      txtype_t = 0 // add arbitrary data to the chain
	Permission txtype_t = 1 // allows others to add to blockchain
)

type transaction_t struct {
	txtype txtype_t
	data   []byte
}

func (tx *transaction_t) Marshal() (b []byte) {
	b = append(b, byte(Entry))
	b = append(b, tx.data...)
	return b
}

func (tx *transaction_t) Unmarshal(b []byte) error {
	tx.txtype = txtype_t(b[0])
	tx.data = b[1:]
	if len(tx.data) != len(b)-1 {
		return errors.New("unmarshaling unsuccessful")
	}
	return nil
}

func NewTx_Entry(data []byte) (tx transaction_t) {
	tx.txtype = Entry
	tx.data = data
	return tx
}

func NewTx_Permission(n uint32, pubKey []byte) (tx transaction_t) {
	tx.txtype = Permission
	b := make([]byte, 4)
	binary.BigEndian.PutUint32(b, n)
	tx.data = append(tx.data, b...)
	tx.data = append(tx.data, pubKey...)
	return tx
}

func (tx *transaction_t) ParseTx_Permission() (n uint32, validator []byte, err error) {
	if tx.txtype != Permission {
		return 0, []byte(""), errors.New("not a permission transaction")
	}
	return binary.BigEndian.Uint32(tx.data[0:4]), tx.data[4:], nil
}

func TestTransaction() {
	tx := NewTx_Entry([]byte("Bryan"))
	fmt.Printf("%x", tx.Marshal())
}

package main

import (
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
)

func StartServer() {
	fname := path.Join(BLOCKCHAIN_DIR, MAIN_CHAIN_NAME+".json")
	if _, err := os.Stat(fname); errors.Is(err, os.ErrNotExist) {
		blockchain.Init(MAIN_CHAIN_NAME)
		blockchain.Save()
	} else {
		bc, err := LoadChain(MAIN_CHAIN_NAME)
		if err != nil {
			panic(err)
		}
		blockchain = bc
	}

	http.HandleFunc("/", welcome)
	http.HandleFunc("/headers", headers)
	http.HandleFunc("/tip", tip)
	http.HandleFunc("/submit", submit)

	http.ListenAndServe(":8090", nil)
}

func printReqInfo(r *http.Request) {
	fmt.Printf("server: %s /\n", r.Method)
	fmt.Printf("server: query id: %s\n", r.URL.Query().Get("id"))
	fmt.Printf("server: content-type: %s\n", r.Header.Get("content-type"))
	fmt.Printf("server: headers:\n")
	for headerName, headerValue := range r.Header {
		fmt.Printf("\t%s = %s\n", headerName, strings.Join(headerValue, ", "))
	}

	// reqBody, err := ioutil.ReadAll(r.Body)
	// if err != nil {
	// 	fmt.Printf("server: could not read request body: %s\n", err)
	// }
	// fmt.Printf("server: request body: %s\n", reqBody)
}

func welcome(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "Welcome  to my permissioned blockchain!\n")
}

func headers(w http.ResponseWriter, req *http.Request) {
	for name, headers := range req.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v\n", name, h)
		}
	}
}

func tip(w http.ResponseWriter, req *http.Request) {
	block := blockchain.GetTip()
	fmt.Fprintf(w, "%x\n", block.GetHash())
}

func submit(w http.ResponseWriter, req *http.Request) {

	var block_data []byte
	var b []byte
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Received data: %s\r\n", b)

	block_data, err = hex.DecodeString(string(b))
	if err != nil {
		panic(err)
	}

	_, err = Unmarshal(block_data)
	if err != nil {
		fmt.Fprintf(w, "Invalid block.\n")
	} else {
		fmt.Fprintf(w, "Valid block.\n")
	}
}

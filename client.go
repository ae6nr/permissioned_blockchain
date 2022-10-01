package main

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func SubmitEntry(server_url string, entry []byte, id identity_t) error {

	// get the tip
	resp, err := http.Get(server_url + "/tip")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var tip_data, tip_hash []byte
	resp.Body.Read(tip_data)
	_, err = hex.Decode(tip_hash, tip_data)
	if err != nil {
		return err
	}

	// create the block
	tx := NewTx_Entry([]byte(entry))
	block, err := NewBlock(GetCurrentTimestamp(), tip_hash, tx, id)
	if err != nil {
		return err
	}
	block.Print()

	// submit the block
	requestUrl := server_url + "/submit"

	client := http.Client{
		Timeout: 10 * time.Second,
	}

	s := hex.EncodeToString(block.Marshal())
	req, err := http.NewRequest(http.MethodPost, requestUrl, strings.NewReader(s))
	if err != nil {
		return err
	}

	resp2, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp2.Body.Close()

	post_resp, err := ioutil.ReadAll(resp2.Body)
	if err != nil {
		return err
	}

	fmt.Printf("Response (%d): %s\r\n", len(post_resp), post_resp)

	return nil
}

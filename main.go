package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

func usage(message string) {
	fmt.Fprintln(os.Stderr, message)
	fmt.Fprintln(os.Stderr)
	fmt.Fprintf(os.Stderr, "Usage: %s <url> <username> <password>\n", os.Args[0])
	os.Exit(1)
}

func main() {
	if len(os.Args) < 4 {
		usage("Not enough args")
	}

	var urlString, user, pass = os.Args[1], os.Args[2], os.Args[3]

	var u, err = url.Parse(urlString)
	if err != nil {
		usage(fmt.Sprintf("Invalid URL %q: %s", urlString, err))
	}

	u.User = url.UserPassword(user, pass)

	for {
		run(u)
	}
}

func run(u *url.URL) {
	var data = bytes.NewBufferString(`{"id":0,"method":"getblocktemplate","params":[{"rules":["segwit"]}],"Header":null}`)
	var err = checkBlockHeight(u, data)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to POST to URL %q: %s", u.String(), err)
		time.Sleep(time.Second)
		return
	}

	time.Sleep(time.Millisecond * 100)
}

var lastHeight int64

func checkBlockHeight(u *url.URL, data io.Reader) error {
	var r, err = http.Post(u.String(), "text/plain", data)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	var body []byte
	body, err = io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	var resp map[string]interface{}
	err = json.Unmarshal(body, &resp)
	var height = int64(resp["result"].(map[string]interface{})["height"].(float64))
	if height == lastHeight {
		return nil
	}

	lastHeight = height
	log.Printf("New block template %d received; writing files", height)
	var fheight *os.File
	fheight, err = os.Create("blockheight")
	if err != nil {
		return err
	}
	defer fheight.Close()
	_, err = fmt.Fprintf(fheight, "%d", height)
	if err != nil {
		return err
	}

	var fjson *os.File
	fjson, err = os.Create("getblocktemplate.json")
	if err != nil {
		return err
	}
	defer fjson.Close()
	_, err = io.Copy(fjson, bytes.NewBuffer(body))
	return err
}

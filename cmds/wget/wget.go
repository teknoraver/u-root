// Copyright 2012-2017 the u-root Authors. All rights reserved
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Wget reads one file from a url and writes to stdout.
//
// Synopsis:
//     wget URL
//
// Description:
//     Returns a non-zero code on failure.
//
// Notes:
//     There are a few differences with GNU wget:
//     - Upon error, the return value is always 1.
//     - The protocol (http/https) is mandatory.
//
// Example:
//     wget http://google.com/ > e100.html
package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
)

func wget(arg string, w io.Writer) error {
	resp, err := http.Get(arg)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return fmt.Errorf("non-200 HTTP status: %d", resp.StatusCode)
	}
	_, err = io.Copy(w, resp.Body)
	return err
}

func usage() {
	log.Printf("Usage: %s [ARGS] URL\n", os.Args[0])
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Parse()
	if flag.NArg() != 1 {
		usage()
	}

	argURL := flag.Arg(0)
	if argURL == "" {
		log.Fatalln("Empty URL")
	}

	url, err := url.Parse(argURL)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	fileName := "index.html"
	if url.Path != "" && url.Path[len(url.Path)-1] != '/' {
		fileName = path.Base(url.Path)
	}

	file, err := os.Create(fileName)
	if err != nil {
		log.Fatalf("%v\n", err)
	}

	if err := wget(argURL, file); err != nil {
		log.Fatalf("%v\n", err)
	}
}

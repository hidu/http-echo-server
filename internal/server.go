// Copyright(C) 2024 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2024/5/24

package internal

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

var addr = flag.String("l", ":8088", "http server listen at")
var resp = flag.String("resp", "", "default response")
var contentType = flag.String("ct", "", "default response content type")
var logHeader = flag.Bool("lh", false, "log with http headers")
var statusCode = flag.Int("status", 0, "default response status code")

func Run() {
	flag.Parse()
	log.Println("http-echo-server listen at", *addr)
	router()
	err := http.ListenAndServe(*addr, nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func router() {
	http.HandleFunc("/", logRequest(index))
	http.HandleFunc("/help", logRequest(help))
	http.HandleFunc("/chunk", logRequest(chunk))
	http.HandleFunc("/status", logRequest(status))
	http.HandleFunc("/cal/sum", logRequest(sum))
}

func init() {
	d := flag.Usage
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "http echo server")
		d()
		fmt.Fprintln(os.Stderr, "")
		fmt.Fprintln(os.Stderr, "site:", "https://github.com/hidu/http-echo-server")
		fmt.Fprintln(os.Stderr, helpMsg)
	}
}

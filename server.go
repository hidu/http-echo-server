/**
*for http network test
 */
package main

import (
	"encoding/json"
	"encoding/xml"
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"
)

var addr = flag.String("addr", ":8088", "http server listen at")
var resp = flag.String("resp", "", "default response")
var contentType = flag.String("content_type", "", "default response content type")

func main() {
	flag.Parse()
	http.HandleFunc("/", HelloServer)
	http.HandleFunc("/help", HelpServer)
	fmt.Println("start http server at:", *addr)
	err := http.ListenAndServe(*addr, nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

var reqID uint64

type resData struct {
	ID      uint64 `json:"id"`
	Request string `json:"request'`
}
type Datas struct {
	ResData []*resData
}

func HelloServer(w http.ResponseWriter, req *http.Request) {
	start := time.Now()

	id := atomic.AddUint64(&reqID, 1)

	defer func() {
		_used := fmt.Sprintf("%.4f", time.Now().Sub(start).Seconds()*1000)
		log.Println(id, req.Method, req.URL.RequestURI(), _used)
	}()

	item := new(resData)
	item.ID = id
	dump, err := httputil.DumpRequest(req, true)

	if err != nil {
		item.Request = "error:" + err.Error()
	} else {
		item.Request = string(dump)
	}

	sleep := getIntVal(req, "sleep")
	if sleep > 0 {
		time.Sleep(time.Duration(sleep) * time.Millisecond)
	}

	httpCode := getIntVal(req, "http_code")
	if httpCode > 0 {
		w.WriteHeader(httpCode)
	}
	reqContentType := req.FormValue("content_type")

	repeatNum := getIntVal(req, "repeat")
	if repeatNum == 0 {
		repeatNum = 1
	}

	datas := new(Datas)

	for i := 0; i < repeatNum; i++ {
		datas.ResData = append(datas.ResData, item)
	}

	if req.FormValue("broken") != "" {
		hj, ok := w.(http.Hijacker)
		if ok {
			conn, _, _ := hj.Hijack()
			conn.Close()
			return
		}
	}

	var dataBf []byte

	rt := req.FormValue("type")
	if rt == "" {
		rt = *contentType
	}

	switch rt {
	case "json":
		reqContentType = "application/json"
		dataBf, _ = json.MarshalIndent(datas, "", " ")
	case "xml":
		reqContentType = "text/xml"
		dataBf, _ = xml.MarshalIndent(datas, "", " ")
	default:
		dataBf = []byte(fmt.Sprintf("%q", datas))
	}

	if *resp != "" {
		dataBf = []byte(*resp)
	}

	if reqContentType != "" {
		w.Header().Set("Content-Type", reqContentType)
	}

	w.Write(dataBf)

}

func getIntVal(req *http.Request, key string) int {
	val := req.FormValue(key)
	if val == "" {
		return 0
	}
	n, err := strconv.Atoi(val)
	if err != nil {
		return 0
	}
	return n
}

var helpMsg = `
query/form params:
	sleep        : sleep ms, eg: sleep=100
	http_code    : http status code, eg: http_code=500
	content_type : content type, eg: content_type=text/html;charset=utf-8
	repeat       : repeat content times, eg: repeat=10
	broken       : broken this connect, eg: broken=1
	type         : data output type, allow: [json,xml], eg: type=json

visit url example:
	http://{host}/?sleep=100
	http://{host}/?sleep=100&http_code=500&repeat=1
	`

func HelpServer(w http.ResponseWriter, req *http.Request) {
	help := strings.Replace(helpMsg, "{host}", req.Host, -1)
	w.Write([]byte(help))
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

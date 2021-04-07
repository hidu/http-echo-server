package main

import (
	"bytes"
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
var logHeader = flag.Bool("log_header", false, "")

func main() {
	flag.Parse()
	http.HandleFunc("/", logRequest(index))
	http.HandleFunc("/help", logRequest(help))
	http.HandleFunc("/status", logRequest(status))
	http.HandleFunc("/cal/sum", logRequest(sum))
	fmt.Println("start http server at:", *addr)
	err := http.ListenAndServe(*addr, nil)

	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func logRequest(h http.HandlerFunc) func(w http.ResponseWriter, req *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		num := atomic.AddInt64(&connecting, 1)
		go func() {
			<-req.Context().Done()
			atomic.AddInt64(&connecting, -1)
		}()

		defer func() {
			cost := fmt.Sprintf("%.4f", time.Now().Sub(start).Seconds()*1000)
			fs := []interface{}{
				num, req.RemoteAddr, req.Method, req.RequestURI, cost,
			}
			if *logHeader {
				fs = append(fs, req.Header)
			}
			log.Println(fs...)
		}()

		h(w, req)
	}
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
	http://{host}/cal/sum?ids=123,456
	`

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

var reqID uint64

type resData struct {
	ID      uint64 `json:"id"`
	Request string `json:"request'`
}
type Datas struct {
	ResData []*resData
}

func (d *Datas) Bytes() []byte {
	var buf bytes.Buffer
	for _, item := range d.ResData {
		buf.WriteString(fmt.Sprintf("ID=%d\n\n", item.ID))
		buf.WriteString(item.Request)
		buf.WriteString("\n\n\n")
	}
	return buf.Bytes()
}

var connecting int64

func index(w http.ResponseWriter, req *http.Request) {
	id := atomic.AddUint64(&reqID, 1)
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
		reqContentType = "text/plain"
		dataBf = datas.Bytes()
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

func help(w http.ResponseWriter, req *http.Request) {
	help := strings.Replace(helpMsg, "{host}", req.Host, -1)
	w.Write([]byte(help))
}

func status(w http.ResponseWriter, req *http.Request) {
	str := fmt.Sprintf("connecting=%d", atomic.LoadInt64(&connecting))
	w.Write([]byte(str))
}

type sumResult struct {
	ErrNo int
	Msg   string
	Data  struct {
		Sum int
	}
}

func (sr *sumResult) Bytes() []byte {
	bf, _ := json.Marshal(sr)
	return bf
}

func sum(w http.ResponseWriter, req *http.Request) {
	ids := req.URL.Query().Get("ids")
	if ids == "" {
		ret := &sumResult{
			ErrNo: 400,
			Msg:   "ids empty",
		}
		w.WriteHeader(400)
		w.Write(ret.Bytes())
		return
	}
	var sum int
	for i, idStr := range strings.Split(ids, ",") {
		id, err := strconv.Atoi(idStr)
		if err != nil {
			ret := &sumResult{
				ErrNo: 400,
				Msg:   fmt.Sprintf("ids[%d]=%q not int", i, idStr),
			}
			w.WriteHeader(400)
			w.Write(ret.Bytes())
			return
		}
		sum += id
	}
	ret := &sumResult{
		ErrNo: 0,
		Msg:   "success",
		Data: struct{ Sum int }{
			Sum: sum,
		},
	}
	w.Write(ret.Bytes())
}

// Copyright(C) 2024 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2024/5/24

package internal

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"net/http"
	"net/http/httputil"
	"time"
)

func index(w http.ResponseWriter, req *http.Request) {
	item := new(resData)
	item.ID = req.Header.Get("rid")
	dump, err := httputil.DumpRequest(req, true)

	if err != nil {
		item.Request = "error:" + err.Error()
	} else {
		item.Request = string(dump)
	}

	sleep := getIntVal(req, "sleep")
	if sleep > 0 {
		if !ctxSleep(req.Context(), time.Duration(sleep)*time.Millisecond) {
			return
		}
	}

	httpCode := getIntVal(req, "http_code")
	if httpCode == 0 {
		httpCode = *statusCode
	}
	if httpCode > 0 {
		addLogKV(req.Context(), "RespStatus", httpCode)
		w.WriteHeader(httpCode)
	}
	reqContentType := req.FormValue("content_type")

	repeatNum := getIntVal(req, "repeat")
	if repeatNum == 0 {
		repeatNum = 1
	}

	if req.FormValue("broken") != "" {
		hj, ok := w.(http.Hijacker)
		if ok {
			conn, _, _ := hj.Hijack()
			conn.Close()
			return
		}
	}

	rt := req.FormValue("type")
	if rt == "" {
		rt = *contentType
	}

	var content []byte
	if *resp == "" {
		datas := new(Datas)
		for i := 0; i < repeatNum; i++ {
			datas.ResData = append(datas.ResData, item)
		}
		switch rt {
		case "json":
			reqContentType = "application/json"
			content, _ = json.MarshalIndent(datas, "", " ")
		case "xml":
			reqContentType = "text/xml"
			content, _ = xml.MarshalIndent(datas, "", " ")
		default:
			reqContentType = "text/plain"
			content = datas.Bytes()
		}
	} else {
		content = []byte(*resp)
	}

	if reqContentType != "" {
		w.Header().Set("Content-Type", reqContentType)
	}

	w.Write(content[:1])
	sleepAh := getIntVal(req, "sleep_ah")
	if sleepAh > 0 {
		ctxSleep(req.Context(), time.Duration(sleepAh)*time.Millisecond)
	}
	w.Write(content[1:])
}

type resData struct {
	ID      string
	Request string
}

type Datas struct {
	ResData []*resData
}

func (d *Datas) Bytes() []byte {
	var buf bytes.Buffer
	for _, item := range d.ResData {
		buf.WriteString(fmt.Sprintf("ID=%s\n\n", item.ID))
		buf.WriteString(item.Request)
		buf.WriteString("\n\n\n")
	}
	return buf.Bytes()
}

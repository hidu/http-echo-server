// Copyright(C) 2024 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2024/5/24

package internal

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

func chunk(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	httpCode := getIntVal(req, "http_code")
	if httpCode == 0 {
		httpCode = *statusCode
	}
	if httpCode > 0 {
		w.WriteHeader(httpCode)
	}
	repeat := getIntVal(req, "repeat")
	if repeat == 0 {
		repeat = 100
	}
	hf := w.(http.Flusher)
	for i := 0; i < repeat; i++ {
		fmt.Fprintln(w, "Hello", i)
		hf.Flush()
		if !ctxSleep(req.Context(), time.Second) {
			return
		}
	}
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

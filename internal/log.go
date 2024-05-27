// Copyright(C) 2024 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2024/5/24

package internal

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"slices"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

var remoteAddrs sync.Map
var requestID atomic.Int64
var connecting atomic.Int64

func logRequest(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		start := time.Now()
		num := connecting.Add(1)
		id := requestID.Add(1)

		logKVStart := map[string]any{
			"ID":         id,
			"Connecting": num,
			"LogID":      getLogID(req),
			"Remote":     req.RemoteAddr,
			"Method":     req.Method,
			"URI":        req.RequestURI,
		}
		req = req.WithContext(withLogContext(req.Context(), logKVStart))

		remoteAddrs.Store(req.RemoteAddr, time.Now())

		req.Header.Set("rid", strconv.FormatInt(id, 10))

		defer func() {
			connecting.Add(-1)
			remoteAddrs.Delete(req.RemoteAddr)

			cost := fmt.Sprintf("%.4f", time.Since(start).Seconds()*1000)
			logKVEnd := map[string]any{
				"Cost":   cost,
				"CtxErr": req.Context().Err(),
			}
			if *logHeader {
				logKVEnd["Headers"] = req.Header
			}
			printLog(req.Context(), logKVEnd, "Request Done")
		}()

		handler(w, req)
	}
}

func getLogID(req *http.Request) string {
	if id := req.URL.Query().Get("logid"); len(id) > 0 {
		return id
	}
	if id := req.Header.Get("X_BD_LOGID"); len(id) > 0 {
		return id
	}
	return ""
}

type ctxKey int8

const ctxKeyLog ctxKey = iota

func withLogContext(ctx context.Context, kv map[string]any) context.Context {
	lf := &logFields{}
	lf.AddMap(kv)
	return context.WithValue(ctx, ctxKeyLog, lf)
}

func addLogKV(ctx context.Context, key string, value any) {
	ctx.Value(ctxKeyLog).(*logFields).Add(key, value)
}

func addLogKVs(ctx context.Context, kv map[string]any) {
	ctx.Value(ctxKeyLog).(*logFields).AddMap(kv)
}

func printLog(ctx context.Context, kv map[string]any, msg string) {
	ctx.Value(ctxKeyLog).(*logFields).print(kv, msg)
}

type logFields struct {
	fields []keyValue
	mux    sync.RWMutex
}

func (lf *logFields) Add(key string, value any) {
	lf.mux.Lock()
	lf.fields = append(lf.fields, keyValue{K: key, V: value})
	lf.mux.Unlock()
}

func (lf *logFields) AddMap(kv map[string]any) {
	lf.mux.Lock()
	for k, v := range kv {
		lf.fields = append(lf.fields, keyValue{K: k, V: v})
	}
	lf.mux.Unlock()
}

func (lf *logFields) print(kv map[string]any, msg string) {
	lf.mux.RLock()
	fields := slices.Clone(lf.fields)
	lf.mux.RUnlock()
	for k, v := range kv {
		fields = append(fields, keyValue{K: k, V: v})
	}
	data := make([]any, 0, len(fields)+1)
	data = append(data, fmt.Sprintf("%q", msg))
	for _, item := range fields {
		data = append(data, fmt.Sprintf("%s=%#v", item.K, item.V))
	}
	log.Println(data...)
}

type keyValue struct {
	K string
	V any
}

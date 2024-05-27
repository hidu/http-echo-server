// Copyright(C) 2024 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2024/5/24

package internal

import (
	"encoding/json"
	"net/http"
	"time"
)

func status(w http.ResponseWriter, req *http.Request) {
	addrs := make(map[string]string, 10)
	remoteAddrs.Range(func(key, value any) bool {
		start := value.(time.Time)
		addrs[key.(string)] = time.Since(start).String()
		return true
	})
	data := map[string]any{
		"Connecting":   connecting.Load(),
		"RequestID":    requestID.Load(),
		"RemoteTotal":  len(addrs),
		"RemoteDetail": addrs,
	}
	bf, _ := json.Marshal(data)
	w.Header().Set("Content-Type", "application/json")
	w.Write(bf)
}

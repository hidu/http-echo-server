// Copyright(C) 2024 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2024/5/24

package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

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
	var total int
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
		total += id
	}
	ret := &sumResult{
		ErrNo: 0,
		Msg:   "success",
		Data: struct{ Sum int }{
			Sum: total,
		},
	}
	w.Write(ret.Bytes())
}

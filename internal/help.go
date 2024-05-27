// Copyright(C) 2024 github.com/fsgo  All Rights Reserved.
// Author: hidu <duv123@gmail.com>
// Date: 2024/5/24

package internal

import (
	"net/http"
	"strings"
)

const helpMsg = `
query/form params:
	sleep        : sleep n ms before response header, eg: sleep=100
	sleep_ah     : sleep n ms after response header before body, eg: sleep_ah=100
	http_code    : http status code, eg: http_code=500
	content_type : content type, eg: content_type=text/html;charset=utf-8
	repeat       : repeat content times, eg: repeat=10
	broken       : broken this connect, eg: broken=1
	type         : data output type, allow: [json,xml], eg: type=json

visit url:
	http://{host}/?sleep=100
	http://{host}/?sleep=100&http_code=500&repeat=1
	http://{host}/cal/sum?ids=123,456
	http://{host}/chunk?&http_code=500&repeat=1
	http://{host}/status
	`

func help(w http.ResponseWriter, req *http.Request) {
	str := strings.ReplaceAll(helpMsg, "{host}", req.Host)
	w.Write([]byte(str))
}

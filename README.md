# http-echo-server

you can use it for test your application。  
such as test 404、500、403 exceptions,
and also network broken exception and so on.  

## Install
```
go get -u github.com/hidu/http-echo-server
```

## Exec
```
http-echo-server -addr :8080
```

visit :
```
http://{host}/?sleep=100&http_code=500&repeat=1
```


## Params
```
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
```
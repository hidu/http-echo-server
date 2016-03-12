# http-echo-server

you can use it for test your application,such as test 404,500,403 exceptions.  
and also network broken exception and so on.  


    http params:
    
    sleep        : sleep ms,eg:100
    http_code    : http status code, eg:500
    content_type : content type, eg: text/html;chatset=utf-8
    repeat       : repeat content times, eg:10
    broken       : broken this connect,eg borken=1
    type         : data output type,allow:[json,xml] 
    
    eg:
    http://{host}/?sleep=100
    http://{host}/?sleep=100&http_code=500&repeat=1
    http://{host}/you/path/example?sleep=100&http_code=500&repeat=1&type=json
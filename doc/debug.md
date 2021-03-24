### 开启debug
```
package main

import (
    "log"
    "net/http"
    _ "net/http/pprof" //引入就好了
    "regexp"
)
#使用http启动获取pprof收集内容的接口
http.HandleFunc("/",nil)
_ := http.ListenAndServe(":9999", nil)
```


### 收集数据
```
go tool pprof http://localhost:9999/pprof/profile
```

### 将数据变成火焰图
```
go tool pprof -http=:1234 [步骤【收集数据】中生成的文件]，
```

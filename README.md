# MgCall Go微服务调用客户端

## 安装
```shell script
go get -u github.com/maczh/mgcall
```

## 使用说明
+ 在配置中必须指定服务注册中心类型 go.discovery，不指定为nacos
+ 在配置引用go.config.used中指定nacos或consul
+ 配置服务器中必须要有相应的nacos或consul配置文件

## 使用范例
```go
import (
    "github.com/maczh/mgcall"
    "github.com/maczh/logs"
    "github.com/maczh/utils"
    "github.com/maczh/gintool/result"
)
...

func AddDocument(database, table string, doc interface{}, searchFields []string) result.Result {
	params := make(map[string]string)
	params["database"] = database
	params["table"] = table
	params["doc"] = utils.ToJSON(doc)
	if searchFields != nil && len(searchFields) > 0 {
		params["searchFields"] = utils.ToJSON(searchFields)
	}
	resp, err := mgcall.Call("serviceName", "/doc/add", params)
	if err != nil {
		logs.Error("微服务{}{}调用异常:{}", "serviceName", "/doc/add", err.Error())
		return *result.Error(-1, err.Error())
	}
	var res result.Result
	utils.FromJSON(resp, &res)
	return res
}

```
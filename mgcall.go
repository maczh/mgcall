package mgcall

import (
	"errors"
	"github.com/levigross/grequests"
	"github.com/maczh/logs"
	"github.com/maczh/mgcache"
	config "github.com/maczh/mgconfig"
	"github.com/maczh/mgtrace"
	"time"
)

//微服务调用其他服务的接口
func Call(service string, uri string, params map[string]string) (string, error) {
	host := ""
	if mgcache.OnGetCache("nacos").IsExist(service) {
		h, _ := mgcache.OnGetCache("nacos").Value(service)
		host = h.(string)
	} else {
		discovery := config.GetConfigString("go.discovery")
		if discovery == "" {
			discovery = "nacos"
		}
		switch discovery {
		case "nacos":
			host = config.GetNacosServiceURL(service)
		case "consul":
			host = config.GetConsulServiceURL(service)
		}
		if host != "" {
			mgcache.OnGetCache("nacos").Add(service, host, 5*time.Minute)
		} else {
			return "", errors.New("微服务获取" + service + "服务主机IP端口失败")
		}
	}
	url := host + uri
	logs.Debug("Nacos微服务请求:{}\n请求参数:{}", url, params)
	header := map[string]string{"X-Request-Id": mgtrace.GetRequestId()}
	resp, err := grequests.Post(url, &grequests.RequestOptions{
		Data:    params,
		Headers: header,
	})
	logs.Debug("Nacos微服务返回结果:{}", resp.String())
	if err != nil {
		return "", err
	} else {
		return resp.String(), err
	}
}

//微服务调用其他服务的接口,带header
func CallWithHeader(service string, uri string, params map[string]string, header map[string]string) (string, error) {
	host := ""
	if mgcache.OnGetCache("nacos").IsExist(service) {
		h, _ := mgcache.OnGetCache("nacos").Value(service)
		host = h.(string)
	} else {
		discovery := config.GetConfigString("go.discovery")
		if discovery == "" {
			discovery = "nacos"
		}
		switch discovery {
		case "nacos":
			host = config.GetNacosServiceURL(service)
		case "consul":
			host = config.GetConsulServiceURL(service)
		}
		if host != "" {
			mgcache.OnGetCache("nacos").Add(service, host, 5*time.Minute)
		} else {
			return "", errors.New("微服务获取" + service + "服务主机IP端口失败")
		}
	}
	url := host + uri
	header["X-Request-Id"] = mgtrace.GetRequestId()
	logs.Debug("Nacos微服务请求:{}\n请求参数:{}\n请求头:{}", url, params, header)
	resp, err := grequests.Post(url, &grequests.RequestOptions{
		Data:    params,
		Headers: header,
	})
	logs.Debug("Nacos微服务返回结果:{}", resp.String())
	if err != nil {
		return "", err
	} else {
		return resp.String(), err
	}
}

//微服务调用其他服务的接口,带文件
func CallWithFiles(service string, uri string, params map[string]string, files []grequests.FileUpload) (string, error) {
	host := ""
	if mgcache.OnGetCache("nacos").IsExist(service) {
		h, _ := mgcache.OnGetCache("nacos").Value(service)
		host = h.(string)
	} else {
		discovery := config.GetConfigString("go.discovery")
		if discovery == "" {
			discovery = "nacos"
		}
		switch discovery {
		case "nacos":
			host = config.GetNacosServiceURL(service)
		case "consul":
			host = config.GetConsulServiceURL(service)
		}
		if host != "" {
			mgcache.OnGetCache("nacos").Add(service, host, 5*time.Minute)
		} else {
			return "", errors.New("微服务获取" + service + "服务主机IP端口失败")
		}
	}
	url := host + uri
	header := map[string]string{"X-Request-Id": mgtrace.GetRequestId()}
	logs.Debug("Nacos微服务请求:{}\n请求参数:{}", url, params)
	resp, err := grequests.Post(url, &grequests.RequestOptions{
		Data:    params,
		Files:   files,
		Headers: header,
	})
	logs.Debug("Nacos微服务返回结果:{}", resp.String())
	if err != nil {
		return "", err
	} else {
		return resp.String(), err
	}
}

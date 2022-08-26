package mgcall

import (
	"errors"
	"fmt"
	"github.com/levigross/grequests"
	"github.com/maczh/logs"
	"github.com/maczh/mgcache"
	"github.com/maczh/mgconfig"
	"github.com/maczh/mgerr"
	"github.com/maczh/mgtrace"
	"net/url"
	"strings"
	"time"
)

func RestfulWithHeader(method, service string, uri string, pathparams map[string]string, header map[string]string, body interface{}) (string, error) {
	host, err := getHostFromCache(service)
	group := "DEFAULT_GROUP"
	if err != nil || host == "" {
		discovery := mgconfig.GetConfigString("go.discovery")
		if discovery == "" {
			discovery = "nacos"
		}
		switch discovery {
		case "nacos":
			group, host = mgconfig.GetNacosServiceURL(service)
			if host != "" && !mgcache.OnGetCache("nacos").IsExist("nacos:subscribe:"+service) {
				subscribeNacos(service, group)
				mgcache.OnGetCache("nacos").Add("nacos:subscribe:"+service, "true", 0)
			}
		case "consul":
			host = mgconfig.GetConsulServiceURL(service)
		}
		if host != "" {
			mgcache.OnGetCache("nacos").Add(service, host, 5*time.Minute)
		} else {
			return "", errors.New("微服务获取" + service + "服务主机IP端口失败")
		}
	}
	for k, v := range pathparams {
		strings.ReplaceAll(uri, fmt.Sprintf("{%s}", k), url.PathEscape(v))
	}
	url := host + uri
	logs.Debug("Nacos微服务请求:{}\n请求参数:{}", url, body)
	header["X-Request-Id"] = mgtrace.GetRequestId()
	header["X-Lang"] = mgerr.GetCurrentLanguage()
	header["Content-Type"] = "application/json"
	var resp *grequests.Response
	switch method {
	case "GET":
		resp, err = grequests.Get(url, &grequests.RequestOptions{
			Headers:            header,
			InsecureSkipVerify: true,
			JSON:               body,
		})
	case "POST":
		resp, err = grequests.Post(url, &grequests.RequestOptions{
			Headers:            header,
			InsecureSkipVerify: true,
			JSON:               body,
		})
	case "DELETE":
		resp, err = grequests.Delete(url, &grequests.RequestOptions{
			Headers:            header,
			InsecureSkipVerify: true,
			JSON:               body,
		})
	case "PUT":
		resp, err = grequests.Put(url, &grequests.RequestOptions{
			Headers:            header,
			InsecureSkipVerify: true,
			JSON:               body,
		})
	case "OPTIONS":
		resp, err = grequests.Options(url, &grequests.RequestOptions{
			Headers:            header,
			InsecureSkipVerify: true,
			JSON:               body,
		})
	case "HEAD":
		resp, err = grequests.Head(url, &grequests.RequestOptions{
			Headers:            header,
			InsecureSkipVerify: true,
			JSON:               body,
		})
	}
	logs.Debug("Nacos微服务返回结果:{}", resp.String())
	if err != nil {
		mgcache.OnGetCache("nacos").Delete(service)
		discovery := mgconfig.GetConfigString("go.discovery")
		if discovery == "" {
			discovery = "nacos"
		}
		switch discovery {
		case "nacos":
			group, host = mgconfig.GetNacosServiceURL(service)
			if host != "" && !mgcache.OnGetCache("nacos").IsExist("nacos:subscribe:"+service) {
				subscribeNacos(service, group)
				mgcache.OnGetCache("nacos").Add("nacos:subscribe:"+service, "true", 0)
			}
		case "consul":
			host = mgconfig.GetConsulServiceURL(service)
		}
		if host != "" {
			mgcache.OnGetCache("nacos").Add(service, host, 5*time.Minute)
		} else {
			return "", errors.New("微服务获取" + service + "服务主机IP端口失败")
		}
		url = host + uri
		switch method {
		case "GET":
			resp, err = grequests.Get(url, &grequests.RequestOptions{
				Headers:            header,
				InsecureSkipVerify: true,
				JSON:               body,
			})
		case "POST":
			resp, err = grequests.Post(url, &grequests.RequestOptions{
				Headers:            header,
				InsecureSkipVerify: true,
				JSON:               body,
			})
		case "DELETE":
			resp, err = grequests.Delete(url, &grequests.RequestOptions{
				Headers:            header,
				InsecureSkipVerify: true,
				JSON:               body,
			})
		case "PUT":
			resp, err = grequests.Put(url, &grequests.RequestOptions{
				Headers:            header,
				InsecureSkipVerify: true,
				JSON:               body,
			})
		case "OPTIONS":
			resp, err = grequests.Options(url, &grequests.RequestOptions{
				Headers:            header,
				InsecureSkipVerify: true,
				JSON:               body,
			})
		case "HEAD":
			resp, err = grequests.Head(url, &grequests.RequestOptions{
				Headers:            header,
				InsecureSkipVerify: true,
				JSON:               body,
			})
		}
		logs.Debug("Nacos微服务返回结果:{}", resp.String())
		if err != nil {
			return "", err
		} else {
			return resp.String(), nil
		}
	} else {
		return resp.String(), err
	}
}

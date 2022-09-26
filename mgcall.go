package mgcall

import (
	"encoding/json"
	"errors"
	"github.com/levigross/grequests"
	"github.com/maczh/logs"
	"github.com/maczh/mgcache"
	"github.com/maczh/mgconfig"
	"github.com/maczh/mgerr"
	"github.com/maczh/mgtrace"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

//微服务调用其他服务的接口
func Call(service string, uri string, params map[string]string) (string, error) {
	return CallWithHeader(service, uri, params, map[string]string{})
}

func Get(service string, uri string, params map[string]string) (string, error) {
	return GetWithHeader(service, uri, params, map[string]string{})
}

func GetWithHeader(service string, uri string, params map[string]string, header map[string]string) (string, error) {
	host, err := getHostFromCache(service)
	group := "DEFAULT_GROUP"
	if err != nil || host == "" {
		discovery := mgconfig.GetConfigString("go.discovery")
		if discovery == "" {
			discovery = "nacos"
		}
		switch discovery {
		case "nacos":
			host, group = mgconfig.GetNacosServiceURL(service)
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
	url := host + uri
	logs.Debug("Nacos微服务请求:{}\n请求参数:{}", url, params)
	header["X-Request-Id"] = mgtrace.GetRequestId()
	header["X-Lang"] = mgerr.GetCurrentLanguage()
	resp, err := grequests.Get(url, &grequests.RequestOptions{
		Params:             params,
		Headers:            header,
		InsecureSkipVerify: true,
	})
	logs.Debug("Nacos微服务返回结果:{}", resp.String())
	if err != nil {
		mgcache.OnGetCache("nacos").Delete(service)
		discovery := mgconfig.GetConfigString("go.discovery")
		if discovery == "" {
			discovery = "nacos"
		}
		switch discovery {
		case "nacos":
			host, group = mgconfig.GetNacosServiceURL(service)
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
		resp, err = grequests.Get(url, &grequests.RequestOptions{
			Params:             params,
			Headers:            header,
			InsecureSkipVerify: true,
		})
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

//微服务调用其他服务的接口,带header
func CallWithHeader(service string, uri string, params map[string]string, header map[string]string) (string, error) {
	host, err := getHostFromCache(service)
	group := "DEFAULT_GROUP"
	if err != nil || host == "" {
		discovery := mgconfig.GetConfigString("go.discovery")
		if discovery == "" {
			discovery = "nacos"
		}
		switch discovery {
		case "nacos":
			host, group = mgconfig.GetNacosServiceURL(service)
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
	url := host + uri
	header["X-Request-Id"] = mgtrace.GetRequestId()
	header["X-Lang"] = mgerr.GetCurrentLanguage()
	logs.Debug("Nacos微服务请求:{}\n请求参数:{}\n请求头:{}", url, params, header)
	resp, err := grequests.Post(url, &grequests.RequestOptions{
		Data:               params,
		Headers:            header,
		InsecureSkipVerify: true,
	})
	logs.Debug("Nacos微服务返回结果:{}", resp.String())
	if err != nil {
		mgcache.OnGetCache("nacos").Delete(service)
		discovery := mgconfig.GetConfigString("go.discovery")
		if discovery == "" {
			discovery = "nacos"
		}
		switch discovery {
		case "nacos":
			host, group = mgconfig.GetNacosServiceURL(service)
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
		resp, err = grequests.Post(url, &grequests.RequestOptions{
			Data:    params,
			Headers: header,
		})
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

func CallWithFiles(service string, uri string, params map[string]string, files []grequests.FileUpload) (string, error) {
	return CallWithFilesHeader(service, uri, params, files, map[string]string{})
}

//微服务调用其他服务的接口,带文件
func CallWithFilesHeader(service string, uri string, params map[string]string, files []grequests.FileUpload, header map[string]string) (string, error) {
	host, err := getHostFromCache(service)
	group := "DEFAULT_GROUP"
	if err != nil || host == "" {
		discovery := mgconfig.GetConfigString("go.discovery")
		if discovery == "" {
			discovery = "nacos"
		}
		switch discovery {
		case "nacos":
			host, group = mgconfig.GetNacosServiceURL(service)
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
	url := host + uri
	header["X-Request-Id"] = mgtrace.GetRequestId()
	header["X-Lang"] = mgerr.GetCurrentLanguage()
	logs.Debug("Nacos微服务请求:{}\n请求参数:{}", url, params)
	resp, err := grequests.Post(url, &grequests.RequestOptions{
		Data:               params,
		Files:              files,
		Headers:            header,
		InsecureSkipVerify: true,
	})
	logs.Debug("Nacos微服务返回结果:{}", resp.String())
	if err != nil {
		mgcache.OnGetCache("nacos").Delete(service)
		discovery := mgconfig.GetConfigString("go.discovery")
		if discovery == "" {
			discovery = "nacos"
		}
		switch discovery {
		case "nacos":
			host, group = mgconfig.GetNacosServiceURL(service)
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
		resp, err = grequests.Post(url, &grequests.RequestOptions{
			Data:               params,
			Headers:            header,
			InsecureSkipVerify: true,
		})
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

func getHostFromCache(serviceName string) (string, error) {
	h, _ := mgcache.OnGetCache("nacos").Value(serviceName)
	if h == nil {
		logs.Debug("{}服务无缓存", serviceName)
		return "", errors.New("无此服务缓存")
	} else {
		hosts := strings.Split(h.(string), ",")
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		return hosts[r.Intn(len(hosts))], nil
	}
}

func subscribeNacos(serviceName, groupName string) {
	logs.Debug("Nacos微服务订阅服务名:{}", serviceName)
	if groupName == "" {
		groupName = "DEFAULT_GROUP"
	}
	err := mgconfig.Nacos.Subscribe(&vo.SubscribeParam{
		ServiceName: serviceName,
		Clusters:    []string{"DEFAULT"},
		GroupName:   groupName,
		SubscribeCallback: func(services []model.SubscribeService, err error) {
			subscribeNacosCallback(services, err)
		},
	})
	if err != nil {
		logs.Error("Nacos订阅错误:{}", err.Error())
	}
}

func subscribeNacosCallback(services []model.SubscribeService, err error) {
	logs.Debug("Nacos回调:{}", services)
	if err != nil {
		logs.Error("Nacos订阅回调错误:{}", err.Error())
		return
	}
	if services == nil || len(services) == 0 {
		logs.Error("Nacos订阅回调服务列表为空")
		return
	}
	servicesMap := make(map[string]string)
	for _, s := range services {
		protocal := "http://"
		if s.Metadata != nil && s.Metadata["ssl"] == "true" {
			protocal = "https://"
		}
		if servicesMap[s.ServiceName] == "" {
			servicesMap[s.ServiceName] = protocal + s.Ip + ":" + strconv.Itoa(int(s.Port))
		} else {
			servicesMap[s.ServiceName] = servicesMap[s.ServiceName] + "," + protocal + s.Ip + ":" + strconv.Itoa(int(s.Port))
		}
	}
	for serviceName, host := range servicesMap {
		mgcache.OnGetCache("nacos").Delete(serviceName)
		mgcache.OnGetCache("nacos").Add(serviceName, host, 5*time.Minute)
	}
}

func toJSON(o interface{}) string {
	j, err := json.Marshal(o)
	if err != nil {
		return "{}"
	} else {
		js := string(j)
		js = strings.Replace(js, "\\u003c", "<", -1)
		js = strings.Replace(js, "\\u003e", ">", -1)
		js = strings.Replace(js, "\\u0026", "&", -1)
		return js
	}
}

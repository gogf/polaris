// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ratelimiter

import (
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/polarismesh/polaris-go/api"
	"strings"
	"sync"
)

var (
	limit, err    = api.NewLimitAPI()
	namespace     string
	service       string
	limitFailFunc = func(r *ghttp.Request) {
		r.Response.WriteExit(`{"code":500,"message":"资源不足"}`)
	}
	MatchLabelMap = map[string]map[string]string{}
	mu            sync.RWMutex
)

// RegisterByUriLabel .
func RegisterByUriLabel(labelMap map[string]string, limitExceededFunc ...func(r *ghttp.Request)) error {
	if len(labelMap) == 0 {
		return errors.New("labelMap不能为空")
	}
	if len(limitExceededFunc) != 0 {
		limitFailFunc = limitExceededFunc[0]
	}
	mu.Lock()
	defer mu.Unlock()
	for uri, label := range labelMap {
		labels, err := parseLabels(label)
		if err != nil {
			return err
		}
		MatchLabelMap[uri] = labels
	}
	return nil
}

// RateLimit .
func RateLimit(r *ghttp.Request) {
	uri := r.RequestURI
	mu.RLock()
	defer mu.RUnlock()
	// 能够直接精确匹配
	if label, ok := MatchLabelMap[uri]; ok {
		getQuotaResult(r, label)
	}
	// 遍历所有注册label，进行labelMatch检查，满足的最长的path胜出

}

func getQuotaResult(r *ghttp.Request, label map[string]string) {
	param := api.NewQuotaRequest()
	param.SetLabels(label)
	param.SetNamespace(namespace)
	param.SetService(service)
	getQuota, err := limit.GetQuota(param)
	if err != nil {
		// gf 带有错误回收，只是中断本次请求
		panic(err)
	}
	if getQuota.Get().Code == api.QuotaResultOk {
		r.Middleware.Next()
	}
	limitFailFunc(r)
}

// RegisterByHook .
func RegisterByHook(r *ghttp.Server, limitExceededFunc func(r *ghttp.Request), labelMap map[string]string) error {
	if err != err {
		return errors.New(fmt.Sprintf("fail to create consumerAPI, err is %v", err))
	}
	for pattern, labelsStr := range labelMap {
		label, err := parseLabels(labelsStr)
		if err != nil {
			return err
		}
		param := api.NewQuotaRequest()
		param.SetLabels(label)
		param.SetNamespace(namespace)
		param.SetService(service)
		r.BindHookHandler(pattern, ghttp.HookBeforeServe, func(r *ghttp.Request) {
			getQuota, err := limit.GetQuota(param)
			if err != nil {
				// gf 带有错误回收，只是中断本次请求
				panic(err)
			}
			if getQuota.Get().Code == api.QuotaResultOk {
				r.Middleware.Next()
			}
			limitExceededFunc(r)
		})
	}
	return nil
}

//解析标签列表
func parseLabels(labelsStr string) (map[string]string, error) {
	strLabels := strings.Split(labelsStr, ",")
	labels := make(map[string]string, len(strLabels))
	for _, strLabel := range strLabels {
		if len(strLabel) < 1 {
			continue
		}
		labelKv := strings.Split(strLabel, ":")
		if len(labelKv) != 2 {
			return nil, fmt.Errorf("invalid kv pair str %s", strLabel)
		}
		labels[labelKv[0]] = labelKv[1]
	}
	return labels, nil
}

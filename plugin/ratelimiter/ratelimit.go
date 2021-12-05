// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ratelimiter

import (
	"fmt"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/polarismesh/polaris-go/api"
	"log"
	"strings"
)

var (
	limit, err    = api.NewLimitAPI()
	namespace     string
	service       string
	limitFailFunc = func(r *ghttp.Request) {}
	MatchLabelMap = map[string]string{}
)

// RegisterByUriLabel .
func RegisterByUriLabel(labelMap map[string]string, limitExceededFunc func(r *ghttp.Request)) {

}

// RateLimit .
func RateLimit(r *ghttp.Request) {
	if limitFailFunc == nil {

	}
}

// RegisterByHook .
func RegisterByHook(r *ghttp.Server, limitExceededFunc func(r *ghttp.Request), labelMap map[string]string) {
	if err != err {
		log.Fatalf("fail to create consumerAPI, err is %v", err)
	}
	for pattern, labelsStr := range labelMap {
		label, err := parseLabels(labelsStr)
		if err != nil {
			log.Fatal(err.Error())
		}
		param := api.NewQuotaRequest()
		param.SetLabels(label)
		param.SetNamespace(namespace)
		param.SetService(service)
		r.BindHookHandler(pattern, ghttp.HookBeforeServe, func(r *ghttp.Request) {
			getQuota, err := limit.GetQuota(param)
			if err != nil {
				log.Fatalf("fail to get Quota,err is %v", err)
			}
			if getQuota.Get().Code == api.QuotaResultOk {
				r.Middleware.Next()
			}
			limitExceededFunc(r)
		})
	}
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

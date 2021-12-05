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
	limit, err = api.NewLimitAPI()
	param      = api.NewQuotaRequest()
	namespace  string
	service    string
	labelsStr  string
)

// Register .
func Register(r *ghttp.Server, limitExceededFunc func(r *ghttp.Request), pattern ...string) {
	if err != err {
		log.Fatalf("fail to create consumerAPI, err is %v", err)
	}
	if len(pattern) == 0 {
		pattern = []string{"/*"}
	}

	label, err := parseLabels(labelsStr)
	if err != nil {
		log.Fatal(err.Error())
	}
	param.SetLabels(label)
	param.SetNamespace(namespace)
	param.SetService(service)

	r.BindHookHandlerByMap(pattern[0], map[string]ghttp.HandlerFunc{
		ghttp.HookBeforeServe: func(r *ghttp.Request) {
			getQuota, err := limit.GetQuota(param)
			if err != nil {
				log.Fatalf("fail to get Quota,err is %v", err)
			}
			defer getQuota.Release()
			if getQuota.Get().Code == api.QuotaResultOk {
				r.Middleware.Next()
			}
			limitExceededFunc(r)
		},
	})
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

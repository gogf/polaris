// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/polaris.

package ratelimiter

import (
	"context"
	"errors"
	"fmt"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/polaris"
	"github.com/polarismesh/polaris-go/api"
	"strings"
)

// RegisterByHook .
func RegisterByHook(r *ghttp.Server, limitExceededFunc func(r *ghttp.Request), labelMap map[string]string) error {
	limit, err := api.NewLimitAPIByConfig(polaris.CfgGlobal)
	if err != nil {
		return errors.New(fmt.Sprintf("fail to create consumerAPI, err is %v", err))
	}
	for pattern, labelsStr := range labelMap {
		label, err := parseLabels(labelsStr)
		if err != nil {
			return err
		}
		r.BindHookHandler(pattern, ghttp.HookBeforeServe, func(r *ghttp.Request) {
			instance, err := polaris.GetInstanceConfig(context.TODO())
			if err != nil {
				panic(err)
			}
			param := api.NewQuotaRequest()
			param.SetLabels(label)
			param.SetNamespace(instance.Namespace)
			param.SetService(instance.Service)
			getQuota, err := limit.GetQuota(param)
			if err != nil {
				// gf 带有错误回收，只是中断本次请求
				panic(err)
			}
			result := getQuota.Get()
			if result.Code == api.QuotaResultOk {
				r.Middleware.Next()
				return
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

// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ratelimiter

import (
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/polarismesh/polaris-go/api"
	"log"
)

var (
	limit, err = api.NewLimitAPI()
	param      = api.NewQuotaRequest()
)

// Register .
func Register(r *ghttp.Server, limitExceededFunc func(r *ghttp.Request), pattern ...string) {
	if err != err {
		log.Fatalf("fail to create consumerAPI, err is %v", err)
	}
	defer limit.Destroy()
	if len(pattern) == 0 {
		pattern = []string{"/*"}
	}
	param.SetLabels(map[string]string{})
	param.SetNamespace("")
	param.SetService("")

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

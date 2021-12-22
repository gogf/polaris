// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/polaris.

package main

import (
	"context"
	"fmt"
	"log"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/polaris"
	"github.com/gogf/polaris/plugin/ratelimiter"
)

func main() {
	// init config
	ctx := context.Background()
	adapter, err := gcfg.NewAdapterFile("config.yaml")
	if err != nil {
		g.Log().Fatal(ctx, "boot init g cfg.NewAdapterFile error:", err)
	}
	g.Cfg().SetAdapter(adapter)
	s := g.Server()

	err = polaris.InitConfigPolaris()
	if err != nil {
		g.Log().Fatal(context.TODO(), err.Error())
	}
	// register
	s.Plugin(polaris.GfPolarisPlugin{
		Listener: func(config string) {
			g.Log().Print(ctx, "Polaris register success")
		},
	})

	// router
	s.BindHandler("/", func(r *ghttp.Request) {
		fmt.Println("success.")
	})

	// init rate limit
	limitedFunc := func(r *ghttp.Request) {
		fmt.Println("limit.")
	}
	err = ratelimiter.RegisterByHook(s, limitedFunc, map[string]string{
		"/*": "*:* ; method:*",
	})
	if err != nil {
		log.Fatalf("init fail,this error is %v", err)
	}

	s.Run()
}

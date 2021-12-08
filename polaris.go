// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/polaris.

package polaris

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gipv4"
	"github.com/gogf/gf/v2/os/gcfg"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gtimer"
	"github.com/polarismesh/polaris-go/api"
	"github.com/polarismesh/polaris-go/pkg/config"
)

const (
	// Version polaris plugin version
	Version = "0.0.1"
	// DefaultNamespace name space
	DefaultNamespace = "default"
	// DefaultService Default service name
	DefaultService = "GoGF-polaris"
	// DefaultPort default service port
	DefaultPort = 8199
	// TTL .The timeout period, if the node wants to call heartbeat to report,
	// it must be filled in, otherwise it will be 400141 error code, unit: second
	TTL = 2
)

var (
	instance  *api.InstanceRegisterRequest
	cfgGlobal config.Configuration
	polaris   *Polaris
)

// Init Polaris plug-in initialization.
func Init() error {
	ctx := gctx.New()
	initConfigPolaris(ctx)
	// Consumer(ctx)
	return nil
}

// Deregister .Remove registration
func Deregister() {
	var (
		ctx      = gctx.New()
		provider = NewProvider(ctx)
	)

	defer provider.Destroy()
	deregisterRequest := &api.InstanceDeRegisterRequest{}
	deregisterRequest.Service = instance.Service
	deregisterRequest.Namespace = instance.Namespace
	deregisterRequest.Host = instance.Host
	deregisterRequest.Port = instance.Port
	deregisterRequest.ServiceToken = instance.ServiceToken
	if err := provider.Deregister(deregisterRequest); nil != err {
		g.Log().Fatal(ctx, "fail to deregister instance, err is %v", err)
	}
	g.Log().Info(ctx, "Deregister end")
}

// initConfigPolaris init
func initConfigPolaris(ctx context.Context) {
	var cfg = g.Cfg()

	polarisConfig(ctx, cfg)
	// set logger dir
	api.SetLoggersDir(polaris.Config.LoggerPath)

	cfgGlobal = globalPolarisConfig(ctx, cfg)

	g.Log().Info(ctx, "initConfigPolaris config:", cfgGlobal)

	// 实行注册
	register(ctx)
	if polaris.Config.IsHeartbeat > 0 {
		gtimer.SetInterval(ctx, 5*time.Second, func(ctx context.Context) {
			// heartbeat report
			Heartbeat(ctx)
		})
	}
	g.Log().Info(ctx, "api.initConfigPolaris end")
}

// Consumer .获取服务列表信息
func Consumer(ctx context.Context) {
	consumer, err := api.NewConsumerAPIByConfig(cfgGlobal)
	if nil != err {
		g.Log().Fatalf(ctx, "fail to create consumerAPI, err is %v", err)
	}
	defer consumer.Destroy()

	request := &api.GetAllInstancesRequest{}
	request.Namespace = instance.Namespace
	request.Service = instance.Service
	resp, err := consumer.GetAllInstances(request)
	for i, inst := range resp.Instances {
		g.Log().Printf(ctx, "instance %d is %s:%d\n", i, inst.GetHost(), inst.GetPort())
	}
}

// register .provider register
func register(ctx context.Context) {
	var provider = NewProvider(ctx)
	// before process exits
	defer provider.Destroy()

	g.Log().Debug(ctx, "register request start params:", polaris.Instance)
	resp, err := provider.Register(polaris.Instance)
	if nil != err {
		g.Log().Fatal(ctx, "provider.register params:", resp, " fail reason err:", err)
	}
}

// Heartbeat .heartbeat report
func Heartbeat(ctx context.Context) {
	var provider = NewProvider(ctx)
	// before process exits
	defer provider.Destroy()

	request := &api.InstanceHeartbeatRequest{}
	request.Namespace = polaris.Instance.Namespace
	request.Service = polaris.Instance.Service
	request.Host = polaris.Instance.Host
	request.Port = polaris.Instance.Port
	if err := provider.Heartbeat(request); err != nil {
		g.Log().Error(ctx, "provider Heartbeat params:", request, " fail reason err:", err)
	}
	g.Log().Info(ctx, "provider Heartbeat end ")
}

// NewProvider . create Provider
func NewProvider(ctx context.Context) api.ProviderAPI {
	var provider, err = api.NewProviderAPIByConfig(cfgGlobal)
	if nil != err {
		g.Log().Fatal(ctx, "NewProvider api.NewProviderAPIByConfig fail err:", err)
	}
	return provider
}

// fillDefaults 完善远程默认链接
func fillDefaults() {
	if polaris.Instance.Namespace == "" {
		polaris.Instance.Namespace = DefaultNamespace
	}
	if polaris.Instance.Service == "" {
		polaris.Instance.Service = DefaultService
	}
	if polaris.Instance.Port < 1 {
		polaris.Instance.Port = DefaultPort
	}
	if polaris.Instance.Host == "" {
		polaris.Instance.Host, _ = gipv4.GetIntranetIp()
	}
	instance = polaris.Instance
}

// globalPolarisConfig global Polaris Config
func globalPolarisConfig(ctx context.Context, cfg *gcfg.Config) config.Configuration {
	v, err := cfg.Get(ctx, "global")
	if err != nil {
		g.Log().Fatal(ctx, "GoFrame config get global fail error:", err)
	}

	if v.IsNil() || v.IsEmpty() {
		g.Log().Fatal(ctx, "GoFrame config get global is not exits")
	}

	// 获取配置信息
	var globalConfig = new(config.GlobalConfigImpl)
	globalConfig.Init()
	if err = v.Struct(&globalConfig); err != nil {
		g.Log().Fatal(ctx, "Struct ServerConnector error:", err)
	}
	globalConfig.ServerConnector.Init()
	globalConfig.System.Init()
	globalConfig.StatReporter.Init()

	g.Log().Debug(ctx, "Struct init end globalConfig:", globalConfig)
	cfgGlobal := &config.ConfigurationImpl{}
	cfgGlobal.Init()
	cfgGlobal.SetDefault()
	cfgGlobal.Global = globalConfig
	if len(polaris.Config.BackupPath) > 0 {
		cfgGlobal.Consumer = new(config.ConsumerConfigImpl)
		cfgGlobal.Consumer.Init()
		cfgGlobal.Consumer.LocalCache.Init()
		cfgGlobal.Consumer.LocalCache.SetPersistDir(polaris.Config.BackupPath)
	}
	return cfgGlobal
}

// polarisConfig set config Polaris
func polarisConfig(ctx context.Context, cfg *gcfg.Config) {
	v, err := cfg.Get(ctx, "polaris")
	if err != nil {
		g.Log().Fatal(ctx, "GoFrame config get polaris fail error:", err)
	}

	if v.IsNil() || v.IsEmpty() {
		g.Log().Fatal(ctx, "GoFrame config get polaris is not exist")
	}
	g.Log().Debug(ctx, "polaris map config:", v)
	// 获取配置信息
	if err = v.Struct(&polaris); err != nil {
		g.Log().Fatal(ctx, "error:", err)
	}
	g.Log().Debug(ctx, "polaris struct config:", polaris)

	fillDefaults()
}

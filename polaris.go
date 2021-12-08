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
	TTL int = 2
)

// InitConfigPolaris Polaris plug-in initialization.
func InitConfigPolaris() error {
	var (
		cfg = g.Cfg()
		ctx = gctx.New()
	)

	polarisConfig(ctx, cfg)
	// set logger dir
	api.SetLoggersDir(polaris.Config.LoggerPath)
	cfgGlobal = globalPolarisConfig(ctx, cfg)

	g.Log().Debug(ctx, "InitConfigPolaris config:", cfgGlobal)

	// Perform registration operation
	register(ctx)
	// determine whether to report the heartbeat
	if polaris.Config.IsHeartbeat > 0 {
		gtimer.SetInterval(ctx, 5*time.Second, func(ctx context.Context) {
			// heartbeat report
			heartbeat(ctx)
		})
	}
	g.Log().Info(ctx, "InitConfigPolaris end")
	return nil
}

// Deregister anti registration service
func Deregister() {
	var (
		ctx      = gctx.New()
		provider = provider(ctx)
	)
	defer provider.Destroy()

	deregisterRequest := &api.InstanceDeRegisterRequest{}
	deregisterRequest.Service = polaris.Instance.Service
	deregisterRequest.Namespace = polaris.Instance.Namespace
	deregisterRequest.Host = polaris.Instance.Host
	deregisterRequest.Port = polaris.Instance.Port
	deregisterRequest.ServiceToken = polaris.Instance.ServiceToken
	if err := provider.Deregister(deregisterRequest); nil != err {
		g.Log().Fatal(ctx, "fail to deregister instance, err is %v", err)
	}
	g.Log().Info(ctx, "Deregister end")
}

// register registration service
func register(ctx context.Context) {
	var provider = provider(ctx)
	// before process exits
	defer provider.Destroy()

	g.Log().Debug(ctx, "register request start params:", polaris.Instance)
	request := &api.InstanceRegisterRequest{}
	request.Service = polaris.Instance.Service
	request.Namespace = polaris.Instance.Namespace
	request.Host = polaris.Instance.Host
	request.Port = polaris.Instance.Port
	request.ServiceToken = polaris.Instance.ServiceToken
	if polaris.Instance.TTL < TTL {
		polaris.Instance.TTL = TTL
	}
	if len(polaris.Instance.Version) > 0 {
		request.Version = &polaris.Instance.Version
	}
	request.SetTTL(polaris.Instance.TTL)
	request.SetHealthy(true)
	resp, err := provider.Register(request)
	if nil != err {
		g.Log().Fatal(ctx, "provider.register params:", resp, " fail reason err:", err)
	}
	g.Log().Info(ctx, "provider.register end")
}

// heartbeat .heartbeat report
func heartbeat(ctx context.Context) {
	var provider = provider(ctx)
	// before process exits
	defer provider.Destroy()

	request := &api.InstanceHeartbeatRequest{}
	request.Namespace = polaris.Instance.Namespace
	request.Service = polaris.Instance.Service
	request.Host = polaris.Instance.Host
	request.Port = polaris.Instance.Port
	if err := provider.Heartbeat(request); err != nil {
		g.Log().Error(ctx, "provider heartbeat params:", request, " fail reason err:", err)
	}
	g.Log().Info(ctx, "provider heartbeat end ")
}

// provider . create Provider
func provider(ctx context.Context) api.ProviderAPI {
	var provider, err = api.NewProviderAPIByConfig(cfgGlobal)
	if nil != err {
		g.Log().Fatal(ctx, "provider api.NewProviderAPIByConfig fail err:", err)
	}
	return provider
}

// fillDefaults improve the remote default link
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

	// Initialize configuration information
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
	// Get configuration information
	if err = v.Struct(&polaris); err != nil {
		g.Log().Fatal(ctx, "error:", err)
	}
	g.Log().Debug(ctx, "polaris struct config:", polaris)

	fillDefaults()
}

// Consumer .Get service list information
func Consumer(ctx context.Context) {
	consumer, err := api.NewConsumerAPIByConfig(cfgGlobal)
	if nil != err {
		g.Log().Fatalf(ctx, "fail to create consumerAPI, err is %v", err)
	}
	defer consumer.Destroy()

	request := &api.GetAllInstancesRequest{}
	request.Namespace = polaris.Instance.Namespace
	request.Service = polaris.Instance.Service
	resp, err := consumer.GetAllInstances(request)
	for i, inst := range resp.Instances {
		g.Log().Printf(ctx, "instance %d is %s:%d\n", i, inst.GetHost(), inst.GetPort())
	}
}

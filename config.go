// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/polaris.

package polaris

import (
	"time"

	"github.com/gogf/gf/v2/os/gctx"
	"github.com/polarismesh/polaris-go/api"
	"github.com/polarismesh/polaris-go/pkg/config"
)

var (
	cfgGlobal   config.Configuration
	polaris     *Polaris
	apiProvider api.ProviderAPI
	err         error
	ctx         = gctx.New()
)

// Polaris .
type Polaris struct {
	// Config polaris configuration
	Config *ConfigPolaris `json:"config"`
	// Instance Service configuration information
	Instance *InstanceRequest `json:"server"`
}

// ConfigPolaris config
type ConfigPolaris struct {
	// 是否开启心跳上报
	IsHeartbeat uint `json:"isHeartbeat"`
	// 日志目录
	LoggerPath string `json:"loggerPath"`
	// 备份目录
	BackupPath string `json:"backupPath"`
}

// InstanceRequest .
type InstanceRequest struct {
	// Service Required, service name
	Service string
	// ServiceToken Required, service access token
	ServiceToken string
	// Namespace Required, namespace
	Namespace string
	// Host Required, service monitoring host, support IPv 6 address
	Host string
	// Port Required, service instance monitor port
	Port int

	// The following fields are optional, the default nil means that the client is not configured, and the server configuration is used
	// Protocol Service Agreement
	Protocol string
	// Weight Service weight, default 100, range 0-10000
	Weight int
	// Priority Instance priority, the default is 0, the smaller the value, the higher the priority
	Priority int
	// Version Instance provides service version number
	Version string
	// Metadata user defined metadata information
	Metadata map[string]string
	// Healthy Whether the service instance is healthy, the default is healthy
	Healthy bool
	// Isolate Whether the service instance is isolated, not isolated by default
	Isolate bool
	// TTL timeout, if the node wants to call heartbeat to report,
	// it must be filled in, otherwise it will be 400141 error code, unit: second
	TTL int

	// Timeout Optional, single query timeout time, the default is to directly obtain the global timeout configuration
	// The user's total maximum timeout time is (1+RetryCount) Timeout
	Timeout time.Duration
	// RetryCount Optional, the number of retries, the global timeout configuration is directly obtained by default
	RetryCount int
}

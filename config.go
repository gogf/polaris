// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/polaris.

package polaris

import "github.com/polarismesh/polaris-go/api"

// ConfigPolaris config
type ConfigPolaris struct {
	// 是否开启心跳上报
	IsHeartbeat uint `json:"isHeartbeat"`
	// 日志目录
	LoggerPath string `json:"loggerPath"`
	// 备份目录
	BackupPath string `json:"backupPath"`
}

// Polaris .
type Polaris struct {
	Config *ConfigPolaris `json:"config"`
	// 服务配置信息
	Instance *api.InstanceRegisterRequest `json:"server"`
}

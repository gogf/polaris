// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/polaris.

package polaris

import (
	"fmt"

	"github.com/gogf/gf/v2/net/ghttp"
)

// ConfigListener .
type ConfigListener func(config string)

var configListener ConfigListener = func(config string) {

}

// GfPolarisPlugin .
type GfPolarisPlugin struct {
	Listener ConfigListener
}

// Name Plugin name
func (p GfPolarisPlugin) Name() string {
	return "gf-polaris"
}

// Author  website of author
func (p GfPolarisPlugin) Author() string {
	return "github.com/gogf/polaris"
}

// Version gf Polaris plugin version
func (p GfPolarisPlugin) Version() string {
	return Version
}

// Description desc of plugin
func (p GfPolarisPlugin) Description() string {
	return "GoFrame and Polaris"
}

// Install plugin installation
func (p GfPolarisPlugin) Install(s *ghttp.Server) error {
	fmt.Println("GoFrame-polaris the plugin is being installed...")
	configListener = p.Listener
	fmt.Printf("configListener: %s", configListener)
	return InitConfigPolaris()
}

// Remove plugin removal
func (p GfPolarisPlugin) Remove() error {
	Deregister()
	fmt.Println("GoFrame-polaris plugin removed。")
	return nil
}

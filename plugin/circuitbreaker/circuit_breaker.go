// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package circuitbreaker

import (
	"github.com/gogf/gf/v2/net/ghttp"
)

func Register(r *ghttp.Server, methodFunc func(r *ghttp.Request), pattern ...string) {

}

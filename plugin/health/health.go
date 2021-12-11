// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/polaris.

package health

import (
	"log"

	"github.com/polarismesh/polaris-go/api"
)

const (
	host          = "127.0.0.1"
	startPort     = 2001
	instanceCount = 5
)

// Register .
func Register() {
	var (
		namespace string
		service   string
		token     string

		consumer, err = api.NewConsumerAPI()
	)
	if nil != err {
		log.Fatalf("fail to create consumerAPI, err is %v", err)
	}
	defer consumer.Destroy()

	// share one context with the consumerAPI
	provider := api.NewProviderAPIByContext(consumer.SDKContext())

	registerRequest := &api.InstanceRegisterRequest{}
	registerRequest.Service = service
	registerRequest.Namespace = namespace
	registerRequest.Host = host
	registerRequest.Port = startPort
	registerRequest.ServiceToken = token
	registerRequest.SetHealthy(true)
	resp, err := provider.Register(registerRequest)
	if nil != err {
		log.Fatalf("fail to register instance, err is %v", err)
	}
	log.Printf("register instance response: instanceId %s", resp.InstanceID)
}

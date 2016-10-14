// Copyright 2016 Tamás Gulácsi
//
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package wsclient

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc/credentials"
)

var _ = credentials.PerRPCCredentials(basicAuthCreds{})

type basicAuthCreds struct {
	username, password string
}

func NewBasicAuth(username, password string) basicAuthCreds {
	return basicAuthCreds{username: username, password: password}
}
func (ba basicAuthCreds) RequireTransportSecurity() bool { return true }
func (ba basicAuthCreds) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	return map[string]string{
		"authorization": ba.username + ":" + ba.password,
	}, nil
}

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

const BasicAuthKey = "authorization-basic"

func WithBasicAuth(ctx context.Context, username, password string) context.Context {
	return context.WithValue(ctx, BasicAuthKey, username+":"+password)
}

var _ = credentials.PerRPCCredentials(basicAuthCreds{})

type basicAuthCreds struct {
	up string
}

func NewBasicAuth(username, password string) basicAuthCreds {
	return basicAuthCreds{up: username + ":" + password}
}
func (ba basicAuthCreds) RequireTransportSecurity() bool { return true }
func (ba basicAuthCreds) GetRequestMetadata(ctx context.Context, uri ...string) (map[string]string, error) {
	var up string
	if upI := ctx.Value(BasicAuthKey); upI != nil {
		up = upI.(string)
	}
	if up == "" {
		up = ba.up
	}
	return map[string]string{"authorization": up}, nil
}

// vim: se noet fileencoding=utf-8:

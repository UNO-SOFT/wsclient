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

// Package wsclient provide helpers for calling UNO-SOFT gRPC server.
package wsclient

import (
	"github.com/pkg/errors"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// DialOpts renders the dial options for calling a gRPC server.
//
// * prefix is inserted before the standard request path - if your server serves on different path.
// * caFile is the PEM file with the server's CA.
// * serverHostOverride is to override the CA's host.
func DialOpts(
	prefix, caFile, serverHostOverride string,
	Log func(keyvals ...interface{}) error,
) ([]grpc.DialOption, error) {
	dialOpts := make([]grpc.DialOption, 2, 5)
	dialOpts[0] = grpc.WithCompressor(grpc.NewGZIPCompressor())
	dialOpts[1] = grpc.WithDecompressor(grpc.NewGZIPDecompressor())

	if prefix != "" || Log != nil {
		if Log == nil {
			Log = func(keyvals ...interface{}) error { return nil }
		}
		dialOpts = append(dialOpts,
			grpc.WithStreamInterceptor(
				func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
					Log("method", method)
					return streamer(ctx, desc, cc, prefix+method, opts...)
				}),
			grpc.WithUnaryInterceptor(
				func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
					Log("method", method)
					return invoker(ctx, prefix+method, req, reply, cc, opts...)
				}),
		)
	}
	if caFile == "" {
		dialOpts = append(dialOpts, grpc.WithInsecure())
	} else {
		creds, err := credentials.NewClientTLSFromFile(caFile, serverHostOverride)
		if err != nil {
			return dialOpts, errors.Wrapf(err, "%q,%q", caFile, serverHostOverride)
		}
		dialOpts = append(dialOpts, grpc.WithTransportCredentials(creds))
	}

	return dialOpts, nil
}

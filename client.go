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
func DialOpts(prefix, caFile, serverHostOverride string) ([]grpc.DialOption, error) {
	dialOpts := make([]grpc.DialOption, 2, 5)
	dialOpts[0] = grpc.WithCompressor(grpc.NewGZIPCompressor())
	dialOpts[1] = grpc.WithDecompressor(grpc.NewGZIPDecompressor())

	if prefix != "" {
		dialOpts = append(dialOpts,
			grpc.WithStreamInterceptor(
				func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
					return streamer(ctx, desc, cc, prefix+method, opts...)
				}),
			grpc.WithUnaryInterceptor(
				func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
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

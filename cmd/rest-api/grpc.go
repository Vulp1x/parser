package main

//
// import (
// 	"context"
// 	"net"
// 	"net/url"
// 	"sync"
//
// 	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
// 	"github.com/inst-api/parser/internal/mw"
// 	"github.com/inst-api/parser/pkg/api"
// 	"github.com/inst-api/parser/pkg/logger"
// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/reflection"
// )
//
// // handleGRPCServer starts configures and starts a gRPC server on the given
// // URL. It shuts down the server if any error is received in the error channel.
// func handleGRPCServer(ctx context.Context, u *url.URL, src api.ParserServer, wg *sync.WaitGroup, errc chan error, debug bool) {
//
// 	// Initialize gRPC server with the middleware.
// 	srv := grpc.NewServer(
// 		grpcmiddleware.WithUnaryServerChain(
// 			mw.UnaryRequestID(),
// 			mw.UnaryServerLog(),
// 			mw.Recover,
// 		),
// 	)
//
// 	// Register the servers.
// 	api.RegisterParserServer(srv, src)
//
// 	for svc, info := range srv.GetServiceInfo() {
// 		for _, m := range info.Methods {
// 			logger.Infof(ctx, "serving gRPC method %s", svc+"/"+m.Name)
// 		}
// 	}
//
// 	// Register the server reflection service on the server.
// 	// See https://grpc.github.io/grpc/core/md_doc_server-reflection.html.
// 	reflection.Register(srv)
//
// 	// 	mux, err := swaggway.NewMux(desc, a.publicServer, a.opts.ServeMuxOpts...)
// 	// 	if err != nil {
// 	// 		logger.Fatalf(context.Background(), "error while init gateway:%s", err.Error())
// 	// 	}
// 	//
// 	// 	if err := desc.RegisterGateway(context.Background(), mux); err != nil {
// 	// 		logger.Fatalf(context.Background(), "error while register gateway:%s", err.Error())
// 	// 	}
// 	// }
//
// 	defer wg.Done()
//
// 	// Start gRPC server in a separate goroutine.
// 	go func() {
// 		lis, err := net.Listen("tcp", u.Host)
// 		if err != nil {
// 			errc <- err
// 		}
// 		logger.Infof(ctx, "gRPC server listening on %q", u.Host)
// 		errc <- srv.Serve(lis)
// 	}()
//
// 	<-ctx.Done()
// 	logger.Infof(ctx, "shutting down gRPC server at %q", u.Host)
// 	srv.Stop()
// }

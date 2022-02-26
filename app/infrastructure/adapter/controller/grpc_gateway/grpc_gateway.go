package grpc_gateway

import (
	"context"
	"fmt"
	"io/fs"
	"mime"
	"net/http"
	"strconv"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/timsolov/ms-users/app/infrastructure/logger"
	"github.com/timsolov/ms-users/third_party"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// RegisterServiceHandlerFunc func to register gRPC service handler.
type RegisterServiceHandlerFunc func(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error

// Run runs the gRPC-Gateway on the gatewayAddr using gprcClient connection
// as underlying gRPC client connection to gRPC server started before.
func Run(ctx context.Context, log logger.Logger, gatewayAddr string, grpcClient *grpc.ClientConn, services []RegisterServiceHandlerFunc) chan error {
	errCh := make(chan error, 1)

	gwmux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(headerMatcher),
		runtime.WithForwardResponseOption(httpResponseModifier),
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.HTTPBodyMarshaler{
			Marshaler: &runtime.JSONPb{
				MarshalOptions: protojson.MarshalOptions{
					UseProtoNames:   true,
					EmitUnpopulated: true,
				},
				UnmarshalOptions: protojson.UnmarshalOptions{
					DiscardUnknown: true,
				},
			},
		}),
	)

	for _, fn := range services {
		err := fn(ctx, gwmux, grpcClient)
		if err != nil {
			errCh <- fmt.Errorf("failed to register gateway: %w", err)
			return errCh
		}
	}

	oa := getOpenAPIHandler()

	const openApiPath = "/api/swagger/"

	gwServer := &http.Server{
		Addr: gatewayAddr,
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if strings.HasPrefix(r.URL.Path, openApiPath) {
				http.StripPrefix(openApiPath, oa).ServeHTTP(w, r)
				return
			}
			// loggerMw := LogRequest(log)
			// loggerMw(gwmux).ServeHTTP(w, r)
			gwmux.ServeHTTP(w, r)
		}),
	}

	log.Infof("Serving gRPC-Gateway http://%s", gatewayAddr)
	log.Infof("Serving OpenAPI Documentation on http://%s%s", gatewayAddr, openApiPath)
	go func() {
		errCh <- gwServer.ListenAndServe()
		if err := gwServer.Shutdown(ctx); err != nil {
			log.Errorf("Shutdown gRPC-Gateway error: %s", err)
		} else {
			log.Infof("Shutdown gRPC-Gateway")
		}
	}()

	return errCh
}

// getOpenAPIHandler serves an OpenAPI UI.
// Adapted from https://github.com/philips/grpc-gateway-example/blob/a269bcb5931ca92be0ceae6130ac27ae89582ecc/cmd/serve.go#L63
func getOpenAPIHandler() http.Handler {
	mime.AddExtensionType(".svg", "image/svg+xml")
	// Use subdirectory in embedded files
	subFS, err := fs.Sub(third_party.OpenAPI, "OpenAPI")
	if err != nil {
		panic("couldn't create sub filesystem: " + err.Error())
	}
	return http.FileServer(http.FS(subFS))
}

func httpResponseModifier(ctx context.Context, w http.ResponseWriter, p proto.Message) error {
	md, ok := runtime.ServerMetadataFromContext(ctx)
	if !ok {
		return nil
	}

	// set http status code
	if vals := md.HeaderMD.Get("x-http-code"); len(vals) > 0 {
		code, err := strconv.Atoi(vals[0])
		if err != nil {
			return err
		}
		// delete the headers to not expose any grpc-metadata in http response
		delete(md.HeaderMD, "x-http-code")
		delete(w.Header(), "Grpc-Metadata-X-Http-Code")
		w.WriteHeader(code)
	}

	return nil
}

func headerMatcher(key string) (string, bool) {
	switch key {
	case "X-User-Id":
		return key, true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}

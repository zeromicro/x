package http_test

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest"
	xerrors "github.com/zeromicro/x/errors"
	xhttp "github.com/zeromicro/x/http"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type message struct {
	Name string `json:"name" xml:"name"`
}

func Example_jsonBaseResponse() {
	data := []byte(`{"name":"JsonBaseResponse.example","port":8080}`)
	var serverConf rest.RestConf
	if err := conf.LoadFromJsonBytes(data, &serverConf); err != nil {
		logx.Must(err)
	}

	server, err := rest.NewServer(serverConf)
	if err != nil {
		logx.Must(err)
	}

	server.AddRoutes([]rest.Route{
		{
			Method: http.MethodGet,
			Path:   "/code/msg",
			Handler: func(writer http.ResponseWriter, request *http.Request) {
				// expected output:
				// {"code":1,"msg":"dummy error"}
				xhttp.JsonBaseResponse(writer, xerrors.New(1, "dummy error"))
			},
		},
		{
			Method: http.MethodGet,
			Path:   "/grpc/status",
			Handler: func(writer http.ResponseWriter, request *http.Request) {
				// expected output:
				// {"code":0,"msg":"ok"}
				xhttp.JsonBaseResponse(writer, status.New(codes.OK, "ok"))
			},
		},
		{
			Method: http.MethodGet,
			Path:   "/error",
			Handler: func(writer http.ResponseWriter, request *http.Request) {
				// expected output:
				// {"code":-1,"msg":"dummy error"}
				xhttp.JsonBaseResponse(writer, errors.New("dummy error"))
			},
		},
		{
			Method: http.MethodGet,
			Path:   "/struct",
			Handler: func(writer http.ResponseWriter, request *http.Request) {
				// expected output:
				// {"code":0,"msg":"ok","data":{"name":"anyone"}}
				xhttp.JsonBaseResponse(writer, message{
					Name: "anyone",
				})
			},
		},
	})

	defer server.Stop()
	fmt.Printf("Starting server at %s:%d...\n", serverConf.Host, serverConf.Port)
	server.Start()
}

func Example_xmlBaseResponse() {
	data := []byte(`{"name":"JsonBaseResponse.example","port":8080}`)
	var serverConf rest.RestConf
	if err := conf.LoadFromJsonBytes(data, &serverConf); err != nil {
		logx.Must(err)
	}

	server, err := rest.NewServer(serverConf)
	if err != nil {
		logx.Must(err)
	}

	server.AddRoutes([]rest.Route{
		{
			Method: http.MethodGet,
			Path:   "/code/msg",
			Handler: func(writer http.ResponseWriter, request *http.Request) {
				// expected output:
				// <xml version="1.0" encoding="UTF-8"><code>1</code><msg>dummy error</msg></xml>
				xhttp.XmlBaseResponse(writer, xerrors.New(1, "dummy error"))
			},
		},
		{
			Method: http.MethodGet,
			Path:   "/grpc/status",
			Handler: func(writer http.ResponseWriter, request *http.Request) {
				// expected output:
				// <xml version="1.0" encoding="UTF-8"><code>0</code><msg>ok</msg></xml>
				xhttp.XmlBaseResponse(writer, status.New(codes.OK, "ok"))
			},
		},
		{
			Method: http.MethodGet,
			Path:   "/error",
			Handler: func(writer http.ResponseWriter, request *http.Request) {
				// expected output:
				// <xml version="1.0" encoding="UTF-8"><code>-1</code><msg>dummy error</msg></xml>
				xhttp.XmlBaseResponse(writer, errors.New("dummy error"))
			},
		},
		{
			Method: http.MethodGet,
			Path:   "/struct",
			Handler: func(writer http.ResponseWriter, request *http.Request) {
				// expected output:
				// <xml version="1.0" encoding="UTF-8"><code>0</code><msg>ok</msg><data><name>anyone</name></data></xml>
				xhttp.XmlBaseResponse(writer, message{Name: "anyone"})
			},
		},
	})

	defer server.Stop()
	fmt.Printf("Starting server at %s:%d...\n", serverConf.Host, serverConf.Port)
	server.Start()
}

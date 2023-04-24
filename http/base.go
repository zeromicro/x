package http

import (
	"context"
	"encoding/xml"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"github.com/zeromicro/x/errors"
	"google.golang.org/grpc/status"
)

// BaseResponse is the base response struct.
type BaseResponse[T any] struct {
	// Code represents the business code, not the http status code.
	Code int `json:"code" xml:"code"`
	// Msg represents the business message, if Code = BusinessCodeOK,
	// and Msg is empty, then the Msg will be set to BusinessMsgOk.
	Msg string `json:"msg" xml:"msg"`
	// Data represents the business data.
	Data T `json:"data,omitempty" xml:"data,omitempty"`
}

type baseXmlResponse[T any] struct {
	XMLName  xml.Name `xml:"xml"`
	Version  string   `xml:"version,attr"`
	Encoding string   `xml:"encoding,attr"`
	BaseResponse[T]
}

// JsonBaseResponse writes v into w with http.StatusOK.
func JsonBaseResponse(w http.ResponseWriter, v any) {
	httpx.OkJson(w, wrapBaseResponse(v))
}

// JsonBaseResponseCtx writes v into w with http.StatusOK.
func JsonBaseResponseCtx(ctx context.Context, w http.ResponseWriter, v any) {
	httpx.OkJsonCtx(ctx, w, wrapBaseResponse(v))
}

// XmlBaseResponse writes v into w with http.StatusOK.
func XmlBaseResponse(w http.ResponseWriter, v any) {
	OkXml(w, wrapXmlBaseResponse(v))
}

// XmlBaseResponseCtx writes v into w with http.StatusOK.
func XmlBaseResponseCtx(ctx context.Context, w http.ResponseWriter, v any) {
	OkXmlCtx(ctx, w, wrapXmlBaseResponse(v))
}

func wrapXmlBaseResponse(v any) baseXmlResponse[any] {
	base := wrapBaseResponse(v)
	return baseXmlResponse[any]{
		Version:      xmlVersion,
		Encoding:     xmlEncoding,
		BaseResponse: base,
	}
}

func wrapBaseResponse(v any) BaseResponse[any] {
	var resp BaseResponse[any]
	switch data := v.(type) {
	case *errors.CodeMsg:
		resp.Code = data.Code
		resp.Msg = data.Msg
	case errors.CodeMsg:
		resp.Code = data.Code
		resp.Msg = data.Msg
	case *status.Status:
		resp.Code = int(data.Code())
		resp.Msg = data.Message()
	case error:
		resp.Code = BusinessCodeError
		resp.Msg = data.Error()
	default:
		resp.Code = BusinessCodeOK
		resp.Msg = BusinessMsgOk
		resp.Data = v
	}

	return resp
}

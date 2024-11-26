package http

import (
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/x/search"
	"net/http"
)

type WildcardHandler struct {
	methodTree map[string]*search.Tree
}

func NewWildcardHandler(routes []rest.Route) (*WildcardHandler, error) {
	methodTree := make(map[string]*search.Tree)
	for _, route := range routes {
		tree, ok := methodTree[route.Method]
		if !ok {
			tree = search.NewTree()
		}
		if err := tree.Add(route.Path, route.Handler); err != nil {
			return nil, err
		}
		methodTree[route.Method] = tree
	}

	return &WildcardHandler{
		methodTree: methodTree,
	}, nil
}

func (w *WildcardHandler) WildcardOption() rest.RunOption {
	return rest.WithNotFoundHandler(w.notFoundHandler())
}

func (w *WildcardHandler) notFoundHandler() http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		tree, ok := w.methodTree[request.Method]
		if !ok {
			writer.WriteHeader(http.StatusNotFound)
			return
		}
		result, ok := tree.Search(request.URL.Path)
		if !ok {
			writer.WriteHeader(http.StatusNotFound)
			return
		}
		result.Item.(http.HandlerFunc)(writer, request)
		return
	})
}

package http

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/rest/pathvar"
)

type TestRequest struct {
	Id    int64  `path:"id"`
	Name  string `form:"name"`
	Age   int64  `json:"age"`
	Token string `header:"token"`
}

func (t TestRequest) Validate(_ *http.Request, data any) error {
	val, ok := data.(*TestRequest)
	if !ok {
		return errors.New("data is not TestRequest")
	}

	if *val == t {
		return nil
	}
	return fmt.Errorf("expected %v, got %v", t, val)
}

func TestParse(t *testing.T) {
	var m = map[string]int64{"age": 20}
	data, err := json.Marshal(m)
	assert.NoError(t, err)

	r, err := http.NewRequest(http.MethodPost, "http://localhost:8080/mock?name=foo", bytes.NewBuffer(data))
	assert.NoError(t, err)
	r = pathvar.WithVars(r, map[string]string{"id": "123"})
	r.Header.Set("Content-Type", "application/json")
	r.Header.Set("token", "mock_token")

	var req TestRequest
	err = Parse(r, &req)
	assert.NoError(t, err)
	assert.Equal(t, TestRequest{
		Id:    123,
		Name:  "foo",
		Age:   20,
		Token: "mock_token",
	}, req)
}

package http

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/jinzhu/copier"
	"github.com/zeromicro/go-zero/core/mapping"
	"github.com/zeromicro/go-zero/rest/httpx"
)

const (
	pathKey   = "path"
	formKey   = "form"
	jsonKey   = "json"
	headerKey = "header"
)

var (
	errValueNotSettable     = errors.New("value is not settable")
	errInstanceNotInterface = errors.New("value can not interface")
	tagKeys                 = []string{pathKey, formKey, jsonKey, headerKey}
)

// The Parse function parses the request into the given value. if the value implements
// the httpx.Validator interface, it will be validated, or else it will be parsed by httpx.Parse.
// NOTE: if the value is httpx.Validator, all tags of go-zero will be overridden by optional tag
// to skip the validation in go-zero, it causes the go-zero validator to be invalid.
func Parse(r *http.Request, v any) error {
	validate, ok := v.(httpx.Validator)
	if !ok {
		return httpx.Parse(r, v)
	}

	valueType := reflect.TypeOf(v)
	if valueType.Kind() != reflect.Ptr {
		return errValueNotSettable
	}
	elemType := mapping.Deref(valueType)
	if elemType.Kind() != reflect.Struct {
		return httpx.Parse(r, v)
	}

	numField := elemType.NumField()
	var fields = make([]reflect.StructField, 0, numField)
	for i := 0; i < numField; i++ {
		field := makeOptional(elemType.Field(i))
		fields = append(fields, field)
	}

	reflectStruct := reflect.StructOf(fields)
	newInstance := reflect.New(reflectStruct)
	if !newInstance.CanInterface() {
		return errInstanceNotInterface
	}

	if err := httpx.Parse(r, newInstance.Interface()); err != nil {
		return err
	}

	if err := copier.Copy(v, newInstance.Interface()); err != nil {
		return err
	}

	return validate.Validate(r, v)
}

func makeOptional(field reflect.StructField) reflect.StructField {
	var tags []string
	for _, tagKey := range tagKeys {
		tag, ok := field.Tag.Lookup(tagKey)
		if !ok {
			continue
		}

		tags = append(tags, fmt.Sprintf(`%s:"%s,optional"`, tagKey, tag))
	}

	field.Tag = reflect.StructTag(strings.Join(tags, " "))
	return field
}

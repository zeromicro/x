package errors

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	ast := assert.New(t)
	c := New(1, "test")
	cm, ok := c.(*CodeMsg)
	ast.True(ok)
	ast.NotNil(cm)
	ast.Equal(int(1), cm.Code)
	ast.Equal("test", cm.Msg)
}

func TestCodeMsg_Error(t *testing.T) {
	ast := assert.New(t)
	c := New(1, "test")
	cm, ok := c.(*CodeMsg)
	ast.True(ok)
	ast.NotNil(cm)
	ast.NotEmpty(cm.Error())
}

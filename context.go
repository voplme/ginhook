package ginhook

import (
	"io"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

type Context struct {
	*gin.Context
	isSend  bool
	repData map[string]interface{}
	Token   string
}

func (c *Context) GetToken() string {
	c.Token = c.Request.Header.Get("token")
	if c.Token == "" {
		c.Token = c.Query("tk")
	}
	return c.Token
}

func (c *Context) GetData() ([]byte, error) {
	defer c.Request.Body.Close()
	return io.ReadAll(c.Request.Body)
}

func (c *Context) GetMapData() map[string]interface{} {
	if c.repData == nil {
		c.repData = map[string]interface{}{}
		if err := c.GetJSON(&c.repData); err != nil {
			return nil
		}
	}
	return c.repData
}

func (c *Context) GetKey(key string) interface{} {
	c.GetMapData()
	if v, ok := c.repData[key]; ok {
		return v
	}
	return nil
}

func (c *Context) GetKeyString(key string, err string) string {
	v := c.GetKey(key)
	if v == nil {
		ThrowError(err)
	}
	if data, ok := v.(string); ok {
		return data
	}
	ThrowError(err)
	return ""
}

func (c *Context) GetKeyBool(key string, err string) bool {
	v := c.GetKey(key)
	if v == nil {
		ThrowError(err)
	}
	if data, ok := v.(bool); ok {
		return data
	}
	ThrowError(err)
	return false
}

func (c *Context) GetKeyInt(key string, err string) int64 {
	v := c.GetKey(key)
	if v == nil {
		ThrowError(err)
	}
	switch m := v.(type) {
	case int:
		return int64(m)
	case uint32:
		return int64(m)
	case int64:
		return int64(m)
	case float64:
		return int64(m)
	case uint64:
		return int64(m)
	case uint8:
		return int64(m)
	case int8:
		return int64(m)
	}
	ThrowError(err)
	return 0
}

func (c *Context) GetJSON(v any) error {
	if err := c.ShouldBindWith(v, binding.JSON); err != nil {
		ThrowError(err.Error())
		return err
	}
	return nil
}

func (c *Context) Result(data any, code int, msg any) {
	if c.isSend {
		return
	}
	c.isSend = true
	c.JSON(200, gin.H{
		"code": code,
		"msg":  msg,
		"data": data,
	})
	c.Abort()
}

func (c *Context) Success(data any) {
	c.Result(data, CODE_OK, nil)
}

func (c *Context) Fail(msg any) {
	if m, ok := msg.(*Exception); ok {
		c.Result(m.Data, m.Code, m.Msg)
		return
	}
	if m, ok := msg.(error); ok {
		c.Result(nil, CODE_ERROR, m.Error())
		return
	}
	c.Result(nil, CODE_ERROR, msg)
}

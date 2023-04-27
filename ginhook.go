package ginhook

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/gin-contrib/gzip"
	"github.com/gin-gonic/gin"
)

type HttpContext interface {
	Do(ctx *Context) error
}

type ApiNameCb func(s string) string
type RouterGroup struct {
	*gin.RouterGroup
	ApiName ApiNameCb
}

func snakeString(s string) string {
	data := make([]byte, 0, len(s)*2)
	j := false
	num := len(s)
	for i := 0; i < num; i++ {
		d := s[i]
		// or通过ASCII码进行大小写的转化
		// 65-90（A-Z），97-122（a-z）
		//判断如果字母为大写的A-Z就在前面拼接一个_
		if i > 0 && d >= 'A' && d <= 'Z' && j {
			data = append(data, '_')
		}
		if d != '_' {
			j = true
		}
		data = append(data, d)
	}
	//ToLower把大写字母统一转小写
	return strings.ToLower(string(data[:]))
}

// 自动绑定
func (c *RouterGroup) API(routes HttpContext) *RouterGroup {
	r := reflect.TypeOf(routes)
	for i := 0; i < r.NumMethod(); i++ {
		m := r.Method(i)
		if m.Type.NumIn() == 1 && m.Name != "Do" {
			func(funname string, i int, obj reflect.Type) {
				GApi := c.RouterGroup.POST
				if funname[0:3] == "Get" {
					funname = funname[3:]
					GApi = c.RouterGroup.GET
				} else if funname[0:4] == "Post" {
					funname = funname[4:]
				}
				if c.ApiName == nil {
					c.ApiName = snakeString
				}
				GApi(c.ApiName(funname), func(ctx *gin.Context) {
					ctx2 := &Context{
						Context: ctx,
					}
					defer Try(func(data Exception) {
						ctx2.Result(data.Data, data.Code, data.Msg)
					})
					result := reflect.New(obj.Elem()).Interface()

					if a9, ok := result.(HttpContext); ok {

						if e := a9.Do(ctx2); e != nil {
							ctx2.Fail(e)
							return
						}
					} else {
						ctx2.Fail("拒绝访问")
						return
					}
					fun := reflect.ValueOf(result).Method(i)
					funResult := fun.Call([]reflect.Value{})
					l := len(funResult)
					if l == 0 {
						return
					}
					if l == 1 {
						res := funResult[0].Interface()
						if e, ok := res.(error); ok {
							ctx2.Fail(e)
							return
						} else {
							ctx2.Success(res)
							return
						}
					}
					if l == 2 {
						e := funResult[1].Interface()
						if e != nil {
							ctx2.Fail(e)
							return
						}
						ctx2.Success(funResult[0].Interface())
					}
				})
			}(m.Name, i, r)
		}

	}
	return c
}

func NewRouterGroup(g *gin.RouterGroup, cb ApiNameCb) *RouterGroup {
	//在没有指定的情况下，使用默认值
	if cb == nil {
		cb = snakeString
	}
	return &RouterGroup{
		RouterGroup: g,
		ApiName:     cb,
	}
}

func New(log bool) *gin.Engine {
	r := gin.New()
	if log {
		r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
			// your custom format
			return fmt.Sprintf("%s - [%s] \"%s %s %d %s %s %d\"\n",
				param.ClientIP,
				param.TimeStamp.Format("2006/01/02 15:04:05"),
				param.Method,
				param.Request.Proto,
				param.StatusCode,
				param.Latency,
				param.Path,
				param.BodySize,
			)
		}))
	}
	r.Use(gzip.Gzip(gzip.DefaultCompression))
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*") // 可将将 * 替换为指定的域名
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
		c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization, Token, Authencation,TimeStamp")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(200)
		}
	})
	return r
}

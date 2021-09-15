package middlewares

import (
	"bytes"
	"encoding/json"

	"github.com/gin-gonic/gin"
)

type BodyWriter struct {
	gin.ResponseWriter
	bodyBuf *bytes.Buffer
}

func (w BodyWriter) Write(b []byte) (int, error) {
	return w.bodyBuf.Write(b)
}
func NewResponseMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		resBodyWriter := BodyWriter{
			bodyBuf:        &bytes.Buffer{},
			ResponseWriter: c.Writer,
		}

		c.Writer = resBodyWriter

		c.Next()

		lastErr := c.Errors.Last()

		var data interface{}
		if err := json.Unmarshal(resBodyWriter.bodyBuf.Bytes(), &data); err != nil {
			data = resBodyWriter.bodyBuf.String()
		}
		//TODO need to change this to jsend json format
		var res interface{}

		code := resBodyWriter.Status()

		raw := c.GetBool("raw")
		if raw {
			c.Writer = resBodyWriter.ResponseWriter
			c.AbortWithStatusJSON(resBodyWriter.Status(), data)
			return
		}

		if code >= 200 && code < 300 {
			res = gin.H{
				"status": "success",
				"data":   data,
			}
		} else {
			if lastErr == nil && code == 404 { // not found API
				res = gin.H{
					"status":  "fail",
					"code":    code,
					"message": "Not found",
				}
			} else if code >= 400 && code < 500 {
				res = gin.H{
					"status":  "fail",
					"code":    code,
					"message": lastErr.Error(),
				}
			} else {
				res = gin.H{
					"status":  "error",
					"code":    code,
					"message": lastErr.Error(),
				}
			}
		}

		c.Writer = resBodyWriter.ResponseWriter
		c.AbortWithStatusJSON(resBodyWriter.Status(), res)
	}
}

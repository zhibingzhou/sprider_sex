package thread

import "github.com/gin-gonic/gin"

func GetKey(c *gin.Context) {
	c.Writer.WriteString("6cae6e508ae9a8c9")
}

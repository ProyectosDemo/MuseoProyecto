package middleware
import (
	"fmt"
	"net/http"
	"github.com/gin-gonic/gin"
)

func PanicButton(err error) {
	if err != nil {
		panic(err)
	}
}

func ErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer handlePanic(c)
		c.Next()
	}
}

func handlePanic(c *gin.Context) {
	r := recover()
	if r != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": fmt.Sprint(r),
		})
		c.Abort()
	}
}
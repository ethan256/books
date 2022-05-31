package router

import (
	"github.com/gin-gonic/gin"

	"github.com/ethan256/books/internal/app/book"
)

func SetApiRouter(r *gin.Engine, h book.BookHandler) {
	v1 := r.Group("/books/v1")
	{
		v1.GET("/book/:name", h.DesribeBookInfo())
		v1.GET("/book", h.FindBooks())
		v1.POST("/book/:name", h.SaveBook())
		v1.PATCH("/book/:name", h.UpdateBook())
	}
}

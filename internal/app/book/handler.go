package book

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/ethan256/books/internal/code"
)

type BookHandler interface {
	DesribeBookInfo() gin.HandlerFunc
	FindBooks() gin.HandlerFunc
	UpdateBook() gin.HandlerFunc
	SaveBook() gin.HandlerFunc
}

var _ BookHandler = (*bookHandler)(nil)

type bookHandler struct {
	service BookService
}

// DesribeBookInfo 查看书籍详情
// @Summary 查看书籍详情
// @Description 查看书籍详情
// @Accept application/json
// @Produce json
// @Success 200 {object} Response
// @Failure 200 {object} Failure
// @Failure 400 {object} Failure
// @Router /books/v1/book/:name [get]
// @Security LoginToken
func (h *bookHandler) DesribeBookInfo() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")
		if name == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest,
				BuildFailure(code.ParamBindError, "name 不能为空"),
			)
			return
		}

		book, err := h.service.FindBookByName(c.Request.Context(), name)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK,
				BuildFailure(code.BookKindNotExistError, err.Error()),
			)
			return
		}
		c.JSON(http.StatusOK, &Response{0, "Success", book})
	}
}

// FindBooks 查找同种类的书籍
// @Summary 查找同种类的书籍
// @Description 查找同种类的书籍
// @Accept application/json
// @Produce json
// @Success 200 {object} detailResponse
// @Failure 400 {object} code.Failure
// @Router /books/v1/book [get]
// @Security LoginToken
func (h *bookHandler) FindBooks() gin.HandlerFunc {
	return func(c *gin.Context) {
		var kind int
		if err := c.ShouldBindJSON(&kind); err != nil || kind == 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest,
				BuildFailure(code.ParamBindError, "参数kind必选"),
			)
			return
		}

		books, err := h.service.ListBooksByKind(c.Request.Context(), kind)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusOK,
				BuildFailure(code.BookKindNotExistError, "参数kind必选"),
			)
			return
		}
		c.JSON(http.StatusOK, &Response{0, "Success", books})
	}
}

// SaveBook 保存书籍详
// @Summary 保存书籍详
// @Description 保存书籍详
// @Accept application/json
// @Produce json
// @Success 200 {object} Response
// @Failure 400 {object} Failure
// @Router /books/v1/book/:name [post]
// @Security LoginToken
func (h *bookHandler) SaveBook() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")
		if name == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest,
				BuildFailure(code.ParamBindError, "name 不能为空"),
			)
			return
		}

		var book Book
		if err := c.ShouldBindJSON(&book); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest,
				BuildFailure(code.ParamBindError, err.Error()),
			)
			return
		}

		if err := h.service.SaveBook(c.Request.Context(), &book); err != nil {
			c.AbortWithStatusJSON(http.StatusOK, BuildFailure(code.BookSaveError, err.Error()))
			return
		}
		c.JSON(http.StatusOK, &Response{0, "Success", nil})
	}
}

// SaveBook 更新书籍信息
// @Summary 更新书籍信息
// @Description 更新书籍信息
// @Accept application/json
// @Produce json
// @Success 200 {object} Response
// @Failure 400 {object} Failure
// @Router /books/v1/book/:name [patch]
// @Security LoginToken
func (h *bookHandler) UpdateBook() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Param("name")
		if name == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest,
				BuildFailure(code.ParamBindError, "name 不能为空"),
			)
			return
		}

		var book Book
		if err := c.ShouldBindJSON(&book); err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest,
				BuildFailure(code.ParamBindError, err.Error()),
			)
			return
		}

		if err := h.service.UpdateBook(c.Request.Context(), name, &book); err != nil {
			c.AbortWithStatusJSON(http.StatusOK, BuildFailure(code.BookSaveError, err.Error()))
			return
		}
		c.JSON(http.StatusOK, &Response{0, "Success", nil})
	}
}

func NewBookHandler(svc BookService) BookHandler {
	return &bookHandler{service: svc}
}

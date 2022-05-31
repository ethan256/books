package main

import (
	"github.com/gin-gonic/gin"

	"github.com/ethan256/books/configs"
	"github.com/ethan256/books/internal/app/book"
	"github.com/ethan256/books/internal/pkg/mysql"
	"github.com/ethan256/books/internal/pkg/redis"
	"github.com/ethan256/books/internal/router"
)

func main() {
	// 初始化Config
	if err := configs.InitConfig(); err != nil {
		panic(err)
	}

	// 初始化 MySQL 连接
	repo, err := mysql.New()
	if err != nil {
		panic(err)
	}

	// 初始化 redis 连接
	rdb, err := redis.NewClient()
	if err != nil {
		panic(err)
	}

	bookRepo := book.NewBookRepo(repo.GetDB(), rdb)
	bookSvc := book.NewBookService(bookRepo)
	bookHandler := book.NewBookHandler(bookSvc)

	// http 请求处理
	r := gin.Default()
	router.SetApiRouter(r, bookHandler)
	r.Run(configs.Get().Host)
}

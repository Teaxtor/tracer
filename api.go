package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type Api struct {
	server *http.Server
}

func NewApi(config Config) *Api {
	r := gin.New()

	r.GET("/health", healthEndpoint)
	r.GET("trace", traceEndpoint)

	return &Api{
		server: &http.Server{
			Addr: ":" + config.Port,
			Handler: r,
		},
	}
}

func (a *Api) Start() error {
	fmt.Println("starting server at ", a.server.Addr)

	return a.server.ListenAndServe()
}

func (a *Api) Stop () error {
	fmt.Println("stopping server")

	return a.server.Close()
}

func (c *gin.Context) healthEndpoint () {
	c.JSON(http.StatusOK, gin.H{})
}

func (c *gin.Context) traceEndpoint () {
	// form validation

	// pass to chrome

	// return response

	c.JSON(http.StatusOK, gin.H{})
}
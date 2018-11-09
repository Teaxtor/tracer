package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"tracer/pkg"
)

type App struct {
	server *http.Server
	tracer *pkg.Tracer
}

type Form struct {
	Proxy   string `form:"proxy" binding:"omitempty"`
	Url     string `form:"url" binding:"required"`
}

func New(config Config) *App {
	app := &App{}

	app.tracer = pkg.New(config.Browser, config.ProxyInfo, config.RemotePort)

	r := gin.New()

	r.GET("/health", healthEndpoint())
	r.GET("trace", traceEndpoint(app.tracer))

	app.server = &http.Server{
		Addr: ":" + strconv.Itoa(config.Port),
		Handler: r,
	}

	return app
}

func (a *App) Start() error {
	fmt.Println("starting server at ", a.server.Addr)

	return a.server.ListenAndServe()
}

func (a *App) Stop () {
	fmt.Println("stopping server")
	a.server.Close()

	fmt.Println("stopping tracer")
	a.tracer.Stop()
}

func healthEndpoint () gin.HandlerFunc {
	return func (c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{})
	}
}

func traceEndpoint (tracer *pkg.Tracer) gin.HandlerFunc {
	return func (c *gin.Context) {
		var form Form

		err := c.ShouldBindQuery(&form)

		if err != nil {
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		config := pkg.TraceConfig{
			Url:     form.Url,
			Proxy:   form.Proxy,
		}

		// pass to chrome
		result, err := tracer.Trace(config)

		// return response
		if err != nil {
			c.JSON(http.StatusInternalServerError, err)
			return
		}

		c.JSON(http.StatusOK, result)
	}
}
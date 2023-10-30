package main

// serve as a http server

import (
	"fmt"

	_ "embed"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/cobra"
)

type (
	httpServer struct {
		port string
	}
)

var (
	serveCmd = &cobra.Command{
		Use:   "serve --port=8000",
		Short: "starts a http server",
		Run: func(cmd *cobra.Command, args []string) {
			server.start()
		},
	}
	server = &httpServer{}
)

func init() {
	serveCmd.PersistentFlags().StringVarP(&server.port,
		"port", "p", "8000", "listening port",
	)
	rootCmd.AddCommand(serveCmd)
}

func (s *httpServer) start() error {
	e := echo.New()
	e.Use(middleware.Recover())

	// handler goes here
	e.GET("/", s.index)

	fmt.Println("serving on port: " + s.port)
	e.Logger.Fatal(e.Start(":" + s.port))
	return nil
}

//go:embed index.html
var indexHtml string

func (s *httpServer) index(c echo.Context) error {
	return c.HTML(200, indexHtml)
}

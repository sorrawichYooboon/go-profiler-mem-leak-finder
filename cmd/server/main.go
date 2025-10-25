package main

import (
	_ "go-profiler-mem-leak-finder/docs"
	"go-profiler-mem-leak-finder/internal/handler"
	"net/http"
	_ "net/http/pprof"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

// @title Go Profiler Memory Leak Finder API
// @version 1.0
// @description This is a sample server for a memory leak finder.

// @host localhost:8080
// @BasePath /
func main() {
	e := echo.New()

	// Register pprof handlers
	e.GET("/debug/pprof/*", echo.WrapHandler(http.DefaultServeMux))

	e.GET("/start-leak/:type", handler.StartLeak)
	e.GET("/leak-test", handler.LeakTest)
	e.GET("/safe-test", handler.SafeTest)
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}

package handler

import (
	"fmt"
	"go-profiler-mem-leak-finder/internal/domain"
	"go-profiler-mem-leak-finder/internal/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

// StartLeak godoc
// @Summary Start a memory leak
// @Description Start a memory leak of a specific type
// @Tags leaks
// @Accept  json
// @Produce  json
// @Param type path string true "Leak type" Enums(slice, channel, mutex)
// @Success 200 {string} string "Leak started"
// @Router /start-leak/{type} [get]
func StartLeak(c echo.Context) error {
	switch c.Param("type") {
	case "slice":
		service.RunLeakerJob1_SimpleSlice()
		return c.String(http.StatusOK, "Slice leak started")
	case "channel":
		service.RunLeakerJob2_ChannelGoroutineBlock()
		return c.String(http.StatusOK, "Channel leak started")
	case "mutex":
		service.RunLeakerJob3_MutexMap()
		return c.String(http.StatusOK, "Mutex leak started")
	default:
		return c.String(http.StatusNotFound, "Unknown leak type")
	}
}

// LeakTest godoc
// @Summary Check the status of the leaks
// @Description Check the status of the leaks
// @Tags leaks
// @Accept  json
// @Produce  json
// @Success 200 {string} string "Leak status"
// @Router /leak-test [get]
func LeakTest(c echo.Context) error {
	msg := fmt.Sprintf("LeakyStore1 (Slice): %d MB\n", len(domain.LeakyStore1))
	msg += fmt.Sprintf("LeakyStore3 (Map): %d MB\n", len(domain.LeakyStore3))
	return c.String(http.StatusOK, msg)
}

// SafeTest godoc
// @Summary Start a safe goroutine
// @Description Start a safe goroutine that does not leak
// @Tags leaks
// @Accept  json
// @Produce  json
// @Success 200 {string} string "Safe goroutine started"
// @Router /safe-test [get]
func SafeTest(c echo.Context) error {
	service.RunSafeJob()
	return c.String(http.StatusOK, "Safe goroutine started")
}

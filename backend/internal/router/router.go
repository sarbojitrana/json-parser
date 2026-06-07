package router

import (
	"parser/internal/handler"

	"github.com/labstack/echo/v4"
)

func Register(
	e *echo.Echo,
	h *handler.Handler,
) {

	api := e.Group("/api")

	api.POST(
		"/spreadsheet",
		h.UploadSpreadsheet,
	)

	api.POST(
		"/workflow",
		h.UploadWorkflow,
	)

	api.GET(
		"/applications/:appl_id",
		h.GetApplication,
	)

	api.DELETE(
		"/applications/:appl_id",
		h.DeleteApplication,
	)
}
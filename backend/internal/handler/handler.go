package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"

	"parser/internal/model"
	"parser/internal/repository"
	"parser/internal/utils"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	repo *repository.Repository
}

func New(repo *repository.Repository) *Handler {
	return &Handler{
		repo: repo,
	}
}

func (h *Handler) GetApplication(c echo.Context) error {

	applID, err := strconv.ParseInt(
		c.Param("appl_id"),
		10,
		64,
	)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{
				"error": "invalid application id",
			},
		)
	}

	serviceID, err := strconv.ParseInt(
		c.QueryParam("service_id"),
		10,
		64,
	)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{
				"error": "invalid service id",
			},
		)
	}

	rootType := c.QueryParam("root_type")

	if rootType != "initiated" &&
		rootType != "execution" {

		return c.JSON(
			http.StatusBadRequest,
			map[string]string{
				"error": "root_type must be initiated or execution",
			},
		)
	}

	event, err := h.repo.GetWorkflowEvent(
		c.Request().Context(),
		applID,
		serviceID,
		rootType,
	)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{
				"error": err.Error(),
			},
		)
	}

	if event == nil {
		return c.JSON(
			http.StatusNotFound,
			map[string]string{
				"error": "workflow event not found",
			},
		)
	}

	mappings, err := h.repo.GetMappingsByServiceID(
		c.Request().Context(),
		serviceID,
	)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{
				"error": err.Error(),
			},
		)
	}

	var payload any

	if err := json.Unmarshal(
		event.RawPayload,
		&payload,
	); err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{
				"error": err.Error(),
			},
		)
	}

	resolver := utils.NewResolver(
		mappings,
	)

	resolved := resolver.Resolve(
		payload,
	)

	return c.JSON(
		http.StatusOK,
		map[string]any{
			"application_id": event.ApplID,
			"service_id":     event.ServiceID,
			"root_type":      event.RootType,
			"payload":        resolved,
		},
	)
}

func (h *Handler) DeleteApplication(c echo.Context) error {

	applID, err := strconv.ParseInt(
		c.Param("appl_id"),
		10,
		64,
	)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{
				"error": "invalid application id",
			},
		)
	}

	serviceID, err := strconv.ParseInt(
		c.QueryParam("service_id"),
		10,
		64,
	)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{
				"error": "invalid service id",
			},
		)
	}

	rootType := c.QueryParam("root_type")

	if rootType != "initiated" &&
		rootType != "execution" {

		return c.JSON(
			http.StatusBadRequest,
			map[string]string{
				"error": "root_type must be initiated or execution",
			},
		)
	}

	err = h.repo.DeleteApplication(
		c.Request().Context(),
		applID,
		serviceID,
		rootType,
	)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{
				"error": err.Error(),
			},
		)
	}

	return c.JSON(
		http.StatusOK,
		map[string]string{
			"message": "application deleted",
		},
	)
}

func (h *Handler) UploadSpreadsheet(c echo.Context) error {

	serviceID, err := strconv.ParseInt(
		c.FormValue("service_id"),
		10,
		64,
	)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{
				"error": "invalid service id",
			},
		)
	}

	serviceName := c.FormValue("service_name")

	if serviceName == "" {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{
				"error": "service name is required",
			},
		)
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{
				"error": "spreadsheet file is required",
			},
		)
	}

	src, err := file.Open()
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{
				"error": err.Error(),
			},
		)
	}
	defer src.Close()

	tmpFile, err := os.CreateTemp("", "*.xlsx")
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{
				"error": err.Error(),
			},
		)
	}
	defer tmpFile.Close()
	defer os.Remove(tmpFile.Name())

	if _, err := io.Copy(tmpFile, src); err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{
				"error": err.Error(),
			},
		)
	}

	sheet, err := utils.LoadSpreadsheet(
		tmpFile.Name(),
	)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{
				"error": err.Error(),
			},
		)
	}

	err = h.repo.CreateService(
		c.Request().Context(),
		model.Service{
			ServiceID:   serviceID,
			ServiceName: serviceName,
		},
	)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{
				"error": err.Error(),
			},
		)
	}

	for _, attr := range sheet.Attributes {

		err := h.repo.CreateMapping(
			c.Request().Context(),
			model.ServiceMapping{
				ServiceID:   serviceID,
				SectionName: attr.SectionName,
				SectionID:   attr.SectionID,
				FieldID:     attr.AttributeID,
				FieldName:   attr.Label,
				InputType:   attr.InputType,
				FieldSetID:  attr.FieldSetID,
			},
		)

		if err != nil {
			return c.JSON(
				http.StatusInternalServerError,
				map[string]string{
					"error": err.Error(),
				},
			)
		}
	}

	return c.JSON(
		http.StatusCreated,
		map[string]string{
			"message": "spreadsheet uploaded",
		},
	)
}

func (h *Handler) UploadWorkflow(c echo.Context) error {

	rootType := c.FormValue("root_type")

	if rootType != "initiated" &&
		rootType != "execution" {

		return c.JSON(
			http.StatusBadRequest,
			map[string]string{
				"error": "root_type must be initiated or execution",
			},
		)
	}

	file, err := c.FormFile("file")
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{
				"error": "workflow file is required",
			},
		)
	}

	src, err := file.Open()
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{
				"error": err.Error(),
			},
		)
	}
	defer src.Close()

	rawJSON, err := io.ReadAll(src)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{
				"error": err.Error(),
			},
		)
	}

	parser := utils.NewParser()

	events, err := parser.Parse(
		rawJSON,
		rootType,
	)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{
				"error": err.Error(),
			},
		)
	}

	for _, event := range events {

		_, err := h.repo.CreateWorkflowEvent(
			c.Request().Context(),
			event,
		)

		if err != nil {
			return c.JSON(
				http.StatusInternalServerError,
				map[string]string{
					"error": err.Error(),
				},
			)
		}
	}

	return c.JSON(
		http.StatusCreated,
		map[string]any{
			"message": "workflow uploaded",
			"count":   len(events),
		},
	)
}

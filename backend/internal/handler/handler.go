package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"

	"parser/internal/model"
	"parser/internal/repository"
	"parser/internal/utils"

	"github.com/labstack/echo/v4"
)

type Handler struct {
	repo *repository.Repository
}

const (
	LogLevelError = "ERROR"
	LogLevelInfo  = "INFO"
)

func New(repo *repository.Repository) *Handler {
	return &Handler{
		repo: repo,
	}
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

	if rootType != "initiated_data" &&
		rootType != "execution_data" {

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

		_ = h.repo.CreateLog(
			c.Request().Context(),
			model.Log{
				Level:   LogLevelError,
				Source:  "DeleteApplication",
				Message: err.Error(),
			},
		)

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

func (h *Handler) UploadWorkflow(c echo.Context) error {

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
		_ = h.repo.CreateLog(
			c.Request().Context(),
			model.Log{
				Level:   LogLevelError,
				Source:  "UploadWorkflow",
				Message: err.Error(),
			},
		)

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

		_ = h.repo.CreateLog(
			c.Request().Context(),
			model.Log{
				Level:   LogLevelError,
				Source:  "UploadWorkflow",
				Message: err.Error(),
			},
		)

		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{
				"error": err.Error(),
			},
		)
	}

	var payload utils.Payload

	if err := json.Unmarshal(
		rawJSON,
		&payload,
	); err != nil {

		_ = h.repo.CreateLog(
			c.Request().Context(),
			model.Log{
				Level:   LogLevelError,
				Source:  "UploadWorkflow",
				Message: err.Error(),
			},
		)

		return c.JSON(
			http.StatusBadRequest,
			map[string]string{
				"error": err.Error(),
			},
		)
	}

	attributeIDs := make(
		map[int64]model.AttributeIDs,
	)

	validServices := make(
		map[int64]bool,
	)

	serviceGroups := make(
		map[int64]bool,
	)

	for _, item := range payload.InitiatedData {

		serviceID := utils.ToInt64(
			item["service_id"],
		)

		if _, ok := validServices[serviceID]; ok {
			continue
		}
		serviceGroupID,_ := strconv.ParseInt(strconv.Itoa(int(serviceID))[:4], 10, 64 )

		exists, ok := serviceGroups[serviceGroupID]

		if !ok {

			exists, err = h.repo.ServiceGroupExists(
				c.Request().Context(),
				serviceGroupID,
			)
			if err != nil {

				_ = h.repo.CreateLog(
					c.Request().Context(),
					model.Log{
						Level:   LogLevelError,
						Source:  "UploadWorkflow",
						Message: err.Error(),
					},
				)

				return c.JSON(
					http.StatusInternalServerError,
					map[string]string{
						"error": err.Error(),
					},
				)
			}

			serviceGroups[serviceGroupID] = exists
		}

		validServices[serviceID] = exists

		if !exists {
			continue
		}

		if _, ok := attributeIDs[serviceGroupID]; ok {
			continue
		}

		ids, err := h.repo.GetAttributeIDs(
			c.Request().Context(),
			serviceGroupID,
		)
		if err != nil {

			_ = h.repo.CreateLog(
				c.Request().Context(),
				model.Log{
					Level:   LogLevelError,
					Source:  "UploadWorkflow",
					Message: err.Error(),
				},
			)

			return c.JSON(
				http.StatusInternalServerError,
				map[string]string{
					"error": err.Error(),
				},
			)
		}

		attributeIDs[serviceGroupID] = ids
	}

	parser := utils.NewParser()

	result, err := parser.Parse(
		rawJSON,
		attributeIDs,
	)
	if err != nil {

		_ = h.repo.CreateLog(
			c.Request().Context(),
			model.Log{
				Level:   LogLevelError,
				Source:  "UploadWorkflow",
				Message: err.Error(),
			},
		)

		return c.JSON(
			http.StatusBadRequest,
			map[string]string{
				"error": err.Error(),
			},
		)
	}

	skipped := 0
	workflowsCount := 0
	applicationsCount := 0

	for _, event := range result.WorkflowEvents {

		if !validServices[event.ServiceID] {
			skipped++
			continue
		}

		_, err := h.repo.CreateWorkflowEvent(
			c.Request().Context(),
			event,
		)
		if err != nil {

			_ = h.repo.CreateLog(
				c.Request().Context(),
				model.Log{
					Level:   LogLevelError,
					Source:  "UploadWorkflow",
					Message: err.Error(),
				},
			)

			return c.JSON(
				http.StatusInternalServerError,
				map[string]string{
					"error": err.Error(),
				},
			)
		}

		workflowsCount++
	}

	for _, app := range result.Application {

		if !validServices[app.ServiceID] {
			continue
		}

		err := h.repo.CreateApplication(
			c.Request().Context(),
			app,
		)
		if err != nil {

			_ = h.repo.CreateLog(
				c.Request().Context(),
				model.Log{
					Level:   LogLevelError,
					Source:  "UploadWorkflow",
					Message: err.Error(),
				},
			)

			return c.JSON(
				http.StatusInternalServerError,
				map[string]string{
					"error": err.Error(),
				},
			)
		}

		applicationsCount++
	}

	return c.JSON(
		http.StatusCreated,
		map[string]any{
			"message":      "workflow uploaded",
			"events":       workflowsCount,
			"applications": applicationsCount,
			"skipped":      skipped,
		},
	)
}

func (h *Handler) UploadSpreadsheet(c echo.Context) error {

	serviceGroupID, err := strconv.ParseInt(
		c.FormValue("service_group_id"),
		10,
		64,
	)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{
				"error": "invalid service group id",
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

		_ = h.repo.CreateLog(
			c.Request().Context(),
			model.Log{
				Level:   LogLevelError,
				Source:  "UploadSpreadsheet",
				Message: err.Error(),
			},
		)

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
		_ = h.repo.CreateLog(
			c.Request().Context(),
			model.Log{
				Level:   LogLevelError,
				Source:  "UploadSpreadsheet",
				Message: err.Error(),
			},
		)
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
		_ = h.repo.CreateLog(
			c.Request().Context(),
			model.Log{
				Level:   LogLevelError,
				Source:  "UploadSpreadsheet",
				Message: err.Error(),
			},
		)
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

		_ = h.repo.CreateLog(
			c.Request().Context(),
			model.Log{
				Level:   LogLevelError,
				Source:  "UploadSpreadsheet",
				Message: err.Error(),
			},
		)

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
			ServiceGroupID: serviceGroupID,
			ServiceName:    serviceName,
		},
	)
	if err != nil {
		_ = h.repo.CreateLog(
			c.Request().Context(),
			model.Log{
				Level:   LogLevelError,
				Source:  "UploadSpreadsheet",
				Message: err.Error(),
			},
		)
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
				ServiceGroupID: serviceGroupID,
				SectionName:    attr.SectionName,
				SectionID:      attr.SectionID,
				FieldID:        attr.AttributeID,
				FieldName:      attr.Label,
				InputType:      attr.InputType,
				FieldSetID:     attr.FieldSetID,
			},
		)

		if err != nil {
			_ = h.repo.CreateLog(
				c.Request().Context(),
				model.Log{
					Level:   LogLevelError,
					Source:  "UploadSpreadsheet",
					Message: err.Error(),
				},
			)
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

func (h *Handler) GetApplications(
	c echo.Context,
) error {

	fromStr := c.QueryParam("from")
	toStr := c.QueryParam("to")

	if fromStr == "" || toStr == "" {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{
				"error": "from and to dates are required",
			},
		)
	}

	from, err := time.Parse(
		"2006-01-02",
		fromStr,
	)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{
				"error": "invalid from date",
			},
		)
	}

	to, err := time.Parse(
		"2006-01-02",
		toStr,
	)
	if err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{
				"error": "invalid to date",
			},
		)
	}

	from = time.Date(
		from.Year(),
		from.Month(),
		from.Day(),
		0, 0, 0, 0,
		time.UTC,
	)

	to = time.Date(
		to.Year(),
		to.Month(),
		to.Day(),
		23, 59, 59, 999999999,
		time.UTC,
	)

	page := 1
	if value := c.QueryParam("page"); value != "" {

		page, err = strconv.Atoi(value)
		if err != nil || page < 1 {
			return c.JSON(
				http.StatusBadRequest,
				map[string]string{
					"error": "invalid page",
				},
			)
		}
	}

	limit := 10
	if value := c.QueryParam("limit"); value != "" {

		limit, err = strconv.Atoi(value)
		if err != nil || limit < 1 {
			return c.JSON(
				http.StatusBadRequest,
				map[string]string{
					"error": "invalid limit",
				},
			)
		}

		if limit > 100 {
			limit = 100
		}
	}

	applications, err := h.repo.GetApplications(
		c.Request().Context(),
		from,
		to,
		page,
		limit,
	)
	if err != nil {

		_ = h.repo.CreateLog(
			c.Request().Context(),
			model.Log{
				Level:   LogLevelError,
				Source:  "GetApplications",
				Message: err.Error(),
			},
		)

		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{
				"error": err.Error(),
			},
		)
	}

	return c.JSON(
		http.StatusOK,
		applications,
	)
}

func (h *Handler) GetLogs(
	c echo.Context,
) error {

	logs, err := h.repo.GetLogs(
		c.Request().Context(),
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
		logs,
	)
}

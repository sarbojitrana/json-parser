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

type ParseResult struct {
	WorkflowEvents []model.WorkflowEvent
	InitiatedApps  []model.ApplicationInitiated
	ExecutionApps  []model.ApplicationExecution
}

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

	if rootType != "initiated_data" &&
		rootType != "execution_data" {

		return c.JSON(
			http.StatusBadRequest,
			map[string]string{
				"error": "root_type must be initiated or execution",
			},
		)
	}

	events, err := h.repo.GetWorkflowEvents(
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
				Source:  "GetApplication",
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

	if events == nil {
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

		_ = h.repo.CreateLog(
			c.Request().Context(),
			model.Log{
				Level:   LogLevelError,
				Source:  "GetApplication",
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

	resolver := utils.NewResolver(
		mappings,
	)

	resolvedEvents := make(
		[]map[string]any,
		0,
		len(events),
	)

	for _, event := range events {

		var payload any

		if err := json.Unmarshal(
			event.RawPayload,
			&payload,
		); err != nil {

			_ = h.repo.CreateLog(
				c.Request().Context(),
				model.Log{
					Level:   LogLevelError,
					Source:  "GetApplication",
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

		resolved := resolver.Resolve(
			payload,
		)

		resolvedEvents = append(
			resolvedEvents,
			map[string]any{
				"id":            event.ID,
				"task_name":     event.TaskName,
				"action_no":     event.ActionNo,
				"task_type":     event.TaskType,
				"received_time": event.ReceivedTime,
				"executed_time": event.ExecutedTime,
				"payload":       resolved,
			},
		)
	}

	return c.JSON(
		http.StatusOK,
		map[string]any{
			"application_id": applID,
			"service_id":     serviceID,
			"root_type":      rootType,
			"events":         resolvedEvents,
		},
	)
}

func (h *Handler) GetApplicationAction(
	c echo.Context,
) error {

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

	if c.Param("action_no") == "" {
		application, err := h.repo.GetApplicationInitiated(
			c.Request().Context(),
			applID,
			serviceID,
		)
		if err != nil {

			_ = h.repo.CreateLog(
				c.Request().Context(),
				model.Log{
					Level:   LogLevelError,
					Source:  "GetApplicationAction",
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

		if application == nil {
			return c.JSON(
				http.StatusNotFound,
				map[string]string{
					"error": "application not found",
				},
			)
		}
		return c.JSON(
			http.StatusOK,
			application,
		)
	} else {
		actionNo, err := strconv.Atoi(
			c.Param("action_no"),
		)
		if err != nil {
			return c.JSON(
				http.StatusInternalServerError,
				map[string]string{
					"error": err.Error(),
				},
			)
		}
		application, err := h.repo.GetApplicationExecution(
			c.Request().Context(),
			applID,
			serviceID,
			actionNo,
		)
		if err != nil {

			_ = h.repo.CreateLog(
				c.Request().Context(),
				model.Log{
					Level:   LogLevelError,
					Source:  "GetApplicationAction",
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

		if application == nil {
			return c.JSON(
				http.StatusNotFound,
				map[string]string{
					"error": "application not found",
				},
			)
		}

		return c.JSON(
			http.StatusOK,
			application,
		)
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

	parser := utils.NewParser()

	result, err := parser.Parse(rawJSON)
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
	initiatedCount := 0
	executionCount := 0
	validServices := make(map[int64]bool)

	for _, event := range result.WorkflowEvents {

		serviceID := event.ServiceID

		exists, seen := validServices[serviceID]

		if !seen {
			serviceGroupID := serviceID / 1000
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

			validServices[serviceID] = exists

		}

		if !exists {
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

	for _, app := range result.InitiatedApps {

		if !validServices[app.ServiceID] {
			continue
		}

		err := h.repo.CreateApplicationInitiated(
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
		initiatedCount++
	}

	for _, app := range result.ExecutionApps {

		if !validServices[app.ServiceID] {
			continue
		}

		err := h.repo.CreateApplicationExecution(
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
		executionCount++
	}

	return c.JSON(
		http.StatusCreated,
		map[string]any{
			"message":   "workflow uploaded",
			"events":    workflowsCount,
			"initiated": initiatedCount,
			"execution": executionCount,
			"skipped":   skipped,
		},
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

package http

import (
	"net/http"
	"prakarsa-app/delivery/middleware"
	"prakarsa-app/domain"
	"prakarsa-app/transport/request"
	"prakarsa-app/utils"

	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/labstack/echo/v4"
)

type InstitutionHandler struct {
	InstitutionUC domain.InstitutionUsecase
}

// NewInstitutionHandler will initialize the institution resources endpoint
func NewInstitutionHandler(e *echo.Echo, middleware *middleware.Middleware, institutionUC domain.InstitutionUsecase) {
	handler := &InstitutionHandler{
		InstitutionUC: institutionUC,
	}

	apiV1 := e.Group("/api/v1")
	apiV1.POST("/institutions", handler.Create)
	apiV1.PATCH("/institutions/:id", handler.Update)
	apiV1.DELETE("/institutions/:id", handler.Delete)
}

func (h *InstitutionHandler) Create(c echo.Context) error {
	ctx := c.Request().Context()
	var req request.CreateInstitutionReq

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewUnprocessableEntityError(err.Error()))
	}

	if err := req.Validate(); err != nil {
		errVal := err.(validation.Errors)
		return c.JSON(http.StatusBadRequest, utils.NewInvalidInputError(errVal))
	}

	if res, err := h.InstitutionUC.Create(ctx, &req); err != nil {
		return c.JSON(utils.ParseHttpError(err))
	} else {
		return c.JSON(http.StatusCreated, map[string]interface{}{
			"message": "Institution successfully created",
			"data":    res,
		})
	}
}

func (h *InstitutionHandler) Update(c echo.Context) error {
	ctx := c.Request().Context()
	var req request.UpdateInstitutionReq

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewUnprocessableEntityError(err.Error()))
	}

	if err := req.Validate(); err != nil {
		errVal := err.(validation.Errors)
		return c.JSON(http.StatusBadRequest, utils.NewInvalidInputError(errVal))
	}

	if err := h.InstitutionUC.Update(ctx, &req); err != nil {
		return c.JSON(utils.ParseHttpError(err))
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Institution successfully updated",
	})
}

func (h *InstitutionHandler) Delete(c echo.Context) error {
	ctx := c.Request().Context()
	var req request.DeleteInstitutionReq

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusUnprocessableEntity, utils.NewUnprocessableEntityError(err.Error()))
	}

	if err := req.Validate(); err != nil {
		errVal := err.(validation.Errors)
		return c.JSON(http.StatusBadRequest, utils.NewInvalidInputError(errVal))
	}

	if rowsAffected, err := h.InstitutionUC.Delete(ctx, &req); err != nil {
		return c.JSON(utils.ParseHttpError(err))
	} else {
		return c.JSON(http.StatusOK, map[string]interface{}{
			"message": "Institution successfully deleted",
			"data": map[string]int64{
				"rows_affected": rowsAffected,
			},
		})
	}
}

package http

import (
	"context"
	"net/http"
	"strconv"

	"github.com/fidellr/jastip/backend/rover/models"

	"github.com/fidellr/jastip/backend/rover"
	"github.com/labstack/echo"
)

type contentHandler struct {
	service rover.ContentUsecase
}

type contentRequirements func(d *contentHandler)

func ContentService(service rover.ContentUsecase) contentRequirements {
	return func(d *contentHandler) {
		d.service = service
	}
}

func NewContentHandler(e *echo.Echo, reqs ...contentRequirements) {
	handler := new(contentHandler)
	for _, req := range reqs {
		req(handler)
	}

	e.POST("/content/create", handler.CreateScreenContent)
	e.GET("/content/:screen_name", handler.GetScreenContent)
	e.GET("/contents", handler.FetchContent)
	e.PUT("/content/:content_id", handler.UpdateByContentID)
}

func (h *contentHandler) CreateScreenContent(c echo.Context) (err error) {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	content := new(models.Screen)
	if err = c.Bind(content); err != nil {
		return rover.ConstraintErrorf("%s", err.Error())
	}

	err = h.service.CreateScreenContent(ctx, content)
	if err != nil {
		return rover.ConstraintErrorf("%s", err.Error())
	}

	return c.JSON(http.StatusCreated, content)
}

func (h *contentHandler) FetchContent(c echo.Context) (err error) {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var num int
	if c.QueryParam("num") != "" {
		num, err = strconv.Atoi(c.QueryParam("num"))
		if err != nil {
			return rover.ConstraintErrorf("%s", err.Error())
		}
	}

	filter := rover.Filter{
		Cursor:   c.QueryParam("cursor"),
		Num:      num,
		RoleName: c.QueryParam("rolve"),
	}

	contents, nextCursors, err := h.service.FetchContent(ctx, &filter)
	if err != nil {
		return rover.ConstraintErrorf("%s", err.Error())
	}

	c.Response().Header().Set("X-Cursor", nextCursors)
	return c.JSON(http.StatusOK, contents)
}

func (h *contentHandler) GetScreenContent(c echo.Context) (err error) {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	content := new(models.Screen)
	content, err = h.service.GetContentByScreen(ctx, c.Param("screen_name"))
	if err != nil {
		return rover.ConstraintErrorf("%s", err.Error())
	}

	return c.JSON(http.StatusOK, content)
}

func (h *contentHandler) UpdateByContentID(c echo.Context) (err error) {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	contentID := c.Param("content_id")
	content := new(models.Screen)
	if err = c.Bind(content); err != nil {
		return rover.ConstraintErrorf("%s", err.Error())
	}

	err = h.service.UpdateByContentID(ctx, contentID, content)
	if err != nil {
		return rover.ConstraintErrorf("%s", err.Error())
	}

	return c.JSON(http.StatusOK, true)
}

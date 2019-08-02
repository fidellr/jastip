package http

import (
	"context"
	"net/http"
	"strconv"

	"github.com/fidellr/jastip_way/backend/uranus/models"

	"github.com/fidellr/jastip_way/backend/uranus"
	"github.com/labstack/echo"
)

type userHandler struct {
	service uranus.UserAccountUsecase
}

func (h *userHandler) CreateAccount(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	u := new(models.UserAccount)
	if err := c.Bind(u); err != nil {
		return uranus.ConstraintErrorf("%s", err.Error())
	}

	err := h.service.CreateUserAccount(ctx, u)
	if err != nil {
		return uranus.ConstraintErrorf("%s", err.Error())
	}

	return c.JSON(http.StatusCreated, u)
}

func (h *userHandler) Fetch(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	var num int
	if c.QueryParam("num") != "" {
		var err error
		num, err = strconv.Atoi(c.QueryParam("num"))
		if err != nil {
			return uranus.ConstraintErrorf("%s", err.Error())
		}
	}

	filter := uranus.Filter{
		Cursor: c.QueryParam("cursor"),
		Num:    num,
	}

	users, nextCursors, err := h.service.Fetch(ctx, &filter)
	if err != nil {
		return uranus.ConstraintErrorf("%s", err.Error())
	}

	c.Response().Header().Set("X-Cursor", nextCursors)
	return c.JSON(http.StatusOK, users)
}

func (h *userHandler) GetUserByUUID(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	uuid := c.Param("id")
	user, err := h.service.GetUserByUUID(ctx, uuid)
	if err != nil {
		return uranus.ConstraintErrorf("%s", err.Error())
	}

	return c.JSON(http.StatusOK, user)
}

func (h *userHandler) SuspendAccount(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	uuid := c.Param("id")
	isSuspended, err := h.service.SuspendAccount(ctx, uuid)
	if !isSuspended || err != nil {
		return uranus.ConstraintErrorf("failed to suspend account %t: %s", isSuspended, err.Error())
	}

	return c.JSON(http.StatusOK, isSuspended)
}

func (h *userHandler) RemoveAccount(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	uuid := c.Param("id")
	isDeleted, err := h.service.RemoveAccount(ctx, uuid)
	if !isDeleted || err != nil {
		return uranus.ConstraintErrorf("failed to remove account %t: %s", isDeleted, err.Error())
	}

	return c.JSON(http.StatusOK, isDeleted)
}

type userRequirements func(d *userHandler)

func UserService(service uranus.UserAccountUsecase) userRequirements {
	return func(d *userHandler) {
		d.service = service
	}
}

func NewUserHandler(e *echo.Echo, reqs ...userRequirements) {
	handler := new(userHandler)

	for _, req := range reqs {
		req(handler)
	}

	e.POST("/user/create", handler.CreateAccount)
	e.GET("/user", handler.Fetch)
	e.GET("/user/:id", handler.GetUserByUUID)
	e.POST("/user/:id", handler.SuspendAccount)
	e.GET("/user/:id", handler.RemoveAccount)
}

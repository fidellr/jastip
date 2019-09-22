package http

import (
	"context"
	"net/http"
	"strconv"

	"github.com/labstack/echo"

	"github.com/fidellr/jastip/backend/uranus"
	"github.com/fidellr/jastip/backend/uranus/models"
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
		Cursor:   c.QueryParam("cursor"),
		Num:      num,
		RoleName: c.QueryParam("role"),
	}

	users, nextCursors, err := h.service.Fetch(ctx, &filter)
	if err != nil {
		return uranus.ConstraintErrorf("%s", err.Error())
	}

	c.Response().Header().Set("X-Cursor", nextCursors)
	return c.JSON(http.StatusOK, users)
}

func (h *userHandler) GetUserByID(c echo.Context) error {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	id := c.Param("id")
	user, err := h.service.GetUserByID(ctx, id)
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

func (h *userHandler) UpdateUserByID(c echo.Context) (err error) {
	ctx := c.Request().Context()
	if ctx == nil {
		ctx = context.Background()
	}

	id := c.Param("id")
	u := new(models.UserAccount)
	if err = c.Bind(u); err != nil {
		return uranus.ConstraintErrorf("%s", err.Error())
	}

	err = h.service.UpdateUserByID(ctx, id, u)
	if err != nil {
		return uranus.ConstraintErrorf("%s", err.Error())
	}

	return c.JSON(http.StatusOK, true)
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
	e.GET("/user/:id", handler.GetUserByID)
	e.POST("/user/suspend/:id", handler.SuspendAccount)
	e.DELETE("/user/:id", handler.RemoveAccount)
	e.PUT("/user/:id", handler.UpdateUserByID)
}

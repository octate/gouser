package handler

import (
	"context"
	"errors"
	"fmt"
	"gouser/er"
	"gouser/pkg/user"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	_pg "github.com/go-pg/pg/v10"
	"github.com/sirupsen/logrus"
)

type UserHandler struct {
	log *logrus.Logger

	userService *user.Service
}

func newUserHandler(
	log *logrus.Logger,
	userService *user.Service,
) *UserHandler {
	return &UserHandler{
		log:         log,
		userService: userService,
	}
}

type (
	// UserRes is the response struct of user API
	UserRes struct {
		// User is the total available user.
		Timestamp time.Time `json:"timestamp"`
		Success   bool      `json:"user"`
		Error     string    `json:"error"`
	}
)
type (
	CreateUserRequest struct {
		FirstName      string      `json:"first_name,omitempty"`
		LastName       string      `json:"last_name,omitempty"`
		Mobile         string      `json:"mobile" binding:"required"`
		ProfilePicture string      `json:"profile_picture,omitempty"`
		DOB            *time.Time  `form:"dob" time_format:"2006-01-02" binding:"required"`
		Metadata       interface{} `json:"metadata,omitempty"`
	}
	Response struct {
		Success bool             `json:"success"`
		Message string           `json:"message,omitempty"`
		Data    interface{}      `json:"data,omitempty"`
		Meta    *user.Pagination `json:"meta,omitempty"`
	}
)

func (h *UserHandler) CreateUser(c *gin.Context) {
	var (
		err  error
		now  = time.Now()
		dCtx = context.Background()
		req  = CreateUserRequest{}
		res  = &Response{}
	)
	defer func() {
		if err != nil {
			c.Error(err)
			return
		}
	}()
	if err = c.ShouldBind(&req); err != nil {
		err = er.New(err, er.InvalidRequestBody).SetStatus(http.StatusUnprocessableEntity)
		return
	}
	user := &user.User{
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Mobile:         req.Mobile,
		ProfilePicture: req.ProfilePicture,
		DOB:            req.DOB,
		CreatedAt:      &now,
		UpdatedAt:      &now,
		Metadata:       req.Metadata,
	}
	_, ePrr := h.userService.FetchByMobileNumber(dCtx, req.Mobile)
	switch ePrr {
	case _pg.ErrNoRows:

		ePrr := h.userService.CreateUser(dCtx, user)
		if ePrr != nil {
			h.log.WithFields(logrus.Fields{
				"request":    req,
				"error":      err,
				"req.Mobile": req.Mobile,
			}).Info("error inserting user")
			return
		}
		res.Data = user
		res.Success = true
	case nil:
		h.log.Info("rider already exist with mobile no :", req.Mobile)
		err = er.New(err, er.UserAlreadyExists).SetStatus(http.StatusUnprocessableEntity)
		return
	default:
		h.log.Info("error while fetching data from database", err.Error())
		err = er.New(err, er.UncaughtException).SetStatus(http.StatusUnprocessableEntity)
		res.Message = err.Error()
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *UserHandler) FetchUserByID(c *gin.Context) {
	var (
		err  error
		dCtx = context.Background()
		res  = &Response{}
	)
	defer func() {
		if err != nil {
			c.Error(err)
			return
		}
	}()
	userIDstr, ok := c.Params.Get("user_id")
	if !ok {
		err = errors.New("user_id empty in param")
		h.log.Info(err.Error())
		err = er.New(err, er.UncaughtException).SetStatus(http.StatusUnprocessableEntity)
		res.Message = err.Error()
		return
	}
	userID, err := strconv.Atoi(fmt.Sprint(userIDstr))
	if err != nil {
		h.log.Info("error while converting string to int: " + err.Error())
		return
	}
	user, ePrr := h.userService.FetchUserByID(dCtx, userID)
	switch ePrr {
	case _pg.ErrNoRows, nil:
		res.Data = user
		res.Success = true
	default:
		h.log.Info("error while fetching data from database", err.Error())
		err = er.New(err, er.UncaughtException).SetStatus(http.StatusUnprocessableEntity)
		res.Message = err.Error()
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *UserHandler) FetchAllUsers(c *gin.Context) {
	var (
		err  error
		dCtx = context.Background()
		req  = &user.UserRequest{}
		res  = &Response{}
	)
	defer func() {
		if err != nil {
			c.Error(err)
			return
		}
	}()
	if err = c.ShouldBind(&req); err != nil {
		err = er.New(err, er.InvalidRequestBody).SetStatus(http.StatusUnprocessableEntity)
		res.Message = err.Error()
		return
	}
	users, pagination, ePrr := h.userService.FetchAllUsers(dCtx, req)
	switch ePrr {
	case _pg.ErrNoRows, nil:
		res.Data = users
		res.Meta = &pagination
		res.Success = true
	default:
		h.log.Info("error while fetching data from database", ePrr.Error())
		err = er.New(ePrr, er.UncaughtException).SetStatus(http.StatusUnprocessableEntity)
		res.Message = ePrr.Error()
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *UserHandler) UpdateUser(c *gin.Context) {
	var (
		err  error
		now  = time.Now()
		dCtx = context.Background()
		req  = CreateUserRequest{}
		res  = &Response{}
	)
	defer func() {
		if err != nil {
			c.Error(err)
			return
		}
	}()
	if err = c.ShouldBind(&req); err != nil {
		err = er.New(err, er.InvalidRequestBody).SetStatus(http.StatusUnprocessableEntity)
		return
	}
	userIDstr, ok := c.Params.Get("user_id")
	if !ok {
		err = errors.New("user_id empty in param")
		h.log.Info(err.Error())
		err = er.New(err, er.UncaughtException).SetStatus(http.StatusUnprocessableEntity)
		res.Message = err.Error()
		return
	}
	userID, err := strconv.Atoi(fmt.Sprint(userIDstr))
	if err != nil {
		h.log.Info("error while converting string to int: " + err.Error())
		return
	}

	user := &user.User{
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		Mobile:         req.Mobile,
		ProfilePicture: req.ProfilePicture,
		DOB:            req.DOB,
		UpdatedAt:      &now,
		Metadata:       req.Metadata,
	}
	savedUser, err := h.userService.FetchUserByID(dCtx, userID)
	switch err {
	case _pg.ErrNoRows:
		h.log.Info("error while fetching data from database", err.Error())
		err = er.New(err, er.UncaughtException).SetStatus(http.StatusUnprocessableEntity)
		res.Message = err.Error()
		return
	case nil:
		user.ID = savedUser.ID
		err = h.userService.UpdateUser(dCtx, user)
		if err != nil {
			return
		}
	default:
		h.log.Info("error while fetching data from database", err.Error())
		err = er.New(err, er.UncaughtException).SetStatus(http.StatusUnprocessableEntity)
		res.Message = err.Error()
		return
	}
	res.Data = user
	res.Success = true
	c.JSON(http.StatusOK, res)
}

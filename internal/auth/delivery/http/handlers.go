package http

import (
	"companies-service/config"
	"companies-service/internal/auth"
	"companies-service/internal/models"
	"companies-service/pkg/httphelper"
	"companies-service/pkg/logger"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
)

// Auth handlers
type authHandlers struct {
	cfg         *config.Config
	authService auth.Service
	logger      logger.Logger
}

// NewAuthHandlers Auth handlers constructor
func NewAuthHandlers(
	cfg *config.Config, authService auth.Service, log logger.Logger,
) auth.Handlers {
	return &authHandlers{cfg: cfg, authService: authService, logger: log}
}

// Register godoc
// @Summary Register new user
// @Description register new user, returns user and token
// @Tags Auth
// @Accept json
// @Produce json
// @Success 201 {object} models.User
// @Router /auth/register [post]
func (h *authHandlers) Register(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContext(c, "auth.Register")
	defer span.Finish()

	user := &models.User{}
	if err := httphelper.ReadRequest(c, user); err != nil {
		httphelper.ErrResponseWithLog(c, h.logger, err)
		return
	}

	createdUser, err := h.authService.Register(ctx, user)
	if err != nil {
		httphelper.ErrResponseWithLog(c, h.logger, err)
		return
	}

	c.JSON(http.StatusCreated, createdUser)
}

// Login godoc
// @Summary Login new user
// @Description login user, returns user
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} models.User
// @Router /auth/login [post]
func (h *authHandlers) Login(c *gin.Context) {
	type Login struct {
		Email    string `json:"email" db:"email" validate:"omitempty,lte=60,email"`
		Password string `json:"password,omitempty" db:"password" validate:"required,gte=6"`
	}
	span, ctx := opentracing.StartSpanFromContext(c, "auth.Login")
	defer span.Finish()

	login := &Login{}
	if err := httphelper.ReadRequest(c, login); err != nil {
		httphelper.ErrResponseWithLog(c, h.logger, err)
		return
	}

	userWithToken, err := h.authService.Login(ctx, &models.User{
		Email:    login.Email,
		Password: login.Password,
	})
	if err != nil {
		httphelper.ErrResponseWithLog(c, h.logger, err)
		return
	}

	c.JSON(http.StatusOK, userWithToken)
}

// Update godoc
// @Summary Update user
// @Description update existing user
// @Tags Auth
// @Accept json
// @Param id path int true "user_id"
// @Produce json
// @Success 200 {object} models.User
// @Router /auth/{id} [put]
func (h *authHandlers) Update(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContext(c, "authHandlers.Update")
	defer span.Finish()

	uuid, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		httphelper.ErrResponseWithLog(c, h.logger, err)
		return
	}

	user := &models.User{}
	user.UserID = uuid

	if err = httphelper.ReadRequest(c, user); err != nil {
		httphelper.ErrResponseWithLog(c, h.logger, err)
		return
	}

	updatedUser, err := h.authService.Update(ctx, user)
	if err != nil {
		httphelper.ErrResponseWithLog(c, h.logger, err)
		return
	}

	c.JSON(http.StatusOK, updatedUser)
}

// GetUserByID godoc
// @Summary get user by id
// @Description get string by ID
// @Tags Auth
// @Accept  json
// @Produce  json
// @Param id path int true "user_id"
// @Success 200 {object} models.User
// @Failure 500 {object} httphelper.RestError
// @Router /auth/{id} [get]
func (h *authHandlers) GetUserByID(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContext(c, "authHandlers.GetUserByID")
	defer span.Finish()

	uID, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		httphelper.ErrResponseWithLog(c, h.logger, err)
		return
	}

	user, err := h.authService.GetByID(ctx, uID)
	if err != nil {
		httphelper.ErrResponseWithLog(c, h.logger, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

// Delete
// @Summary Delete user account
// @Description some description
// @Tags Auth
// @Accept json
// @Param id path int true "user_id"
// @Produce json
// @Success 200 {string} string	"ok"
// @Failure 500 {object} httphelper.RestError
// @Router /auth/{id} [delete]
func (h *authHandlers) Delete(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContext(c, "authHandlers.Delete")
	defer span.Finish()

	uuid, err := uuid.Parse(c.Param("user_id"))
	if err != nil {
		httphelper.ErrResponseWithLog(c, h.logger, err)
		return
	}

	if err = h.authService.Delete(ctx, uuid); err != nil {
		httphelper.ErrResponseWithLog(c, h.logger, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		uuid.String(): "Deleted",
	})
}

// GetMe godoc
// @Summary Get user by id
// @Description Get current user by id
// @Tags Auth
// @Accept json
// @Produce json
// @Success 200 {object} models.User
// @Failure 500 {object} httphelper.RestError
// @Router /auth/me [get]
func (h *authHandlers) GetMe(c *gin.Context) {
	span, _ := opentracing.StartSpanFromContext(c, "authHandlers.GetMe")
	defer span.Finish()

	user, ok := c.MustGet("user").(*models.User)
	if !ok {
		httphelper.ErrResponseWithLog(c, h.logger, httphelper.ErrUnauthorized)
		return
	}

	c.JSON(http.StatusOK, user)
}

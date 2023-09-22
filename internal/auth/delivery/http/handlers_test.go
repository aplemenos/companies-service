package http

import (
	"companies-service/config"
	"companies-service/internal/auth/mock"
	"companies-service/internal/models"
	"companies-service/pkg/converter"
	"companies-service/pkg/logger"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// performRequest is a helper function to perform a test request
func performRequest(r http.Handler, method, path, body string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Content-Type", gin.MIMEJSON)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func TestAuthHandlers_Register(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		Logger: config.Logger{
			Development: true,
		},
	}

	apiLogger := logger.NewApiLogger(cfg)
	mockAuthService := mock.NewMockService(ctrl)
	handlers := NewAuthHandlers(cfg, mockAuthService, apiLogger)

	gender := "male"
	user := &models.User{
		FirstName: "Nick",
		LastName:  "Papapdopoulos",
		Email:     "nick.pap@gmail.com",
		Password:  "123456",
		Gender:    &gender,
	}

	userUUID := uuid.New()
	userWithToken := &models.UserWithToken{
		User: &models.User{
			UserID: userUUID,
		},
	}

	mockAuthService.EXPECT().Register(gomock.Any(), gomock.Eq(user)).Return(userWithToken, nil)

	// Define the test route
	router := gin.Default()
	router.POST("/api/v1/auth/register", handlers.Register)

	buf, err := converter.AnyToBytesBuffer(user)
	require.NoError(t, err)
	require.NotNil(t, buf)
	require.Nil(t, err)
	w := performRequest(router, "POST", "/api/v1/auth/register", buf.String())

	// Assert the response status code (HTTP 201 Created)
	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestAuthHandlers_Login(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		Logger: config.Logger{
			Development: true,
		},
	}

	apiLogger := logger.NewApiLogger(cfg)
	mockAuthService := mock.NewMockService(ctrl)
	handlers := NewAuthHandlers(cfg, mockAuthService, apiLogger)

	// Define a mock login request
	type Login struct {
		Email    string `json:"email" db:"email" validate:"omitempty,lte=60,email"`
		Password string `json:"password,omitempty" db:"password" validate:"required,gte=6"`
	}

	loginRequest := &Login{
		Email:    "nick.pap@example.com",
		Password: "123456",
	}

	// Define a mock user and userWithToken for the response
	mockUser := &models.User{
		Email:    loginRequest.Email,
		Password: loginRequest.Password,
	}

	mockUserWithToken := &models.UserWithToken{
		User: &models.User{
			Email:    loginRequest.Email,
			Password: loginRequest.Password,
		},
	}

	// Expect the AuthService's Login method to be called with the mockUser
	mockAuthService.EXPECT().Login(gomock.Any(), gomock.Eq(mockUser)).Return(mockUserWithToken, nil)

	// Define the test route
	router := gin.Default()
	router.POST("/api/v1/auth/login", handlers.Login)

	buf, err := converter.AnyToBytesBuffer(loginRequest)
	require.NoError(t, err)
	require.NotNil(t, buf)
	require.Nil(t, err)

	// Perform the request
	w := performRequest(router, "POST", "/api/v1/auth/login", buf.String())

	// Assert the response status code (HTTP 200 OK in this case)
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse the JSON response body
	var response models.UserWithToken
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
}

func TestAuthHandlers_Update(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		Logger: config.Logger{
			Development: true,
		},
	}

	apiLogger := logger.NewApiLogger(cfg)
	mockAuthService := mock.NewMockService(ctrl)
	handlers := NewAuthHandlers(cfg, mockAuthService, apiLogger)

	// Define a mock user update request
	userID := uuid.New()
	userUpdateRequest := &models.User{
		FirstName: "Nikolas",
		LastName:  "Papadopoulakis",
	}

	// Expect the AuthService's Update method to be called with the user update request
	mockAuthService.EXPECT().Update(gomock.Any(),
		gomock.Eq(userUpdateRequest)).Return(userUpdateRequest, nil)

	// Define the test route
	router := gin.Default()
	router.PUT("/api/v1/auth/update/:user_id", handlers.Update)

	buf, err := converter.AnyToBytesBuffer(userUpdateRequest)
	require.NoError(t, err)
	require.NotNil(t, buf)
	require.Nil(t, err)

	// Perform the request with the user ID as a path parameter
	w := performRequest(router, "PUT", fmt.Sprintf("/api/v1/auth/update/%s", userID.String()),
		buf.String())

	// Assert the response status code (HTTP 200 OK in this case)
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse the JSON response body
	var response models.User
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)

	// Implement assertions for the response body as needed
	assert.Equal(t, userUpdateRequest.FirstName, response.FirstName)
	assert.Equal(t, userUpdateRequest.LastName, response.LastName)
}

func TestAuthHandlers_GetUserByID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		Logger: config.Logger{
			Development: true,
		},
	}

	apiLogger := logger.NewApiLogger(cfg)
	mockAuthService := mock.NewMockService(ctrl)
	handlers := NewAuthHandlers(cfg, mockAuthService, apiLogger)

	// Define a mock user ID
	userID := uuid.New()

	// Define a mock user for the response
	mockUser := &models.User{
		UserID: userID,
	}

	// Expect the AuthService's GetByID method to be called with the mock user ID
	mockAuthService.EXPECT().GetByID(gomock.Any(), gomock.Eq(userID)).Return(mockUser, nil)

	// Define the test route
	router := gin.Default()
	router.GET("/api/v1/auth/user/:user_id", handlers.GetUserByID)

	// Perform the request with the user ID as a path parameter
	w := performRequest(router, "GET", fmt.Sprintf("/api/v1/auth/user/%s", userID.String()), "")

	// Assert the response status code (HTTP 200 OK in this case)
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse the JSON response body
	var response models.User
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)
}

func TestAuthHandlers_Delete(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		Logger: config.Logger{
			Development: true,
		},
	}

	apiLogger := logger.NewApiLogger(cfg)
	mockAuthService := mock.NewMockService(ctrl)
	handlers := NewAuthHandlers(cfg, mockAuthService, apiLogger)

	// Define a mock user ID
	userID := uuid.New()

	// Expect the AuthService's Delete method to be called with the mock user ID
	mockAuthService.EXPECT().Delete(gomock.Any(), gomock.Eq(userID)).Return(nil)

	// Define the test route
	router := gin.Default()
	router.DELETE("/api/v1/auth/delete/:user_id", handlers.Delete)

	// Perform the request with the user ID as a path parameter
	w := performRequest(router, "DELETE",
		fmt.Sprintf("/api/v1/auth/delete/%s", userID.String()), "")

	// Assert the response status code (HTTP 200 OK in this case)
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse the JSON response body
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)

	// Perform assertions on the response, if needed
	assert.Equal(t, "Deleted", response[userID.String()])
}

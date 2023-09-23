package http

import (
	"companies-service/config"
	"companies-service/internal/companies/mock"
	"companies-service/internal/companies/service"
	"companies-service/internal/models"
	"companies-service/pkg/converter"
	"companies-service/pkg/logger"
	"encoding/json"
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

func TestCompanyHandlers_Create(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		KafkaTopics: config.KafkaTopics{
			CompanyCreated: config.TopicConfig{
				TopicName: "company_created",
			},
		},
	}

	apiLogger := logger.NewApiLogger(nil)
	mockCompanyService := mock.NewMockService(ctrl)
	companiesService := service.NewCompaniesService(nil, mockCompanyService, apiLogger)

	mockKafka := mock.NewMockKafka(ctrl)
	handlers := NewCompaniesHandlers(cfg, companiesService, mockKafka, apiLogger)

	// Create a test request with a JSON body
	company := &models.Company{
		CompanyName:        "Test Company",
		CompanyDescription: "It is a new test company",
		AmountOfEmployees:  10,
		Registered:         true,
		CompanyType:        "NonProfit",
	}

	mockCompanyService.EXPECT().Create(gomock.Any(), gomock.Any()).Return(company, nil)

	mockKafka.EXPECT().PublishMessage(gomock.Any(), gomock.Any()).Return(nil)

	// Define the test route
	router := gin.Default()
	router.POST("/api/v1/companies", handlers.Create)

	buf, err := converter.AnyToBytesBuffer(company)
	require.NoError(t, err)
	require.NotNil(t, buf)
	require.Nil(t, err)
	w := performRequest(router, "POST", "/api/v1/companies", buf.String())

	// Assert the response status code (HTTP 201 Created)
	assert.Equal(t, http.StatusCreated, w.Code)

	// Parse the JSON response body
	var response models.Company
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)

	// Implement assertions for the response body as needed
	assert.Equal(t, company.CompanyName, response.CompanyName)
	assert.Equal(t, company.CompanyDescription, response.CompanyDescription)
	assert.Equal(t, company.AmountOfEmployees, response.AmountOfEmployees)
	assert.Equal(t, company.Registered, response.Registered)
	assert.Equal(t, company.CompanyType, response.CompanyType)
}

func TestCompanyHandlers_Update(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		KafkaTopics: config.KafkaTopics{
			CompanyUpdated: config.TopicConfig{
				TopicName: "company_updated",
			},
		},
	}

	apiLogger := logger.NewApiLogger(nil)
	mockCompanyService := mock.NewMockService(ctrl)
	companiesService := service.NewCompaniesService(nil, mockCompanyService, apiLogger)

	mockKafka := mock.NewMockKafka(ctrl)
	handlers := NewCompaniesHandlers(cfg, companiesService, mockKafka, apiLogger)

	// Create a test request with a JSON body
	companyID := uuid.New()
	company := &models.Company{
		CompanyName:        "Test Company",
		CompanyDescription: "It is a test company to update",
		AmountOfEmployees:  10,
		Registered:         true,
		CompanyType:        "NonProfit",
	}

	mockCompanyService.EXPECT().Update(gomock.Any(), gomock.Any()).Return(company, nil)

	mockKafka.EXPECT().PublishMessage(gomock.Any(), gomock.Any()).Return(nil)

	// Define the test route
	router := gin.Default()
	router.PATCH("/api/v1/companies/:company_id", handlers.Update)

	buf, err := converter.AnyToBytesBuffer(company)
	require.NoError(t, err)
	require.NotNil(t, buf)
	require.Nil(t, err)
	w := performRequest(router, "PATCH", "/api/v1/companies/"+companyID.String(), buf.String())

	// Assert the response status code (HTTP 200 OK)
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse the JSON response body
	var response models.Company
	err = json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)

	// Implement assertions for the response body as needed
	assert.Equal(t, company.CompanyName, response.CompanyName)
	assert.Equal(t, company.CompanyDescription, response.CompanyDescription)
	assert.Equal(t, company.AmountOfEmployees, response.AmountOfEmployees)
	assert.Equal(t, company.Registered, response.Registered)
	assert.Equal(t, company.CompanyType, response.CompanyType)
}

func TestCompanyHandlers_Delete(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		KafkaTopics: config.KafkaTopics{
			CompanyDeleted: config.TopicConfig{
				TopicName: "company_deleted",
			},
		},
	}

	apiLogger := logger.NewApiLogger(nil)
	mockCompanyService := mock.NewMockService(ctrl)
	companiesService := service.NewCompaniesService(nil, mockCompanyService, apiLogger)

	mockKafka := mock.NewMockKafka(ctrl)
	handlers := NewCompaniesHandlers(cfg, companiesService, mockKafka, apiLogger)

	// Define the company ID for deletion
	companyID := uuid.New()

	mockCompanyService.EXPECT().Delete(gomock.Any(), companyID).Return(nil)

	mockKafka.EXPECT().PublishMessage(gomock.Any(), gomock.Any()).Return(nil)

	// Define the test route
	router := gin.Default()
	router.DELETE("/api/v1/companies/:company_id", handlers.Delete)

	w := performRequest(router, "DELETE", "/api/v1/companies/"+companyID.String(), "")

	// Assert the response status code (HTTP 200 OK)
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse the JSON response body
	var response map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)

	// Perform assertions on the response, if needed
	assert.Equal(t, "Deleted", response[companyID.String()])
}

func TestCompanyHandlers_GetByID(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	apiLogger := logger.NewApiLogger(nil)
	mockCompanyService := mock.NewMockService(ctrl)
	companiesService := service.NewCompaniesService(nil, mockCompanyService, apiLogger)

	handlers := NewCompaniesHandlers(nil, companiesService, nil, apiLogger)

	// Define a test company ID
	companyID := uuid.New()

	// Create a test request with a JSON body
	mockResponse := &models.Company{
		CompanyID:          companyID,
		CompanyName:        "Test Company",
		CompanyDescription: "It is a test company",
	}

	mockCompanyService.EXPECT().GetByID(gomock.Any(), companyID).Return(mockResponse, nil)

	// Define the test route
	router := gin.Default()
	router.GET("/api/v1/companies/:company_id", handlers.GetByID)

	w := performRequest(router, "GET", "/api/v1/companies/"+companyID.String(), "")

	// Assert the response status code (HTTP 200 OK)
	assert.Equal(t, http.StatusOK, w.Code)

	// Parse the JSON response body
	var response models.Company
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.Nil(t, err)

	// Implement assertions for the response body as needed
	assert.Equal(t, mockResponse.CompanyName, response.CompanyName)
	assert.Equal(t, mockResponse.CompanyDescription, response.CompanyDescription)
}

package http

import (
	"companies-service/config"
	"companies-service/internal/companies"
	"companies-service/internal/models"
	"companies-service/pkg/httphelper"
	"companies-service/pkg/kafka"
	"companies-service/pkg/logger"
	"companies-service/pkg/tracing"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/opentracing/opentracing-go"
	kafkaMessages "github.com/segmentio/kafka-go"
)

// Companies handlers
type companiesHandlers struct {
	cfg            *config.Config
	companyService companies.Service
	kafkaProducer  kafka.Producer
	logger         logger.Logger
}

// NewCompaniesHandlers Companies Handlers Constructor
func NewCompaniesHandlers(
	cfg *config.Config,
	companyService companies.Service,
	kafkaProducer kafka.Producer,
	logger logger.Logger,
) companies.Handlers {
	return &companiesHandlers{cfg: cfg, companyService: companyService,
		kafkaProducer: kafkaProducer, logger: logger}
}

// Create
// @Summary Create a new company
// @Description create a new company
// @Tags Companies
// @Accept  json
// @Produce  json
// @Success 201 {object} models.Company
// @Failure 400 {object} httphelper.RestErr
// @Failure 500 {object} httphelper.RestErr
// @Router /companies [post]
func (h *companiesHandlers) Create(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContext(c, "companiesHandlers.Create")
	defer span.Finish()

	company := &models.Company{}

	if err := httphelper.SanitizeRequest(c, company); err != nil {
		httphelper.ErrResponseWithLog(c, h.logger, err)
		return
	}

	createdCompany, err := h.companyService.Create(ctx, company)
	if err != nil {
		httphelper.ErrResponseWithLog(c, h.logger, err)
		return
	}

	// Prepare the kafka event
	msgBytes, err := json.Marshal(createdCompany)
	if err != nil {
		httphelper.ErrResponseWithLog(c, h.logger, err)
		return
	}

	message := kafkaMessages.Message{
		Topic:   h.cfg.KafkaTopics.CompanyCreated.TopicName,
		Value:   msgBytes,
		Time:    time.Now().UTC(),
		Headers: tracing.GetKafkaTracingHeadersFromSpanCtx(span.Context()),
	}

	// Publish a company created event to the kafka broker
	err = h.kafkaProducer.PublishMessage(ctx, message)
	if err != nil {
		httphelper.ErrResponseWithLog(c, h.logger, err)
		return
	}

	c.JSON(http.StatusCreated, createdCompany)
}

// Update
// @Summary Update a company
// @Description update a company
// @Tags Companies
// @Accept  json
// @Produce  json
// @Param id path int true "company_id"
// @Success 200 {object} models.Company
// @Failure 400 {object} httphelper.RestErr
// @Failure 404 {object} httphelper.RestErr
// @Failure 500 {object} httphelper.RestErr
// @Router /companies/{id} [patch]
func (h *companiesHandlers) Update(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContext(c, "companiesHandlers.Update")
	defer span.Finish()

	companyID, err := uuid.Parse(c.Param("company_id"))
	if err != nil {
		httphelper.ErrResponseWithLog(c, h.logger, err)
		return
	}

	company := &models.Company{}

	if err := httphelper.SanitizeRequest(c, company); err != nil {
		httphelper.ErrResponseWithLog(c, h.logger, err)
		return
	}

	company.CompanyID = companyID

	updatedCompany, err := h.companyService.Update(ctx, company)
	if err != nil {
		httphelper.ErrResponseWithLog(c, h.logger, err)
		return
	}

	// Prepare the kafka event
	msgBytes, err := json.Marshal(updatedCompany)
	if err != nil {
		httphelper.ErrResponseWithLog(c, h.logger, err)
		return
	}

	message := kafkaMessages.Message{
		Topic:   h.cfg.KafkaTopics.CompanyUpdated.TopicName,
		Value:   msgBytes,
		Time:    time.Now().UTC(),
		Headers: tracing.GetKafkaTracingHeadersFromSpanCtx(span.Context()),
	}

	// Publish a company updated event to the kafka broker
	err = h.kafkaProducer.PublishMessage(ctx, message)
	if err != nil {
		httphelper.ErrResponseWithLog(c, h.logger, err)
		return
	}

	c.JSON(http.StatusOK, updatedCompany)
}

// Delete
// @Summary Delete a company
// @Description delete a company
// @Tags Companies
// @Accept  json
// @Produce  json
// @Param id path int true "company_id"
// @Success 200 {string} string	"ok"
// @Failure 404 {object} httphelper.RestErr
// @Failure 500 {object} httphelper.RestErr
// @Router /companies/{id} [delete]
func (h *companiesHandlers) Delete(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContext(c, "companiesHandlers.Delete")
	defer span.Finish()

	companyID, err := uuid.Parse(c.Param("company_id"))
	if err != nil {
		httphelper.ErrResponseWithLog(c, h.logger, err)
		return
	}

	if err = h.companyService.Delete(ctx, companyID); err != nil {
		httphelper.ErrResponseWithLog(c, h.logger, err)
		return
	}

	// Prepare the kafka event
	msg := `{CompanyID: ` + companyID.String() + `}`
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		httphelper.ErrResponseWithLog(c, h.logger, err)
		return
	}

	message := kafkaMessages.Message{
		Topic:   h.cfg.KafkaTopics.CompanyDeleted.TopicName,
		Value:   msgBytes,
		Time:    time.Now().UTC(),
		Headers: tracing.GetKafkaTracingHeadersFromSpanCtx(span.Context()),
	}

	// Publish a company deleted event to the kafka broker
	err = h.kafkaProducer.PublishMessage(ctx, message)
	if err != nil {
		httphelper.ErrResponseWithLog(c, h.logger, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		companyID.String(): "Deleted",
	})
}

// GetByID
// @Summary Get a company
// @Description Get a company by id
// @Tags Companies
// @Accept  json
// @Produce  json
// @Param id path int true "company_id"
// @Success 200 {string} string	"ok"
// @Failure 404 {object} httphelper.RestErr
// @Failure 500 {object} httphelper.RestErr
// @Router /companies/{id} [get]
func (h *companiesHandlers) GetByID(c *gin.Context) {
	span, ctx := opentracing.StartSpanFromContext(c, "companiesHandlers.GetByID")
	defer span.Finish()

	companyID, err := uuid.Parse(c.Param("company_id"))
	if err != nil {
		httphelper.ErrResponseWithLog(c, h.logger, err)
		return
	}

	company, err := h.companyService.GetByID(ctx, companyID)
	if err != nil {
		httphelper.ErrResponseWithLog(c, h.logger, err)
		return
	}

	c.JSON(http.StatusOK, company)
}

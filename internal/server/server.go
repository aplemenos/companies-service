package server

import (
	"companies-service/config"
	kafkaClient "companies-service/pkg/kafka"
	"companies-service/pkg/logger"
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
)

const (
	maxHeaderBytes = 1 << 20
	ctxTimeout     = 5
)

// Server struct
type Server struct {
	gin         *gin.Engine
	cfg         *config.Config
	db          *sqlx.DB
	redisClient *redis.Client
	kafkaConn   *kafka.Conn
	logger      logger.Logger
}

// NewServer constructor
func NewServer(
	cfg *config.Config, db *sqlx.DB, redisClient *redis.Client, logger logger.Logger,
) *Server {
	return &Server{gin: gin.New(), cfg: cfg, db: db, redisClient: redisClient, logger: logger}
}

func (s *Server) Run() error {
	ctx, shutdown := context.WithTimeout(context.Background(), ctxTimeout*time.Second)
	defer shutdown()

	server := &http.Server{
		Addr:           s.cfg.Server.Port,
		ReadTimeout:    time.Second * s.cfg.Server.ReadTimeout,
		WriteTimeout:   time.Second * s.cfg.Server.WriteTimeout,
		MaxHeaderBytes: maxHeaderBytes,
		Handler:        s.gin,
	}

	go func() {
		s.logger.Infof("Starting Server on PORT: %s", s.cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatalf("Error ListenAndServe: %s", err)
		}
	}()

	if err := s.connectKafkaBrokers(ctx); err != nil {
		return errors.Wrap(err, "s.connectKafkaBrokers")
	}
	defer s.kafkaConn.Close()

	if s.cfg.Kafka.InitTopics {
		s.initKafkaTopics(ctx)
	}

	kafkaProducer := kafkaClient.NewProducer(s.logger, s.cfg.Kafka.Brokers)
	defer kafkaProducer.Close() // nolint: errcheck
	s.logger.Info("Kafka connected")

	if err := s.MapHandlers(kafkaProducer); err != nil {
		return err
	}

	s.runHealthCheck()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	s.logger.Info("shutdown server  ...")

	// Shut downs gracefully the server
	if err := server.Shutdown(ctx); err != nil {
		s.logger.Error(err)
		return err
	}

	s.logger.Info("shutting down gracefully")
	return nil
}

func (s *Server) connectKafkaBrokers(ctx context.Context) error {
	kafkaConn, err := kafkaClient.NewKafkaConn(ctx, s.cfg.Kafka)
	if err != nil {
		return errors.Wrap(err, "kafka.NewKafkaCon")
	}

	s.kafkaConn = kafkaConn

	brokers, err := kafkaConn.Brokers()
	if err != nil {
		return errors.Wrap(err, "kafkaConn.Brokers")
	}

	s.logger.Infof("kafka connected to brokers: %+v", brokers)

	return nil
}

func (s *Server) initKafkaTopics(ctx context.Context) {
	controller, err := s.kafkaConn.Controller()
	if err != nil {
		s.logger.Warnf("kafkaConn.Controller", err)
		return
	}

	controllerURI := net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port))
	s.logger.Infof("kafka controller uri: %s", controllerURI)

	conn, err := kafka.DialContext(ctx, "tcp", controllerURI)
	if err != nil {
		s.logger.Warnf("initKafkaTopics.DialContext", err)
		return
	}
	defer conn.Close() // nolint: errcheck

	s.logger.Infof("established new kafka controller connection: %s", controllerURI)

	companyCreatedTopic := kafka.TopicConfig{
		Topic:             s.cfg.KafkaTopics.CompanyCreated.TopicName,
		NumPartitions:     s.cfg.KafkaTopics.CompanyCreated.Partitions,
		ReplicationFactor: s.cfg.KafkaTopics.CompanyCreated.ReplicationFactor,
	}

	companyUpdatedTopic := kafka.TopicConfig{
		Topic:             s.cfg.KafkaTopics.CompanyUpdated.TopicName,
		NumPartitions:     s.cfg.KafkaTopics.CompanyUpdated.Partitions,
		ReplicationFactor: s.cfg.KafkaTopics.CompanyUpdated.ReplicationFactor,
	}

	companyDeletedTopic := kafka.TopicConfig{
		Topic:             s.cfg.KafkaTopics.CompanyDeleted.TopicName,
		NumPartitions:     s.cfg.KafkaTopics.CompanyDeleted.Partitions,
		ReplicationFactor: s.cfg.KafkaTopics.CompanyDeleted.ReplicationFactor,
	}

	if err := conn.CreateTopics(
		companyCreatedTopic,
		companyUpdatedTopic,
		companyDeletedTopic,
	); err != nil {
		s.logger.Warn("kafkaConn.CreateTopics", err)
		return
	}

	s.logger.Infof("kafka topics created or already exists: %+v",
		[]kafka.TopicConfig{companyCreatedTopic, companyUpdatedTopic, companyDeletedTopic})
}

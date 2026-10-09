package main

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/OleksUMD/ecommerce_api/internal/config"
	"github.com/OleksUMD/ecommerce_api/internal/models"
	"github.com/OleksUMD/ecommerce_api/internal/notifications"
	"github.com/OleksUMD/ecommerce_api/internal/providers"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-aws/sqs"
	"github.com/ThreeDotsLabs/watermill/message"
)

func main() {
	log.Println("Starting notification service...")

	ctx := context.Background()

	cfg, err := config.Load()

	if err != nil {
		log.Fatalf("Failed to load config: %w", err)
	}

	emailConfig := &notifications.SMTPConfig{
		Host:     cfg.SMTP.Host,
		Port:     cfg.SMTP.Port,
		Username: cfg.SMTP.Username,
		Password: cfg.SMTP.Password,
		From:     cfg.SMTP.From,
	}

	emailNotifier := notifications.NewEmailNotifier(emailConfig)
	awsConfig, err := providers.CreateAWSConfig(ctx, cfg.AWS.S3Endpoint, cfg.AWS.Region)
	if err != nil {
		log.Fatalf("Failed to create AWS config: %w", err)
	}

	logger := watermill.NewStdLogger(false, false)
	subscriber, err := sqs.NewSubscriber(sqs.SubscriberConfig{
		AWSConfig: awsConfig,
	}, logger)
	if err != nil {
		log.Fatalf("Failed to initialize SQS subscriber: %w", err)
	}
	messages, err := subscriber.Subscribe(ctx, cfg.AWS.EventQueueName)
	if err != nil {
		if errClose := subscriber.Close(); errClose != nil {
			log.Fatalf("Failed to subscribe to the queue and close the subscriber: %w\n%w", err, errClose)
		}
		log.Fatalf("Failed to subscribe to the queue: %w", err)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	log.Println("Notification service started and waiting for messages...")

	for {
		select {
		case msg := <-messages:
			if err := processMessage(msg, emailNotifier); err != nil {
				log.Printf("Error processing message: %w", err)
				msg.Nack()
			} else {
				msg.Ack()
			}
		case <-sigChan:
			log.Println("Shutting down the notification service...")
			if err := subscriber.Close(); err != nil {
				log.Fatalf("Failed to close the subscriber: %w\n%w", err)
			}
			return
		}
	}
}

func processMessage(msg *message.Message, notifier *notifications.EmailNotifier) error {
	eventType := msg.Metadata.Get("event_type")
	switch eventType {
	case notifications.UserLoggedIn:
		return handleUserLoggedIn(msg, notifier)
	default:
		log.Printf("Unknown event type: %s", eventType)
		return nil
	}
}

func handleUserLoggedIn(msg *message.Message, notifier *notifications.EmailNotifier) error {
	var user models.User
	if err := json.Unmarshal(msg.Payload, &user); err != nil {
		return err
	}
	userName := user.FirstName + " " + user.LastName
	if userName == " " {
		userName = "user"
	}
	log.Printf("Sending login notification to %s", userName)
	return notifier.SendLoginNotification(user.Email, userName)
}

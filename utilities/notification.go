package utilities

import (
	"context"
	"fmt"
	"log"
	"schoolPayment/constants"
	request "schoolPayment/dtos/request"
	"time"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	"google.golang.org/api/option"
)

type FirebaseClient struct {
	client *messaging.Client
}

// Initialize Firebase client
func NewFirebaseClient(credentialsFile string) (*FirebaseClient, error) {
	opt := option.WithCredentialsFile(credentialsFile)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, err
	}

	client, err := app.Messaging(context.Background())
	if err != nil {
		return nil, err
	}

	return &FirebaseClient{client: client}, nil
}

// Send notification to a specific token
func (fc *FirebaseClient) SendNotificationPayment(token, title, body, nis, studentName, transactionID, notificationType, redirectUrl, userId string) error {
	if fc.client == nil {
		return fmt.Errorf(constants.MessageErrorFirebaseClientNotInitialized)
	}
	timestamp := time.Now().Format("20060102150405")
	id := timestamp + transactionID

	message := &messaging.Message{
		Token: token,
		Data: map[string]string{
			"id":               id,
			"title":            title,
			"body":             body,
			"nis":              nis,
			"image":            "",
			"type":             "",
			"announcementID":   "",
			"studentName":      studentName,
			"transactionID":    transactionID,
			"notificationType": notificationType,
			"redirectUrl":      redirectUrl,
			"userID":           userId,
		},
	}

	response, err := fc.client.Send(context.Background(), message)
	if err != nil {
		// Log error with details
		if messaging.IsInvalidArgument(err) {
			log.Printf(constants.MessageErrorInvalidArgument, err)
		} else if messaging.IsSenderIDMismatch(err) {
			log.Printf(constants.MessageErrorSenderIdMismatch, err)
		} else {
			log.Printf(constants.MessageErrorSendingMessage, err)
		}
	}

	log.Printf(constants.MessageSuccessSendMessage, response)
	return nil
}

func (fc *FirebaseClient) SendNotificationDummy(token string, request *request.DummyNotifRequest, userId string) error {
	var id string

	if fc.client == nil {
		return fmt.Errorf(constants.MessageErrorFirebaseClientNotInitialized)
	}

	timestamp := time.Now().Format("20060102150405")
	if request.Type == "transaction" {
		id = timestamp + request.AnnouncementId
	} else {
		id = timestamp + request.TransactionID
	}

	message := &messaging.Message{
		Token: token,
		Data: map[string]string{
			"id":               id,
			"title":            request.Title,
			"body":             request.Body,
			"nis":              request.Nis,
			"image":            request.Image,
			"type":             request.Type,
			"announcementID":   request.AnnouncementId,
			"studentName":      request.StudentName,
			"transactionID":    request.TransactionID,
			"notificationType": request.NotificationType,
			"redirectUrl":      request.RedirectUrl,
			"userID":           userId,
		},
	}

	response, err := fc.client.Send(context.Background(), message)
	if err != nil {
		// Log error with details
		if messaging.IsInvalidArgument(err) {
			log.Printf(constants.MessageErrorInvalidArgument, err)
		} else if messaging.IsSenderIDMismatch(err) {
			log.Printf(constants.MessageErrorSenderIdMismatch, err)
		} else {
			log.Printf(constants.MessageErrorSendingMessage, err)
		}
	}

	log.Printf(constants.MessageSuccessSendMessage, response)
	return nil
}

func (fc *FirebaseClient) SendingAnouncementNotification(token, title, body, image, typeAnnouncement, announcementID string) error {
	if fc.client == nil {
		return fmt.Errorf(constants.MessageErrorFirebaseClientNotInitialized)
	}

	timestamp := time.Now().Format("20060102150405")
	id := timestamp + announcementID
	message := &messaging.Message{
		Token: token,
		Data: map[string]string{
			"id":               id,
			"title":            title,
			"body":             body,
			"image":            image,
			"type":             typeAnnouncement,
			"announcementID":   announcementID,
			"nis":              "",
			"studentName":      "",
			"transactionID":    "",
			"notificationType": "information",
			"redirectUrl":      "",
			"userID":           "",
		},
	}

	response, err := fc.client.Send(context.Background(), message)
	if err != nil {
		// Log error with details
		if messaging.IsInvalidArgument(err) {
			log.Printf(constants.MessageErrorInvalidArgument, err)
		} else if messaging.IsSenderIDMismatch(err) {
			log.Printf(constants.MessageErrorSenderIdMismatch, err)
		} else {
			log.Printf(constants.MessageErrorSendingMessage, err)
		}
	}

	log.Printf(constants.MessageSuccessSendMessage, response)
	return nil
}

package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"mms-backend/config"
	"mms-backend/models"
	"mms-backend/utils"
)

// PushService handles push notifications
type PushService struct {
	config *config.Config
}

// NewPushService creates a new push service
func NewPushService(cfg *config.Config) *PushService {
	return &PushService{
		config: cfg,
	}
}

// FCMPayload represents Firebase Cloud Messaging payload
type FCMPayload struct {
	To           string                 `json:"to"`
	Notification FCMNotification        `json:"notification"`
	Data         map[string]interface{} `json:"data,omitempty"`
	Priority     string                 `json:"priority"`
}

// FCMNotification represents FCM notification data
type FCMNotification struct {
	Title string `json:"title"`
	Body  string `json:"body"`
	Sound string `json:"sound"`
}

// SendMessageNotification sends a push notification for a new message
func (s *PushService) SendMessageNotification(receiver *models.User, senderName, messagePreview string) error {
	if receiver.DeviceToken == "" {
		return fmt.Errorf("no device token for user")
	}

	title := utils.T(receiver.Language, "message_received", senderName)
	body := messagePreview

	switch receiver.Platform {
	case "android":
		return s.sendFCM(receiver.DeviceToken, title, body, map[string]interface{}{
			"type":        "message",
			"sender_name": senderName,
		})
	case "ios":
		return s.sendAPNS(receiver.DeviceToken, title, body, map[string]interface{}{
			"type":        "message",
			"sender_name": senderName,
		})
	default:
		// Try FCM as default
		return s.sendFCM(receiver.DeviceToken, title, body, map[string]interface{}{
			"type":        "message",
			"sender_name": senderName,
		})
	}
}

// SendGroupMessageNotification sends a push notification for a new group message
func (s *PushService) SendGroupMessageNotification(receiver *models.User, groupName, senderName, messagePreview string) error {
	if receiver.DeviceToken == "" {
		return fmt.Errorf("no device token for user")
	}

	title := groupName
	body := fmt.Sprintf("%s: %s", senderName, messagePreview)

	switch receiver.Platform {
	case "android":
		return s.sendFCM(receiver.DeviceToken, title, body, map[string]interface{}{
			"type":        "group_message",
			"group_name":  groupName,
			"sender_name": senderName,
		})
	case "ios":
		return s.sendAPNS(receiver.DeviceToken, title, body, map[string]interface{}{
			"type":        "group_message",
			"group_name":  groupName,
			"sender_name": senderName,
		})
	default:
		return s.sendFCM(receiver.DeviceToken, title, body, map[string]interface{}{
			"type":        "group_message",
			"group_name":  groupName,
			"sender_name": senderName,
		})
	}
}

// sendFCM sends a notification via Firebase Cloud Messaging
func (s *PushService) sendFCM(deviceToken, title, body string, data map[string]interface{}) error {
	if s.config.Push.FCMServerKey == "" {
		log.Println("FCM server key not configured, skipping push notification")
		return nil
	}

	payload := FCMPayload{
		To: deviceToken,
		Notification: FCMNotification{
			Title: title,
			Body:  body,
			Sound: "default",
		},
		Data:     data,
		Priority: "high",
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://fcm.googleapis.com/fcm/send", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "key="+s.config.Push.FCMServerKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("FCM request failed with status: %d", resp.StatusCode)
		return fmt.Errorf("FCM request failed with status: %d", resp.StatusCode)
	}

	log.Println("FCM notification sent successfully")
	return nil
}

// sendAPNS sends a notification via Apple Push Notification Service
func (s *PushService) sendAPNS(deviceToken, title, body string, data map[string]interface{}) error {
	// Note: This is a simplified implementation
	// In production, you should use a proper APNs library like github.com/sideshow/apns2
	if s.config.Push.APNSKeyID == "" {
		log.Println("APNs key not configured, skipping push notification")
		return nil
	}

	// TODO: Implement proper APNs integration with JWT authentication
	// For now, we'll just log that we would send an APNs notification
	log.Printf("Would send APNs notification to device: %s, title: %s, body: %s", deviceToken, title, body)
	
	return nil
}


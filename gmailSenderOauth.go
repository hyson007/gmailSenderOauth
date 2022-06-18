package gmailSenderOauth

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

var GmailService *gmail.Service

var (
	ClientID     string
	ClientSecret string
	AccessToken  string
	RefreshToken string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	ClientID = os.Getenv("CLIENT_ID")
	ClientSecret = os.Getenv("CLIENT_SECRET")
	AccessToken = os.Getenv("ACCESS_TOKEN")
	RefreshToken = os.Getenv("REFRESH_TOKEN")
}

func OAuthGmailService() {
	config := oauth2.Config{
		ClientID:     ClientID,
		ClientSecret: ClientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  "http://localhost",
	}

	token := oauth2.Token{
		AccessToken:  AccessToken,
		RefreshToken: RefreshToken,
		TokenType:    "Bearer",
		Expiry:       time.Now(),
	}

	var tokenSource = config.TokenSource(context.Background(), &token)

	srv, err := gmail.NewService(context.Background(), option.WithTokenSource(tokenSource))
	if err != nil {
		log.Printf("Unable to retrieve Gmail client: %v", err)
	}

	GmailService = srv
	if GmailService != nil {
		fmt.Println("Email service is initialized")
	}
}

func SendEmailOAUTH2(to string, kind string, data interface{}) (bool, error) {

	var message gmail.Message
	var body string

	emailTo := "To: " + to + "\r\n"
	subject := "Subject: " + "Welcome to subscribe earth quake notification" + "\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"

	if kind == "registration" {
		body = "Hello,\n\n" + "You have subscribed to earth quake notification.\n\n" + "Thank you for using our service.\n\n" + "Earth Quake Notification Service\n" + "Kindly check the URL to confirm your subscription\n" +
			"http://localhost:8080/verification/" + data.(string) + "\n\n" + "Regards,\n" + "Earth Quake Notification Service"
	} else if kind == "notification" {
		body = "Hello,\n\n" + "We have detected suspected earthquake in your interested area\n" + "Details: \n" +
			data.(string) + "\n\n" + "Regards,\n" + "Earth Quake Notification Service"
	}

	msg := []byte(emailTo + subject + mime + "\n" + body)
	message.Raw = base64.URLEncoding.EncodeToString(msg)

	// Send the message
	_, err := GmailService.Users.Messages.Send("me", &message).Do()
	if err != nil {
		return false, err
	}
	return true, nil
}

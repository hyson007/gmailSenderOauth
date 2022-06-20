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
	if ClientID == "" || ClientSecret == "" || AccessToken == "" || RefreshToken == "" {
		log.Fatal("Error loading .env file")
	}

}

func NewOAuthGmailService() *gmail.Service {
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

	fmt.Println("Email service is initialized")
	return srv
}

func SendEmailOAUTH2(svc *gmail.Service, to string, subject string, body string) (bool, error) {

	var message gmail.Message
	emailTo := "To: " + to + "\r\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"

	msg := []byte(emailTo + subject + mime + "\n" + body)
	message.Raw = base64.URLEncoding.EncodeToString(msg)

	// Send the message
	_, err := svc.Users.Messages.Send("me", &message).Do()
	if err != nil {
		log.Println("Error sending message: " + err.Error())
		return false, err
	}
	return true, nil
}

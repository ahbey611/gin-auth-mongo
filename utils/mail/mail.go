package mail

import (
	// "context"
	"errors"
	"fmt"
	"gin-auth-mongo/utils/consts"
	"io"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

// var ch chan *gomail.Message
var ch chan *EmailMessage

// var ctx = context.Background()

// EmailMessage contains email message and related Redis parameters
type EmailMessage struct {
	Message    *gomail.Message
	Email      string
	ResultChan chan error
}

type SendResult struct {
	Error error
}

type VerificationMethod string

const (
	VerificationMethodEmail VerificationMethod = "EMAIL"
	VerificationMethodPhone VerificationMethod = "PHONE"
)

type VerificationFormType string

const (
	VerificationFormTypeLink VerificationFormType = "LINK"
	VerificationFormTypeCode VerificationFormType = "CODE"
)

type VerificationRequestType string

const (
	VerificationRequestTypeRegister      VerificationRequestType = "REGISTER"
	VerificationRequestTypeResetPassword VerificationRequestType = "RESET"
)

func InitMail() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	SMTPServer := os.Getenv("SMTP_SERVER")
	SMTPUser := os.Getenv("SMTP_USER")
	SMTPPassword := os.Getenv("SMTP_PASS")
	SMTPPort := os.Getenv("SMTP_PORT")
	SMTPPortInt, err := strconv.Atoi(SMTPPort)
	if err != nil {
		log.Fatal("Error converting SMTP_PORT to int")
	}

	ch = make(chan *EmailMessage)

	go func() {
		d := gomail.NewDialer(SMTPServer, SMTPPortInt, SMTPUser, SMTPPassword)

		var s gomail.SendCloser
		var err error
		open := false
		for {
			select {
			case m, ok := <-ch:
				if !ok {
					return
				}

				// if not open, open the connection
				if !open {

					// try to open the connection
					s, err = d.Dial()
					if err != nil {
						log.Println("Error dialing", err)
						m.ResultChan <- err
						continue
					} else {
						open = true
						log.Println("SMTP connection opened")
					}
				}
				if err := gomail.Send(s, m.Message); err != nil {
					log.Print(err)
					m.ResultChan <- err
					continue
				}

				log.Printf("Email sent to %s\n", m.Email)

				m.ResultChan <- nil

			// Close the connection to the SMTP server if no email was sent in
			// the last 30 seconds.
			case <-time.After(30 * time.Second):
				if open {
					if s != nil {

						err := s.Close()
						if err != nil && err != io.EOF {
							log.Println("Unexpected error closing SMTP connection:", err)
						} else {

							log.Println("SMTP connection closed (or already closed by server)")
						}
						open = false
					}
				}
			}
		}
	}()
}

func SendVerificationEmail(email string, username string, formType VerificationFormType, requestType VerificationRequestType, message string, expiry string) *SendResult {

	if email == "" || username == "" || formType == "" || requestType == "" || message == "" {
		return &SendResult{
			Error: errors.New("invalid email, username, form type, request type, or message"),
		}
	}

	resultChan := make(chan error)
	m := gomail.NewMessage()
	SMTPFromAddress := os.Getenv("SMTP_FROM_ADDRESS")
	SMTPFromName := os.Getenv("SMTP_FROM_NAME")
	m.SetHeader("From", SMTPFromAddress)
	m.SetHeader("To", email)
	m.SetAddressHeader("Cc", SMTPFromAddress, SMTPFromName)
	m.SetHeader("Subject", SMTPFromName+" Email Verification")
	var content string

	switch formType {
	case VerificationFormTypeLink:
		content = GetVerificationLinkContent(email, username, requestType, message, expiry)
	case VerificationFormTypeCode:
		content = GetVerificationCodeContent(email, username, requestType, message, expiry)
	default:
		return &SendResult{
			Error: errors.New("invalid form type"),
		}
	}

	if content == "" {
		return &SendResult{
			Error: errors.New("invalid form type"),
		}
	}

	m.SetBody("text/html", content)

	emailMessage := &EmailMessage{
		Message:    m,
		Email:      email,
		ResultChan: resultChan,
	}
	ch <- emailMessage
	err := <-resultChan
	close(resultChan)
	return &SendResult{
		Error: err,
	}
}

func GetVerificationLinkContent(email string, username string, requestType VerificationRequestType, link string, expiry string) string {
	switch requestType {
	case VerificationRequestTypeRegister:
		return fmt.Sprintf(EmailRegisterLinkTemplate, username, link, link, consts.VERIFY_EMAIL_REGISTER_LINK_EXPIRY, expiry)
	case VerificationRequestTypeResetPassword:
		return fmt.Sprintf(PasswordResetLinkTemplate, username, link, link, consts.VERIFY_EMAIL_RESET_PWD_LINK_EXPIRY, expiry)
	}
	return ""
}

func GetVerificationCodeContent(email string, username string, requestType VerificationRequestType, code string, expiry string) string {
	switch requestType {
	case VerificationRequestTypeRegister:
		return fmt.Sprintf(EmailRegisterCodeTemplate, username, code, consts.VERIFY_EMAIL_REGISTER_CODE_EXPIRY, expiry)
	case VerificationRequestTypeResetPassword:
		return fmt.Sprintf(PasswordResetCodeTemplate, username, code, consts.VERIFY_EMAIL_RESET_PWD_CODE_EXPIRY, expiry)
	}
	return ""
}

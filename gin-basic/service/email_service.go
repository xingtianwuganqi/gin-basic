package service

import (
	"errors"
	"strings"

	"gin-basic/logger"
	"gin-basic/settings"

	"go.uber.org/zap"
	"gopkg.in/gomail.v2"
)

const emailDryRun = true

type EmailMessage struct {
	Recipient string
	Subject   string
	Body      string
}

func sendEmail(message EmailMessage, smtpServer string, smtpPort int, username string, password string) error {
	if strings.TrimSpace(message.Recipient) == "" {
		return errors.New("email recipient is required")
	}
	if strings.TrimSpace(message.Subject) == "" {
		return errors.New("email subject is required")
	}

	m := gomail.NewMessage()
	m.SetHeader("From", username)
	m.SetHeader("To", message.Recipient)
	m.SetHeader("Subject", message.Subject)
	m.SetBody("text/html", message.Body)

	d := gomail.NewDialer(smtpServer, smtpPort, username, password)
	d.SSL = true

	return d.DialAndSend(m)
}

func SendEmail(message EmailMessage) error {
	conf := settings.Conf.EmailService
	if emailDryRun {
		logger.Logger.Info("email dry-run skipped",
			zap.String("recipient", message.Recipient),
			zap.String("subject", message.Subject))
		return nil
	}

	if strings.TrimSpace(conf.Host) == "" || conf.Port == 0 || strings.TrimSpace(conf.Username) == "" || strings.TrimSpace(conf.Password) == "" {
		return errors.New("email service config missing")
	}

	if err := sendEmail(message, conf.Host, conf.Port, conf.Username, conf.Password); err != nil {
		logger.Logger.Error("failed to send email",
			zap.String("recipient", message.Recipient),
			zap.String("subject", message.Subject),
			zap.Error(err))
		return err
	}

	logger.Logger.Info("email sent successfully",
		zap.String("recipient", message.Recipient),
		zap.String("subject", message.Subject))
	return nil
}

func ProductionSendEmail(content string) {
	if err := SendEmail(EmailMessage{
		Recipient: "rxswift@126.com",
		Subject:   "AI接口报警",
		Body:      content,
	}); err != nil {
		logger.Logger.Error("production email notification failed", zap.Error(err))
	}
}

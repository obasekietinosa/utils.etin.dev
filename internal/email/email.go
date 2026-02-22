package email

import (
	"fmt"
	"net/smtp"
	"os"
	"strings"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	From     string
	To       string
	CC       []string
}

func LoadConfig() (*Config, error) {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	username := os.Getenv("SMTP_USERNAME")
	password := os.Getenv("SMTP_PASSWORD")
	from := os.Getenv("SMTP_FROM_EMAIL")
	to := os.Getenv("SMTP_TO_EMAIL")
	cc := os.Getenv("SMTP_CC_EMAILS")

	if host == "" || port == "" || username == "" || password == "" || from == "" || to == "" {
		return nil, fmt.Errorf("missing required SMTP environment variables")
	}

	var ccList []string
	if cc != "" {
		ccList = strings.Split(cc, ",")
		for i := range ccList {
			ccList[i] = strings.TrimSpace(ccList[i])
		}
	}

	return &Config{
		Host:     host,
		Port:     port,
		Username: username,
		Password: password,
		From:     from,
		To:       to,
		CC:       ccList,
	}, nil
}

func SendEmail(cfg *Config, subject, body string) error {
	auth := smtp.PlainAuth("", cfg.Username, cfg.Password, cfg.Host)

	to := []string{cfg.To}
	to = append(to, cfg.CC...)

	// RFC 822 style headers
	headers := make(map[string]string)
	headers["From"] = cfg.From
	headers["To"] = cfg.To
	if len(cfg.CC) > 0 {
		headers["Cc"] = strings.Join(cfg.CC, ", ")
	}
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/plain; charset=\"utf-8\""

	headerStr := ""
	for k, v := range headers {
		headerStr += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	msg := []byte(headerStr + "\r\n" + body)

	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	err := smtp.SendMail(addr, auth, cfg.From, to, msg)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

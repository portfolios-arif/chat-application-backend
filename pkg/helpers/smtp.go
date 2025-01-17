package helpers

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/smtp"
	"os"
	"strings"
	"time"
)

type smtpConfig struct {
	Host   string
	Port   string
	Sender string
	Email  string
	Pass   string
}

func SendEmailOTP(ctx context.Context, otp, recipient string) (string, error) {
	config := smtpConfig{
		Host:   os.Getenv("SMTP_HOST"),
		Port:   os.Getenv("SMTP_PORT"),
		Sender: os.Getenv("SMTP_SENDER"),
		Email:  os.Getenv("SMTP_AUTH_EMAIL"),
		Pass:   os.Getenv("SMTP_AUTH_PASSWORD"),
	}

	to := []string{recipient}
	smtpAddr := config.Host + ":" + config.Port

	// Configure TLS
	tlsConfig := &tls.Config{
		ServerName: config.Host,
		MinVersion: tls.VersionTLS12,
	}

	// Create a dialer with timeout
	d := &net.Dialer{
		Timeout: 10 * time.Second,
	}

	// First connect without TLS
	conn, err := d.DialContext(ctx, "tcp", smtpAddr)
	if err != nil {
		return "", fmt.Errorf("connection failed: %w", err)
	}
	defer conn.Close()

	c, err := smtp.NewClient(conn, config.Host)
	if err != nil {
		return "", fmt.Errorf("failed to create SMTP client: %w", err)
	}
	defer c.Close()

	// Start TLS
	if err = c.StartTLS(tlsConfig); err != nil {
		return "", fmt.Errorf("start TLS failed: %w", err)
	}

	// Now authenticate after TLS is established
	auth := smtp.PlainAuth("", config.Email, config.Pass, config.Host)
	if err = c.Auth(auth); err != nil {
		return "", fmt.Errorf("authentication failed: %w", err)
	}

	if err = c.Mail(config.Email); err != nil {
		return "", fmt.Errorf("setting sender failed: %w", err)
	}

	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return "", fmt.Errorf("setting recipient failed: %w", err)
		}
	}

	w, err := c.Data()
	if err != nil {
		return "", fmt.Errorf("getting data writer failed: %w", err)
	}

	body := "From: " + config.Sender + "\r\n" +
		"To: " + strings.Join(to, ",") + "\r\n" +
		"Subject: OTP Code\r\n" +
		"Content-Type: text/html; charset=utf-8\r\n" +
		"\r\n" +
		"OTP Anda adalah: <b>" + otp + "</b>|\r\n\n" +
		"OTP ini akan hangus dalam 3 menit. Harap segera memasukkan OTP."

	_, err = w.Write([]byte(body))
	if err != nil {
		return "", fmt.Errorf("writing body failed: %w", err)
	}

	err = w.Close()
	if err != nil {
		return "", fmt.Errorf("closing data writer failed: %w", err)
	}

	err = c.Quit()
	if err != nil {
		return "", fmt.Errorf("QUIT command failed: %w", err)
	}

	return "success", nil
}

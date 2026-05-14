// internal/service/smtp_sender.go
package service

import (
	"fmt"
	"net/smtp"
)

// SMTPSender - реальная отправка писем через SMTP
type SMTPSender struct {
	host     string
	port     int
	username string
	password string
	from     string
	fromName string
}

// NewSMTPSender создаёт новый экземпляр SMTP отправителя
func NewSMTPSender(host string, port int, username, password, from, fromName string) *SMTPSender {
	return &SMTPSender{
		host:     host,
		port:     port,
		username: username,
		password: password,
		from:     from,
		fromName: fromName,
	}
}

// Send отправляет email через SMTP
func (s *SMTPSender) Send(to, subject, body string) error {
	// Формируем заголовки письма
	headers := make(map[string]string)
	headers["From"] = fmt.Sprintf("%s <%s>", s.fromName, s.from)
	headers["To"] = to
	headers["Subject"] = subject
	headers["MIME-Version"] = "1.0"
	headers["Content-Type"] = "text/html; charset=\"UTF-8\""

	// Собираем полное сообщение
	message := ""
	for key, value := range headers {
		message += fmt.Sprintf("%s: %s\r\n", key, value)
	}
	message += "\r\n" + body

	// Адрес SMTP сервера
	addr := fmt.Sprintf("%s:%d", s.host, s.port)

	// Аутентификация (если требуется)
	var auth smtp.Auth
	if s.username != "" && s.password != "" {
		auth = smtp.PlainAuth("", s.username, s.password, s.host)
	}

	// Отправка письма
	err := smtp.SendMail(
		addr,
		auth,
		s.from,
		[]string{to},
		[]byte(message),
	)

	if err != nil {
		return fmt.Errorf("failed to send email to %s: %w", to, err)
	}

	return nil
}

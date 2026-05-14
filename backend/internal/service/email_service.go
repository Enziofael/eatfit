// internal/service/email_service.go
package service

import (
	"log"
	"strings"

	"github.com/Enziofael/eatfit/backend/internal/config"
)

// EmailService - сервис для отправки email уведомлений
type EmailService struct {
	cfg    *config.Config
	mailer EmailSender
}

// EmailSender - интерфейс для отправки писем
type EmailSender interface {
	Send(to, subject, body string) error
}

// Шаблоны писем (константы с плейсхолдерами {{CODE}} и {{LOGIN}})
const verificationTemplate = `<!DOCTYPE html>
<html>
<head>
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            margin: 0;
            padding: 0;
            background-color: #f5f5f5;
        }
        .container {
            max-width: 600px;
            margin: 20px auto;
            background: white;
            border-radius: 12px;
            overflow: hidden;
            box-shadow: 0 2px 8px rgba(0,0,0,0.1);
        }
        .header {
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
            padding: 30px;
            text-align: center;
            color: white;
        }
        .header h1 {
            margin: 0;
            font-size: 28px;
        }
        .content {
            padding: 30px;
        }
        .code-box {
            background: #f8f9fa;
            border: 2px dashed #667eea;
            border-radius: 8px;
            padding: 20px;
            text-align: center;
            margin: 20px 0;
        }
        .code {
            font-size: 32px;
            font-weight: bold;
            letter-spacing: 8px;
            color: #667eea;
            font-family: 'Courier New', monospace;
        }
        .footer {
            background: #f8f9fa;
            padding: 20px;
            text-align: center;
            color: #6c757d;
            font-size: 12px;
        }
        .warning {
            color: #dc3545;
            font-size: 12px;
            margin-top: 20px;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🏋️ Welcome to Eatfit!</h1>
        </div>
        <div class="content">
            <h2>Verify Your Email Address</h2>
            <p>Thank you for registering! To complete your registration, please use the verification code below:</p>
            
            <div class="code-box">
                <p style="margin: 0 0 10px 0; color: #6c757d;">Your verification code:</p>
                <div class="code">{{CODE}}</div>
            </div>
            
            <p>This code will expire in <strong>15 minutes</strong>.</p>
            <p>If you didn't create an Eatfit account, please ignore this email.</p>
            
            <p class="warning">⚠️ Never share this code with anyone. Our team will never ask for your verification code.</p>
        </div>
        <div class="footer">
            <p>© 2025 Eatfit. All rights reserved.</p>
            <p>This is an automated message, please do not reply to this email.</p>
        </div>
    </div>
</body>
</html>`

const welcomeTemplate = `<!DOCTYPE html>
<html>
<head>
    <style>
        body {
            font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
            margin: 0;
            padding: 0;
            background-color: #f5f5f5;
        }
        .container {
            max-width: 600px;
            margin: 20px auto;
            background: white;
            border-radius: 12px;
            overflow: hidden;
            box-shadow: 0 2px 8px rgba(0,0,0,0.1);
        }
        .header {
            background: linear-gradient(135deg, #11998e 0%, #38ef7d 100%);
            padding: 30px;
            text-align: center;
            color: white;
        }
        .header h1 {
            margin: 0;
            font-size: 28px;
        }
        .content {
            padding: 30px;
        }
        .welcome-text {
            font-size: 18px;
            line-height: 1.6;
            color: #333;
        }
        .features {
            display: grid;
            grid-template-columns: 1fr 1fr;
            gap: 15px;
            margin: 30px 0;
        }
        .feature {
            padding: 15px;
            background: #f8f9fa;
            border-radius: 8px;
            text-align: center;
        }
        .feature-icon {
            font-size: 24px;
            margin-bottom: 8px;
        }
        .feature-text {
            font-size: 14px;
            color: #555;
        }
        .footer {
            background: #f8f9fa;
            padding: 20px;
            text-align: center;
            color: #6c757d;
            font-size: 12px;
        }
    </style>
</head>
<body>
    <div class="container">
        <div class="header">
            <h1>🎉 Welcome to Eatfit!</h1>
        </div>
        <div class="content">
            <p class="welcome-text">
                Hey <strong>{{LOGIN}}</strong>! 👋<br>
                Your email has been verified successfully. You're all set to start your fitness journey!
            </p>
            
            <div class="features">
                <div class="feature">
                    <div class="feature-icon">📊</div>
                    <div class="feature-text">Track Calories</div>
                </div>
                <div class="feature">
                    <div class="feature-icon">🏋️</div>
                    <div class="feature-text">Workout Plans</div>
                </div>
                <div class="feature">
                    <div class="feature-icon">🍽️</div>
                    <div class="feature-text">Meal Recipes</div>
                </div>
                <div class="feature">
                    <div class="feature-icon">👥</div>
                    <div class="feature-text">Community</div>
                </div>
            </div>
            
            <p style="text-align: center; color: #666;">
                Start tracking your nutrition and achieving your fitness goals today!
            </p>
        </div>
        <div class="footer">
            <p>© 2025 Eatfit. All rights reserved.</p>
            <p>This is an automated message, please do not reply to this email.</p>
        </div>
    </div>
</body>
</html>`

// NewEmailService создаёт новый экземпляр EmailService
func NewEmailService(cfg *config.Config) *EmailService {
	var mailer EmailSender

	if cfg.AppEnv == "development" {
		mailer = NewSMTPSender(
			cfg.SMTPHost,
			cfg.SMTPPort,
			cfg.SMTPUser,
			cfg.SMTPPass,
			cfg.EmailFrom,
			cfg.EmailFromName,
		)
		log.Printf("Email service initialized with SMTP (MailHog): %s:%d", cfg.SMTPHost, cfg.SMTPPort)
	} else {
		mailer = NewSMTPSender(
			cfg.SMTPHost,
			cfg.SMTPPort,
			cfg.SMTPUser,
			cfg.SMTPPass,
			cfg.EmailFrom,
			cfg.EmailFromName,
		)
		log.Printf("Email service initialized with SMTP: %s:%d", cfg.SMTPHost, cfg.SMTPPort)
	}

	return &EmailService{
		cfg:    cfg,
		mailer: mailer,
	}
}

// SendVerificationEmail отправляет письмо с кодом подтверждения
func (s *EmailService) SendVerificationEmail(to string, code string) error {
	subject := "Eatfit - Email Verification"
	body := strings.Replace(verificationTemplate, "{{CODE}}", code, 1)
	return s.mailer.Send(to, subject, body)
}

// SendWelcomeEmail отправляет приветственное письмо после подтверждения email
func (s *EmailService) SendWelcomeEmail(to string, login string) error {
	subject := "Welcome to Eatfit! 🎉"
	body := strings.Replace(welcomeTemplate, "{{LOGIN}}", login, 1)
	return s.mailer.Send(to, subject, body)
}

// MockEmailSender - моковая реализация для разработки
type MockEmailSender struct{}

// Send логирует письмо вместо реальной отправки
func (m *MockEmailSender) Send(to, subject, body string) error {
	log.Printf(strings.Repeat("=", 50))
	log.Printf("MOCK EMAIL")
	log.Printf("To: %s", to)
	log.Printf("Subject: %s", subject)
	log.Printf("Body: %s", body)
	log.Printf(strings.Repeat("=", 50))
	return nil
}

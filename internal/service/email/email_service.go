// Package service
package service

import (
	"fmt"
	"net/smtp"
	"strings"

	"future-letter/internal/config"
	"future-letter/internal/models"
)

// EmailService Struct yang menangani pengiriman email
type EmailService struct {
	cfg *config.Config
}

// NewEmailService membuat instance email service baru
func NewEmailService(cfg *config.Config) *EmailService {
	return &EmailService{
		cfg: cfg,
	}
}

// SendCapsuleEmail mengirim capsule via email ke user
// Parameter :
//   - user : data user yang akan menerima email
//   - capsule : data capsule yang akan dikirim
func (s *EmailService) SendCapsuleEmail(user *models.User, capsule *models.Capsule) error {
	// Judul email
	subject := fmt.Sprintf("Time Capsule: %s", capsule.Title)

	// Email body / isi email
	body := s.comploseEmailBody(user, capsule)

	auth := smtp.PlainAuth(
		"",
		s.cfg.Email.SMTPUsername,
		s.cfg.Email.SMTPPassword,
		s.cfg.Email.SMTPHost,
	)

	// Mmembuat message
	message := []byte(
		"From: " + s.cfg.Email.SMTPFrom + "\r\n" +
			"To: " + user.Email + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: Text/html; charset=UTF-8\r\n" +
			"\r\n" +
			body + "\r\n",
	)

	addr := fmt.Sprintf("%s:%s", s.cfg.Email.SMTPHost, s.cfg.Email.SMTPPort)
	err := smtp.SendMail(
		addr,
		auth,
		s.cfg.Email.SMTPUsername,
		[]string{user.Email},
		message,
	)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}

func (s *EmailService) comploseEmailBody(user *models.User, capsule *models.Capsule) string {
	// Format category dan mood jika ada
	category := "Not specified"
	if capsule.Category.Valid {
		category = capsule.Category.String
	}

	mood := "Not specified"
	if capsule.Mood.Valid {
		mood = capsule.Mood.String
	}

	html := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
</head>
<body style="font-family: Arial, sans-serif; line-height: 1.6; color: #333; max-width: 600px; margin: 0 auto; padding: 20px;">
    <!-- Header -->
    <div style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); padding: 30px; text-align: center; border-radius: 10px 10px 0 0;">
        <h1 style="color: white; margin: 0; font-size: 28px;">⏰ Your Time Capsule Has Arrived!</h1>
    </div>
    
    <!-- Content -->
    <div style="background: #f9f9f9; padding: 30px; border-radius: 0 0 10px 10px; border: 1px solid #e0e0e0;">
        <p style="font-size: 16px; margin-bottom: 20px;">
            Hi <strong>%s</strong>,
        </p>
        
        <p style="font-size: 16px; margin-bottom: 20px;">
            Remember this? You wrote this message to your future self on <strong>%s</strong>:
        </p>
        
        <!-- Capsule Card -->
        <div style="background: white; padding: 25px; border-radius: 8px; border-left: 4px solid #667eea; margin: 20px 0; box-shadow: 0 2px 4px rgba(0,0,0,0.1);">
            <h2 style="color: #667eea; margin-top: 0; font-size: 22px;">%s</h2>
            <div style="background: #f5f5f5; padding: 20px; border-radius: 5px; margin: 15px 0;">
                <p style="margin: 0; white-space: pre-wrap; font-size: 15px; line-height: 1.8;">%s</p>
            </div>
            
            <!-- Metadata -->
            <div style="margin-top: 20px; padding-top: 20px; border-top: 1px solid #e0e0e0;">
                <p style="margin: 5px 0; font-size: 14px; color: #666;">
                    📁 <strong>Category:</strong> %s
                </p>
                <p style="margin: 5px 0; font-size: 14px; color: #666;">
                    😊 <strong>Mood when written:</strong> %s
                </p>
                <p style="margin: 5px 0; font-size: 14px; color: #666;">
                    📅 <strong>Written on:</strong> %s
                </p>
            </div>
        </div>
        
        <!-- Reflection Prompt -->
        <div style="background: #fff3cd; border: 1px solid #ffc107; padding: 15px; border-radius: 8px; margin: 20px 0;">
            <p style="margin: 0; font-size: 14px;">
                💭 <strong>Take a moment to reflect:</strong><br>
                How have you grown since you wrote this? Have you achieved your goals?
            </p>
        </div>
        
        <p style="font-size: 14px; color: #666; margin-top: 30px;">
            Keep moving forward! 🚀
        </p>
    </div>
    
    <!-- Footer -->
    <div style="text-align: center; padding: 20px; font-size: 12px; color: #999;">
        <p>This is an automated message from Future Self Reminders</p>
        <p>You received this because you scheduled a time capsule to be sent today.</p>
    </div>
</body>
</html>
`

	// Replace placeholder dengan data sebenarnya
	html = fmt.Sprintf(html,
		escapeHTML(user.Name),
		capsule.CreatedAt.Format("January 2, 2006"),
		escapeHTML(capsule.Title),
		escapeHTML(capsule.Message),
		escapeHTML(category),
		escapeHTML(mood),
		capsule.CreatedAt.Format("January 1, 2006 at 3:04 PM"),
	)

	return html
}

func (s *EmailService) SendTestEmail(toEmail string) error {
	subject := "Test email from Future Self Reminders"

	body := `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
</head>
<body style="font-family: Arial, sans-serif; padding: 20px;">
    <h2 style="color: #667eea;">✅ Email Configuration Test</h2>
    <p>If you're reading this, your SMTP configuration is working correctly!</p>
    <p>You can now send time capsules to your future self.</p>
    <hr>
    <p style="font-size: 12px; color: #999;">
        This is a test email from Future Self Reminders API
    </p>
</body>
</html>
`

	auth := smtp.PlainAuth("", s.cfg.Email.SMTPUsername, s.cfg.Email.SMTPPassword, s.cfg.Email.SMTPHost)

	message := []byte(
		"From: " + s.cfg.Email.SMTPFrom + "\r\n" +
			"To: " + toEmail + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: Text/html; charset=UTF-8\r\n" +
			"\r\n" +
			body + "\r\n",
	)

	addr := fmt.Sprintf("%s:%s", s.cfg.Email.SMTPHost, s.cfg.Email.SMTPPort)
	err := smtp.SendMail(
		addr,
		auth,
		s.cfg.Email.SMTPUsername,
		[]string{toEmail},
		message,
	)
	if err != nil {
		return fmt.Errorf("failed to send test email: %w", err)
	}

	return nil
}

func escapeHTML(s string) string {
	// Map karakter yang perlu direplace
	replacements := map[string]string{
		"&":  "&amp;",
		"<":  "&lt;",
		">":  "&gt;",
		"\"": "&quot;",
		"'":  "&#39",
	}

	// Replace semua karakter spesial
	for old, new := range replacements {
		s = strings.ReplaceAll(s, old, new)
	}

	return s
}

func (s *EmailService) SendWelcomeEmail(user *models.User) error {
	subject := "Welcome to Future Self Reminders"

	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
</head>
<body style="font-family: Arial, sans-serif; padding: 20px; max-width: 600px; margin: 0 auto;">
    <div style="background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); padding: 30px; text-align: center; border-radius: 10px;">
        <h1 style="color: white; margin: 0;">Welcome to Future Self Reminders! 🎉</h1>
    </div>
    
    <div style="padding: 30px; background: #f9f9f9; border-radius: 0 0 10px 10px;">
        <p>Hi <strong>%s</strong>,</p>
        
        <p>Thank you for joining Future Self Reminders! We're excited to help you send messages to your future self.</p>
        
        <h3 style="color: #667eea;">What's next?</h3>
        <ul>
            <li>Create your first time capsule</li>
            <li>Set goals and reflect on them in the future</li>
            <li>Track your personal growth journey</li>
        </ul>
        
        <p>Start by creating your first time capsule today!</p>
        
        <p style="margin-top: 30px;">Happy time traveling! 🚀</p>
    </div>
    
    <div style="text-align: center; padding: 20px; font-size: 12px; color: #999;">
        <p>Future Self Reminders - Your personal time capsule service</p>
    </div>
</body>
</html>
`, escapeHTML(user.Name))

	auth := smtp.PlainAuth("", s.cfg.Email.SMTPUsername, s.cfg.Email.SMTPPassword, s.cfg.Email.SMTPHost)
	message := []byte(
		"From: " + s.cfg.Email.SMTPFrom + "\r\n" +
			"To: " + user.Email + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-Version: 1.0\r\n" +
			"Content-Type: Text/html; charset=UTF-8\r\n" +
			"\r\n" +
			body + "\r\n",
	)

	addr := fmt.Sprintf("%s:%s", s.cfg.Email.SMTPHost, s.cfg.Email.SMTPPort)
	err := smtp.SendMail(addr, auth, s.cfg.Email.SMTPFrom, []string{user.Email}, message)
	if err != nil {
		// Tidak mereturn error karena registrasi tetap berhasil walau welcome email gagal
		fmt.Printf("Warning: Failed to send welcome email: %v\n", err)
		return nil
	}

	return nil
}

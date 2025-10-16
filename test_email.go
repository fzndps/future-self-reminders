package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	"future-letter/internal/config"
	"future-letter/internal/models"
	emailService "future-letter/internal/service/email"
)

// ==========================================
// MAIN FUNCTION FOR EMAIL TESTING
// ==========================================
// File ini untuk testing email service
// PENTING: Pastikan konfigurasi SMTP di .env sudah benar!
func main() {
	fmt.Println("📧 Testing Email Service...")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// ==========================================
	// 1. LOAD CONFIGURATION
	// ==========================================
	fmt.Println("\n📋 Loading configuration...")
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}
	fmt.Println("✅ Configuration loaded")

	// Tampilkan SMTP config (hide password)
	fmt.Println("\n🔧 SMTP Configuration:")
	fmt.Printf("   Host: %s\n", cfg.Email.SMTPHost)
	fmt.Printf("   Port: %d\n", cfg.Email.SMTPPort)
	fmt.Printf("   Username: %s\n", cfg.Email.SMTPUsername)
	fmt.Printf("   From: %s\n", cfg.Email.SMTPFrom)

	// ==========================================
	// 2. INITIALIZE EMAIL SERVICE
	// ==========================================
	fmt.Println("\n📧 Initializing email service...")
	emailService := emailService.NewEmailService(cfg)
	fmt.Println("✅ Email service initialized")

	// ==========================================
	// 3. TEST 1: SEND TEST EMAIL
	// ==========================================
	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("TEST 1: Send Test Email")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Ganti dengan email kamu untuk testing
	testEmail := "mkhabib47@gmail.com" // Kirim ke diri sendiri

	fmt.Printf("📤 Sending test email to: %s\n", testEmail)
	err = emailService.SendTestEmail(testEmail)
	if err != nil {
		fmt.Printf("❌ Failed to send test email: %v\n", err)
		fmt.Println("\n💡 Troubleshooting tips:")
		fmt.Println("   1. Pastikan SMTP credentials di .env benar")
		fmt.Println("   2. Untuk Gmail, gunakan App Password bukan password biasa")
		fmt.Println("   3. Cara buat App Password: https://support.google.com/accounts/answer/185833")
		fmt.Println("   4. Pastikan 'Less secure app access' diaktifkan (jika pakai password biasa)")
		return
	}
	fmt.Println("✅ Test email sent successfully!")
	fmt.Println("   Check your inbox (and spam folder)")

	// ==========================================
	// 4. TEST 2: SEND WELCOME EMAIL
	// ==========================================
	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("TEST 2: Send Welcome Email")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Buat dummy user untuk testing
	dummyUser := &models.User{
		ID:       1,
		Name:     "khabib",
		Email:    testEmail,
		Timezone: "Asia/Jakarta",
	}

	fmt.Printf("📤 Sending welcome email to: %s\n", dummyUser.Email)
	err = emailService.SendWelcomeEmail(dummyUser)
	if err != nil {
		fmt.Printf("❌ Failed to send welcome email: %v\n", err)
		return
	}
	fmt.Println("✅ Welcome email sent successfully!")

	// ==========================================
	// 5. TEST 3: SEND CAPSULE EMAIL
	// ==========================================
	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("TEST 3: Send Time Capsule Email")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Buat dummy capsule untuk testing
	dummyCapsule := &models.Capsule{
		ID:             1,
		UserID:         1,
		Title:          "My First Time Capsule",
		Message:        "Hi future me!\n\nI hope you are doing great and have achieved all your goals.\n\nRemember to stay positive and keep learning!\n\nBest regards,\nPast You",
		DueDate:        time.Now(),
		DeliveryMethod: "email",
		Status:         "pending",
		Category: sql.NullString{
			String: "personal",
			Valid:  true,
		},
		Mood: sql.NullString{
			String: "motivated",
			Valid:  true,
		},
		CreatedAt: time.Now().AddDate(0, -6, 0), // 6 bulan yang lalu
	}

	fmt.Printf("📤 Sending capsule email to: %s\n", dummyUser.Email)
	fmt.Printf("   Title: %s\n", dummyCapsule.Title)
	err = emailService.SendCapsuleEmail(dummyUser, dummyCapsule)
	if err != nil {
		fmt.Printf("❌ Failed to send capsule email: %v\n", err)
		return
	}
	fmt.Println("✅ Capsule email sent successfully!")

	// ==========================================
	// SUMMARY
	// ==========================================
	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("✨ All email tests completed successfully!")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("\n📊 Test Summary:")
	fmt.Println("   ✅ Test email sent")
	fmt.Println("   ✅ Welcome email sent")
	fmt.Println("   ✅ Time capsule email sent")
	fmt.Println("\n💡 Check your email inbox for 3 emails!")
	fmt.Println("   (Don't forget to check spam folder)")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}

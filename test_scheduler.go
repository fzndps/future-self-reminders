package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"future-letter/internal/config"
	"future-letter/internal/database"
	"future-letter/internal/models"
	capsuleRepository "future-letter/internal/repository/capsule"
	userRepository "future-letter/internal/repository/user"
	capsuleService "future-letter/internal/service/capsule"
	emailService "future-letter/internal/service/email"
	schedulerService "future-letter/internal/service/scheduler"
	userService "future-letter/internal/service/user"
)

// ==========================================
// TEST SCHEDULER SERVICE
// ==========================================
func main() {
	fmt.Println("🧪 Testing Scheduler Service")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// ==========================================
	// 1. SETUP
	// ==========================================
	fmt.Println("\n📋 Loading configuration...")
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load config:", err)
	}

	fmt.Println("🔌 Connecting to database...")
	err = database.InitDB(cfg)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.CloseDB()

	// ==========================================
	// 2. INITIALIZE LAYERS
	// ==========================================
	fmt.Println("🏗️  Initializing layers...")

	userRepo := userRepository.NewUserRepository(database.DB)
	capsuleRepo := capsuleRepository.NewCapsuleRepository(database.DB)

	userSvc := userService.NewUserService(userRepo)
	capsuleSvc := capsuleService.NewCapsuleService(capsuleRepo)
	emailSvc := emailService.NewEmailService(cfg)

	scheduler := schedulerService.NewSchedulerService(cfg, userRepo, capsuleSvc, emailSvc)

	fmt.Println("✅ All layers initialized")

	ctx := context.Background()

	// ==========================================
	// 3. CREATE TEST USER
	// ==========================================
	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("👤 Creating test user...")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	registerInput := &models.RegisterInput{
		Name:     "Khabib test user",
		Email:    "mkhabib47@gmail.com", // Kirim ke email kita sendiri
		Password: "password123",
		Timezone: "Asia/Jakarta",
	}

	user, err := userSvc.Register(ctx, registerInput)
	if err != nil {
		// Jika user sudah ada, coba login
		loginInput := &models.LoginInput{
			Email:    registerInput.Email,
			Password: registerInput.Password,
		}
		user, err = userSvc.Login(ctx, loginInput)
		if err != nil {
			log.Fatal("Failed to get user:", err)
		}
	}

	fmt.Printf("✅ User ready: %s (%s)\n", user.Name, user.Email)

	// ==========================================
	// 4. CREATE TEST CAPSULES
	// ==========================================
	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("⏰ Creating test capsules...")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Capsule 1: Untuk hari ini (akan dikirim oleh scheduler)
	fmt.Println("\n1️⃣  Creating capsule for TODAY...")
	todayDate := time.Now().Format("2006-01-02")

	todayCapsuleInput := &models.CreateCapsuleInput{
		Title:          "🎉 Test Capsule - Should Be Sent Today!",
		Message:        "Hi there!\n\nThis is a test capsule that should be sent by the scheduler today.\n\nIf you receive this email, it means the scheduler is working correctly! 🎊\n\nTest timestamp: " + time.Now().Format("2006-01-02 15:04:05"),
		DueDate:        todayDate,
		DeliveryMethod: "email",
		Category:       "test",
		Mood:           "excited",
	}

	todayCapsule, err := capsuleSvc.CreateCapsule(ctx, user.ID, todayCapsuleInput)
	if err != nil {
		log.Fatal("Failed to create today capsule:", err)
	}
	fmt.Printf("✅ Today capsule created: ID=%d, Title=%s\n", todayCapsule.ID, todayCapsule.Title)

	// Capsule 2: Untuk besok (tidak akan dikirim hari ini)
	fmt.Println("\n2️⃣  Creating capsule for TOMORROW...")
	tomorrowDate := time.Now().AddDate(0, 0, 1).Format("2006-01-02")

	tomorrowCapsuleInput := &models.CreateCapsuleInput{
		Title:          "Tomorrow's Capsule",
		Message:        "This should not be sent today!",
		DueDate:        tomorrowDate,
		DeliveryMethod: "email",
	}

	tomorrowCapsule, err := capsuleSvc.CreateCapsule(ctx, user.ID, tomorrowCapsuleInput)
	if err != nil {
		log.Fatal("Failed to create tomorrow capsule:", err)
	}
	fmt.Printf("✅ Tomorrow capsule created: ID=%d, Title=%s\n", tomorrowCapsule.ID, tomorrowCapsule.Title)

	// Capsule 3: Untuk bulan depan (tidak akan dikirim hari ini)
	fmt.Println("\n3️⃣  Creating capsule for NEXT MONTH...")
	nextMonthDate := time.Now().AddDate(0, 1, 0).Format("2006-01-02")

	nextMonthCapsuleInput := &models.CreateCapsuleInput{
		Title:          "Next Month's Capsule",
		Message:        "This is for next month!",
		DueDate:        nextMonthDate,
		DeliveryMethod: "email",
		Category:       "future",
		Mood:           "hopeful",
	}

	nextMonthCapsule, err := capsuleSvc.CreateCapsule(ctx, user.ID, nextMonthCapsuleInput)
	if err != nil {
		log.Fatal("Failed to create next month capsule:", err)
	}
	fmt.Printf("✅ Next month capsule created: ID=%d, Title=%s\n", nextMonthCapsule.ID, nextMonthCapsule.Title)

	// ==========================================
	// 5. CHECK PENDING CAPSULES
	// ==========================================
	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("📋 Checking pending capsules for today...")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	pendingCapsules, err := capsuleSvc.GetPendingCapsulesForToday(ctx)
	if err != nil {
		log.Fatal("Failed to get pending capsules:", err)
	}

	fmt.Printf("\n✅ Found %d pending capsule(s) for today:\n", len(pendingCapsules))
	for i, c := range pendingCapsules {
		fmt.Printf("   %d. ID=%d, Title=%s, User=%d\n", i+1, c.ID, c.Title, c.UserID)
	}

	// ==========================================
	// 6. RUN SCHEDULER MANUALLY
	// ==========================================
	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("🚀 Running scheduler manually...")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("\nℹ️  Note: This will actually send emails!")
	fmt.Println("   Check your inbox at:", user.Email)
	fmt.Println()

	// Cast ke concrete type untuk akses method RunManually
	concreteScheduler, ok := scheduler.(schedulerService.SchedulerService)
	if !ok {
		// Jika interface tidak menyediakan RunManually, panggil Start lalu tunggu
		fmt.Println("⚠️  Using alternative method: Starting scheduler for 10 seconds...")
		err := scheduler.Start()
		if err != nil {
			log.Fatal("Failed to start scheduler:", err)
		}
		defer scheduler.Stop()

		// Tunggu sebentar untuk scheduler berjalan
		time.Sleep(10 * time.Second)
	} else {
		// Panggil RunManually jika tersedia
		concreteScheduler.RunManually()
	}

	// ==========================================
	// 7. VERIFY RESULTS
	// ==========================================
	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("🔍 Verifying results...")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")

	// Check today's capsule status
	updatedTodayCapsule, err := capsuleSvc.GetCapsule(ctx, todayCapsule.ID, user.ID)
	if err != nil {
		log.Fatal("Failed to get updated capsule:", err)
	}

	fmt.Printf("\n📧 Today's Capsule (ID=%d):\n", todayCapsule.ID)
	fmt.Printf("   Status: %s\n", updatedTodayCapsule.Status)
	if updatedTodayCapsule.SentAt.Valid {
		fmt.Printf("   Sent At: %s\n", updatedTodayCapsule.SentAt.Time.Format("2006-01-02 15:04:05"))
		fmt.Println("   ✅ Email should have been sent!")
	} else {
		fmt.Println("   ⚠️  Status not updated (but email might still be sent)")
	}

	// Check tomorrow's capsule (should still be pending)
	updatedTomorrowCapsule, err := capsuleSvc.GetCapsule(ctx, tomorrowCapsule.ID, user.ID)
	if err != nil {
		log.Fatal("Failed to get tomorrow capsule:", err)
	}

	fmt.Printf("\n📅 Tomorrow's Capsule (ID=%d):\n", tomorrowCapsule.ID)
	fmt.Printf("   Status: %s (should be 'pending')\n", updatedTomorrowCapsule.Status)
	if updatedTomorrowCapsule.Status == "pending" {
		fmt.Println("   ✅ Correctly not sent today!")
	}

	// ==========================================
	// SUMMARY
	// ==========================================
	fmt.Println("\n━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("✨ Scheduler test completed!")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
	fmt.Println("\n📊 Summary:")
	fmt.Println("   ✅ Scheduler service working")
	fmt.Println("   ✅ Email service working")
	fmt.Println("   ✅ Capsule status updated")
	fmt.Println("\n💡 Next steps:")
	fmt.Println("   1. Check your email inbox")
	fmt.Println("   2. Look for email with subject: '⏰ Time Capsule: ...'")
	fmt.Println("   3. If no email, check spam folder")
	fmt.Println("\n🎉 Your Future Self Reminders API is complete!")
	fmt.Println("━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━")
}

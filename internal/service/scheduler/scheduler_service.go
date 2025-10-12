// Package service scheduler
package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"future-letter/internal/config"
	repository "future-letter/internal/repository/user"
	capsule "future-letter/internal/service/capsule"
	email "future-letter/internal/service/email"

	"github.com/robfig/cron/v3"
)

// SchedulerService interface
type SchedulerService interface {
	Start() error
	Stop()
	RunManually()
}

// schedulerService struct implementation
type schedulerService struct {
	cfg            *config.Config
	cron           *cron.Cron
	userRepo       repository.UserRepository
	capsuleService capsule.CapsuleService
	emailService   *email.EmailService
}

// NewSchedulerService instance baru SchedulerService
// inject semua service yang dibutuhkan
func NewSchedulerService(
	cfg *config.Config,
	userRepo repository.UserRepository,
	capsuleService capsule.CapsuleService,
	emailService *email.EmailService,
) SchedulerService {
	// Load timezone dari config
	location, err := time.LoadLocation(cfg.Schedular.Timezone)
	if err != nil {
		log.Printf("Failed to load timezone %s, using UTC: %v", cfg.Schedular.Timezone, err)
		location = time.UTC
	}

	// Buat cron dengan timezone
	cronScheduler := cron.New(
		cron.WithSeconds(),
		cron.WithLocation(location),
	)

	return &schedulerService{
		cfg:            cfg,
		cron:           cronScheduler,
		userRepo:       userRepo,
		capsuleService: capsuleService,
		emailService:   emailService,
	}
}

// Start untuk memulai scheduler
func (s *schedulerService) Start() error {
	log.Println("Start scheduler service...")

	// Register cron job
	// addFunc untuk menambahkan job ke scheduler
	_, err := s.cron.AddFunc(s.cfg.Schedular.CronExpression, func() {
		// Fungsi akan dijalankan sesuai cron expression
		log.Println("Schedular running: checking pending capsuless...")

		// jalankan job untuk kirtim capsules
		s.processPendingCapsules()
	})
	if err != nil {
		return fmt.Errorf("failed to add cron job: %w", err)
	}

	// Start cron scheduler menjalankan scheduler di background (goroutine)
	s.cron.Start()

	log.Printf("Schedular started with expression: %s", s.cfg.Schedular.CronExpression)

	log.Printf("Timezone: %s", s.cfg.Schedular.Timezone)

	log.Println("Schedular is running in background...")

	return nil
}

func (s *schedulerService) Stop() {
	log.Println("Stopping scheduler service...")

	// Menghentikan scheduler dan menunggu semua job selesai
	ctx := s.cron.Stop()

	// Tunggu semua job selesai (dengan TO 30s)
	select {
	case <-ctx.Done():
		log.Println("Schedular stopped gracefully")
	case <-time.After(30 * time.Second):
		log.Panicln("Scheduler stop timeout after 30 seconds")
	}
}

func (s *schedulerService) processPendingCapsules() {
	// Buat context dengan timeout
	// Agar job tidak berjalan selamanya / loop
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// GetPendingForToday mendapat capsule yang pending pada hari ini
	capsules, err := s.capsuleService.GetPendingCapsulesForToday(ctx)
	if err != nil {
		log.Printf("Failed to get pending capsules: %v", err)
		return
	}

	// Jika tidak ada capsule yang jatuh tempo
	if len(capsules) == 0 {
		log.Println("No pending capsule for today")
		return
	}

	log.Printf("Found %d pendiing capsule(S) to send", len(capsules))

	// Proses setiap capsule
	successCount := 0
	failCount := 0

	for _, capsule := range capsules {
		// Dapatkan user data untuk email
		user, err := s.userRepo.GetByID(ctx, capsule.ID)
		if err != nil {
			log.Printf("Failed to get user %d for capsule %d: %v", capsule.UserID, capsule.ID, err)
			failCount++
			continue
		}

		// Kirim email
		log.Printf("Sending capsule %d to %s (%s)", capsule.ID, user.Name, user.Email)

		err = s.emailService.SendCapsuleEmail(user, &capsule)
		if err != nil {
			log.Printf("Failed to send email for capsule %d: %v", capsule.ID, err)
			failCount++
			continue
		}

		// Tandai jika sudah dikirim
		err = s.capsuleService.MarkCapsulesAsSent(ctx, capsule.ID)
		if err != nil {
			log.Printf("Email send but failed to update status for capsule %d: %v", capsule.ID, err)
			// email sudah terkirim tetapi status di database belum terupdate
		}

		log.Printf("Capsule %d sent successfully to %s", capsule.ID, user.Email)
		successCount++
	}

	// Ringkasan Log
	log.Printf("Success : %d capsules", successCount)

	log.Printf("Failed : %d capsules", failCount)

	log.Printf("Total processed : %d capsules", len(capsules))
}

func (s *schedulerService) RunManually() {
	log.Println("Running scheduler manually for testing...")
	s.processPendingCapsules()
}

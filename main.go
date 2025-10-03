package main

import (
	"fmt"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Gagal memuat env")
	}

	// Data dari login SMTP gmail
	from := "futureletter.app@gmail.com"
	password := os.Getenv("PASSWORD")

	// penerima
	to := []string{"deanurindasalsadil@gmail.com"}

	msg := []byte("Subject: Test future letter\r\n" + "\r\n" + "Pesan ini terkirim otomatis dari program golang")

	// Alamat server gmail
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"

	// autentikasi
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Kirim
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, msg)
	if err != nil {
		fmt.Println("Gagal mengirim pesan : ", err)
		return
	}
}

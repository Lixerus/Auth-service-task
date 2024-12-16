package services

import (
	"log"
	"net/smtp"
)

func SendEmail() {
	from := "from@example.io"
	password := "password"
	to := []string{
		"dest@email.io",
	}
	smtpHost := "smtp.gmail.com"
	smtpPort := "587"
	message := []byte("Warning! Your credentials on test site might be used!")
	// send(from, password, smtpHost, smtpPort, to, message) // tries to send email to smtp server
	mocksend(from, password, smtpHost, smtpPort, to, message)
}

func send(from, password, smtpHost, smtpPort string, to []string, message []byte) {
	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		log.Println(err)
	}
}

func mocksend(from, password, smtpHost, smtpPort string, to []string, message []byte) {
	log.Println("Using credentials "+from, password, ". Using smtp host "+smtpHost+smtpPort)
	log.Println("Sending email to ", to, "With message "+string(message))
}

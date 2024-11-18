package email

import (
	"log"
	"net/mail"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
	"github.com/scorredoira/email"
)

func SendEmail(emaill string, filename string) error {

	err := godotenv.Load("./.env")
	if err != nil {
		log.Println("Couldn't load environment variables: ", err)
		return err
	}

	from := os.Getenv("EMAIL")
	password := os.Getenv("EMAIL_PASS")

	to := []string{emaill}

	m := email.NewMessage("آگهی ها", "آگهی های فیلتر شده توسط شما")
	m.From = mail.Address{Name: "crawler", Address: from}
	m.To = to

	// Attach csv file
	if err := m.Attach(filename); err != nil {
		log.Println("Could't attach file: ", err)
		return err
	}

	// SMTP server configuration
	host := os.Getenv("SMTP")
	port := os.Getenv("SMTP_PORT")
	address := host + ":" + port

	auth := smtp.PlainAuth("", from, password, host)

	// send the email
	if err := email.Send(address, auth, m); err != nil {
		log.Println("Couldn't send the email: ", err)
		return err
	}
	return nil

}

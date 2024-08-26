package smtp

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"log"
	"mime"
	"net"
	"net/mail"
	"net/smtp"
)

// New - конструктор для создания SMTP с настроенными параметрами.
func New(host, port, name, email, login, password string) *SMTP {
	return &SMTP{
		host:     host,
		port:     port,
		name:     name,
		email:    email,
		login:    login,
		password: password,
	}
}

// SMTP - структура для отправки писем.
type SMTP struct {
	host     string
	port     string
	name     string
	email    string
	login    string
	password string
}

// Метод отправляет текстовое письмо указанному адресату.
func (s *SMTP) Text(name, email, subject, message string) error {
	return s.sendEmail(name, email, subject, "text/plain", []byte(message))
}

// Метод отправляет HTML письмо указанному адресату.
func (s *SMTP) HTML(name, email, subject string, message []byte) error {
	return s.sendEmail(name, email, subject, "text/html", message)
}

func (s *SMTP) sendEmail(name, email, subject, mime string, message []byte) error {
	from, to := s.addresses(name, email)
	auth := smtp.PlainAuth("", s.login, s.password, s.host)
	client, err := s.smtp()
	if err != nil {
		return err
	}

	err = client.Auth(auth)
	if err != nil {
		return err
	}

	err = client.Mail(from.Address)
	if err != nil {
		return err
	}

	err = client.Rcpt(to.Address)
	if err != nil {
		return err
	}

	wr, err := client.Data()
	defer func() {
		err := wr.Close()

		if err != nil {
			log.Println(err)
		}
	}()
	if err != nil {
		return err
	}

	_, _ = wr.Write(s.message(from, to, subject, mime, message))
	return client.Quit()
}

func (m *SMTP) addresses(name, email string) (from, to mail.Address) {
	from = mail.Address{Name: m.name, Address: m.email}
	to = mail.Address{Name: name, Address: email}
	return
}

func (m *SMTP) message(from, to mail.Address, subject, mimeType string, message []byte) []byte {
	buffer := new(bytes.Buffer)

	buffer.WriteString(fmt.Sprintf("From: %s\r\n", from.String()))
	buffer.WriteString(fmt.Sprintf("To: %s\r\n", to.String()))
	buffer.WriteString(fmt.Sprintf("Subject: %s\r\n", mime.QEncoding.Encode("utf-8", subject)))
	buffer.WriteString(fmt.Sprintf("MIME-Version: 1.0;\r\nContent-Type: %s; charset=\"UTF-8\"\r\n\r\n", mimeType))
	buffer.Write(message)

	return buffer.Bytes()
}

func (m *SMTP) smtp() (*smtp.Client, error) {
	servername := m.host + ":" + m.port
	host, _, _ := net.SplitHostPort(servername)

	tlsconfig := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         m.host,
	}

	// Вызов tcp соедиения
	conn, err := tls.Dial("tcp", servername, tlsconfig)
	if err != nil {
		return nil, err
	}

	return smtp.NewClient(conn, host)
}

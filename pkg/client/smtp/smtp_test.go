package smtp_test

import (
	"testing"

	"github.com/ReanSn0w/gokit/pkg/app"
	"github.com/ReanSn0w/gokit/pkg/client/smtp"
)

var (
	opts = smtp.Config{}
)

func init() {
	_, err := app.LoadConfiguration("Test SMTP Client", "unknown", &opts)
	if err != nil {
		panic(err)
	}
}

func Test_SendPlainMail(t *testing.T) {
	err := opts.NewSMTP().Text(
		"Дмитрий Папков",
		"papkovda@me.com",
		"Test mail",
		"Текст сообщения",
	)

	if err != nil {
		t.Log(err)
		t.Fail()
	}
}

func Test_SendHTMLMail(t *testing.T) {
	err := opts.NewSMTP().HTML(
		"Дмитрий Папков",
		"papkovda@me.com",
		"Test mail",
		[]byte("<html><body><h1>hello, world</h1></body></html>"),
	)

	if err != nil {
		t.Log(err)
		t.Fail()
	}
}

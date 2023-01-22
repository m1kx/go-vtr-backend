package notify

import (
	"context"
	"fmt"
	"os"

	"github.com/m1kx/go-vtr-backend/pkg/config"
	"github.com/mailgun/mailgun-go/v4"
)

func SendMail(msg string, mail string) error {

	url := os.Getenv("MAIL_URL")
	key := os.Getenv("MAIL_KEY")

	mg := mailgun.NewMailgun(url, key)
	m := mg.NewMessage(
		"Plan Info <no-reply@mg.mikadev.tech>",
		"Neuer Plan",
		msg,
		mail,
	)

	mg.SetAPIBase("https://api.eu.mailgun.net/v3")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_, id, err := mg.Send(ctx, m)

	if err != nil {
		fmt.Println("ERR")
		return err
	}

	fmt.Printf("Sent mail to %s%s%s (id: %s%s%s)\n", config.Yellow, mail, config.Reset, config.Yellow, id, config.Reset)
	return nil
}

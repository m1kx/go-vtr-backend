package notify

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/SherClockHolmes/webpush-go"
	"github.com/m1kx/go-vtr-backend/pkg/config"
)

func Send(msg string, subscr string, mail string) {
	priv := os.Getenv("PRIV")
	pub := os.Getenv("PUB")

	s := &webpush.Subscription{}
	err := json.Unmarshal([]byte(subscr), s)

	resp, err := webpush.SendNotification([]byte(msg), s, &webpush.Options{
		Subscriber:      mail,
		VAPIDPublicKey:  pub,
		VAPIDPrivateKey: priv,
		TTL:             30,
	})
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Sent PWA notification to %s%s%s \n", config.Yellow, mail, config.Reset)
	defer resp.Body.Close()
}

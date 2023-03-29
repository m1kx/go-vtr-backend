package notify

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/SherClockHolmes/webpush-go"
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
	defer resp.Body.Close()
}

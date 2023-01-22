package notify

import (
	"bytes"
	"net/http"

	"github.com/m1kx/go-vtr-backend/pkg/utils/structs"
)

func SendPerRequest(req_opt *structs.HttpReq) error {
	var body []byte
	if req_opt.METHOD == "POST" {
		body = []byte(req_opt.PROPS)
	}
	url := req_opt.BASE_URL
	req, err := http.NewRequest(req_opt.METHOD, url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	if req_opt.METHOD == "POST" {
		req.Header.Set("Content-Type", "application/json")
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	return nil
}

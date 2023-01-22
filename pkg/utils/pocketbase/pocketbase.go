package pocketbase

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/m1kx/go-vtr-backend/pkg/utils/structs"
)

// authenticate as admin and return token
func auth() (string, error) {
	// request admin login to pocketbase
	url := "http://127.0.0.1:8090/api/admins/auth-with-password"
	admin_mail := os.Getenv("ADMIN_MAIL")
	password := os.Getenv("ADMIN_PASS")
	var body = []byte(fmt.Sprintf(`{"identity": "%s", "password": "%s"}`, admin_mail, password))
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println("error")
		return "", err
	}
	if res.StatusCode != 200 {
		return "", errors.New(fmt.Sprintf("Pocketbase Admin auth failed with code: %d", res.StatusCode))
	}
	defer res.Body.Close()

	// read and format the response from pocketbase
	bodyres, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var main structs.AuthResponse
	err = json.Unmarshal(bodyres, &main)
	if err != nil {
		return "", err
	}

	//return the admin token
	return main.TOKEN, nil
}

// retrieve all users from pocketbase users table
func GetAllUsers() ([]structs.User, error) {
	token, err := auth()
	if err != nil {
		return nil, err
	}
	url := "http://127.0.0.1:8090/api/collections/users/records"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", token)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != 200 {
		return nil, errors.New(fmt.Sprintf("Pocketbase User fetch failed with code: %d", res.StatusCode))
	}
	defer res.Body.Close()
	bodyres, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var main structs.UserResponse
	err = json.Unmarshal(bodyres, &main)
	if err != nil {
		return nil, err
	}

	return main.ITEMS, nil
}

// update user in pocketbase
func EditField(identifier string, id string, data interface{}) error {
	token, err := auth()
	if err != nil {
		return err
	}
	url := fmt.Sprintf("http://127.0.0.1:8090/api/collections/users/records/%s", id)
	value_send := data
	switch v := data.(type) {
	case string:
		value_send = fmt.Sprintf(`"%v"`, v)
	default:
		value_send = v
	}
	var body = []byte(fmt.Sprintf(`{"%s": %v}`, identifier, value_send))
	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", token)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != 200 {
		return errors.New(fmt.Sprintf("Pocketbase Field Update failed with code: %d", res.StatusCode))
	}
	defer res.Body.Close()
	return nil
}

package pocketbase

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v5"
	"github.com/m1kx/go-vtr-backend/pkg/utils/health"
	"github.com/m1kx/go-vtr-backend/pkg/utils/structs"
	pb "github.com/pocketbase/pocketbase"
	"github.com/pocketbase/pocketbase/core"
)

var app *pb.PocketBase

// authenticate as admin and return token
func auth() (string, error) {
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
		return "", err
	}
	if res.StatusCode != 200 {
		return "", errors.New(fmt.Sprintf("Pocketbase Admin auth failed with code: %d", res.StatusCode))
	}
	defer res.Body.Close()

	bodyres, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	var main structs.AuthResponse
	err = json.Unmarshal(bodyres, &main)
	if err != nil {
		return "", err
	}

	return main.TOKEN, nil
}

func Start() {
	app = pb.New()

	app.OnBeforeServe().Add(func(e *core.ServeEvent) error {
		e.Router.AddRoute(
			echo.Route{
				Method:  http.MethodGet,
				Path:    "/app/api/health",
				Handler: health.Health,
			},
		)
		return nil
	})

	if err := app.Start(); err != nil {
		log.Fatal(err)
	}
}

func GetAllUsers() ([]structs.User, error) {
	var users []structs.User
	err := app.Dao().ConcurrentDB().NewQuery("SELECT * FROM users").All(&users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func ApplyPoints() {
	start := time.Now()
	users, err := GetAllUsers()
	if err != nil {
		fmt.Println("Error occured while getting users to apply points, trying again in 10s:")
		fmt.Println(err)
		time.Sleep(time.Second * 10)
		ApplyPoints()
	}
	for i := 0; i < len(users); i++ {
		if users[i].H_Score == 0 {
			continue
		}
		err = user_points(users[i])
		if err != nil {
			fmt.Println("Error occured while apply points, trying again in 10s:")
			fmt.Println(err)
			time.Sleep(time.Second * 10)
			ApplyPoints()
		}
	}
	fmt.Printf("Successfully added todays score to SCORE in %dms\n", time.Since(start).Milliseconds())
}

func user_points(user structs.User) error {
	record, err := app.Dao().FindRecordById("users", user.Id)
	if err != nil {
		return err
	}
	record.Set("h_score", 0)
	record.Set("score", user.H_Score+user.Score)
	err = app.Dao().SaveRecord(record)
	return err
}

func EditField(identifier string, id string, collection string, data interface{}) error {
	record, err := app.Dao().FindRecordById(collection, id)
	if err != nil {
		return err
	}
	record.Set(identifier, data)
	err = app.Dao().SaveRecord(record)
	return err
}

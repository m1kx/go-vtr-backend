package structs

import "github.com/pocketbase/pocketbase/tools/types"

// struct for auth req response
type AuthResponse struct {
	Admin struct {
		ID      string `json:"id"`
		CREATED string `json:"created"`
		UPDATED string `json:"updated"`
		AVATAR  int    `json:"avatar"`
		EMAIL   string `json:"email"`
	}
	TOKEN string `json:"token"`
}

type User struct {
	Id, Username, Email, Subjects, Class, H_Hash, M_Hash, Notifications string
	NewUpdate, NewVersion, Verified                                     bool
	Score, H_Score                                                      int
	ReqInfo                                                             types.JsonMap `db:"reqinfo"`
}

type HealthResponse struct {
	Status, Last_Words string
}

// request props for sending messag per request
type HttpReq struct {
	METHOD   string
	PROPS    string
	BASE_URL string
}

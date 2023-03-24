package structs

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

// struct for user database
type UserResponse struct {
	PAGE       int    `json:"page"`
	PERPAGE    int    `json:"perPage"`
	TOTALITEMS int    `json:"totalItems"`
	TOTALPAGES int    `json:"totalPages"`
	ITEMS      []User `json:"items"`
}

// struct for auth user
type User struct {
	COLLECTIONID    string      `json:"collectionId"`
	COLLECTIONNAME  string      `json:"collectionName"`
	CREATED         string      `json:"created"`
	EMAIL           string      `json:"email"`
	EMAILVISIBILITY bool        `json:"emailVisibility"`
	H_HASH          string      `json:"h_hash"`
	M_HASH          string      `json:"m_hash"`
	UPDATE          bool        `json:"update"`
	SUBJECTS        string      `json:"subjects"`
	CLASS           string      `json:"class"`
	ID              string      `json:"id"`
	UPDATED         string      `json:"updated"`
	USERNAME        string      `json:"username"`
	VERIFIED        bool        `json:"verified"`
	SCORE           int         `json:"score"`
	NEW_VERSION     bool        `json:"new_version"`
	REQINFO         RequestInfo `json:"reqinfo"`
}

// struct for webhook info
type RequestInfo struct {
	URL     string
	METHOD  string
	INFOFMT string
}

// request props for sending messag per request
type HttpReq struct {
	METHOD   string
	PROPS    string
	BASE_URL string
}

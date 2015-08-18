package core

type logoutRequest struct {
	Email string `json:"email"`
}

type updateRequest struct {
	Name string `json:"name"`
}

type updateAppRequest struct {
	Name string `json:"name"`
	OS   string `json:"os"`
}

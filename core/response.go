package core

//NewDeveloperResponse developer created response
type NewDeveloperResponse struct {
	ObjectID    string `json:"objectID"`
	AccessToken string `json:"accessToken"`
}

type loginResponse struct {
	ObjectID    string `json:"objectID"`
	AccessToken string `json:"accessToken"`
	Status      string `json:"status"`
}

type verifyLogin struct {
	Status string `json:"status"`
}

type logoutResponse struct {
	Status string `json:"status"`
}

type queryDevResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type deleteResponse struct {
	Status string `json:"status"`
}

type okResponse struct {
	Status string `json:"status"`
}

//AppResponse app auth response
type AppResponse struct {
	APIToken  string `json:"apitoken"`
	APISecret string `json:"apisecret"`
	Name      string `json:"name"`
	OS        string `json:"os"`
}

type allAppResponse struct {
	Responses []AppResponse `json:"responses"`
}

//SaveFileResponse file is saved
type SaveFileResponse struct {
	FileID   string `json:"fileID"`
	FileName string `json:"fileName"`
}

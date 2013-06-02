package hamster

/*http response*/

//source: http://code.google.com/p/goweb/wiki/BuildingAPIs#Data_representation_-_What_does_the_response_looks_like?
type Response struct {

	// The context of the request that initiated this response
	C string `json:"c"`

	// The HTTP Status code of this response
	S int `json:"s"`

	// The data (if any) for this response
	D interface{} `json:"d"`

	// A list of any errors that occurred while processing
	// the response
	E string `json:"e"`
}

type ResponseData struct {
	Response Response `json:"response"`
}

type NewDeveloperResponse struct {
	ObjectId    string `json:"object_id"`
	AccessToken string `json:"access_token"`
}

type LoginResponse struct {
	ObjectId    string `json:"object_id"`
	AccessToken string `json:"access_token"`
	Status      string `json:"status"`
}

type VerifyLogin struct {
	Status string `json:"status"`
}

type LogoutResponse struct {
	Status string `json:"status"`
}

type QueryDevResponse struct {
	Name  string `json:"name"`
	Email string `json:"email"`
}

type DeleteResponse struct {
	Status string `json:"status"`
}

type AppResponse struct {
	ApiToken  string `json:"apitoken"`
	ApiSecret string `json:"apisecret"`
	Name      string `json:"name"`
	OS        string `json:"os"`
}

type AllAppResponse struct {
	Responses []AppResponse `json:"responses"`
}

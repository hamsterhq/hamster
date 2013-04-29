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

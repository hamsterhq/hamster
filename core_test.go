package main

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/hamsterhq/hamster/core"
	"github.com/kr/fernet"
)

/*Tests ending with OK is a positive test, while others are negative tests
*TODO: write more negative tests.
 */

var (
	host        = "http://localhost:" + os.Getenv("SERVER_PORT")
	contentType = "application/json"
	accessToken = ""
	devID       = ""
	apiToken    = ""
	apiSecret   = ""
	objectID    = ""
	objectName  = ""
	fileID      = ""
	fileName    = ""
)

func testHTTPRequest(verb string, resource string, body string) (*http.Response, error) {
	client := &http.Client{Transport: &http.Transport{DisableKeepAlives: true}}
	r, _ := http.NewRequest(verb, fmt.Sprintf("%s%s", host, resource), strings.NewReader(body))
	r.Header.Add("Content-Type", contentType)
	return client.Do(r)
}
func testHTTPRequestWithHeaders(verb string, resource string, body string, header map[string]string) (*http.Response, error) {

	client := &http.Client{Transport: &http.Transport{DisableKeepAlives: true}}
	r, _ := http.NewRequest(verb, fmt.Sprintf("%s%s", host, resource), strings.NewReader(body))
	r.Header.Add("Content-Type", contentType)
	for key, value := range header {
		r.Header.Add(key, value)
	}
	return client.Do(r)

}

func testPostPng(resource string, fileReader io.Reader, header map[string]string) (*http.Response, error) {
	client := &http.Client{Transport: &http.Transport{DisableKeepAlives: true}}
	r, _ := http.NewRequest("POST", fmt.Sprintf("%s%s", host, resource), fileReader)
	r.Header.Add("Content-Type", "image/png")
	for key, value := range header {
		r.Header.Add(key, value)
	}
	return client.Do(r)
}

func testHTTP(verb string, resource string, header map[string]string) (*http.Response, error) {
	client := &http.Client{Transport: &http.Transport{DisableKeepAlives: true}}
	r, _ := http.NewRequest(verb, fmt.Sprintf("%s%s", host, resource), nil)

	for key, value := range header {
		r.Header.Add(key, value)
	}
	return client.Do(r)

}

func testServer(f func(s *core.Server)) {

	server := core.NewServer("../hamster.toml")
	//server.Quiet()
	server.ListenAndServe()
	defer server.Shutdown()
	f(server)
}

func TestCreateDeveloperHeaderEmpty(t *testing.T) {
	testServer(func(s *core.Server) {
		res, err := testHTTPRequest("POST", "/api/v1/developers/", `{"name":"adnaan","email":"badr.adnaan@gmail.com","password":"mypassword"}`)
		if err != nil {
			t.Fatalf("Unable to create developer: %v", err)

		} else {

			body, _ := ioutil.ReadAll(res.Body)

			if res.StatusCode == 200 {
				t.Fatalf("able to create developer: %v", string(body))
			}

		}

	})

}

func TestCreateDeveloperHeaderNoTime(t *testing.T) {
	testServer(func(s *core.Server) {
		headers := make(map[string]string)
		headers["X-Access-Token"] = "Z0FBQUFBQlJvMVAtdnNqS1c2dkNNVlRJSjF6Q2x4LW5YaElCRWVvZ00yRE1UaU9nc0huU0hMVUVYRGNoX2ZzUHBQczhZSk9yaTJXOHNpZWl6R21RSmp4SnlPSVJNTDF2TWc9PQ=="
		res, err := testHTTPRequestWithHeaders("POST", "/api/v1/developers/", `{"name":"adnaan","email":"badr.adnaan@gmail.com","password":"mypassword"}`, headers)
		if err != nil {
			t.Fatalf("Unable to create developer: %v", err)

		} else {

			body, _ := ioutil.ReadAll(res.Body)

			if res.StatusCode == 200 {
				t.Fatalf("able to create developer: %v", string(body))
			}

		}

	})

}

func TestCreateDeveloperHeaderWithTimeOK(t *testing.T) {
	testServer(func(s *core.Server) {
		headers := make(map[string]string)
		k := fernet.MustDecodeKeys("YI1ZYdopn6usnQ/5gMAHg8+pNh6D0DdaJkytdoLWUj0=")
		tok, err := fernet.EncryptAndSign([]byte("mysharedtoken"), k[0])
		if err != nil {
			t.Fatalf("fernet encryption failed %v\n", err)
		}
		stok := base64.URLEncoding.EncodeToString(tok)
		headers["X-Access-Token"] = stok
		res, err := testHTTPRequestWithHeaders("POST", "/api/v1/developers/", `{"name":"adnaan","email":"badr.adnaan@gmail.com","password":"mypassword"}`, headers)
		if err != nil {
			t.Fatalf("Unable to create developer: %v", err)

		} else {

			body, _ := ioutil.ReadAll(res.Body)

			if res.StatusCode != 200 {
				t.Fatalf("Unable to create developer: %v", string(body))
			}

			response := core.NewDeveloperResponse{}
			err := json.Unmarshal(body, &response)
			if err != nil {
				t.Fatalf("fail to parse body: %v", string(body))
			}

			accessToken = response.AccessToken
			devID = response.ObjectID

		}

	})

}

func TestCreateDeveloperEmailUnique(t *testing.T) {
	testServer(func(s *core.Server) {
		headers := make(map[string]string)
		k := fernet.MustDecodeKeys("YI1ZYdopn6usnQ/5gMAHg8+pNh6D0DdaJkytdoLWUj0=")
		tok, err := fernet.EncryptAndSign([]byte("mysharedtoken"), k[0])
		if err != nil {
			t.Fatalf("fernet encryption failed %v\n", err)
		}
		stok := base64.URLEncoding.EncodeToString(tok)
		headers["X-Access-Token"] = stok
		res, err := testHTTPRequestWithHeaders("POST", "/api/v1/developers/", `{"name":"adnaan","email":"badr.adnaan@gmail.com","password":"mypassword"}`, headers)
		if err != nil {
			t.Fatalf("email unique failed %v", err)

		} else {

			body, _ := ioutil.ReadAll(res.Body)

			if res.StatusCode != 500 {
				t.Fatalf("able to create developer: %v", string(body))
			}
		}

	})

}

func TestCreateDeveloperEmailExists(t *testing.T) {
	testServer(func(s *core.Server) {
		headers := make(map[string]string)
		k := fernet.MustDecodeKeys("YI1ZYdopn6usnQ/5gMAHg8+pNh6D0DdaJkytdoLWUj0=")
		tok, err := fernet.EncryptAndSign([]byte("mysharedtoken"), k[0])
		if err != nil {
			t.Fatalf("fernet encryption failed %v\n", err)
		}
		stok := base64.URLEncoding.EncodeToString(tok)
		headers["X-Access-Token"] = stok

		res, err := testHTTPRequestWithHeaders("POST", "/api/v1/developers/", `{"name":"adnaan"}`, headers)
		if err != nil {
			t.Fatalf("email exists failed %v", err)

		} else {
			body, _ := ioutil.ReadAll(res.Body)

			if res.StatusCode != 500 {
				t.Fatalf("able to create developer: %v", string(body))

			}
		}

	})

}

func TestLoginOK(t *testing.T) {
	testServer(func(s *core.Server) {

		headers := make(map[string]string)
		userpass := base64.StdEncoding.EncodeToString([]byte("badr.adnaan@gmail.com:mypassword"))
		headers["Authorization"] = "Basic " + userpass

		//headers["X-Access-Token"] = accessToken
		//make request
		client := &http.Client{Transport: &http.Transport{DisableKeepAlives: true}}
		r, _ := http.NewRequest("POST", fmt.Sprintf("%s%s", host, "/api/v1/developers/login/"), nil)

		for key, value := range headers {
			r.Header.Add(key, value)
		}

		res, err := client.Do(r)

		//res, err := testHTTP("GET", "/api/v1/developers/login/", headers)
		if err != nil {
			t.Fatalf("login failed! %v", err)

		} else {
			body, _ := ioutil.ReadAll(res.Body)

			if res.StatusCode != 200 {
				t.Fatalf("unable to login: %v", string(body))
			}
		}

	})

}

func TestQueryDeveloperOK(t *testing.T) {

	testServer(func(s *core.Server) {

		headers := make(map[string]string)
		headers["X-Access-Token"] = accessToken

		url := "/api/v1/developers/" + devID
		//make request
		res, err := testHTTPRequestWithHeaders("GET", url, ``, headers)

		if err != nil {
			t.Fatalf("unable to query: %v , %v", url, err)

		} else {

			body, _ := ioutil.ReadAll(res.Body)

			if res.StatusCode != 200 {
				t.Fatalf("unable to query: %v , %v", url, string(body))
			}

			//fmt.Printf("query response: %v \n", string(body))

		}

	})

}

func TestUpdateDeveloperOK(t *testing.T) {

	testServer(func(s *core.Server) {

		headers := make(map[string]string)
		headers["X-Access-Token"] = accessToken

		url := "/api/v1/developers/" + devID
		//make request
		res, err := testHTTPRequestWithHeaders("PUT", url, `{"name":"adnaan badr"}`, headers)

		if err != nil {
			t.Fatalf("unable to query: %v , %v", url, err)

		} else {

			body, _ := ioutil.ReadAll(res.Body)

			if res.StatusCode != 200 {
				t.Fatalf("unable to query: %v , %v", url, string(body))
			}

			//fmt.Printf("query response: %v ", string(body))

		}

	})

}

func TestCreateAppOK(t *testing.T) {

	testServer(func(s *core.Server) {

		headers := make(map[string]string)
		headers["X-Access-Token"] = accessToken

		url := "/api/v1/developers/" + devID + "/apps/"
		//make request
		res, err := testHTTPRequestWithHeaders("POST", url, `{"name":"traverik","os":"android"}`, headers)

		if err != nil {
			t.Fatalf("unable to create app: %v , %v", url, err)

		} else {

			body, _ := ioutil.ReadAll(res.Body)

			if res.StatusCode != 200 {
				t.Fatalf("unable to create app: %v , %v", url, string(body))
			}

			//fmt.Printf("query response: %v ", string(body))

			response := core.AppResponse{}
			err := json.Unmarshal(body, &response)
			if err != nil {
				t.Fatalf("fail to parse body: %v", string(body))
			}

			apiToken = response.APIToken
			apiSecret = response.APISecret

			fmt.Printf("apiToken: %v\n", apiToken)
			fmt.Printf("apiSecret: %v\n", apiSecret)

		}

	})

}

func TestQueryAppOK(t *testing.T) {

	testServer(func(s *core.Server) {

		headers := make(map[string]string)
		headers["X-Access-Token"] = accessToken

		url := "/api/v1/developers/apps/" + apiToken
		//make request
		res, err := testHTTPRequestWithHeaders("GET", url, ``, headers)

		if err != nil {
			t.Fatalf("unable to query: %v , %v", url, err)

		} else {

			body, _ := ioutil.ReadAll(res.Body)

			if res.StatusCode != 200 {
				t.Fatalf("unable to query: %v , %v", url, string(body))
			}

			//fmt.Printf("query response: %v \n", string(body))

		}

	})

}

func TestQueryAllAppsOK(t *testing.T) {

	testServer(func(s *core.Server) {

		headers := make(map[string]string)
		headers["X-Access-Token"] = accessToken

		url := "/api/v1/developers/" + devID + "/apps/"
		//make request
		res, err := testHTTPRequestWithHeaders("GET", url, ``, headers)

		if err != nil {
			t.Fatalf("unable to query: %v , %v", url, err)

		} else {

			body, _ := ioutil.ReadAll(res.Body)

			if res.StatusCode != 200 {
				t.Fatalf("unable to query: %v , %v", url, string(body))
			}

			//fmt.Printf("query response: %v \n", string(body))

		}

	})

}

func TestUpdateAppOK(t *testing.T) {

	testServer(func(s *core.Server) {

		headers := make(map[string]string)
		headers["X-Access-Token"] = accessToken

		url := "/api/v1/developers/apps/" + apiToken
		//make request
		res, err := testHTTPRequestWithHeaders("PUT", url, `{"name":"traverik alpha","os":"iOS"}`, headers)

		if err != nil {
			t.Fatalf("unable to update: %v , %v", url, err)

		} else {

			body, _ := ioutil.ReadAll(res.Body)

			if res.StatusCode != 200 {
				t.Fatalf("unable to update: %v , %v", url, string(body))
			}

			//fmt.Printf("update response: %v ", string(body))

		}

	})

}

type GameScore struct {
	Score      int      `json:"score"`
	PlayerName string   `json:"playerName"`
	Skills     []string `json:"skills"`
}

func getGameScore(score int, name string, t *testing.T) string {
	skills := make([]string, 2)
	skills[0] = "dying"
	skills[1] = "rebirth"

	gs := GameScore{Score: score, PlayerName: name, Skills: skills}
	s, err := json.MarshalIndent(&gs, "", "  ")
	if err != nil {

		t.Fatalf("marshal score error: %v ", err)

	}
	return string(s)

}

func TestCreateObjectOK(t *testing.T) {

	testServer(func(s *core.Server) {
		//test create
		score := getGameScore(1001, "adnaan", t)
		headers := make(map[string]string)
		headers["X-Api-Token"] = apiToken
		headers["X-Api-Secret"] = apiSecret

		url := "/api/v1/objects/" + "GameScore"
		//make request
		res, err := testHTTPRequestWithHeaders("POST", url, string(score), headers)

		if err != nil {
			t.Fatalf("unable to create object: %v , %v", url, err)

		} else {

			body, _ := ioutil.ReadAll(res.Body)

			if res.StatusCode != 200 {
				t.Fatalf("unable to create object: %v , %v", url, string(body))
			}

			//fmt.Printf("create object response: %v \n", string(body))
			var response map[string]interface{}
			err := json.Unmarshal(body, &response)
			if err != nil {
				t.Fatalf("fail to parse body: %v", string(body))
			}

			objectID = response["objectID"].(string)

		}

	})

}

func TestQueryObjectOK(t *testing.T) {

	testServer(func(s *core.Server) {

		headers := make(map[string]string)
		headers["X-Api-Token"] = apiToken
		headers["X-Api-Secret"] = apiSecret

		url := "/api/v1/objects/" + "GameScore/" + objectID
		//make request
		res, err := testHTTPRequestWithHeaders("GET", url, ``, headers)

		if err != nil {
			t.Fatalf("unable to query: %v , %v", url, err)

		} else {

			body, _ := ioutil.ReadAll(res.Body)

			if res.StatusCode != 200 {
				t.Fatalf("unable to query: %v , %v", url, string(body))
			}

			//fmt.Printf("object query response: %v \n", string(body))

		}

	})

}

type GameScoreBatch struct {
	Batch      []GameScore `json:"objects"`
	Operation  string      `json:"__op"`
	ObjectName string      `json:"objectName"`
}

func getGameScoreBatch(baseScore int, baseName string, t *testing.T) string {
	skills := make([]string, 2)
	skills[0] = "dying"
	skills[1] = "rebirth"

	var batch []GameScore

	for i := 0; i < 3; i++ {
		name := fmt.Sprintf(baseName+"%v", i)

		batch = append(batch, GameScore{Score: baseScore + i, PlayerName: name, Skills: skills})

	}

	gs := GameScoreBatch{Operation: "InsertBatch", Batch: batch, ObjectName: "GameScore"}
	s, err := json.MarshalIndent(&gs, "", "  ")
	if err != nil {

		t.Fatalf("marshal score error: %v ", err)

	}
	return string(s)

}

func TestCreateManyObjectsOK(t *testing.T) {

	testServer(func(s *core.Server) {
		//test create
		scores := getGameScoreBatch(1005, "adnaan", t)
		headers := make(map[string]string)
		headers["X-Api-Token"] = apiToken
		headers["X-Api-Secret"] = apiSecret

		url := "/api/v1/objects/batch" + "GameScore"
		//make request
		res, err := testHTTPRequestWithHeaders("POST", url, scores, headers)

		if err != nil {
			t.Fatalf("unable to create objects: %v , %v", url, err)

		} else {

			body, _ := ioutil.ReadAll(res.Body)

			if res.StatusCode != 200 {
				t.Fatalf("unable to create objects: %v , %v", url, string(body))
			}

			//fmt.Printf("create objects response: %v \n", string(body))
			/*var response map[string]interface{}
			err := json.Unmarshal(body, &response)
			if err != nil {
				t.Fatalf("fail to parse body: %v", string(body))
			}

			objectID = response["objectID"].(string)*/

		}

	})

}

func TestQueryObjectsOK(t *testing.T) {

	testServer(func(s *core.Server) {

		headers := make(map[string]string)
		headers["X-Api-Token"] = apiToken
		headers["X-Api-Secret"] = apiSecret

		url := "/api/v1/objects/" + "GameScore"
		//make request
		res, err := testHTTPRequestWithHeaders("GET", url, ``, headers)

		if err != nil {
			t.Fatalf("unable to query: %v , %v", url, err)

		} else {

			body, _ := ioutil.ReadAll(res.Body)

			if res.StatusCode != 200 {
				t.Fatalf("unable to query: %v , %v", url, string(body))
			}

			//fmt.Printf("objects query response: %v \n", string(body))

		}

	})

}

func TestUpdateObjectOK(t *testing.T) {

	testServer(func(s *core.Server) {

		headers := make(map[string]string)
		headers["X-Api-Token"] = apiToken
		headers["X-Api-Secret"] = apiSecret

		url := "/api/v1/objects/" + "GameScore/" + objectID
		//make request
		res, err := testHTTPRequestWithHeaders("PUT", url, `{"playerName":"superman"}`, headers)

		if err != nil {
			t.Fatalf("unable to update: %v , %v", url, err)

		} else {

			body, _ := ioutil.ReadAll(res.Body)

			if res.StatusCode != 200 {
				t.Fatalf("unable to update: %v , %v", url, string(body))
			}

			//fmt.Printf("update object response: %v ", string(body))

		}

	})

}

func TestSaveImageOK(t *testing.T) {
	testServer(func(s *core.Server) {
		//test create
		filePath := ""
		fileName = "gophers.png"
		file, err := os.Open(filePath + fileName)
		if err != nil {
			t.Fatalf("unable to open image: %v", err)
		}

		defer file.Close()
		fileReader := bufio.NewReader(file)

		metadata := make(map[string]interface{})
		metadata["category"] = "screenshot"
		metadata["view"] = "homeview"
		metadata["width"] = 480
		metadata["height"] = 854

		meta, err := json.MarshalIndent(metadata, "", "  ")
		if err != nil {

			t.Fatalf("marshal meta error: %v ", err)

		}
		//fmt.Printf("meta %v\n", string(meta))

		headers := make(map[string]string)

		headers["X-Api-Token"] = apiToken
		headers["X-Api-Secret"] = apiSecret
		headers["X-Meta-Data"] = string(meta)

		url := "/api/v1/files/" + fileName
		//make request
		res, err := testPostPng(url, fileReader, headers)

		if err != nil {
			t.Fatalf("unable to create object: %v , %v", url, err)

		} else {

			body, _ := ioutil.ReadAll(res.Body)

			if res.StatusCode != 200 {
				t.Fatalf("unable to save file: %v , %v", url, string(body))
				return
			}

			//fmt.Printf("save file response: %v \n", string(body))
			response := core.SaveFileResponse{}
			err := json.Unmarshal(body, &response)
			if err != nil {
				t.Fatalf("fail to parse body: %v", string(body))
			}

			fileID = response.FileID
			fileName = response.FileName

		}

	})

}

func TestGetImageOK(t *testing.T) {

	testServer(func(s *core.Server) {

		headers := make(map[string]string)
		headers["X-Api-Token"] = apiToken
		headers["X-Api-Secret"] = apiSecret

		url := "/api/v1/files/" + fileName + "/" + fileID
		//make request
		res, err := testHTTPRequestWithHeaders("GET", url, ``, headers)

		if err != nil {
			t.Fatalf("unable to get image: %v , %v", url, err)

		} else {

			body, _ := ioutil.ReadAll(res.Body)

			if res.StatusCode != 200 {
				t.Fatalf("unable to image: %v , %v", url, string(body))
			}

			//fmt.Printf("object query response: %v \n", string(body))
			//write file
			f, err := os.OpenFile("download.png", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
			if err != nil {
				t.Fatalf("unable to open image ")

			}
			defer f.Close()

			/*// Grab the image data
			var buf bytes.Buffer
			io.Copy(&buf, res.Body)
			//decode
			_, format, err := image.Decode(&buf)
			if err != nil {
				t.Fatalf("unable to decode %v  %v \n ", err, format)
			}*/

			f.Write(body)
		}

	})

}
func TestDeleteObjectOK(t *testing.T) {

	testServer(func(s *core.Server) {
		//test delete
		headers := make(map[string]string)
		headers["X-Api-Token"] = apiToken
		headers["X-Api-Secret"] = apiSecret

		url := "/api/v1/objects/" + "GameScore/" + objectID

		//make request
		res, err := testHTTPRequestWithHeaders("DELETE", url, ``, headers)

		if err != nil {
			t.Fatalf("unable to delete object: %v , %v", url, err)

		} else {

			body, _ := ioutil.ReadAll(res.Body)

			if res.StatusCode != 200 {
				t.Fatalf("unable to delete object: %v , %v", url, string(body))
			}

			//fmt.Printf("object delete response: %v\n ", string(body))

		}

	})

}

func TestDeleteAppOK(t *testing.T) {

	testServer(func(s *core.Server) {
		//test delete

		headers := make(map[string]string)
		headers["X-Access-Token"] = accessToken

		url := "/api/v1/developers/apps/" + apiToken
		//make request
		res, err := testHTTPRequestWithHeaders("DELETE", url, ``, headers)

		if err != nil {
			t.Fatalf("unable to delete app: %v , %v", url, err)

		} else {

			body, _ := ioutil.ReadAll(res.Body)

			if res.StatusCode != 200 {
				t.Fatalf("unable to delete app: %v , %v", url, string(body))
			}

			//fmt.Printf("query response: %v\n ", string(body))

		}

	})

}

func TestDeleteDeveloperOK(t *testing.T) {

	testServer(func(s *core.Server) {
		//test delete

		headers := make(map[string]string)
		headers["X-Access-Token"] = accessToken

		url := "/api/v1/developers/" + devID
		//make request
		res, err := testHTTPRequestWithHeaders("DELETE", url, ``, headers)

		if err != nil {
			t.Fatalf("unable to delete: %v , %v", url, err)

		} else {

			body, _ := ioutil.ReadAll(res.Body)

			if res.StatusCode != 200 {
				t.Fatalf("unable to delete: %v , %v", url, string(body))
			}

			//fmt.Printf("delete response: %v\n ", string(body))

		}

	})

}

func TestLogoutOK(t *testing.T) {
	testServer(func(s *core.Server) {

		headers := make(map[string]string)
		headers["X-Access-Token"] = accessToken

		//make request
		res, err := testHTTPRequestWithHeaders("POST", "/api/v1/developers/logout/", `{"email":"badr.adnaan@gmail.com"}`, headers)

		if err != nil {
			t.Fatalf("unable to logout: %v", err)

		} else {

			body, _ := ioutil.ReadAll(res.Body)

			if res.StatusCode != 200 {
				t.Fatalf("unable to logout: %v", string(body))
			}

		}

	})
}

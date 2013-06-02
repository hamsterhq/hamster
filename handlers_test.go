package hamster

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/kr/fernet"
	"io/ioutil"
	"net/http"
	"testing"
)

var access_token = ""
var dev_id = ""
var api_token = ""
var api_secret = ""
var object_id = ""
var object_name = ""

//developers handlers
func TestCreateDeveloperHeaderEmpty(t *testing.T) {
	testServer(func(s *Server) {
		res, err := testHttpRequest("POST", "/api/v1/developers/", `{"name":"adnaan","email":"badr.adnaan@gmail.com","password":"mypassword"}`)
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
	testServer(func(s *Server) {
		headers := make(map[string]string)
		headers["X-Access-Token"] = "Z0FBQUFBQlJvMVAtdnNqS1c2dkNNVlRJSjF6Q2x4LW5YaElCRWVvZ00yRE1UaU9nc0huU0hMVUVYRGNoX2ZzUHBQczhZSk9yaTJXOHNpZWl6R21RSmp4SnlPSVJNTDF2TWc9PQ=="
		res, err := testHttpRequestWithHeaders("POST", "/api/v1/developers/", `{"name":"adnaan","email":"badr.adnaan@gmail.com","password":"mypassword"}`, headers)
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
	testServer(func(s *Server) {
		headers := make(map[string]string)
		k := fernet.MustDecodeKeys("YI1ZYdopn6usnQ/5gMAHg8+pNh6D0DdaJkytdoLWUj0=")
		tok, err := fernet.EncryptAndSign([]byte("mysharedtoken"), k[0])
		if err != nil {
			t.Fatalf("fernet encryption failed %v\n", err)
		}
		stok := base64.URLEncoding.EncodeToString(tok)
		headers["X-Access-Token"] = stok
		res, err := testHttpRequestWithHeaders("POST", "/api/v1/developers/", `{"name":"adnaan","email":"badr.adnaan@gmail.com","password":"mypassword"}`, headers)
		if err != nil {
			t.Fatalf("Unable to create developer: %v", err)

		} else {

			body, _ := ioutil.ReadAll(res.Body)

			if res.StatusCode != 200 {
				t.Fatalf("unable to create developer: %v", string(body))
			}

			response := NewDeveloperResponse{}
			err := json.Unmarshal(body, &response)
			if err != nil {
				t.Fatalf("fail to parse body: %v", string(body))
			}

			access_token = response.AccessToken
			dev_id = response.ObjectId

		}

	})

}

func TestCreateDeveloperEmailUnique(t *testing.T) {
	testServer(func(s *Server) {
		headers := make(map[string]string)
		k := fernet.MustDecodeKeys("YI1ZYdopn6usnQ/5gMAHg8+pNh6D0DdaJkytdoLWUj0=")
		tok, err := fernet.EncryptAndSign([]byte("mysharedtoken"), k[0])
		if err != nil {
			t.Fatalf("fernet encryption failed %v\n", err)
		}
		stok := base64.URLEncoding.EncodeToString(tok)
		headers["X-Access-Token"] = stok
		res, err := testHttpRequestWithHeaders("POST", "/api/v1/developers/", `{"name":"adnaan","email":"badr.adnaan@gmail.com","password":"mypassword"}`, headers)
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
	testServer(func(s *Server) {
		headers := make(map[string]string)
		k := fernet.MustDecodeKeys("YI1ZYdopn6usnQ/5gMAHg8+pNh6D0DdaJkytdoLWUj0=")
		tok, err := fernet.EncryptAndSign([]byte("mysharedtoken"), k[0])
		if err != nil {
			t.Fatalf("fernet encryption failed %v\n", err)
		}
		stok := base64.URLEncoding.EncodeToString(tok)
		headers["X-Access-Token"] = stok

		res, err := testHttpRequestWithHeaders("POST", "/api/v1/developers/", `{"name":"adnaan"}`, headers)
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
	testServer(func(s *Server) {
		//test login
		host = "http://localhost:8686"
		headers := make(map[string]string)
		userpass := base64.StdEncoding.EncodeToString([]byte("badr.adnaan@gmail.com:mypassword"))
		headers["Authorization"] = "Basic " + userpass

		headers["X-Access-Token"] = access_token
		//make request
		client := &http.Client{Transport: &http.Transport{DisableKeepAlives: true}}
		r, _ := http.NewRequest("POST", fmt.Sprintf("%s%s", host, "/api/v1/developers/login/"), nil)

		for key, value := range headers {
			r.Header.Add(key, value)
		}

		res, err := client.Do(r)

		//res, err := testHttp("GET", "/api/v1/developers/login/", headers)
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

func TestQueryDeveloper(t *testing.T) {

	testServer(func(s *Server) {
		//test login

		headers := make(map[string]string)
		headers["X-Access-Token"] = access_token

		url := "/api/v1/developers/" + dev_id
		//make request
		res, err := testHttpRequestWithHeaders("GET", url, ``, headers)

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

func TestUpdateDeveloper(t *testing.T) {

	testServer(func(s *Server) {
		//test login

		headers := make(map[string]string)
		headers["X-Access-Token"] = access_token

		url := "/api/v1/developers/" + dev_id
		//make request
		res, err := testHttpRequestWithHeaders("PUT", url, `{"name":"adnaan badr"}`, headers)

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

func TestCreateApp(t *testing.T) {

	testServer(func(s *Server) {
		//test login

		headers := make(map[string]string)
		headers["X-Access-Token"] = access_token

		url := "/api/v1/developers/" + dev_id + "/apps/"
		//make request
		res, err := testHttpRequestWithHeaders("POST", url, `{"name":"traverik","os":"android"}`, headers)

		if err != nil {
			t.Fatalf("unable to create app: %v , %v", url, err)

		} else {

			body, _ := ioutil.ReadAll(res.Body)

			if res.StatusCode != 200 {
				t.Fatalf("unable to create app: %v , %v", url, string(body))
			}

			//fmt.Printf("query response: %v ", string(body))

			response := AppResponse{}
			err := json.Unmarshal(body, &response)
			if err != nil {
				t.Fatalf("fail to parse body: %v", string(body))
			}

			api_token = response.ApiToken
			api_secret = response.ApiSecret

		}

	})

}

func TestQueryApp(t *testing.T) {

	testServer(func(s *Server) {
		//test login

		headers := make(map[string]string)
		headers["X-Access-Token"] = access_token

		url := "/api/v1/developers/apps/" + api_token
		//make request
		res, err := testHttpRequestWithHeaders("GET", url, ``, headers)

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

func TestQueryAllApps(t *testing.T) {

	testServer(func(s *Server) {
		//test login

		headers := make(map[string]string)
		headers["X-Access-Token"] = access_token

		url := "/api/v1/developers/" + dev_id + "/apps/"
		//make request
		res, err := testHttpRequestWithHeaders("GET", url, ``, headers)

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

func TestUpdateApp(t *testing.T) {

	testServer(func(s *Server) {
		//test login

		headers := make(map[string]string)
		headers["X-Access-Token"] = access_token

		url := "/api/v1/developers/apps/" + api_token
		//make request
		res, err := testHttpRequestWithHeaders("PUT", url, `{"name":"traverik alpha","os":"iOS"}`, headers)

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

func TestCreateObject(t *testing.T) {

	testServer(func(s *Server) {
		//test create
		score := getGameScore(1001, "adnaan", t)
		headers := make(map[string]string)
		headers["X-Api-Token"] = api_token
		headers["X-Api-Secret"] = api_secret

		url := "/api/v1/objects/" + "GameScore"
		//make request
		res, err := testHttpRequestWithHeaders("POST", url, string(score), headers)

		if err != nil {
			t.Fatalf("unable to create object: %v , %v", url, err)

		} else {

			body, _ := ioutil.ReadAll(res.Body)

			if res.StatusCode != 200 {
				t.Fatalf("unable to create object: %v , %v", url, string(body))
			}

			fmt.Printf("create object response: %v \n", string(body))
			var response map[string]interface{}
			err := json.Unmarshal(body, &response)
			if err != nil {
				t.Fatalf("fail to parse body: %v", string(body))
			}

			object_id = response["object_id"].(string)

		}

	})

}

func TestQueryObject(t *testing.T) {

	testServer(func(s *Server) {
		//test login

		headers := make(map[string]string)
		headers["X-Api-Token"] = api_token
		headers["X-Api-Secret"] = api_secret

		url := "/api/v1/objects/" + "GameScore/" + object_id
		//make request
		res, err := testHttpRequestWithHeaders("GET", url, ``, headers)

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

func TestCreateManyObjects(t *testing.T) {

	scores := getGameScoreBatch(1005, "adnaan", t)
	fmt.Printf("batch: %v\n", scores)

}

func TestQueryObjects(t *testing.T) {

	testServer(func(s *Server) {
		//test login

		headers := make(map[string]string)
		headers["X-Api-Token"] = api_token
		headers["X-Api-Secret"] = api_secret

		url := "/api/v1/objects/" + "GameScore"
		//make request
		res, err := testHttpRequestWithHeaders("GET", url, ``, headers)

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

func TestUpdateObject(t *testing.T) {

	testServer(func(s *Server) {
		//test login
		headers := make(map[string]string)
		headers["X-Api-Token"] = api_token
		headers["X-Api-Secret"] = api_secret

		url := "/api/v1/objects/" + "GameScore/" + object_id
		//make request
		res, err := testHttpRequestWithHeaders("PUT", url, `{"playerName":"superman"}`, headers)

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

func TestDeleteObject(t *testing.T) {

	testServer(func(s *Server) {
		//test delete
		headers := make(map[string]string)
		headers["X-Api-Token"] = api_token
		headers["X-Api-Secret"] = api_secret

		url := "/api/v1/objects/" + "GameScore/" + object_id

		//make request
		res, err := testHttpRequestWithHeaders("DELETE", url, ``, headers)

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

func TestDeleteApp(t *testing.T) {

	testServer(func(s *Server) {
		//test delete

		headers := make(map[string]string)
		headers["X-Access-Token"] = access_token

		url := "/api/v1/developers/apps/" + api_token
		//make request
		res, err := testHttpRequestWithHeaders("DELETE", url, ``, headers)

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

func TestDeleteDeveloper(t *testing.T) {

	testServer(func(s *Server) {
		//test delete

		headers := make(map[string]string)
		headers["X-Access-Token"] = access_token

		url := "/api/v1/developers/" + dev_id
		//make request
		res, err := testHttpRequestWithHeaders("DELETE", url, ``, headers)

		if err != nil {
			t.Fatalf("unable to query: %v , %v", url, err)

		} else {

			body, _ := ioutil.ReadAll(res.Body)

			if res.StatusCode != 200 {
				t.Fatalf("unable to query: %v , %v", url, string(body))
			}

			//fmt.Printf("query response: %v\n ", string(body))

		}

	})

}

func TestLogoutOK(t *testing.T) {
	testServer(func(s *Server) {
		//test login

		headers := make(map[string]string)
		headers["X-Access-Token"] = access_token

		//make request
		res, err := testHttpRequestWithHeaders("POST", "/api/v1/developers/logout/", `{"email":"badr.adnaan@gmail.com"}`, headers)

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

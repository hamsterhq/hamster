package hamster

import (
	"io/ioutil"
	"testing"
)

//developers handlers
func TestCreateDeveloper(t *testing.T) {
	testServer(func(s *Server) {
		res, err := testHttpRequest("POST", "/developers", `{"name":"adnaan","email":"badr1.adnaan@gmail.com","password":"mypassword"}`)
		if err != nil {
			t.Fatalf("Unable to create developer: %v", err)

		} else {

			body, _ := ioutil.ReadAll(res.Body)

			if res.StatusCode != 200 {
				t.Fatalf("unable to create developer: %v", string(body))
			} else {
				s.logger.SetPrefix("test: ")
				s.logger.Printf("response: %+v\n", string(body))
			}

		}

	})

}

func TestCreateDeveloperEmailUnique(t *testing.T) {
	testServer(func(s *Server) {
		res, err := testHttpRequest("POST", "/developers", `{"name":"adnaan","email":"badr.adnaan@gmail.com","password":"mypassword"}`)
		if err != nil {
			t.Fatalf("email unique failed %v", err)

		} else {

			body, _ := ioutil.ReadAll(res.Body)

			if res.StatusCode == 500 {
				t.Fatalf("unable to create developer: %v", string(body))
			} else {
				s.logger.SetPrefix("test: ")
				s.logger.Printf("response: %+v\n", string(body))
			}
		}

	})

}

func TestCreateDeveloperEmailExists(t *testing.T) {
	testServer(func(s *Server) {
		res, err := testHttpRequest("POST", "/developers", `{"name":"adnaan"}`)
		if err != nil {
			t.Fatalf("email exists failed %v", err)

		} else {
			body, _ := ioutil.ReadAll(res.Body)

			if res.StatusCode == 500 {
				s.logger.SetPrefix("test: ")
				s.logger.Printf("response: %+v\n", string(body))

			} else {
				t.Fatalf("able to create developer: %v", string(body))

			}
		}

	})

}

package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Auth struct{}

func (a Auth) DoLogin(body []byte, params url.Values) (string, error) {
	if string(body) == "username=admin&password=1234" {
		return "mocktoken123", nil
	}
	return "", fmt.Errorf("invalid credentials")
}

var myAuth = Auth{}

func main() {
	http.HandleFunc("/login", func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			http.NotFound(res, req)
			return
		}
		params := req.URL.Query()
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			http.Error(res, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		token, err := myAuth.DoLogin(body, params)
		if err == nil {
			res.WriteHeader(http.StatusOK)
			_, _ = res.Write([]byte(fmt.Sprintf("Login successful. Token: %s", token)))
		} else {
			res.WriteHeader(http.StatusUnauthorized)
			_, _ = res.Write([]byte("Unauthorized"))
		}
	})

	fmt.Println("Server starting on port 8080")
	http.ListenAndServe(":8080", nil)
}

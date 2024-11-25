package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

type Auth struct{}

var (
	blacklist   = make(map[string]bool)
	blacklistMu sync.Mutex
)

func (a Auth) DoLogin(body []byte, params url.Values) (string, error) {
	if string(body) == "username=admin&password=1234" {
		return "mocktoken123", nil
	}
	return "", fmt.Errorf("invalid credentials")
}

func (a Auth) authenticate(authHeader string) bool {
	if authHeader == "" {
		return false
	}

	const prefix = "Token "
	if !strings.HasPrefix(authHeader, prefix) {
		return false
	}

	token := strings.TrimPrefix(authHeader, prefix)

	blacklistMu.Lock()
	defer blacklistMu.Unlock()
	if blacklist[token] {
		return false
	}
	return token == "mocktoken123"
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
			res.Header().Set("Authorization", fmt.Sprintf(" Token %s", token))
			res.WriteHeader(http.StatusOK)
			_, _ = res.Write([]byte(fmt.Sprintf("Login successful. Token: %s", token)))
		} else {
			res.WriteHeader(http.StatusUnauthorized)
			_, _ = res.Write([]byte("Unauthorized"))
		}
	})
	http.HandleFunc("/logout", func(res http.ResponseWriter, req *http.Request) {
		if req.Method != "POST" {
			http.NotFound(res, req)
			return
		}
		authHeader := req.Header.Get("Authorization")
		if !myAuth.authenticate(authHeader) {
			res.WriteHeader(http.StatusUnauthorized)
			_, _ = res.Write([]byte("Unauthorized"))
			return
		}

		token := strings.TrimPrefix(authHeader, "Token ")
		blacklistMu.Lock()
		blacklist[token] = true
		blacklistMu.Unlock()

		res.WriteHeader(http.StatusOK)
		_, _ = res.Write([]byte("Logged out successfully"))
	})

	fmt.Println("Server starting on port 8080")
	http.ListenAndServe(":8080", nil)
}

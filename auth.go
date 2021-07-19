package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type session struct {
	sessions map[string]sessionInfo
}

type sessionInfo struct {
	accessToken    string
	refreshToken   string
	expirationDate time.Time
}

func (a *App) InitializeAuth() {
	s := session{}
	a.Router.Use(s.Authenticate)

	a.Router.HandleFunc("/login", s.login).Methods("GET")
	// a.Router.HandleFunc("/callback", s.callback)
	a.Router.Path("/callback").
		Queries("code", "", "state", "").
		HandlerFunc(s.callback)
}

func (s *session) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session := r.Header.Get("session")

		if _, found := s.sessions[session]; found || true { // delete || true
			next.ServeHTTP(w, r)
		} else {
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	})
}

func (s *session) login(w http.ResponseWriter, r *http.Request) {
	session := r.Header.Get("session")
	hasher := sha256.New()
	hasher.Write([]byte(session))
	state := base64.URLEncoding.EncodeToString(hasher.Sum(nil))

	http.Redirect(w, r,
		"https://discord.com/api/oauth2/authorize?"+
			"client_id="+keys.CLIENT_ID+
			"&redirect_uri="+"https://localhost:8000/callback"+
			"&response_type=code"+
			"&scope=identify"+
			"&state="+state+
			"&prompt=consent", http.StatusTemporaryRedirect)
}

func (s *session) callback(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)

	code := mux.Vars(r)["code"]

	var body map[string]string

	println("3")
	println(s.sessions)

	session := body["session"]

	accessToken, refreshToken, err := getAccessToken(code)
	if err != nil {
		s.sessions[session] = sessionInfo{accessToken, refreshToken, time.Now()}
	}
	println("4")
	println(s.sessions)
}

type AuthTokenRequestBody struct {
	client_id     string
	client_secret string
	redirect_url  string
	grant_type    string
	code          string
	scope         string
}

type AuthTokenResponseBody struct {
	access_token  string
	expires_in    string
	refresh_token string
	scope         string
	token_type    string
}

func getAccessToken(code string) (accessToken string, refreshToken string, err error) {
	body := &AuthTokenRequestBody{
		keys.CLIENT_ID,
		keys.CLIENT_SECRET,
		"https://localhost:8000/callback",
		"authorization_code",
		"PkWE7J2CIajabve8uyNqrK5w7SJUJV",
		"identify",
	}

	jsonBytes, _ := json.Marshal(body)

	url := "https://discord.com/api/oauth2/token"
	res, err := http.Post(url, "application/x-www-form-urlencoded", bytes.NewBuffer(jsonBytes))
	if err != nil {
		return "", "", fmt.Errorf("Invalid code")
	}

	var resBody AuthTokenResponseBody
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&resBody); err != nil {
		return "", "", fmt.Errorf("Invalid response body")
	}
	defer res.Body.Close()

	return resBody.access_token, resBody.refresh_token, nil
}

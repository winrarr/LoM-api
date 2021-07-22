package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
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
	s.sessions = make(map[string]sessionInfo)
	a.Router.Use(s.Authenticate)

	a.Router.HandleFunc("/login", s.login).Methods("GET")
	a.Router.Path("/callback").
		Queries("code", "", "state", "").
		HandlerFunc(s.callback).
		Methods("GET")
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
	session, err := r.Cookie("session")
	if err != nil {
		log.Println(err)
		return
	}
	hasher := sha256.New()
	hasher.Write([]byte(session.Value))
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

type callbackRequest struct {
	session string
}

func (s *session) callback(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)

	code := r.FormValue("code")
	state := r.FormValue("state")
	session, err := r.Cookie("session")
	if err != nil {
		log.Println(err)
		return
	}
	println(code, state, session.Value)

	accessToken, refreshToken, err := getAccessToken(code)
	if err != nil {
		log.Println(err)
		return
	}
	s.sessions[session.Value] = sessionInfo{accessToken, refreshToken, time.Now()}
	fmt.Println("map:", s.sessions)
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

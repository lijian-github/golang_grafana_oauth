package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-session/session"
	"log"
	"net/http"
	"os"
)
type Token struct {
	AccessToken string `json:"access_token"`
	TokenType string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	ExpiresInt int `json:"expires_int"`
}

type Data struct{
	Email string `json:"email"`
	Id string `json:"id"`
	Name string `json:"name"`
	Role string `json:"role"`
	Username string `json:"username"`
}

func main() {
	http.HandleFunc("/login", loginHandler)


	http.HandleFunc("/token", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r)
		h:=w.Header()
		//必须设置Content-Type为application-json，否则Grafana后台不会去取access-token，
		//然后grafana前端报login.OAuthLogin(NewTransportWithCode)，日志报oauth2: server response missing access_token
		h.Set("Content-Type","application/json")
		//access_token是必须的
		token :=Token{
			AccessToken:  `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.
eyJpZCI6IjQiLCJuYW1lIjoidmlld2VyIiwidXNlcm5hbWUiOiJ2aWV3ZXIiLCJ3b3JrX2NvZGUiOiIxIiwiZW1haWwiOiIxMTI5ODIwMjk0QHFxLmNvbSJ9.
IIL8sC6tFcMaCcqpmhSxsKg9JyHKoTipFjWNW1YR5wY`,
			TokenType:    "Bearer",
			RefreshToken: "ODg2M2M3MDItMDQ2ZC00YzZiLWExNTktYzQ2NWJkYzI4MDE3.26a7d743c3bfdb70a5f8fe00d99f9a4199c45442",
			ExpiresInt:   6000,
		}
		fmt.Println(token)
		tokenB, _ := json.Marshal(token)
		w.Write(tokenB)
	})

	http.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
		outputHTML(w, r, "./static/grafana.html")
	})

	http.HandleFunc("/userinfo", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r)
		//必须要有Email
		data := Data{Email: "99",Role: "Editor",Username: "99"}
		ret_json,_ := json.Marshal(data)
		w.Write(ret_json)
	})

	log.Println("Server is running at 9096 port.")
	log.Fatal(http.ListenAndServe(":9096", nil))
}


func loginHandler(w http.ResponseWriter, r *http.Request) {
	store, err := session.Start(r.Context(), w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if r.Form == nil {
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
	store.Set("state", r.Form.Get("state"))
	store.Save()
	state, _ := store.Get("state")
	w.Header().Set("Location", r.Form.Get("redirect_uri")+"?state="+fmt.Sprint(state)+"&code=nbayramberdiyev")
	w.WriteHeader(http.StatusFound)
}


func outputHTML(w http.ResponseWriter, req *http.Request, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer file.Close()
	fi, _ := file.Stat()
	http.ServeContent(w, req, file.Name(), fi.ModTime(), file)
}

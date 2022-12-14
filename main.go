package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/livekit/protocol/auth"
)

type TokenHandler struct {
	Token string `json:"token"`
}

func (t *TokenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	apiKey := os.Getenv("apiKey")
	apiSecret := os.Getenv("apiSecret")
	room := "room123"

	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

	if r.Method == "OPTIONS" {
		fmt.Println("OPTION method")
		fmt.Fprintf(w, "option ok")
		return
	}

	if r.Method == "GET" {
		fmt.Println("GET method")
		fmt.Fprintf(w, "get ok")
		return
	}

	//To allocate slice for request body
	length, err := strconv.Atoi(r.Header.Get("Content-Length"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Read Body
	body, err := ioutil.ReadAll(r.Body)

	if err != nil {
		fmt.Println("StatusInternalServerError")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	fmt.Println(length)

	//parse to json
	var jsonBody map[string]string

	fmt.Println(jsonBody)

	err = json.Unmarshal(body, &jsonBody)
	if err != nil {
		fmt.Println(err)
		fmt.Println("Unmarshal:StatusInternalServerError")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//存在する
	if _, ok := jsonBody["identity"]; !ok {
		fmt.Println("identity not exists ")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	token, err := t.getJoinToken(apiKey, apiSecret, room, jsonBody["identity"])
	if err != nil {
		panic(err)
	}
	t.Token = token

	fmt.Println(jsonBody["identity"])

	res, err := json.Marshal(t)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write(res)

}

func (t *TokenHandler) getJoinToken(apiKey, apiSeacret, room, identity string) (string, error) {
	at := auth.NewAccessToken(apiKey, apiSeacret)
	grant := &auth.VideoGrant{
		RoomJoin: true,
		Room:     room,
	}
	at.AddGrant(grant).
		SetIdentity(identity).
		SetValidFor(time.Hour)

	return at.ToJWT()
}

func main() {
	http.Handle("/token", new(TokenHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}

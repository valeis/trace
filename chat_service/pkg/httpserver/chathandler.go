package httpserver

import (
	"chat_app/pkg/redisrepo"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type userReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Client   string `json:"client"`
}

type contactListReq struct {
	Username string `json:"username"`
}

type participantsReq struct {
	User1 string `json:"u1"`
	User2 string `json:"u2"`
}

type response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Total   int         `json:"total,omitempty"`
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	u := &userReq{}
	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		http.Error(w, "error decoding request object", http.StatusBadRequest)
		return
	}

	res := register(u)
	json.NewEncoder(w).Encode(res)
}

func verifyContactHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	u := &userReq{}
	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		http.Error(w, "error decoding request object", http.StatusBadRequest)
		return
	}

	res := verifyContact(u.Username)
	json.NewEncoder(w).Encode(res)
}

func chatHistoryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	participants := &participantsReq{}
	if err := json.NewDecoder(r.Body).Decode(participants); err != nil {
		http.Error(w, "error decoding request object", http.StatusBadRequest)
		return
	}

	u1 := participants.User1
	u2 := participants.User2

	fromTS, toTS := "0", "+inf"

	if r.URL.Query().Get("from-ts") != "" && r.URL.Query().Get("to-ts") != "" {
		fromTS = r.URL.Query().Get("from-ts")
		toTS = r.URL.Query().Get("to-ts")
	}

	res := chatHistory(u1, u2, fromTS, toTS)
	json.NewEncoder(w).Encode(res)
}

func contactListHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	u := &contactListReq{}
	log.Println(u)
	if err := json.NewDecoder(r.Body).Decode(u); err != nil {
		http.Error(w, "error decoding request object", http.StatusBadRequest)
		return
	}

	res := contactList(u.Username)
	json.NewEncoder(w).Encode(res)
}

func register(u *userReq) *response {
	res := &response{Status: true}

	status := redisrepo.IsUserExist(u.Username)
	if status {
		res.Status = false
		res.Message = "Username already taken. Try something else."
		return res
	}

	err := redisrepo.RegisterNewUser(u.Username, u.Password)
	if err != nil {
		res.Status = false
		res.Message = "something went wrong while registering the user. please try againg after sometime"
		return res
	}

	return res
}

func verifyContact(username string) *response {
	res := &response{Status: true}
	status := redisrepo.IsUserExist(username)

	if !status {
		res.Status = false
		res.Message = "invalid username"
	}
	return res
}

func chatHistory(username1, username2, fromTS, toTS string) *response {
	res := &response{}

	fmt.Println(username1, username2)
	if !redisrepo.IsUserExist(username1) || !redisrepo.IsUserExist(username2) {
		res.Message = "incorrect username"
		return res
	}

	chats, err := redisrepo.FetchChatBetween(username1, username2, fromTS, toTS)
	if err != nil {
		log.Println("error in fetch chat between", err)
		res.Message = "unable to fetch chat history. Please try again later."
		return res
	}

	res.Status = true
	res.Data = chats
	res.Total = len(chats)
	return res
}

func contactList(username string) *response {
	res := &response{}

	if !redisrepo.IsUserExist(username) {
		res.Message = "incorrect username"
		return res
	}

	contactList, err := redisrepo.FetchContactList(username)
	if err != nil {
		log.Println("error in fetch contact list of username: ", username, err)
		res.Message = "unable to fetch contact list. please try again later."
		return res
	}

	res.Status = true
	res.Data = contactList
	res.Total = len(contactList)
	return res
}

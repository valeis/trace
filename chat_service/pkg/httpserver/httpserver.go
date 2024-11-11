package httpserver

import (
	"chat_app/pkg/redisrepo"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"net/http"
)

func StartHTTPServer() {
	redisClient := redisrepo.InitialiseRedis()
	defer redisClient.Close()

	redisrepo.CreateFetchChatBetweenIndex()

	r := mux.NewRouter()
	r.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Simple Server")
	}).Methods(http.MethodGet)

	r.HandleFunc("/register", registerHandler).Methods(http.MethodPost)
	r.HandleFunc("/verify-contact", verifyContactHandler).Methods(http.MethodPost)
	r.HandleFunc("/contact-list", contactListHandler).Methods(http.MethodPost)
	r.HandleFunc("/chat-history", chatHistoryHandler).Methods(http.MethodPost)

	handler := cors.Default().Handler(r)
	http.ListenAndServe(":8080", handler)
}

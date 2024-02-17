package spotistats

import (
	"github.com/erfgypO/spotistats/lib/server"
	"log"
	"net/http"
	"os"
)

func runApi() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /auth-url", server.UseAuth(server.HandleGetAuthUrl))
	mux.HandleFunc("GET /auth-redirect", server.HandleAuthRedirect)
	mux.HandleFunc("GET /ping", server.UseAuth(server.HandlePing))
	mux.HandleFunc("POST /sign-up", server.HandleSignUp)
	mux.HandleFunc("POST /sign-in", server.HandleSignIn)
	mux.HandleFunc("GET /user/me", server.UseAuth(server.HandleGetMe))
	err := http.ListenAndServe(os.Getenv("API_URL"), mux)
	if err != nil {
		log.Panicln(err)
	}
}

package oauth

import (
	"fmt"
	"log"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/drive/v3"
)

var (
	oauthConfig *oauth2.Config
)

const (
	oauthGoogleURLAPI = "https://www.googleapis.com/oauth2/v2/userinfo"
)

func init() {
	oauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/callback",
		ClientID:     "481585171843-qlhj01c6hhqtjaoqpbrnj8satngah5vk.apps.googleusercontent.com",
		ClientSecret: "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",
		Scopes:       []string{"https://www.googleapis.com/auth/drive"},
		Endpoint:     google.Endpoint,
	}
}

func main() {
	http.HandleFunc("/", handleMain)
	http.HandleFunc("/login", handleGoogleLogin)
	http.HandleFunc("/callback", handleGoogleCallback)
	fmt.Println(http.ListenAndServe(":8080", nil))
}

func handleMain(w http.ResponseWriter, r *http.Request) {
	var htmlIndex = `<html><body><a href="/login">Google Log In</a></body></html>`
	fmt.Fprintf(w, htmlIndex)
}

func handleGoogleLogin(w http.ResponseWriter, r *http.Request) {
	url := oauthConfig.AuthCodeURL("")
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

func handleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	token, err := oauthConfig.Exchange(r.Context(), code)
	if err != nil {
		log.Fatal(err)
	}

	client := oauthConfig.Client(r.Context(), token)

	srv, err := drive.New(client)
	if err != nil {
		log.Fatalf("erro ao criar serviço do google drive: %v", err)
	}

	file := drive.File{Name: "NomeDoArquivo"}
	_, err = srv.Files.Create(&file).Media(r.Body).Do()
	if err != nil {
		log.Fatalf("erro: %v", err)
	}

	var htmlSuccess = `<html><body>Authentication Successful!</body></html>`
	fmt.Fprintf(w, htmlSuccess)
}

package oauth

import (
	"fmt"
	"net/http"
	"sync"

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

var UserTokens map[int64]string
var mutex sync.Mutex

func init() {
	UserTokens = make(map[int64]string)

	oauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/callback",
		ClientID:     "481585171843-qlhj01c6hhqtjaoqpbrnj8satngah5vk.apps.googleusercontent.com",
		ClientSecret: "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",
		Scopes:       []string{"https://www.googleapis.com/auth/drive"},
		Endpoint:     google.Endpoint,
	}
}

func StartServer() {
	http.HandleFunc("/", handleMain)
	http.HandleFunc("/login", handleGoogleLogin)
	http.HandleFunc("/callback", HandleGoogleCallback)
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

func HandleGoogleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.FormValue("code")
	token, err := oauthConfig.Exchange(r.Context(), code)
	if err != nil {
		http.Error(w, "erro ao trocar codigo por token", http.StatusInternalServerError)
		return
	}

	client := oauthConfig.Client(r.Context(), token)

	srv, err := drive.NewService(client)
	if err != nil {
		http.Error(w, "erro ao criar serviço do google drive", http.StatusInternalServerError)
		return
	}

	file := drive.File{Name: "NomeDoArquivo"}

	_, err = srv.Files.Create(&file).Media(r.Body).Do()
	if err != nil {
		http.Error(w, "erro ao fazer upload do arquivo", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, "autenticaçao e upload bem sucedidos!")
}

func SetToken(userID int64, token string) {
	mutex.Lock()
	defer mutex.Unlock()
	UserTokens[userID] = token
}

func IsLoggedIn(userID int64) bool {
	mutex.Lock()
	defer mutex.Unlock()
	_, ok := UserTokens[userID]
	return ok
}

func GetToken(userID int64) (string, error) {
	mutex.Lock()
	defer mutex.Unlock()

	token, ok := UserTokens[userID]
	if !ok {
		return "", fmt.Errorf("Token do usuario não encontrada, ID: %d", userID)
	}
	return token, nil
}

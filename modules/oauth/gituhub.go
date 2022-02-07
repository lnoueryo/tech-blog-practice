package oauth

import (
	"encoding/json"
	"fmt"
	"helloworld/config"
	"io/ioutil"
	"log"
	"net/http"
)


type Token struct {
	AccessToken string `json:"access_token"`
}

type GithubOAuthInfo struct {
	Token
	Name 	  string `json:"name"`
	AvatarURL string `json:"avatar_url"`
	Email 	  string `json:"email"`
}

func NewGithubOAuthInfo() *GithubOAuthInfo{
    return &GithubOAuthInfo{}
}

func GithubOAuth(w http.ResponseWriter, r *http.Request) (*GithubOAuthInfo, error) {
	// トークンをリクエストする
	user := NewGithubOAuthInfo()
	err := r.ParseForm()
	if err != nil {
		// errorlog.Print("could not parse query: %v", err)
		return user, err
	}
	code := r.FormValue("code")
	err = user.tokenRequest(code)
	if err != nil {
		return user, err
		// errorlog.Println(err)
	}
	err = user.apiRequest()
	if err != nil {
		return user, err
		// errorlog.Print(err)
	}
	return user, err
}

func (u *GithubOAuthInfo)tokenRequest(code string) (err error) {
	clientID := config.ApiKey.GitHubClientId
	clientSecret := config.ApiKey.GitHubSecretId

	// Next, lets for the HTTP request to call the github oauth enpoint
	// to get our access token
	reqURL := fmt.Sprintf("https://github.com/login/oauth/access_token?client_id=%s&client_secret=%s&code=%s", clientID, clientSecret, code)
	req, err := http.NewRequest(http.MethodPost, reqURL, nil)
	if err != nil {
		return err
		// errorlog.Print(err)
	}
	client := &http.Client{}
	header := http.Header{}
	header.Set("Accept", "application/json")
	header.Set("Content-Type", "application/json")
	req.Header = header
	resp, err := client.Do(req)
	if err != nil {
		return err
		// errorlog.Printf("request err: %s", err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
		// errorlog.Printf("request err: %s", err)
	}
	if err = json.Unmarshal(body, &u); err != nil {
		fmt.Println(err)
	}
	return
}

func (u *GithubOAuthInfo)apiRequest() (err error) {
	userAPI := "https://api.github.com/user"

	req, err := http.NewRequest("GET", userAPI, nil)
	if err != nil {
		return err
		// errorlog.Print(err)
	}
	// 取得したアクセストークンをHeaderにセットしてリソースサーバにリクエストを送る
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", u.AccessToken))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		log.Printf("http status code is %d, err: %s", resp.StatusCode, err)
		return err
		// errorlog.Print(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
		// errorlog.Print(err)
	}
	if err = json.Unmarshal(body, &u); err != nil {
		return err
		// errorlog.Print(err)
	}
	return
}
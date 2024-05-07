package github

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/fabiokaelin/f-oauth/config"
)

type (
	GitHubOauthToken struct {
		Access_token string
	}
	GitHubUserResult struct {
		Name  string
		Photo string
		Email string
	}
)

func getGitHubOauthToken(code string) (*GitHubOauthToken, error) {
	const rootURl = "https://github.com/login/oauth/access_token"

	values := url.Values{}
	values.Add("code", code)
	values.Add("client_id", config.GitHubClientID)
	values.Add("client_secret", config.GitHubClientSecret)

	query := values.Encode()

	queryString := fmt.Sprintf("%s?%s", rootURl, bytes.NewBufferString(query))
	req, err := http.NewRequest("POST", queryString, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	client := http.Client{
		Timeout: time.Second * 30,
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("could not retrieve token")
	}

	var resBody bytes.Buffer
	_, err = io.Copy(&resBody, res.Body)
	if err != nil {
		return nil, err
	}

	parsedQuery, err := url.ParseQuery(resBody.String())
	if err != nil {
		return nil, err
	}

	tokenBody := &GitHubOauthToken{
		Access_token: parsedQuery["access_token"][0],
	}

	return tokenBody, nil
}

func getGitHubUser(access_token string) (*GitHubUserResult, error) {
	rootUrl := "https://api.github.com/user"

	req, err := http.NewRequest("GET", rootUrl, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access_token))

	client := http.Client{
		Timeout: time.Second * 30,
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, errors.New("could not retrieve user")
	}

	var resBody bytes.Buffer
	_, err = io.Copy(&resBody, res.Body)
	if err != nil {
		return nil, err
	}

	var GitHubUserRes map[string]interface{}

	if err := json.Unmarshal(resBody.Bytes(), &GitHubUserRes); err != nil {
		return nil, err
	}

	spew.Dump(GitHubUserRes)
	email := ""

	email, ok := GitHubUserRes["email"].(string)

	if !ok {
		fmt.Println("email not found")
		rootUrl := "https://api.github.com/user/emails"

		req, err := http.NewRequest("GET", rootUrl, nil)
		if err != nil {
			return nil, err
		}

		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access_token))

		client := http.Client{
			Timeout: time.Second * 30,
		}

		res, err := client.Do(req)
		if err != nil {
			return nil, err
		}

		if res.StatusCode != http.StatusOK {
			return nil, errors.New("could not retrieve user")
		}

		var resBody bytes.Buffer
		_, err = io.Copy(&resBody, res.Body)
		if err != nil {
			return nil, err
		}
		// [
		// 	{
		// 	  "email": "supermae206@gmail.com",
		// 	  "primary": true,
		// 	  "verified": true,
		// 	  "visibility": "private"
		// 	},
		// 	{
		// 	  "email": "90179796+FabioKaelin@users.noreply.github.com",
		// 	  "primary": false,
		// 	  "verified": true,
		// 	  "visibility": null
		// 	}
		//   ]

		var GitHubUserEmailRes []map[string]interface{}

		if err := json.Unmarshal(resBody.Bytes(), &GitHubUserEmailRes); err != nil {
			return nil, err
		}

		spew.Dump(GitHubUserEmailRes)

		for _, emailRes := range GitHubUserEmailRes {
			spew.Dump(emailRes)
			if emailRes["primary"].(bool) {
				email, ok = emailRes["email"].(string)
				fmt.Println("email found", email)
				if !ok {
					return nil, errors.New("email not found")
				}
				break
			}
		}
	}

	fmt.Println("email", email)

	userBody := &GitHubUserResult{
		Email: email,
		Name:  GitHubUserRes["login"].(string),
		Photo: GitHubUserRes["avatar_url"].(string),
	}

	return userBody, nil
}

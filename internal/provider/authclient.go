package provider

import "net/http"

type AuthClient struct {
	apikey string
}

func NewAuthClient(apikey string) *AuthClient {
	return &AuthClient{
		apikey: apikey,
	}
}

func (a *AuthClient) Do(r *http.Request) (*http.Response, error) {
	r.Header.Set("Authorization", "Api-Key "+a.apikey)

	return http.DefaultClient.Do(r)
}

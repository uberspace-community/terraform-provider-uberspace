package provider

import "net/http"

type authClient struct {
	apikey string
}

func newAuthClient(apikey string) *authClient {
	return &authClient{
		apikey: apikey,
	}
}

func (a *authClient) Do(r *http.Request) (*http.Response, error) {
	r.Header.Set("Authorization", "Api-Key "+a.apikey)

	return http.DefaultClient.Do(r)
}

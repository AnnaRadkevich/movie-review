package client

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/cloudmachinery/movie-reviews/contracts"
	"github.com/go-resty/resty/v2"
)

type Client struct {
	client  *resty.Client
	baseURL string
}

func New(url string) *Client {
	hc := &http.Client{}
	rc := resty.NewWithClient(hc)
	rc.OnAfterResponse(func(client *resty.Client, response *resty.Response) error {
		if response.IsError() {
			herr := contracts.HttpError{}
			_ = json.Unmarshal(response.Body(), &herr)

			return &Error{Code: response.StatusCode(), Message: herr.Message}
		}
		return nil
	})
	return &Client{
		client:  rc,
		baseURL: url,
	}
}

func (c *Client) path(f string, args ...any) string {
	return fmt.Sprintf(c.baseURL+f, args...)
}

func logRequest(client *resty.Client, request *resty.Request) error {
	log.Printf("Request URL: %s", request.URL)
	log.Printf("Request Method: %s", request.Method)
	log.Printf("Request Body: %v", request.Body)
	return nil
}

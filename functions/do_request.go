package functions

import (
	"encoding/json"
	"io"
	"kmuttBot/types/response"
	"kmuttBot/utils/network"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func DoRequest(c *fiber.Ctx, method string, url string, body io.Reader, modifier func(r *http.Request), data interface{}) *response.ErrorInstance {
	// Construct request
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return &response.ErrorInstance{Message: "Cannot construct request", Err: err}
	}

	// Modify request
	if modifier != nil {
		modifier(req)
	}

	// Construct client
	cli := network.NewClient()
	resp, err := cli.Do(req)
	if err != nil {
		return &response.ErrorInstance{Message: "Cannot send request", Err: err}
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	// Read response
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return &response.ErrorInstance{Message: "Cannot read response", Err: err}
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return &response.ErrorInstance{Message: "Status code is not 200"}
	}

	// Parse response
	if err := json.Unmarshal(respBody, data); err != nil {
		return &response.ErrorInstance{Message: "Cannot parse response", Err: err}
	}

	return nil
}

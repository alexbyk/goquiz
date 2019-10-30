package httpapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/alexbyk/goquiz/common/model"
)

// HTTPPublisher implements integrator.Publisher for http endpoint
type HTTPPublisher struct {
	endpoint string
}

// NewHTTPPublisher returns a new publisher
func NewHTTPPublisher(endpoint string) *HTTPPublisher { return &HTTPPublisher{endpoint} }

// Publish posts data to endpoint as a JSON. If response is 200, return nil
func (p *HTTPPublisher) Publish(c *model.Customer) error {
	data, err := json.Marshal(c)
	if err != nil {
		return err
	}

	resp, err := http.Post(p.endpoint, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Recieved status code: %v", resp.StatusCode)
	}

	_, err = ioutil.ReadAll(resp.Body)
	return nil
}

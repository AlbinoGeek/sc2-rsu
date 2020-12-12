package sc2replaystats

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"runtime"
	"time"

	jsoniter "github.com/json-iterator/go"
)

var (
	// Hostname represents the root domain sc2replaystats is hosted at
	Hostname = "sc2replaystats.com"

	// Protocol represents the HTTP protocol we use when communicating with sc2replaystats
	Protocol = "https"

	// APIRoot represents the base URL for requests to the sc2replaystats JSON-ish API
	APIRoot = fmt.Sprintf("%s://api.%s", Protocol, Hostname)

	// WebRoot represents the base URL for requests to the sc2replaystats Website
	WebRoot = fmt.Sprintf("%s://%s", Protocol, Hostname)

	// ClientIdentifier represents the "upload_method" shown to sc2replaystats
	ClientIdentifier = fmt.Sprintf("sc2-rsu-%s", runtime.GOOS)
)

// Client allows you to communicate with the sc2ReplayStats API
type Client struct {
	apikey string
	client *http.Client
}

// New returns an sc2ReplayStats API Client
func New(apikey string) *Client {
	return &Client{
		apikey: apikey,
		client: &http.Client{
			Timeout: time.Second * 3,
		},
	}
}

func (client *Client) doRequest(method, slug, contentType string, data io.Reader) (resp *http.Response, err error) {
	req, err := http.NewRequest(method, fmt.Sprintf("%s/%s", APIRoot, slug), data)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare request: %v", err)
	}

	req.Header.Set("Authorization", client.apikey)
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	resp, err = client.client.Do(req)
	if err != nil {
		err = fmt.Errorf("failed to send request: %v", err)
	}

	return
}

func (client *Client) requestBytes(method, slug, contentType string, data io.Reader) (result []byte, err error) {
	resp, err := client.doRequest(method, slug, contentType, data)
	defer resp.Body.Close()

	result, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("sc2replaystats API returned error: %v", resp.Status)
	}

	return
}

func (client *Client) requestMap(method, slug, contentType string, data io.Reader) (result map[string]string, err error) {
	resp, err := client.doRequest(method, slug, contentType, data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	result = make(map[string]string)
	if err = jsoniter.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decode response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("sc2replaystats API returned error: %v", resp.Status)
	}

	return
}

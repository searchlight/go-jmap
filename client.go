package jmap

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/oauth2"
)

// A JMAP Client
type Client struct {
	// The HttpClient.Client to use for requests. The HttpClient.Client should handle
	// authentication.
	HttpClient *http.Client

	// The JMAP Session Resource Endpoint. If the client detects the Session
	// object needs refetching, it will automatically do so
	SessionEndpoint string

	// the JMAP Session object
	Session *Session
}

// Set the HttpClient to a client which authenticates using the provided
// username and password
func (c *Client) WithBasicAuth(username string, password string) *Client {
	ctx := context.Background()
	auth := username + ":" + password
	t := &oauth2.Token{
		AccessToken: base64.StdEncoding.EncodeToString([]byte(auth)),
		TokenType:   "basic",
	}
	cfg := &oauth2.Config{}
	c.HttpClient = oauth2.NewClient(ctx, cfg.TokenSource(ctx, t))
	return c
}

// Set the HttpClient to a client which authenticates using the provided Access
// Token
func (c *Client) WithAccessToken(token string) *Client {
	ctx := context.Background()
	t := &oauth2.Token{
		AccessToken: token,
		TokenType:   "bearer",
	}
	cfg := &oauth2.Config{}
	c.HttpClient = oauth2.NewClient(ctx, cfg.TokenSource(ctx, t))
	return c
}

// Authenticate authenticates the client and retrieves the Session object.
// Authenticate will be called automatically when Do is called if the Session
// object hasn't already been initialized. Call Authenticate before any requests
// if you need to access information from the Session object prior to the first
// request
func (c *Client) Authenticate() (*Session, error) {
	if c.SessionEndpoint == "" {
		return nil, fmt.Errorf("no session url is set")
	}

	req, err := http.NewRequest("GET", c.SessionEndpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("couldn't authenticate")
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	s := &Session{}
	err = json.Unmarshal(data, s)
	if err != nil {
		return nil, err
	}
	c.Session = s
	return s, nil
}

// Do performs a JMAP request and returns the response
func (c *Client) Do(req *Request) (*Response, error) {
	if c.Session == nil {
		_, err := c.Authenticate()
		if err != nil {
			return nil, err
		}
	}
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	httpReq, err := http.NewRequest("POST", c.Session.APIURL, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	httpResp, err := c.HttpClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	if httpResp.StatusCode != 200 {
		return nil, decodeHttpError(httpResp)
	}

	data, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}
	resp := &Response{}
	err = json.Unmarshal(data, resp)
	if err != nil {
		return nil, fmt.Errorf("error? %v", err)
	}

	return resp, nil
}

// Upload sends binary data to the server and returns blob ID and some
// associated meta-data.
//
// There are some caveats to keep in mind:
// - Server may return the same blob ID for multiple uploads of the same blob.
// - Blob ID may become invalid after some time if it is unused.
// - Blob ID is usable only by the uploader until it is used, even for shared accounts.
func (c *Client) Upload(accountID string, blob io.Reader) (*BlobInfo, error) {
	if c.SessionEndpoint == "" {
		return nil, fmt.Errorf("jmap/client: SessionEndpoint is empty")
	}
	if c.Session == nil {
		_, err := c.Authenticate()
		if err != nil {
			return nil, err
		}
	}

	url := strings.Replace(c.Session.UploadURL, "{accountId}", accountID, -1)
	req, err := http.NewRequest("POST", url, blob)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, decodeHttpError(resp)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	info := &BlobInfo{}
	err = json.Unmarshal(data, info)
	if err != nil {
		return nil, err
	}

	return info, nil
}

// Download downloads binary data by its Blob ID from the server.
func (c *Client) Download(accountID string, blobID string) (io.ReadCloser, error) {
	if c.SessionEndpoint == "" {
		return nil, fmt.Errorf("jmap/client: SessionEndpoint is empty")
	}
	if c.Session == nil {
		_, err := c.Authenticate()
		if err != nil {
			return nil, err
		}
	}

	urlRepl := strings.NewReplacer(
		"{accountId}", accountID,
		"{blobId}", blobID,
		"{type}", "application/octet-stream", // TODO: are any other values necessary?
		"{name}", "filename",
	)
	tgtUrl := urlRepl.Replace(c.Session.DownloadURL)
	req, err := http.NewRequest("GET", tgtUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode/100 != 2 {
		defer resp.Body.Close()
		return nil, decodeHttpError(resp)
	}

	return resp.Body, nil
}

func decodeHttpError(resp *http.Response) error {
	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json" {
		return fmt.Errorf("HTTP %d %s", resp.StatusCode, resp.Status)
	}

	var requestErr error
	if err := json.NewDecoder(resp.Body).Decode(&requestErr); err != nil {
		return fmt.Errorf("HTTP %d %s (failed to decode JSON body: %v)", resp.StatusCode, resp.Status, err)
	}

	return requestErr
}
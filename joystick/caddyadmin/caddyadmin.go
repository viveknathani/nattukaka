package caddyadmin

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const baseURL = "http://localhost:2019"

// Route defines a Caddy HTTP route object
type Route struct {
	ID       string     `json:"@id"`
	Match    []Match    `json:"match"`
	Handle   []Subroute `json:"handle"`
	Terminal bool       `json:"terminal"`
}

// Match defines a Caddy HTTP match object
type Match struct {
	Host []string `json:"host"`
}

// Subroute defines a Caddy HTTP subroute object
type Subroute struct {
	Handler string       `json:"handler"`
	Routes  []InnerRoute `json:"routes"`
}

// InnerRoute defines a Caddy HTTP inner route object
type InnerRoute struct {
	Handle []ReverseProxy `json:"handle"`
}

// ReverseProxy defines a Caddy HTTP reverse proxy object
type ReverseProxy struct {
	Handler   string     `json:"handler"`
	Upstreams []Upstream `json:"upstreams"`
}

// Upstream defines a Caddy HTTP upstream object
type Upstream struct {
	Dial string `json:"dial"`
}

// AddRoute appends a new route to srv0
func AddRoute(id string, hosts []string, upstream string) error {
	route := Route{
		ID:    id,
		Match: []Match{{Host: hosts}},
		Handle: []Subroute{{
			Handler: "subroute",
			Routes: []InnerRoute{{
				Handle: []ReverseProxy{{
					Handler:   "reverse_proxy",
					Upstreams: []Upstream{{Dial: upstream}},
				}},
			}},
		}},
		Terminal: true,
	}

	body, _ := json.Marshal(route)
	resp, err := http.Post(
		baseURL+"/config/apps/http/servers/srv0/routes/",
		"application/json",
		bytes.NewReader(body),
	)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return handleResp(resp)
}

// UpdateRoute replaces a route by @id
func UpdateRoute(id string, hosts []string, upstream string) error {
	route := Route{
		ID:    id,
		Match: []Match{{Host: hosts}},
		Handle: []Subroute{{
			Handler: "subroute",
			Routes: []InnerRoute{{
				Handle: []ReverseProxy{{
					Handler:   "reverse_proxy",
					Upstreams: []Upstream{{Dial: upstream}},
				}},
			}},
		}},
		Terminal: true,
	}

	body, _ := json.Marshal(route)
	req, _ := http.NewRequest("PATCH", baseURL+"/id/"+id, bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return handleResp(resp)
}

// DeleteRoute removes a route by @id
func DeleteRoute(id string) error {
	req, _ := http.NewRequest("DELETE", baseURL+"/id/"+id, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return handleResp(resp)
}

// GetRoute returns a route by @id
func GetRoute(id string) (*Route, error) {
	req, _ := http.NewRequest("GET", baseURL+"/id/"+id, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		body, _ := io.ReadAll(resp.Body)
		var route Route
		err := json.Unmarshal(body, &route)
		if err != nil {
			return nil, err
		}
		return &route, nil
	}
	body, _ := io.ReadAll(resp.Body)
	return nil, fmt.Errorf("caddy error: %d %s", resp.StatusCode, string(body))
}

// handleResp checks for non-2xx status
func handleResp(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}
	body, _ := io.ReadAll(resp.Body)
	return fmt.Errorf("caddy error: %d %s", resp.StatusCode, string(body))
}

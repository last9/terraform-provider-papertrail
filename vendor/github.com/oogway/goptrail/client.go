package goptrail

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"strconv"

	"github.com/pkg/errors"
	"github.com/sethgrid/pester"
)

type DefaultClient struct {
	netLoc string
	token  string
}

var pClient Client

const defaultEndpoint = "https://papertrailapp.com/api/v1"

const PesterRetries = 50

func MakeRestClient() *pester.Client {
	client := pester.New()

	client.KeepLog = false
	client.Concurrency = 1
	client.MaxRetries = PesterRetries
	client.Timeout = 5 * time.Second

	// Use a LinearJitter Backoff algorithm.
	client.Backoff = pester.LinearJitterBackoff
	client.LogHook = func(e pester.ErrEntry) {
		log.Printf(client.FormatError(e))
	}

	return client
}

func NewClient(token string) Client {
	return &DefaultClient{defaultEndpoint, token}
}

func defaultParams() map[string]string {
	return make(map[string]string)
}

func (c *DefaultClient) ListUsers() ([]User, error) {
	users := []User{}
	params := defaultParams()

	err := c.execute("GET", "/users", params, &users)
	return users, err
}

func (c *DefaultClient) ListLogDestinations() ([]LogDestination, error) {
	lds := []LogDestination{}
	params := defaultParams()

	err := c.execute("GET", "/destinations", params, &lds)
	return lds, err
}

func (c *DefaultClient) ListSystems() ([]OutputSystem, error) {
	out := []OutputSystem{}
	params := defaultParams()

	err := c.execute("GET", "/systems", params, &out)

	return out, err
}

func (c *DefaultClient) GetSystem(id string) (*OutputSystem, error) {
	out := OutputSystem{}
	params := defaultParams()
	path := fmt.Sprintf("/systems/%s", id)

	if err := c.execute("GET", path, params, &out); err != nil {
		return nil, err
	}

	return &out, nil
}

func (c *DefaultClient) RegisterSystem(s InputSystem) (OutputSystem, error) {
	params := parseSystemParams(s)
	out := OutputSystem{}

	err := c.execute("POST", "/systems", params, &out)
	return out, err
}

func (c *DefaultClient) UpdateSystem(s InputSystem) error {
	params := parseSystemParams(s)
	var out interface{}

	path := fmt.Sprintf("/systems/%d", s.ID)
	return c.execute("PUT", path, params, &out)
}

func (c *DefaultClient) UnregisterSystem(id string) error {
	params := defaultParams()
	var out interface{}
	path := fmt.Sprintf("/systems/%s", id)
	return c.execute("DELETE", path, params, &out)
}

func (c *DefaultClient) CreateGroup(g Group) (Group, error) {
	params := parseGroupParams(g)
	out := Group{}

	err := c.execute("POST", "/groups", params, &out)
	return out, err
}

func (c *DefaultClient) GetGroup(id string) (Group, error) {
	out := Group{}
	params := defaultParams()
	path := fmt.Sprintf("/groups/%s", id)

	if err := c.execute("GET", path, params, &out); err != nil {
		return out, err
	}

	return out, nil
}

func (c *DefaultClient) ListGroups() ([]Group, error) {
	out := []Group{}
	params := defaultParams()

	err := c.execute("GET", "/groups", params, &out)

	return out, err
}

func (c *DefaultClient) UpdateGroup(g Group) error {
	params := parseGroupParams(g)
	var out interface{}
	path := fmt.Sprintf("/groups/%s", g.ID)

	return c.execute("PUT", path, params, &out)
}

func (c *DefaultClient) DeleteGroup(id string) error {
	var out interface{}
	params := defaultParams()
	path := fmt.Sprintf("/groups/%s", id)

	return c.execute("DELETE", path, params, &out)
}

func (c *DefaultClient) AddSystemToGroup(sID, gID string) error {
	var out interface{}
	params := defaultParams()
	params["group_id"] = gID
	path := fmt.Sprintf("/systems/%s/join", sID)

	return c.execute("POST", path, params, &out)
}

func (c *DefaultClient) RemoveSystemFromGroup(sID, gID string) error {
	var out interface{}
	params := defaultParams()
	params["group_id"] = gID
	path := fmt.Sprintf("/systems/%s/leave", sID)

	return c.execute("POST", path, params, &out)
}

func (c *DefaultClient) execute(method, path string, reqParams map[string]string, respBody interface{}) error {
	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", c.netLoc, path), nil)
	if err != nil {
		return errors.Wrapf(err, "[%v] Error creating reqeust", path)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Papertrail-Token", c.token)

	q := req.URL.Query()
	q.Add("format", "json")

	for k, v := range reqParams {
		q.Add(k, v)
	}
	req.URL.RawQuery = q.Encode()

	client := MakeRestClient()

	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrapf(err, "[%v] Error requesting", path)
	}

	if resp.Body != nil {
		defer resp.Body.Close()
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrapf(err, "[%v] Cannot read body", path)
	}

	errMessage := struct {
		Message string `json:"message"`
	}{}

	if resp.StatusCode >= 400 {
		json.Unmarshal(body, &errMessage)
		return errors.New(errMessage.Message)
	}

	if err := json.Unmarshal(body, respBody); err != nil {
		return errors.Wrapf(err, "[%v] Cannot parse %v as JSON", path, string(body))
	}

	return nil
}

func parseSystemParams(s InputSystem) map[string]string {
	params := defaultParams()
	params["system[name]"] = s.Name

	if s.IpAddress != "" {
		params["system[ip_address]"] = s.IpAddress
	}

	if s.Hostname != "" {
		params["system[hostname]"] = s.Hostname
	}

	if s.DestinationID > 0 {
		params["destination_id"] = strconv.Itoa(s.DestinationID)
	}

	if s.DestinationPort > 0 {
		params["destination_port"] = strconv.Itoa(s.DestinationPort)
	}

	if s.Description != "" {
		params["description"] = s.Description
	}
	return params
}

func parseGroupParams(g Group) map[string]string {
	params := defaultParams()
	params["group[name]"] = g.Name

	if g.SystemWildcard != "" {
		params["group[system_wildcard"] = g.SystemWildcard
	}
	return params
}

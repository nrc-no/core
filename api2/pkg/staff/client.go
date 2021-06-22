package staff

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/nrc-no/core-kafka/pkg/auth"
	"io/ioutil"
	"net/http"
)

type Client struct {
	basePath string
}

func NewClient(basePath string) *Client {
	return &Client{
		basePath: basePath,
	}
}

func (c *Client) List(ctx context.Context, listOptions ListOptions) (*StaffList, error) {
	req, err := http.NewRequest("GET", c.basePath+"/apis/v1/staff", nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	auth.Forward(ctx, req)
	req.Header.Set("Accept", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var list StaffList
	if err := json.Unmarshal(bodyBytes, &list); err != nil {
		return nil, err
	}
	return &list, nil
}

func (c *Client) Get(ctx context.Context, id string) (*Staff, error) {
	req, err := http.NewRequest("GET", c.basePath+"/apis/v1/staff/"+id, nil)
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	auth.Forward(ctx, req)
	req.Header.Set("Accept", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}
	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var t Staff
	if err := json.Unmarshal(bodyBytes, &t); err != nil {
		return nil, err
	}
	return &t, nil
}

func (c *Client) Update(ctx context.Context, team *Staff) (*Staff, error) {
	bodyBytes, err := json.Marshal(team)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("PUT", c.basePath+"/apis/v1/staff/"+team.ID, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	auth.Forward(ctx, req)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}
	responseBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var t Staff
	if err := json.Unmarshal(responseBytes, &t); err != nil {
		return nil, err
	}
	return &t, nil
}

func (c *Client) Create(ctx context.Context, team *Staff) (*Staff, error) {
	bodyBytes, err := json.Marshal(team)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", c.basePath+"/apis/v1/staff", bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, err
	}
	req = req.WithContext(ctx)
	auth.Forward(ctx, req)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}
	responseBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	var t Staff
	if err := json.Unmarshal(responseBytes, &t); err != nil {
		return nil, err
	}
	return &t, nil
}

func (c *Client) Delete(ctx context.Context, id string) error {
	req, err := http.NewRequest("DELETE", c.basePath+"/apis/v1/staff/"+id, nil)
	if err != nil {
		return err
	}
	req = req.WithContext(ctx)
	auth.Forward(ctx, req)
	req.Header.Set("Accept", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", res.StatusCode)
	}
	return nil
}

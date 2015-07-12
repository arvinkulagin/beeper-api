package beeperapi

import (
	"net/http"
	"net/url"
	"bytes"
	"errors"
	"encoding/json"
	"io/ioutil"
)

const (
	scheme = "http"
	topicPrefix = "topic"
	pingPrefix = "ping"
)

type Client struct {
	rurl url.URL
}

func NewClient(host string) (Client, error) {
	c := Client{
		rurl: url.URL{
			Scheme: scheme,
			Host: host,
		},
	}
	if err := c.Ping(); err != nil {
		return c, err
	}
	return c, nil
}

func (c Client) Add(topic string) error {
	out := bytes.NewReader([]byte(topic))
	c.rurl.Path = c.rurl.Path + topicPrefix
	request, err := http.NewRequest("POST", c.rurl.String(), out)
	if err != nil {
		return errors.New("Bad request")
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return errors.New("Сan not connect to " + c.rurl.Host)
	}
	if response.StatusCode != 200 {
		return errors.New("Topic already exists")
	}
	return nil
}

func (c Client) Del(topic string) error {
	c.rurl.Path = c.rurl.Path + topicPrefix + "/" + topic
	req, err := http.NewRequest("DELETE", c.rurl.String(), nil)
	if err != nil {
		return errors.New("Can not make request")
	}
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.New("Сan not connect to " + c.rurl.Host)
	}
	if r.StatusCode != 200 {
		return errors.New("Topic does not exist")
	}
	return nil
}

func (c Client) Pub(topic string, data string) error {
	out := bytes.NewReader([]byte(data))
	c.rurl.Path = c.rurl.Path + topicPrefix + "/" + topic
	request, err := http.NewRequest("POST", c.rurl.String(), out)
	if err != nil {
		return errors.New("Bad request")
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return errors.New("Сan not connect to " + c.rurl.Host)
	}
	if response.StatusCode != 200 {
		return errors.New("Topic does not exist")
	}
	return nil
}

func (c Client) List() ([]string, error) {
	list := []string{}
	c.rurl.Path = c.rurl.Path + topicPrefix
	req, err := http.NewRequest("GET", c.rurl.String(), nil)
	if err != nil {
		return list, errors.New("Can not make request")
	}
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return list, errors.New("Сan not connect to " + c.rurl.Host)
	}
	if r.StatusCode == 500 {
		return list, errors.New("Server error")
	}
	if r.StatusCode != 200 {
		return list, errors.New("Topic does not exist")
	}
	body, err := ioutil.ReadAll(r.Body) // Ограничить количество байт из ReadAll с помощью io.LimitReader
	defer r.Body.Close()
	if err != nil {
		return list, errors.New("Can not read response body")
	}
	json.Unmarshal(body, &list)
	return list, nil
}

func (c Client) Ping() error {
	c.rurl.Path = c.rurl.Path + pingPrefix
	request, err := http.NewRequest("GET", c.rurl.String(), nil)
	if err != nil {
		return errors.New("Bad request")
	}
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return errors.New("Сan not connect to " + c.rurl.Host)
	}
	if response.StatusCode != 200 {
		return errors.New("Service is unavailable")
	}
	return nil
}
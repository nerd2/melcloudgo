package melcloudgo

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-resty/resty/v2"
)

type Options struct {
	Username   string
	Password   string
	Url        string // defaults to DEFAULT_URL
	HttpClient *http.Client
}

const (
	DEFAULT_URL           = "https://app.melcloud.com/Mitsubishi.Wifi.Client/"
	USER_LOGIN_ENDPOINT   = "Login/ClientLogin2"
	LIST_DEVICES_ENDPOINT = "User/ListDevices"
)

type MelCloud interface {
	Login() (*LoginResponse, error)
	ListDevices() ([]ListDevicesResponse, error)
}

func NewMelCloud(options *Options) MelCloud {
	client := resty.New()
	if options == nil {
		options = &Options{}
	}
	if options.HttpClient != nil {
		client = resty.NewWithClient(options.HttpClient)
	}
	if options.Url == "" {
		options.Url = DEFAULT_URL
	}

	return &melCloud{
		client:  client,
		options: options,
	}
}

type melCloud struct {
	options *Options
	token   string
	client  *resty.Client
}

func (n *melCloud) ListDevices() ([]ListDevicesResponse, error) {
	var data []ListDevicesResponse
	err := n.jsonRequest(LIST_DEVICES_ENDPOINT, nil, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (n *melCloud) Login() (*LoginResponse, error) {
	var loginResponse LoginResponse
	err := n.jsonRequest(USER_LOGIN_ENDPOINT, LoginRequest{
		AppVersion:      "1.2.3.4",
		CaptchaResponse: nil,
		Email:           n.options.Username,
		Language:        0,
		Password:        n.options.Password,
		Persist:         true,
	}, &loginResponse)
	if err != nil {
		return nil, err
	}
	n.token = loginResponse.LoginData.ContextKey

	return &loginResponse, nil
}

func (n *melCloud) jsonRequest(endpoint string, request interface{}, response interface{}) error {

	req := n.client.R().SetResult(response).SetHeader("x-mitscontextkey", n.token)

	var resp *resty.Response
	var err error
	if request != nil {
		resp, err = req.SetHeader("Content-Type", "application/json").SetBody(request).Post(n.options.Url + endpoint)
	} else {
		resp, err = req.Get(n.options.Url + endpoint)
	}
	if err != nil {
		return err
	}
	if resp.StatusCode() != http.StatusOK {
		return fmt.Errorf("Unexpected status code in request to %s: %s", endpoint, resp.StatusCode())
	}
	err = json.Unmarshal(resp.Body(), response)
	if err != nil {
		return fmt.Errorf("JSON unmarshal error: %s", err.Error())
	}

	return nil
}

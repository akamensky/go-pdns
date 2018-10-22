package pdns

import (
	"fmt"
	"github.com/dghubble/sling"
	"net/url"
)

const (
	clientBaseURL = "/api/v1/servers/%s%s"
)

type Client struct {
	url    string
	server string
	apiKey string
}

type ClientOptions struct {
	URL    string
	Server string
	ApiKey string
}

type ServerInfo struct {
	Id         string `json:"id"`
	Type       string `json:"type"`
	Url        string `json:"url"`
	Version    string `json:"version"`
	DaemonType string `json:"daemon_type"`
}

func NewClient(options *ClientOptions) (*Client, error) {
	// Make new client
	client := new(Client)

	// Try to parse options url
	parsedUrl, err := url.Parse(options.URL)
	if err != nil {
		return nil, err
	}
	saneUrl := new(url.URL)
	saneUrl.Scheme = parsedUrl.Scheme
	saneUrl.Host = parsedUrl.Host
	saneUrl.Path = "/"
	if saneUrl.Scheme == "" {
		saneUrl.Scheme = "http"
	}
	client.url = saneUrl.String()

	// Set server ID and API key
	client.server = options.Server
	client.apiKey = options.ApiKey

	// Try new connection right away
	_, err = client.GetServerInfo()
	if err != nil {
		return nil, err
	}

	return client, nil
}

func (c *Client) GetServerInfo() (*ServerInfo, error) {
	result := new(ServerInfo)
	f := new(failure)

	_, err := c.getSling().Get(c.getApiPath("")).Receive(result, f)
	if err != nil {
		return nil, err
	}

	if err := f.getError(); err != nil {
		return nil, err
	}

	return result, nil
}

func (c *Client) NewZone(name string, kind zoneKind, nameservers []string) (*Zone, error) {
	zoneReq := new(Zone)
	zoneReq.Name = FQDN(name)
	zoneReq.Kind = kind
	for i, v := range nameservers {
		nameservers[i] = FQDN(v)
	}
	zoneReq.Nameservers = nameservers

	zoneResp := new(Zone)
	f := new(failure)

	_, err := c.getSling().Post(c.getApiPath("/zones")).BodyJSON(zoneReq).Receive(zoneResp, f)
	if err != nil {
		return nil, err
	}

	if err := f.getError(); err != nil {
		return nil, err
	}

	zoneResp.client = c

	return zoneResp, nil
}

func (c *Client) GetZone(name string) (*Zone, error) {
	zone := new(Zone)
	f := new(failure)

	name = FQDN(name)

	_, err := c.getSling().Get(c.getApiPath("/zones/"+name)).Receive(zone, f)
	if err != nil {
		return nil, err
	}

	if err := f.getError(); err != nil {
		return nil, err
	}

	zone.client = c

	return zone, nil
}

func (c *Client) getSling() *sling.Sling {
	sl := sling.New()
	sl.Base(c.url)
	sl.Set("X-API-Key", c.apiKey)

	return sl
}

func (c *Client) getApiPath(path string) string {
	return fmt.Sprintf(clientBaseURL, c.server, path)
}

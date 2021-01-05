package restclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/Drinkey/keyvault/commands/request"
)

const (
	apiVersion = "v1"
)

// URL type represents the URL components
type URL struct {
	Scheme, HostString, Version string
	ResourceName, Path          string
	PortNum                     int
	Query                       map[string]string
}

func (r URL) String() (url string) {
	url = fmt.Sprintf("%s://%s:%d/api/%s/%s",
		r.Scheme, r.HostString, r.PortNum, r.Version, r.ResourceName)
	if len(r.Path) != 0 {
		url = fmt.Sprintf("%s/%s", url, r.Path)
	}
	if len(r.Query) == 0 {
		return
	}
	var params []string
	for k, v := range r.Query {
		params = append(params, fmt.Sprintf("%s=%s", k, v))
	}
	url = fmt.Sprintf("%s?%s", url, strings.Join(params, "&"))
	return
}

type RESTFulClienter interface {
	Url(path string, query ...map[string]string) string
	Read(path string, query ...map[string]string) (*http.Response, error)
	Create(path string, data interface{}) (*http.Response, error)
	Update(path string, data interface{}) (*http.Response, error)
	Delete(path string, query ...map[string]string) (*http.Response, error)
}

type RESTFulClient struct {
	URL
	req      request.Requests
	Insecure bool
	Host     string
	Port     int
}

func (c *RESTFulClient) setResourceName() {
	c.ResourceName = "certificate"
}

func (c *RESTFulClient) initClient() {
	if c.req.Client == nil {
		c.req.InitClient(c.Insecure)
	}
}

func (c *RESTFulClient) generateURL(path string, query ...map[string]string) {
	c.Scheme = "https"
	c.HostString = c.Host
	c.PortNum = c.Port
	c.Version = apiVersion
	c.setResourceName()
	if len(query) > 0 {
		c.Query = query[0]
	}
	if len(path) > 0 {
		c.Path = path
	}
}

func (c RESTFulClient) Url(path string, query ...map[string]string) string {
	c.generateURL(path, query...)
	return c.String()
}

func (c RESTFulClient) Read(path string, query ...map[string]string) (*http.Response, error) {
	c.generateURL(path, query...)
	c.initClient()
	url := c.String()
	return c.req.Get(url)
}

func (c RESTFulClient) Create(path string, data interface{}) (*http.Response, error) {

	c.generateURL(path)
	c.initClient()
	url := c.String()
	log.Printf("restclient.Certificate::Create %s", url)
	cert, err := json.Marshal(data)
	if err != nil {
		log.Printf("failed to encode data to JSON: %s", err)
		panic(err)
	}
	body := bytes.NewReader(cert)
	return c.req.Post(url, body)
}

func (c RESTFulClient) Update(path string, data interface{}) (*http.Response, error) {
	c.generateURL(path)
	c.initClient()
	url := c.String()
	fmt.Println(url)
	body := bytes.NewReader([]byte(""))
	return c.req.Put(url, body)
}

func (c RESTFulClient) Delete(path string, query ...map[string]string) (*http.Response, error) {
	c.generateURL(path, query...)
	c.initClient()
	url := c.String()
	fmt.Println(url)
	return c.req.Delete(url)
}

type Certificate struct {
	RESTFulClient
}

func (c *Certificate) setResourceName() {
	c.ResourceName = "certificate"
}

type Namespace struct {
	RESTFulClient
}

func (c *Namespace) setResourceName() {
	c.ResourceName = "namespace"
}

type Secret struct {
	RESTFulClient
}

func (c *Secret) setResourceName() {
	c.ResourceName = "vault"
}

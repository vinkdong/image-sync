package docker

import (
	"github.com/docker/docker/client"
	"golang.org/x/net/context"
	"github.com/docker/docker/api/types"
	"os"
	"net/http"
	"fmt"
	"strings"
	"net/url"
	"io/ioutil"
	"encoding/json"
	"encoding/base64"
	"github.com/vinkdong/gox/log"
	"github.com/docker/docker/pkg/jsonmessage"
	"github.com/docker/docker/cli/command"
)

const ApiVersion = "v2"

type Docker struct {
	Registry     string
	ApiUrl       string `yaml:"apiUrl"`
	Username     string
	Password     string
	RegistryAuth string
}

func (docker *Docker) Login() error {
	if docker.Username != "" && docker.Password != "" {
		authConfig := types.AuthConfig{
			Username: docker.Username,
			Password: docker.Password,
		}
		encodedJSON, err := json.Marshal(authConfig)
		if err != nil {
			return err
		}
		authStr := base64.URLEncoding.EncodeToString(encodedJSON)
		docker.RegistryAuth = authStr
	}
	return nil
}

func (docker *Docker) pullImage(name, tag string) error {

	ctx := context.Background()
	c, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	if tag == "" {
		tag = "latest"
	}

	pullOpts := types.ImagePullOptions{
		RegistryAuth: docker.RegistryAuth,
	}

	refStr := fmt.Sprintf("%s/%s:%s", docker.Registry, name, tag)
	reader, err := c.ImagePull(ctx, refStr, pullOpts)
	if err != nil {
		return err
	}

	outStream := command.NewOutStream(os.Stdout)
	jsonmessage.DisplayJSONMessagesToStream(reader, outStream, nil)
	return nil
}

func (docker *Docker) tagImage(source, target string) error {
	ctx := context.Background()
	c, err := client.NewEnvClient()
	if err != nil {
		return err
	}
	return c.ImageTag(ctx, source, target)
}

func (docker *Docker) pushImage(name, tag string) error {
	ctx := context.Background()
	c, err := client.NewEnvClient()
	if err != nil {
		return err
	}
	image := fmt.Sprintf("%s/%s:%s", docker.Registry, name, tag)

	pushOpts := types.ImagePushOptions{
		RegistryAuth: docker.RegistryAuth,
	}

	reader, err := c.ImagePush(ctx, image, pushOpts)
	if err != nil {
		return err
	}
	outStream := command.NewOutStream(os.Stdout)
	jsonmessage.DisplayJSONMessagesToStream(reader, outStream, nil)
	return nil
}

type Image struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

func (docker *Docker) listTags(name string) (*Image, error) {
	apiUrl := docker.ApiUrl
	if apiUrl == "" {
		apiUrl = docker.Registry
	}
	uri := fmt.Sprintf("https://%s/%s/%s/%s", apiUrl, ApiVersion, name, "tags/list")
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return nil, err
	}
	resp, err := docker.req(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		log.Infof("registry %s/%s maybe not exists", docker.Registry, name)
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	image := &Image{}

	if err := json.Unmarshal(data, image); err != nil {
		return nil, err
	}
	return image, nil
}

type Auth struct {
	Token string `json:"token"`
}

func (docker *Docker) getToken(bc bearerChallenge) (string, error) {
	realm := bc.values["realm"]
	service := bc.values["service"]
	scope := bc.values["scope"]
	uri := fmt.Sprintf("%s?service=%s&scope=%s", realm, service, scope)
	req, err := http.NewRequest("GET", uri, nil)
	req.SetBasicAuth(docker.Username, docker.Password)
	if err != nil {
		return "", err
	}
	c := http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return "", err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	a := &Auth{}
	json.Unmarshal(data, a)
	return a.Token, nil
}

type bearerChallenge struct {
	values map[string]string
}

const (
	bearerChallengeHeader = "Www-Authenticate"
	bearer                = "Bearer"
	tenantID              = "tenantID"
)

// code from azure
// returns true if the HTTP response contains a bearer challenge
func hasBearerChallenge(resp *http.Response) bool {
	authHeader := resp.Header.Get(bearerChallengeHeader)
	if len(authHeader) == 0 || strings.Index(authHeader, bearer) < 0 {
		return false
	}
	return true
}

// code from helm
func newBearerChallenge(resp *http.Response) (bc bearerChallenge, err error) {
	challenge := strings.TrimSpace(resp.Header.Get(bearerChallengeHeader))
	trimmedChallenge := challenge[len(bearer)+1:]

	// challenge is a set of key=value pairs that are comma delimited
	pairs := strings.Split(trimmedChallenge, ",")
	if len(pairs) < 1 {
		err = fmt.Errorf("challenge '%s' contains no pairs", challenge)
		return bc, err
	}

	bc.values = make(map[string]string)
	for i := range pairs {
		trimmedPair := strings.TrimSpace(pairs[i])
		pair := strings.Split(trimmedPair, "=")
		if len(pair) == 2 {
			// remove the enclosing quotes
			key := strings.Trim(pair[0], "\"")
			value := strings.Trim(pair[1], "\"")

			switch key {
			case "authorization", "authorization_uri":
				// strip the tenant ID from the authorization URL
				asURL, err := url.Parse(value)
				if err != nil {
					return bc, err
				}
				bc.values[tenantID] = asURL.Path[1:]
			default:
				bc.values[key] = value
			}
		}
	}
	return bc, err
}

func (docker *Docker) req(req *http.Request) (*http.Response, error) {
	c := http.Client{}
	resp, err := c.Do(req)
	if err != nil {
		return resp, err
	}
	if hasBearerChallenge(resp) {
		bc, err := newBearerChallenge(resp)
		if err != nil {
			return resp, err
		}
		token, err := docker.getToken(bc)
		if err != nil {
			return resp, err
		}
		req.Header.Set("Authorization", "Bearer "+token)
		return c.Do(req)
	}
	return resp, err
}

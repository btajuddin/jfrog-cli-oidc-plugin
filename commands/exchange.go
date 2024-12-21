package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/jfrog/jfrog-cli-core/v2/plugins/components"
	"io"
	"net/http"
	"strings"
)

type httpClient interface {
	Post(url string, contentType string, body io.Reader) (*http.Response, error)
}

type flagContext interface {
	GetStringFlagValue(flagName string) string
}

func GetExchangeCommand() components.Command {
	return components.Command{
		Name:        "exchange",
		Description: "Exchanges an OIDC token for an access token.",
		Aliases:     []string{"x"},
		Arguments:   []components.Argument{},
		Flags:       getExchangeFlags(),
		EnvVars:     []components.EnvVar{},
		Action: func(c *components.Context) error {
			token, err := exchangeCommand(c, &http.Client{})
			if err == nil {
				println(token)
			}
			return err
		},
	}
}

func getExchangeFlags() []components.Flag {
	return []components.Flag{
		components.StringFlag{
			BaseFlag: components.BaseFlag{
				Name:        "token",
				Description: "OIDC token to exchange (JWT)",
			},
			Mandatory: true,
		},
		components.StringFlag{
			BaseFlag: components.BaseFlag{
				Name:        "server",
				Description: "Server to exchange the token with",
			},
			Mandatory: true,
			HelpValue: "https://<org>.jfrog.io",
		},
		components.StringFlag{
			BaseFlag: components.BaseFlag{
				Name:        "provider",
				Description: "Name of the OIDC provider configured on the server",
			},
			Mandatory: true,
		},
		components.StringFlag{
			BaseFlag: components.BaseFlag{
				Name:        "project",
				Description: "Project to authenticate for or skip to not specify a project",
			},
			Mandatory: false,
		},
	}
}

type payload struct {
	GrantType        string  `json:"grant_type"`
	SubjectTokenType string  `json:"subject_token_type"`
	SubjectToken     string  `json:"subject_token"`
	ProviderName     string  `json:"provider_name"`
	ProjectKey       *string `json:"project_key,omitempty"`
}

type exchangeResponse struct {
	AccessToken string `json:"access_token"`
}

func exchangeCommand(c flagContext, client httpClient) (*string, error) {
	data := payload{
		GrantType:        "urn:ietf:params:oauth:grant-type:token-exchange",
		SubjectTokenType: "urn:ietf:params:oauth:token-type:id_token",
		SubjectToken:     c.GetStringFlagValue("token"),
		ProviderName:     c.GetStringFlagValue("provider"),
	}

	if project := c.GetStringFlagValue("project"); project != "" {
		data.ProjectKey = &project
	}

	baseUrl, _ := strings.CutSuffix(c.GetStringFlagValue("server"), "/")
	exchangeUrl := fmt.Sprintf("%s/access/api/v1/oidc/token", baseUrl)

	encodedData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	resp, err := client.Post(exchangeUrl, "application/json", bytes.NewReader(encodedData))
	if err != nil {
		return nil, err
	}

	jsonResp := new(exchangeResponse)
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(resp.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(buf.Bytes(), jsonResp)
	if err != nil {
		return nil, err
	}

	return &jsonResp.AccessToken, nil
}

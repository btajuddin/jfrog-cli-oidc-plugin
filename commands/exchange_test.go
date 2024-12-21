package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"testing"
)

type mockedHttp struct {
	mock.Mock
}

func (m *mockedHttp) Post(url string, contentType string, body io.Reader) (*http.Response, error) {
	args := m.Called(url, contentType, body)
	return args.Get(0).(*http.Response), args.Error(1)
}

type mockedContext struct {
	mock.Mock
}

func (m *mockedContext) GetStringFlagValue(key string) string {
	args := m.Called(key)
	return args.String(0)
}

func TestSimpleExchange(t *testing.T) {
	domains := []string{"https://test-server.local", "https://test-server.local/"}

	for _, domain := range domains {
		t.Run(domain, func(t *testing.T) {
			token := "this.is.atest"
			provider := "github"
			accessToken := "this-is-an-access-token-from-a-test"
			requestBodyStruct := payload{
				GrantType:        "urn:ietf:params:oauth:grant-type:token-exchange",
				SubjectTokenType: "urn:ietf:params:oauth:token-type:id_token",
				SubjectToken:     token,
				ProviderName:     provider,
			}

			response := &http.Response{
				StatusCode: 200,
				Body:       io.NopCloser(bytes.NewBufferString(fmt.Sprintf(`{ "access_token": "%s" }`, accessToken))),
			}

			requestBodyJson, err := json.Marshal(requestBodyStruct)
			assert.Nil(t, err)

			mockHttp := new(mockedHttp)
			mockHttp.On("Post", "https://test-server.local/access/api/v1/oidc/token", "application/json", bytes.NewReader(requestBodyJson)).Return(response, nil)

			mockContext := new(mockedContext)
			mockContext.On("GetStringFlagValue", "token").Return(token)
			mockContext.On("GetStringFlagValue", "provider").Return(provider)
			mockContext.On("GetStringFlagValue", "project").Return("")
			mockContext.On("GetStringFlagValue", "server").Return(domain)

			retrievedAccessToken, err := exchangeCommand(mockContext, mockHttp)
			assert.NotNil(t, retrievedAccessToken)
			assert.Nil(t, err)
			assert.Equal(t, accessToken, *retrievedAccessToken)
		})
	}
}

func TestSimpleExchangeWithProject(t *testing.T) {
	domain := "https://test-server.local/"
	token := "this.is.atest"
	provider := "github"
	accessToken := "this-is-an-access-token-from-a-test"
	project := "some sort of project"
	requestBodyStruct := payload{
		GrantType:        "urn:ietf:params:oauth:grant-type:token-exchange",
		SubjectTokenType: "urn:ietf:params:oauth:token-type:id_token",
		SubjectToken:     token,
		ProviderName:     provider,
		ProjectKey:       &project,
	}

	response := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(bytes.NewBufferString(fmt.Sprintf(`{ "access_token": "%s" }`, accessToken))),
	}

	requestBodyJson, err := json.Marshal(requestBodyStruct)
	assert.Nil(t, err)

	mockHttp := new(mockedHttp)
	mockHttp.On("Post", "https://test-server.local/access/api/v1/oidc/token", "application/json", bytes.NewReader(requestBodyJson)).Return(response, nil)

	mockContext := new(mockedContext)
	mockContext.On("GetStringFlagValue", "token").Return(token)
	mockContext.On("GetStringFlagValue", "provider").Return(provider)
	mockContext.On("GetStringFlagValue", "project").Return(project)
	mockContext.On("GetStringFlagValue", "server").Return(domain)

	retrievedAccessToken, err := exchangeCommand(mockContext, mockHttp)
	assert.NotNil(t, retrievedAccessToken)
	assert.Nil(t, err)
	assert.Equal(t, accessToken, *retrievedAccessToken)
}

func TestFlags(t *testing.T) {
	flags := getExchangeFlags()
	if len(flags) != 4 {
		t.Error("Wrong number of flags")
	}
}

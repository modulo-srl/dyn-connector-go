package connector

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Connector object
type Connector struct {
	Debug bool

	host string

	clientID     string
	clientSecret string

	persistor   TokensPersistor
	accessToken string

	customHTTPheader map[string]string
	protocolByHost   bool
}

const (
	URLApi         = "/api/"
	URLAccessToken = "auth/token"
)

// NewConnector make a new Dyn Connector
func NewConnector(host, clientID, clientSecret string, persistor TokensPersistor) *Connector {
	conn := Connector{
		host:         host,
		clientID:     clientID,
		clientSecret: clientSecret,
	}

	conn.persistor = persistor
	conn.accessToken = conn.persistor.GetAccessToken()

	return &conn
}

// SetDebug enable or disable debug and internal params
func (c *Connector) SetDebug(enableDebug bool, customHTTPheader map[string]string, protocolByHost bool) {
	c.Debug = enableDebug
	c.customHTTPheader = customHTTPheader
	c.protocolByHost = protocolByHost
}

// Send operation request, passing data and receiving custom data types
func (c *Connector) Send(operation string, dataSend, dataReceive interface{}) *ResponseError {
	if c.Debug {
		fmt.Println("[send request \"" + operation + "\"]")
	}

	err := c.sendRequest(operation, dataSend, dataReceive)
	if err != nil {
		switch err.ID {

		case ErrorUnauthorized:
			// Authentication needed

			if c.Debug {
				fmt.Println("[authentication needed, try to get access token...]")
			}

			err := c.doAuth()
			if err != nil {
				return err
			}

			// Retry request
			if c.Debug {
				fmt.Println("[authentication done, access token \"" + c.accessToken + "\"]")
				fmt.Println("[resend request \"" + operation + "\"]")
			}

			err = c.sendRequest(operation, dataSend, dataReceive)
			if err != nil {
				return err
			}

		default:
			return err
		}
	}

	return nil
}

func (c *Connector) doAuth() *ResponseError {
	c.accessToken = ""

	req := authRequest{
		GrantType:    "client_credentials",
		ClientID:     c.clientID,
		ClientSecret: c.clientSecret,
	}

	resp := authResponse{}

	err := c.sendRequest(URLAccessToken, req, &resp)
	if err != nil {
		return err
	}

	if resp.AccessToken == "" {
		return &ResponseError{
			ID:       resp.ErrorID,
			Reason:   resp.ErrorReason,
			HTTPCode: 200,
		}
	}

	c.accessToken = resp.AccessToken
	c.persistor.SetAccessToken(c.accessToken)
	return nil
}

func (c *Connector) sendRequest(operation string, dataSend, dataReceive interface{}) *ResponseError {
	var err error

	url := ""
	if !c.protocolByHost {
		url = "https://"
	}
	url += c.host + URLApi + operation

	if c.Debug {
		fmt.Println("[make POST to", url, "]")
	}

	client := http.Client{}
	client.Timeout = time.Duration(time.Second * 60) // 60 seconds timeout

	var requestBody []byte
	requestBody, err = json.Marshal(dataSend)
	if err != nil {
		return &ResponseError{
			ID:     ErrorInternal,
			Reason: err.Error(),
		}
	}
	if c.Debug {
		fmt.Println("[raw request", string(requestBody), "]")
	}

	httpReq, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		return &ResponseError{
			ID:     ErrorInternal,
			Reason: err.Error(),
		}
	}

	httpReq.Header.Set("Content-type", "application/json")
	if len(c.customHTTPheader) > 0 {
		for key, value := range c.customHTTPheader {
			httpReq.Header.Set(key, value)
		}
	}

	if operation != URLAccessToken && c.accessToken != "" {
		httpReq.Header.Set("authorization", "Bearer "+c.accessToken)
	}

	httpResp, err := client.Do(httpReq)
	if err != nil {
		return &ResponseError{
			ID:     ErrorNetwork,
			Reason: err.Error(),
		}
	}
	defer httpResp.Body.Close()

	/*if c.Debug {
		body, _ := ioutil.ReadAll(httpResp.Body)
		fmt.Printf("[server raw response: `%s`\n", string(body))
	}*/

	if httpResp.StatusCode != 200 {
		/*if c.Debug {
			fmt.Println("[response error]")
		}*/

		rErr := ResponseError{}
		err = json.NewDecoder(httpResp.Body).Decode(&rErr)
		if err != nil {
			return &ResponseError{
				ID:       ErrorInvalidResponse,
				Reason:   err.Error(),
				HTTPCode: httpResp.StatusCode,
			}
		}

		rErr.HTTPCode = httpResp.StatusCode
		//fmt.Println(rErr)

		return &rErr
	}

	//err := json.Unmarshal(body, &dataReceive)
	err = json.NewDecoder(httpResp.Body).Decode(&dataReceive)
	if err != nil {
		return &ResponseError{
			ID:       ErrorInvalidResponse,
			Reason:   err.Error(),
			HTTPCode: httpResp.StatusCode,
		}
	}

	return nil
}

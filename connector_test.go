package connector

import (
	"fmt"
	"testing"
)

// Tokens persistor
type testPersistor struct {
	accessToken string
}

func (s *testPersistor) GetAccessToken() string {
	fmt.Println("get access token \"" + s.accessToken + "\"")
	return s.accessToken
}

func (s *testPersistor) SetAccessToken(token string) {
	s.accessToken = token
	fmt.Println("set access token \"" + s.accessToken + "\"")
}

// Test data to send and receive
type testData struct {
	KeyA    int64
	SubData testSubData
}
type testSubData struct {
	Key1 int
	Key2 bool
	Key3 string
}

func newConnector() *Connector {
	domain := "test.modulo.srl"
	clientID := "test"
	clientSecret := "test"

	persistor := testPersistor{}
	conn := NewConnector(domain, clientID, clientSecret, &persistor)

	// Enable debugging output (should be disabled in production)
	conn.SetDebug(true, nil, false)

	return conn
}

func TestConnector(t *testing.T) {
	conn := newConnector()

	dataSend := testData{
		KeyA: 1024,
		SubData: testSubData{
			Key1: 1,
			Key2: true,
			Key3: "This is a test",
		},
	}
	dataReceive := testData{}

	/*dataSend := map[string]interface{}{}
	  dataReceive := map[string]interface{}{}
	*/

	err := conn.Send("test/echo", dataSend, &dataReceive)
	if err != nil {
		t.Error(err)
		return
	}
	if dataReceive != dataSend {
		t.Error("dataSend != dataReceive")
		return
	}

	fmt.Println("Server response:", dataReceive)
}

func TestConnectorAuth(t *testing.T) {
	conn := newConnector()

	dataSend := testData{
		KeyA: 1024,
		SubData: testSubData{
			Key1: 1,
			Key2: true,
			Key3: "This is a test",
		},
	}
	dataReceive := testData{}

	/*dataSend := map[string]interface{}{}
	  dataReceive := map[string]interface{}{}
	*/

	err := conn.Send("auth/test/echo", dataSend, &dataReceive)
	if err != nil {
		t.Error(err)
		return
	}
	if dataReceive != dataSend {
		t.Error("dataSend != dataReceive")
		return
	}

	fmt.Println("Server response:", dataReceive)
}

func TestRemoteLogin(t *testing.T) {
	conn := newConnector()

	dataSend := struct {
		TaxCode string `json:"tax_code"`
	}{
		TaxCode: "codice fiscale",
	}
	dataReceive := struct {
		AccessURL string `json:"access_url"`
	}{}

	err := conn.Send("remote-login", dataSend, &dataReceive)
	if err != nil {
		t.Error(err)
		return
	}
	if dataReceive.AccessURL == "" {
		t.Error("access URL not found")
		return
	}

	fmt.Println("Server response:", dataReceive)
}

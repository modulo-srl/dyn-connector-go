package connector

import (
	"fmt"
	"testing"
)

// Session persistor
type testSession struct {
	sessionToken string
}

func (s *testSession) GetSessionToken() string {
	fmt.Println("get session token \"" + s.sessionToken + "\"")
	return s.sessionToken
}

func (s *testSession) SetSessionToken(token string) {
	s.sessionToken = token
	fmt.Println("set session token \"" + s.sessionToken + "\"")
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

	authUID := "test"
	masterToken := "test"

	session := testSession{}
	conn := NewConnector(domain, authUID, masterToken, &session)

	// Enable debugging output (please disable in production)
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

	err := conn.Send("echo", dataSend, &dataReceive)
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

	err := conn.Send("echo/auth", dataSend, &dataReceive)
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

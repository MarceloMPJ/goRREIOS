package request_test

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/MarceloMPJ/goRREIOS/config"
	"github.com/MarceloMPJ/goRREIOS/request"
)

var CorreiosRequest string

const (
	DefaultUser     = "ECT"
	DefaultPassword = "SRO"
	DefaultType     = "L"
	DefaultResult   = "T"
	DefaultLanguage = "101"
)

func Test_NewBuscaEventosRequest(t *testing.T) {
	loadFixtures()

	code := "QI558593122BR"

	t.Run("add default informations to header", func(t *testing.T) {
		expected := newRequest(t, "")
		result, err := request.NewBuscaEventosRequest(code)

		checkEmptyError(t, err)
		checkHeader(t, result, expected)
		checkURL(t, result, expected)
		checkMethod(t, result, expected)
	})

	t.Run("add user, password, typeRequest, result, language and code to request body", func(t *testing.T) {
		expectedBody := requestBody(t, DefaultUser, DefaultPassword, DefaultType, DefaultResult, DefaultLanguage, code)
		expected := newRequest(t, expectedBody)
		result, err := request.NewBuscaEventosRequest(code)

		checkEmptyError(t, err)
		checkBody(t, result, expected)
	})
}

func loadFixtures() {
	CorreiosRequestBytes, err := os.ReadFile("../test/fixtures/correios_request.xml")
	if err != nil {
		panic(err)
	}

	CorreiosRequest = string(CorreiosRequestBytes)
}

func newRequest(t *testing.T, body string) *http.Request {
	t.Helper()

	req, err := http.NewRequest("POST", config.RastroURL(), strings.NewReader(body))

	if err != nil {
		t.Fatalf("unexpected error when create request: %v", err)
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Content-Type", `text/xml; charset="utf-8"`)
	req.Header.Set("User-Agent", "gowsdl/0.1")

	return req
}

func requestBody(t *testing.T, user, password, typeRequest, result, language, code string) string {
	return fmt.Sprintf(CorreiosRequest, user, password, typeRequest, result, language, code)
}

func checkHeader(t *testing.T, result, expected *http.Request) {
	t.Helper()

	if !reflect.DeepEqual(result.Header, expected.Header) {
		t.Errorf("result header: %v\nexpected header: %v", result.Header, expected.Header)
	}
}

func checkBody(t *testing.T, result, expected *http.Request) {
	t.Helper()

	if result.Body == nil {
		t.Fatalf("body of result is nil")
	}

	if expected.Body == nil {
		t.Fatalf("body of expected is nil")
	}

	resultBody, _ := ioutil.ReadAll(result.Body)
	expectedBody, _ := ioutil.ReadAll(expected.Body)

	rexp := regexp.MustCompile("[\t\n ]")

	resultStrBody := string(rexp.ReplaceAll(resultBody, []byte("")))
	expectedStrBody := string(rexp.ReplaceAll(expectedBody, []byte("")))

	if resultStrBody != expectedStrBody {
		t.Errorf("result Body: %v\nexpected Body: %v", resultStrBody, expectedStrBody)
	}
}

func checkURL(t *testing.T, result, expected *http.Request) {
	t.Helper()

	if result.URL.String() != expected.URL.String() {
		t.Errorf("result URL: %s, expected: %s", result.URL, expected.URL)
	}
}

func checkMethod(t *testing.T, result, expected *http.Request) {
	t.Helper()

	if result.Method != expected.Method {
		t.Errorf("result Method: %s, expected: %s", result.Method, expected.Method)
	}
}

func checkEmptyError(t *testing.T, resultErr error) {
	t.Helper()

	if resultErr != nil {
		t.Fatalf("error unexpected: %v\n", resultErr)
	}
}

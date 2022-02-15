package correios_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	correios "github.com/MarceloMPJ/goRREIOS"
	"github.com/MarceloMPJ/goRREIOS/config"
	"github.com/MarceloMPJ/goRREIOS/response"
)

var correiosResponse string
var correiosResponseWithError string

func Test_BuscaEventos(t *testing.T) {
	loadFixtures()

	t.Run("returns events when code is 'AA598971235BR'", func(t *testing.T) {
		code := "AA598971235BR"

		expected := buildDeliveredEvent(code)
		correiosServer := newCorreiosServer(http.StatusOK, responseByCode(code))
		defer correiosServer.Close()

		config.SetRastroURL(correiosServer.URL)
		result, err := correios.BuscaEventos(code)

		checkEmptyError(t, err)
		checkResult(t, result, expected)
	})

	t.Run("returns events when code is 'BA598971235BR'", func(t *testing.T) {
		code := "BA598971235BR"

		expected := buildDeliveredEvent(code)
		correiosServer := newCorreiosServer(http.StatusOK, responseByCode(code))
		defer correiosServer.Close()

		config.SetRastroURL(correiosServer.URL)
		result, err := correios.BuscaEventos(code)

		checkEmptyError(t, err)
		checkResult(t, result, expected)
	})

	t.Run("returns events with error when code is invalid", func(t *testing.T) {
		code := "QI558593122BR"

		expected := []*response.Event{}
		expectedErr := fmt.Errorf(response.CorreiosException, "Objeto não encontrado na base de dados dos Correios.")
		correiosServer := newCorreiosServer(http.StatusOK, responseWithErrorByCode(code))
		defer correiosServer.Close()

		config.SetRastroURL(correiosServer.URL)
		result, err := correios.BuscaEventos(code)

		checkError(t, err, expectedErr)
		checkResult(t, result, expected)
	})

	t.Run("returns error when status code is InternalServerError", func(t *testing.T) {
		code := "BA598971235BR"

		expected := []*response.Event{}
		correiosServer := newCorreiosServer(http.StatusInternalServerError, "")
		defer correiosServer.Close()

		config.SetRastroURL(correiosServer.URL)
		result, err := correios.BuscaEventos(code)

		checkError(t, err, fmt.Errorf(correios.InvalidStatusCode, http.StatusInternalServerError))
		checkResult(t, result, expected)
	})
}

func Test_IntegrationTest_BuscaEventos(t *testing.T) {
	t.Run("returns an array with Event expected", func(t *testing.T) {
		result, err := correios.BuscaEventos("AA598971235BR")

		expected := []*response.Event{
			{
				Tipo:      "PO",
				Status:    1,
				Data:      "04/02/2022",
				Hora:      "13:03",
				Descricao: "Objeto postado",
				Local:     "Agência dos Correios",
				Codigo:    0,
				Cidade:    "SAO PAULO",
				UF:        "SP",
			},
		}

		checkEmptyError(t, err)
		checkResult(t, result, expected)
	})

	t.Run("returns an error when code is incorrect", func(t *testing.T) {
		_, resultErr := correios.BuscaEventos("QI720913955BA")
		expectedError := fmt.Errorf(response.CorreiosException, "Objeto não encontrado na base de dados dos Correios.")

		checkError(t, resultErr, expectedError)
	})
}

func loadFixtures() {
	correiosResponseBytes, err := os.ReadFile("./test/fixtures/correios_response.xml")
	if err != nil {
		panic(err)
	}

	correiosResponseWithErrorBytes, err := os.ReadFile("./test/fixtures/correios_response_error.xml")
	if err != nil {
		panic(err)
	}

	correiosResponse = string(correiosResponseBytes)
	correiosResponseWithError = string(correiosResponseWithErrorBytes)
}

func checkResult(t *testing.T, result, expected []*response.Event) {
	t.Helper()

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("result: %v\nexpected %v", eventsToString(result), eventsToString(expected))
	}
}

func checkEmptyError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("error not expeted: %v", err)
	}
}

func checkError(t *testing.T, result, expected error) {
	t.Helper()

	if result.Error() != expected.Error() {
		t.Errorf("Result error: '%v', expected error: '%v'", result, expected)
	}
}

func buildSendingEvent(code string) (codeEvents []*response.Event) {
	event := response.Event{
		Tipo:      "BDE",
		Status:    23,
		Data:      "18/03/2014",
		Hora:      "18:37",
		Descricao: "Objeto enviado ao remetente",
		Local:     "CTCE MACEIO",
		Codigo:    57060971,
		Cidade:    "MACEIO",
		UF:        "AL",
	}
	codeEvents = append(codeEvents, &event)

	return codeEvents
}

func buildDeliveredEvent(code string) (codeEvents []*response.Event) {
	codeEvents = buildSendingEvent(code)
	event := response.Event{
		Tipo:      "BDE",
		Status:    23,
		Data:      "18/03/2014",
		Hora:      "20:37",
		Descricao: "Objeto entregue",
		Local:     "CTCE MACEIO",
		Codigo:    57060971,
		Cidade:    "MACEIO",
		UF:        "AL",
	}
	codeEvents = append(codeEvents, &event)

	return codeEvents
}

func eventsToString(events []*response.Event) (str string) {
	for _, event := range events {
		str += fmt.Sprint(*event)
	}
	str = fmt.Sprintf("[%s]", str)

	return
}

func responseByCode(code string) string {
	return fmt.Sprintf(correiosResponse, code)
}

func responseWithErrorByCode(code string) string {
	return fmt.Sprintf(correiosResponseWithError, code)
}

func newCorreiosServer(status int, response string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(status)
		res.Write([]byte(response))
	}))
}

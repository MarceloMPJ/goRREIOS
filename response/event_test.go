package response_test

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"testing"

	"github.com/MarceloMPJ/correios/response"
)

var correiosResponseSending string
var correiosResponseDelivered string
var correiosResponseWithError string
var correiosBadResponse string

func Test_XMLToEvents(t *testing.T) {
	loadFixtures()
	code := "QI558593122BR"

	t.Run("with sending event", func(t *testing.T) {
		expected := buildSendingEvent(code)
		result, err := response.XMLToEvents(fmt.Sprintf(correiosResponseSending, code))

		checkEmptyError(t, err)
		checkResponse(t, result, expected)
	})

	t.Run("with sending and delivered event", func(t *testing.T) {
		expected := buildDeliveredEvent(code)
		result, err := response.XMLToEvents(fmt.Sprintf(correiosResponseDelivered, code))

		checkEmptyError(t, err)
		checkResponse(t, result, expected)
	})

	t.Run("with error, tracking code invalid", func(t *testing.T) {
		expected := []*response.Event{}
		result, resultErr := response.XMLToEvents(fmt.Sprintf(correiosResponseWithError, code))
		expectedErr := fmt.Errorf(response.CorreiosException, "Objeto não encontrado na base de dados dos Correios.")

		checkErrors(t, resultErr, expectedErr)
		checkResponse(t, result, expected)
	})

	t.Run("with bad response", func(t *testing.T) {
		expected := []*response.Event{}
		result, resultErr := response.XMLToEvents(fmt.Sprintf(correiosBadResponse, code))

		checkErrors(t, resultErr, errors.New("expected element type <Envelope> but have <Body>"))
		checkResponse(t, result, expected)
	})
}

func loadFixtures() {
	correiosResponseSendingBytes, err := os.ReadFile("../test/fixtures/sending_correios_response.xml")
	if err != nil {
		panic(err)
	}

	correiosResponseDeliveredBytes, err := os.ReadFile("../test/fixtures/delivered_correios_response.xml")
	if err != nil {
		panic(err)
	}

	correiosResponseWithErrorBytes, err := os.ReadFile("../test/fixtures/correios_response_error.xml")
	if err != nil {
		panic(err)
	}

	correiosBadResponseBytes, err := os.ReadFile("../test/fixtures/bad_response.xml")
	if err != nil {
		panic(err)
	}

	correiosResponseSending = string(correiosResponseSendingBytes)
	correiosResponseDelivered = string(correiosResponseDeliveredBytes)
	correiosResponseWithError = string(correiosResponseWithErrorBytes)
	correiosBadResponse = string(correiosBadResponseBytes)
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

func checkResponse(t *testing.T, result, expected []*response.Event) {
	t.Helper()

	if !reflect.DeepEqual(result, expected) {
		t.Errorf("result: %v\nexpected: %v", eventsToString(result), eventsToString(expected))
	}
}

func eventsToString(events []*response.Event) (str string) {
	for _, event := range events {
		str += fmt.Sprint(*event)
	}
	str = fmt.Sprintf("[%s]", str)

	return
}

func checkEmptyError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("error unexpected: %v", err)
	}
}

func checkErrors(t *testing.T, result, expected error) {
	t.Helper()

	if result.Error() != expected.Error() {
		t.Errorf("result error: '%v', expected error: '%v'", result, expected)
	}
}

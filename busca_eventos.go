package correios

import (
	"fmt"
	"io"
	"net/http"

	"github.com/MarceloMPJ/correios/request"
	"github.com/MarceloMPJ/correios/response"
)

const InvalidStatusCode string = "invalid status code: %d"

// BuscaEventos is a function that receive one tracking code and return events for this code
func BuscaEventos(code string) ([]*response.Event, error) {
	resp, err := requestBuscaEventos(code)
	if err != nil {
		return []*response.Event{}, err
	}

	events, err := response.XMLToEvents(resp)
	if err != nil {
		return []*response.Event{}, err
	}

	return events, nil
}

func requestBuscaEventos(code string) (string, error) {
	req, err := request.NewBuscaEventosRequest(code)
	if err != nil {
		return "", err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode > 299 {
		return "", fmt.Errorf(InvalidStatusCode, resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

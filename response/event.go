package response

import (
	"encoding/xml"
	"fmt"
)

const CorreiosException string = "correios exception: %v"

type envelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Body    body     `xml:"Body"`
}

type body struct {
	XMLName      xml.Name     `xml:"Body"`
	BuscaEventos buscaEventos `xml:"buscaEventosResponse"`
}

type buscaEventos struct {
	ReturnResponse returnResponse `xml:"return"`
}

type returnResponse struct {
	Objeto objeto `xml:"objeto"`
}

type objeto struct {
	Erro    string   `xml:"erro"`
	Eventos []*Event `xml:"evento"`
}

// Event is a struct that represent each event of tracking code
type Event struct {
	Tipo      string `xml:"tipo"`
	Status    uint32 `xml:"status"`
	Data      string `xml:"data"`
	Hora      string `xml:"hora"`
	Descricao string `xml:"descricao"`
	Local     string `xml:"local"`
	Codigo    uint32 `xml:"codigo"`
	Cidade    string `xml:"cidade"`
	UF        string `xml:"uf"`
}

// XMLToEvents is a function that receive a csv string and convert to array of Event
// If XML has a error should return this error
func XMLToEvents(xmlresponse string) ([]*Event, error) {
	var envelope envelope

	if err := xml.Unmarshal([]byte(xmlresponse), &envelope); err != nil {
		return []*Event{}, err
	}

	if envelope.Body.BuscaEventos.ReturnResponse.Objeto.Erro != "" {
		return []*Event{}, fmt.Errorf(CorreiosException, envelope.Body.BuscaEventos.ReturnResponse.Objeto.Erro)
	}

	return envelope.Body.BuscaEventos.ReturnResponse.Objeto.Eventos, nil
}

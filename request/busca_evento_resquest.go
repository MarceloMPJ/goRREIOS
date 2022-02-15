package request

import (
	"encoding/xml"
	"net/http"
	"strings"

	"github.com/MarceloMPJ/correios/config"
)

func NewBuscaEventosRequest(code string) (*http.Request, error) {
	bodyXML, err := buildBodyXML(code)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", config.RastroURL(), strings.NewReader(bodyXML))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "*/*")
	req.Header.Set("Content-Type", `text/xml; charset="utf-8"`)
	req.Header.Set("User-Agent", "gowsdl/0.1")

	return req, nil
}

func buildBodyXML(code string) (string, error) {
	requestEnvelope := envelope{
		Tag1: "http://schemas.xmlsoap.org/soap/envelope/",
		Tag2: "http://resource.webservice.correios.com.br/",
		Body: body{
			BuscaEventos: buscaEventos{
				Usuario:   config.User(),
				Senha:     config.Password(),
				Tipo:      config.Type(),
				Resultado: config.Result(),
				Lingua:    config.Language(),
				Objetos:   code,
			},
		},
	}

	xmlByte, err := xml.Marshal(&requestEnvelope)
	if err != nil {
		return "", err
	}

	return string(xmlByte), nil
}

type envelope struct {
	XMLName xml.Name `xml:"soapenv:Envelope"`
	Tag1    string   `xml:"xmlns:soapenv,attr"`
	Tag2    string   `xml:"xmlns:res,attr"`
	Body    body     `xml:"soapenv:Body"`
}

type body struct {
	XMLName      xml.Name     `xml:"soapenv:Body"`
	BuscaEventos buscaEventos `xml:"res:buscaEventos"`
}

type buscaEventos struct {
	Usuario   string `xml:"usuario"`
	Senha     string `xml:"senha"`
	Tipo      string `xml:"tipo"`
	Resultado string `xml:"resultado"`
	Lingua    string `xml:"lingua"`
	Objetos   string `xml:"objetos"`
}

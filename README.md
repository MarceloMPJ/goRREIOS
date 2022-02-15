# goRREIOS
### An implementation to API Correios builded using GoLang

```GOLANG
// Example:
// Get events of 'AA598971235BR'

result, err := correios.BuscaEventos("AA598971235BR")

fmt.Println(result)
/*
[]*response.Event{
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
*/
````

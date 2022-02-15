# goRREIOS
### An implementation to API Correios builded using GoLang

```GOLANG
// Example:
// Get events of 'AA598971235BR'
package main

import (
	"fmt"

	correios "github.com/MarceloMPJ/goRREIOS"
)

func main() {
	events, err := correios.BuscaEventos("AA598971235BR")

	if err != nil {
		fmt.Println(err)
		return
	}

	if len(events) > 0 {
		fmt.Println(*events[0])
    /*
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
    }
    */
	}
}

````

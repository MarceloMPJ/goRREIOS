package config

var (
	rastroURL  = "http://webservice.correios.com.br/service/rastro"
	user       = "ECT"
	password   = "SRO"
	typeResult = "L"
	result     = "T"
	language   = "101"
)

func SetConfig(userConfig, passwordConfig string) {
	user = userConfig
	password = passwordConfig
}

func SetRastroURL(url string) {
	rastroURL = url
}

func RastroURL() string {
	return rastroURL
}

func User() string {
	return user
}

func Password() string {
	return password
}

func Type() string {
	return typeResult
}

func Result() string {
	return result
}

func Language() string {
	return language
}

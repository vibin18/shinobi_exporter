package config

type (
	Opts struct {
		ServerBind string `long:"bind"     env:"SERVER_BIND"   description:"Server address"     default:":8080"`
	}
)

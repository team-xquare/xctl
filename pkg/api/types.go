package api

import "fmt"

type Application struct {
	Name          string
	Host          string
	ImageRegistry string
	ImageTag      string
	ContainerPort int32
	Prefix        string
	Environment   string
	Type          string
}

const (
	Backend  = "backend"
	Frontend = "frontend"
)

func CheckApplicationType(t string) (string, error) {
	switch t {
	case "backend", "be":
		return Backend, nil
	case "frontend", "fe":
		return Frontend, nil
	default:
		return "", fmt.Errorf("cannot use an application type: %s.\n support \"frontend\", \"fe\", \"backend\" \"be\" only", t)

	}
}

const (
	Production = "prod"
	Staging    = "staging"
)

func CheckApplicationEnvironment(e string) (string, error) {
	switch e {
	case "production", "prod":
		return Production, nil
	case "staging", "stag":
		return Staging, nil
	default:
		return "", fmt.Errorf("cannot use an application environment: %s.\n support \"production\", \"prod\", \"staging\" \"stag\" only", e)

	}
}

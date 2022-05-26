package api

const (
	Backend  = "backend"
	Frontend = "frontend"

	Production = "production"
	Staging    = "staging"
)

type Application struct {
	Name          string
	Type          string
	Host          string
	ImageRegistry string
	ImageTag      string
	Environment   string
	ContainerPort int32
	Prefix        string
}

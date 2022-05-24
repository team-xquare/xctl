package api

type Application struct {
	Name          string
	Type          string
	Host          string
	ImageRegistry string
	ImageTag      string
	ContainerPort int32
}

package application

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/go-github/v45/github"
	"github.com/xctl/pkg/api"
)

func (c *applications) Get(ctx context.Context) ([]*api.Application, error) {
	return c.getApplication(ctx, "")
}

func (c *applications) GetFrontend(ctx context.Context) ([]*api.Application, error) {
	return c.getApplication(ctx, api.Frontend)
}

func (c *applications) GetBackend(ctx context.Context) ([]*api.Application, error) {
	return c.getApplication(ctx, api.Backend)
}

func (c *applications) getApplication(ctx context.Context, serviceType string) (apps []*api.Application, err error) {
	directoryContent, err := c.getDirectoryContent(ctx, "resource")
	if err != nil {
		return nil, err
	}
	for _, content := range directoryContent {
		fileContent, err := c.getFileContent(ctx, *content.Path)
		if err != nil {
			return nil, err
		}

		if strings.HasSuffix(*content.Path, ".json") {
			var app *api.Application
			err = json.Unmarshal(fileContent, &app)
			if err != nil {
				return nil, err
			}
			if serviceType == "" {
				apps = append(apps, app)
			} else if strings.HasSuffix(*content.Path, fmt.Sprintf("%s.json", serviceType)) {
				apps = append(apps, app)
			}
		}
	}
	return apps, nil
}

func (c *applications) getDirectoryContent(ctx context.Context, path string) ([]*github.RepositoryContent, error) {
	_, directoryContent, err := c.client.GetContents(ctx, path)
	return directoryContent, err
}

func (c *applications) getFileContent(ctx context.Context, path string) ([]byte, error) {
	fileContent, _, err := c.client.GetContents(ctx, path)
	if err != nil {
		return nil, err
	}
	content, err := fileContent.GetContent()
	return []byte(content), err
}

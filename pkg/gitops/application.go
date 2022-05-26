package gitops

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	"github.com/google/go-github/github"
	"github.com/xctl/pkg/api"
	cmdutil "github.com/xctl/pkg/cmd/util"
)

var appTemplatePath = "template/app-template"

type ApplicationInterface interface {
	Create(ctx context.Context, app *api.Application) error
	Get(ctx context.Context) ([]*api.Application, []*api.Application, error)
}

type applications struct {
	gitops *Gitops
}

func newApplication(c *Gitops) *applications {
	return &applications{c}
}

func (c *applications) Create(ctx context.Context, app *api.Application) error {
	ref, err := c.gitops.client.GetRef(ctx)
	if err != nil {
		return err
	}
	entries, err := c.createTreeEntries(ctx, app)
	if err != nil {
		return err
	}
	tree, err := c.gitops.client.CreateTreeFromEntries(ctx, ref, entries)
	if err != nil {
		return err
	}
	head, err := c.gitops.client.GetHeadFromRef(ctx, ref)
	if err != nil {
		return err
	}
	message := fmt.Sprintf("⚡️ :: %s-%s 애플리케이션 추가", app.Name, app.Type)
	newCommit, err := c.gitops.client.CreateCommit(ctx, message, tree, head.Commit)
	if err != nil {
		return err
	}
	err = c.gitops.client.UpdateRef(ctx, ref, newCommit)
	return err
}

func (c *applications) createTreeEntries(ctx context.Context, app *api.Application) ([]github.TreeEntry, error) {
	var entries = []github.TreeEntry{}
	entries = append(entries, github.TreeEntry{
		Path:    github.String(fmt.Sprintf("resource/app-%s-%s.json", app.Name, app.Type)),
		Mode:    github.String("100644"),
		Type:    github.String("blob"),
		Content: github.String(cmdutil.MarshalObject(app)),
	})

	var buf bytes.Buffer

	err := filepath.Walk(appTemplatePath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.IsDir() {
				return nil
			}
			if runtime.GOOS == "windows" {
				path = strings.Replace(path, "\\", "/", -1)
			}

			tmpl, err := template.ParseFiles(path)
			if err != nil {
				return err
			}

			data := struct {
				Name          string
				Type          string
				Host          string
				ImageRegistry string
				ImageTag      string
				ContainerPort int32
				Environment   string
				Prefix        string
			}{app.Name,
				app.Type,
				app.Host,
				app.ImageRegistry,
				app.ImageTag,
				app.ContainerPort,
				app.Environment,
				app.Prefix}
			err = tmpl.Execute(&buf, &data)
			if err != nil {
				return err
			}

			fileName := strings.Join(strings.Split(path, "/")[2:], "/")
			content := buf.String()
			buf.Reset()
			var entryPath string
			if fileName == "app-template.yaml" {
				entryPath = fmt.Sprintf("apps/templates/%s/%s", app.Type, fmt.Sprintf("app-%s-%s.yaml", app.Name, app.Type))
			} else {
				entryPath = fmt.Sprintf("charts/%s/%s/%s", app.Type, app.Name, fileName)
			}

			entries = append(entries, github.TreeEntry{
				Path:    github.String(entryPath),
				Content: github.String(content),
				Type:    github.String("blob"),
				Mode:    github.String("100644"),
			})
			return nil
		})
	return entries, err
}

func (c *applications) Get(ctx context.Context) (frontend []*api.Application, backend []*api.Application, err error) {
	directoryContent, err := c.getDirectoryContent(ctx, "resource")
	if err != nil {
		return nil, nil, err
	}
	for _, content := range directoryContent {
		fileContent, err := c.getFileContent(ctx, *content.Path)
		if err != nil {
			return nil, nil, err
		}

		var app *api.Application
		err = json.Unmarshal(fileContent, &app)
		if err != nil {
			return nil, nil, err
		}

		if strings.HasSuffix(*content.Path, "frontend.json") {
			frontend = append(frontend, app)
		} else if strings.HasSuffix(*content.Path, "backend.json") {
			backend = append(backend, app)
		}
	}
	return frontend, backend, nil
}

func (c *applications) getDirectoryContent(ctx context.Context, path string) ([]*github.RepositoryContent, error) {
	_, directoryContent, err := c.gitops.client.GetContents(ctx, path)
	return directoryContent, err
}

func (c *applications) getFileContent(ctx context.Context, path string) ([]byte, error) {
	fileContent, _, err := c.gitops.client.GetContents(ctx, path)
	if err != nil {
		return nil, err
	}
	content, err := fileContent.GetContent()
	return []byte(content), err
}

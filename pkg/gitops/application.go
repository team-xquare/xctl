package gitops

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	"github.com/google/go-github/github"
	"github.com/xctl/pkg/api"
)

var appTemplatePath = "template/app-template"

type ApplicationInterface interface {
	Create(ctx context.Context, app *api.Application) error
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
	var entries []github.TreeEntry
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
			}{app.Name, app.Type, app.Host, app.ImageRegistry, app.ImageTag, app.ContainerPort, app.Environment}
			err = tmpl.Execute(&buf, &data)
			if err != nil {
				return err
			}

			fileName := strings.Join(strings.Split(path, "/")[2:], "/")
			content := buf.String()
			var entryPath string
			if fileName == "app-template.yaml" {
				entryPath = fmt.Sprintf("apps/templates/%s/%s", app.Type, fmt.Sprintf("app-%s-%s.yaml", app.Name, app.Type))
			} else {
				entryPath = fmt.Sprintf("charts/%s/%s/%s", app.Type, app.Name, fileName)
			}

			entries = append(entries, github.TreeEntry{
				Path:    github.String(entryPath),
				Content: github.String(content),
				Mode:    github.String("blob"),
				Type:    github.String("100644"),
			})
			return nil
		})
	return entries, err
}

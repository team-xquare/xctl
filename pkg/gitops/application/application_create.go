package application

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"

	"github.com/google/go-github/v45/github"
	cmdutil "github.com/xctl/pkg/cmd/util"
)

func (c *applications) Create(ctx context.Context) error {
	ref, err := c.client.GetRef(ctx)
	if err != nil {
		return err
	}
	entries, err := c.createTreeEntries(ctx)
	if err != nil {
		return err
	}
	tree, err := c.client.CreateTreeFromEntries(ctx, ref, entries)
	if err != nil {
		return err
	}
	head, err := c.client.GetHeadFromRef(ctx, ref)
	if err != nil {
		return err
	}
	message := fmt.Sprintf("⚡️ :: %s-%s 애플리케이션 추가", c.application.Name, c.application.Type)
	newCommit, err := c.client.CreateCommit(ctx, message, tree, head.Commit)
	if err != nil {
		return err
	}
	err = c.client.UpdateRef(ctx, ref, newCommit)
	return err
}

func (c *applications) createTreeEntries(ctx context.Context) ([]*github.TreeEntry, error) {
	var entries = []*github.TreeEntry{}
	entries = append(entries, &github.TreeEntry{
		Path:    github.String(fmt.Sprintf("resource/app-%s-%s.json", c.application.Name, c.application.Type)),
		Mode:    github.String("100644"),
		Type:    github.String("blob"),
		Content: github.String(cmdutil.MarshalObject(c.application)),
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

			err = tmpl.Execute(&buf, &c.application)
			if err != nil {
				return err
			}

			fileName := strings.Join(strings.Split(path, "/")[2:], "/")
			content := buf.String()
			buf.Reset()
			var entryPath string
			if fileName == "app-template.yaml" {
				entryPath = fmt.Sprintf("apps/templates/%s/%s",
					c.application.Type,
					fmt.Sprintf("app-%s-%s.yaml", c.application.Name, c.application.Type))
			} else {
				entryPath = fmt.Sprintf("charts/%s/%s/%s",
					c.application.Type,
					c.application.Name,
					fileName)
			}

			entries = append(entries, &github.TreeEntry{
				Path:    github.String(entryPath),
				Content: github.String(content),
				Type:    github.String("blob"),
				Mode:    github.String("100644"),
			})
			return nil
		})
	return entries, err
}

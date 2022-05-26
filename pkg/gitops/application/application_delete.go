package application

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/google/go-github/v45/github"
)

func (c *applications) Delete(ctx context.Context) error {
	ref, err := c.client.GetRef(ctx)
	if err != nil {
		return err
	}
	entries, err := c.deleteTreeEntries(ctx)
	if err != nil {
		return err
	}
	newTree, err := c.client.CreateTreeFromEntries(ctx, ref, entries)
	if err != nil {
		return err
	}
	head, err := c.client.GetHeadFromRef(ctx, ref)
	if err != nil {
		return err
	}
	message := fmt.Sprintf("⚰️ :: %s-%s 애플리케이션 삭제", c.application.Name, c.application.Type)
	newCommit, err := c.client.CreateCommit(ctx, message, newTree, head.Commit)
	if err != nil {
		return err
	}
	return c.client.UpdateRef(ctx, ref, newCommit)
}

func (c *applications) deleteTreeEntries(ctx context.Context) ([]*github.TreeEntry, error) {
	var entries []*github.TreeEntry
	entries = append(entries, &github.TreeEntry{
		Path:    github.String(fmt.Sprintf("resource/app-%s-%s.json", c.application.Name, c.application.Type)),
		Mode:    github.String("100644"),
		Type:    github.String("blob"),
		SHA:     nil,
		Content: nil,
	})

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

			fileName := strings.Join(strings.Split(path, "/")[2:], "/")
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
				Mode:    github.String("100644"),
				Type:    github.String("blob"),
				SHA:     nil,
				Content: nil,
			})

			return nil
		})
	return entries, err
}

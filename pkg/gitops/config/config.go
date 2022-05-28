package config

import (
	"encoding/json"
	"io/ioutil"
	"os"

	cmdutil "github.com/xctl/pkg/cmd/util"
)

var (
	CredentialDir  = cmdutil.GetHomeDir() + "/.xctl"
	CredentialPath = CredentialDir + "/credential.json"
)

type Credential struct {
	GithubToken string
}

func SetCredential(c *Credential) error {
	err := createCredentialDir()
	if err != nil {
		return err
	}

	f, err := os.Create(CredentialPath)
	if err != nil {
		return err
	}
	defer f.Close()

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	return err
}

func createCredentialDir() error {
	if _, err := os.Stat(CredentialDir); os.IsNotExist(err) {
		err := os.Mkdir(CredentialDir, os.ModeDir)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetCredential() (*Credential, error) {
	f, err := ioutil.ReadFile(CredentialPath)
	if err != nil {
		return nil, err
	}

	var c *Credential
	err = json.Unmarshal(f, &c)
	return c, err
}

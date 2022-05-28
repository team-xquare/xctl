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

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(CredentialPath, data, 0666)
}

func createCredentialDir() error {
	if _, err := os.Stat(CredentialDir); os.IsNotExist(err) {
		err := os.Mkdir(CredentialDir, 0777)
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

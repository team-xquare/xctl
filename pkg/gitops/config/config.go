package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"runtime"
)

var (
	CredentialDir  = getHomeDir() + "/.xctl"
	CredentialPath = CredentialDir + "/credential.json"
)

type Credential struct {
	GithubToken string
}

func getHomeDir() string {
	if runtime.GOOS == "windows" {
		home := os.Getenv("HOMEDRIVE") + os.Getenv("HOMEPATH")
		if home == "" {
			home = os.Getenv("USERPROFILE")
		}
		return home
	} else if runtime.GOOS == "linux" {
		home := os.Getenv("XDG_CONFIG_HOME")
		if home != "" {
			return home
		}
	}
	return os.Getenv("HOME")
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

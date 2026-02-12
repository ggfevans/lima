package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// Credentials stores LinkedIn authentication data.
type Credentials struct {
	Cookie       string `json:"cookie"`
	XLiTrack     string `json:"x_li_track"`
	PageInstance string `json:"page_instance"`
}

// IsEmpty returns true if no credentials are stored.
func (c Credentials) IsEmpty() bool {
	return c.Cookie == ""
}

// CredentialsPath returns the path to the credentials file.
func CredentialsPath() (string, error) {
	dir, err := ConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "credentials.json"), nil
}

// LoadCredentials reads stored credentials from disk.
func LoadCredentials() (Credentials, error) {
	var creds Credentials

	path, err := CredentialsPath()
	if err != nil {
		return creds, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return creds, nil
		}
		return creds, err
	}

	if err := json.Unmarshal(data, &creds); err != nil {
		return creds, err
	}

	return creds, nil
}

// SaveCredentials writes credentials to disk with restricted permissions.
func SaveCredentials(creds Credentials) error {
	dir, err := ConfigDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	path := filepath.Join(dir, "credentials.json")
	data, err := json.MarshalIndent(creds, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}

// ClearCredentials removes stored credentials.
func ClearCredentials() error {
	path, err := CredentialsPath()
	if err != nil {
		return err
	}

	err = os.Remove(path)
	if err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

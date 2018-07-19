package packet

import (
	"fmt"
	"os"
	"strings"

	"github.com/packethost/packngo"
	. "github.com/pharmer/cloud-storage/cloud"
	"github.com/pharmer/cloud-storage/util"
)

const (
	tokenEnv     = "PACKET_API_KEY"
	projectIDEnv = "PACKET_PROJECT_ID"
	apiKey       = "apiKey"
	projectIDKey = "projectID"
)

type TokenSource struct {
	ApiKey string `json:"apiKey"`
}

func getProjectID() (string, error) {
	if id, ok := os.LookupEnv(projectIDEnv); ok && id != "" {
		return strings.TrimSpace(id), nil
	}
	return util.ReadSecretKeyFromFile(SecretDefaultLocation, projectIDKey)
}

func getCredential() (*TokenSource, error) {
	if t, err := util.ReadSecretKeyFromFile(SecretDefaultLocation, apiKey); err == nil {
		return &TokenSource{
			ApiKey: t,
		}, nil
	}

	if f, ok := os.LookupEnv(CredentialFileEnv); ok && f != "" {
		cred, err := util.ReadCredentialFromFile(f, &TokenSource{})
		if err != nil {
			return nil, err
		}
		return cred.(*TokenSource), nil
	}

	if t, ok := os.LookupEnv(tokenEnv); ok && t != "" {
		return &TokenSource{
			ApiKey: strings.TrimSpace(t),
		}, nil
	}

	cred, err := util.ReadCredentialFromFile(CredentialDefaultLocation, &TokenSource{})
	if err != nil {
		return nil, err
	}
	tokenSource := cred.(*TokenSource)
	if tokenSource.ApiKey != "" {
		return tokenSource, nil
	}

	return nil, fmt.Errorf("no credential provided for packet")
}

func (t *TokenSource) getClient() *packngo.Client {
	return packngo.NewClientWithAuth("", t.ApiKey, nil)
}

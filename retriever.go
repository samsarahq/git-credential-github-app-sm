package main

import (
	"bytes"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"os/exec"
	"strings"
)

// interface guard
var _ stscreds.IdentityTokenRetriever = (*commandRetriever)(nil)

// commandRetriever implements AWS's IdentityTokenRetriever for getting the OIDC token for a web identity role
// assumption
type commandRetriever struct {
	command string
}

// GetIdentityToken implements IdentityTokenRetriever
func (c *commandRetriever) GetIdentityToken() ([]byte, error) {
	var buffer bytes.Buffer

	cmdParts := strings.Split(c.command, " ") // This is not ideal, but it should do for now.

	cmd := exec.Command(cmdParts[0], cmdParts[1:]...)
	cmd.Stdout = &buffer
	err := cmd.Run()

	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

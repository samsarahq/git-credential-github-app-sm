package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/samsarahq/github-app-credential-helper/common"
)

func main() {
	secretArn := flag.String("secret-arn", "", "Secret ARN")
	role := flag.String("role", "", "Role ARN if a role should be assumed")
	tokenCommand := flag.String("token-command", "", "OIDC token command if using web identity")
	flag.Parse()

	subcommand := flag.Arg(0)

	// We are ignoring the other operations here and assume any caching it provided by another helper.
	if subcommand == "get" {
		get(secretArn, role, tokenCommand)
	}
}

// get provides the functionality for the "get" command
func get(secretArn *string, role *string, tokenCommand *string) {
	if secretArn == nil || *secretArn == "" {
		log.Fatal("-secret-arn is required")
	}

	// Use the common check to see if we should even run this.
	if !common.ShouldRun() {
		return
	}

	provider := secretsManagerProvider{
		secretArn:    *secretArn,
		role:         role,
		tokenCommand: tokenCommand,
	}

	helper := common.NewAuthenticator(&provider)
	output, err := helper.Authenticate()
	if err != nil {
		log.Fatal(err)
	}

	// The output needs to be pristine, so we print here at the end only if we are sure everything is working.
	fmt.Println(output)
}

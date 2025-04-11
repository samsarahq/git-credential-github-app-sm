package main

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	"github.com/samsarahq/github-app-credential-helper/common"
)

// interface guard
var _ common.SecretProvider = (*secretsManagerProvider)(nil)

type secretsManagerProvider struct {
	role         *string
	tokenCommand *string
	secretArn    string
}

// Credentials retrieves the app credentials from Secrets Manager
func (s *secretsManagerProvider) Credentials() (*common.AppSecret, error) {

	shouldAssume := false
	if s.role != nil && *s.role != "" {
		shouldAssume = true
	}

	hasToken := false
	if s.tokenCommand != nil && *s.tokenCommand != "" {
		hasToken = true
		if !shouldAssume {
			return nil, errors.New("-token-command not allowed when not assuming a role")
		}
	}

	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, err
	}

	// Set up STS credential provider if necessary.
	if shouldAssume {
		stsClient := sts.NewFromConfig(cfg)
		var creds aws.CredentialsProvider
		if hasToken {
			creds = stscreds.NewWebIdentityRoleProvider(stsClient, *s.role, &commandRetriever{command: *s.tokenCommand})
		} else {
			creds = stscreds.NewAssumeRoleProvider(stsClient, *s.role)
		}
		cfg.Credentials = aws.NewCredentialsCache(creds)
	}

	sm := secretsmanager.NewFromConfig(cfg)

	output, err := sm.GetSecretValue(context.Background(), &secretsmanager.GetSecretValueInput{SecretId: &s.secretArn})
	if err != nil {
		return nil, err
	}

	var result common.AppSecret
	err = json.Unmarshal([]byte(*output.SecretString), &result)

	return &result, err
}

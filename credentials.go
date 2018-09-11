package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

var (
	envVariables = map[string]string{
		"AWS_ACCESS_KEY_ID":      ".AccessKeyID",
		"AWS_SECRET_ACCESS_KEY":  ".SecretAccessKey",
		"AWS_SESSION_TOKEN":      ".SessionToken",
		"AWS_SESSION_EXPIRATION": ".TokenExpiration",
		"AWS_ACCOUNT_NAME":       ".AccountName",
		"AWS_ACCOUNT_NUMBER":     ".AccountID",
	}
)

type EnvVariables struct {
	AccountID       string
	AccessKeyID     string
	AccountName     string
	Delimiter       string
	TokenExpiration string
	Prefix          string
	Region          string
	SecretAccessKey string
	SessionToken    string
	Suffix          string
}

func envTemplate() string {
	var templ string = "{{ if .Region }}{{ .Prefix }}AWS_REGION{{ .Delimiter }}{{ .Region }}{{ .Suffix }}{{ .Prefix }}AWS_DEFAULT_REGION{{ .Delimiter }}{{ .Region }}{{ .Suffix }}{{ end }}"

	keys := []string{}
	for env, _ := range envVariables {
		keys = append(keys, env)
	}
	sort.Strings(keys)

	for _, key := range keys {
		templ += fmt.Sprintf(
			"{{ .Prefix }}%s{{ .Delimiter }}{{ %s }}{{ .Suffix }}",
			key,
			envVariables[key],
		)
	}

	return templ
}

// Detect the User Shell
func detectShell() string {
	var userShell string
	switch runtime.GOOS {
	case "windows":
		userShell = "powershell"
	case "linux":
		userShell = "bash"
	case "darwin":
		userShell = "bash"
	}

	return userShell
}

// Generate Session Name
func generateSessionName(roleName, mfaDeviceID string) string {
	userName := strings.Split(mfaDeviceID, "mfa/")[1]
	timestamp := time.Now().Unix()

	return fmt.Sprintf("%s@%s_%d", userName, roleName, timestamp)
}

type AWSCredentials struct {
	AccessKeyID     string
	SecretAccessKey string
}

func (a AWSCredentials) Retrieve() (credentials.Value, error) {
	return credentials.Value{
		AccessKeyID:     a.AccessKeyID,
		SecretAccessKey: a.SecretAccessKey,
		ProviderName:    "tok",
	}, nil
}

func (a AWSCredentials) IsExpired() bool {
	return false
}

type assumeRoleInput struct {
	AWSAccessKeyID     string `required:"true"`
	AWSSecretAccessKey string `required:"true"`
	AWSAccountNumber   string `required:"true"`
	RoleName           string `required:"true"`
	MFADeviceID        string `required:"true"`
	TokenCode          string `required:"true"`
	SessionName        string
	Duration           int
}

func assumeRole(input assumeRoleInput) (*sts.AssumeRoleOutput, error) {
	creds := AWSCredentials{
		AccessKeyID:     input.AWSAccessKeyID,
		SecretAccessKey: input.AWSSecretAccessKey,
	}

	svc := sts.New(session.New(
		&aws.Config{
			Credentials: credentials.NewCredentials(&creds),
		},
	))
	roleArn := fmt.Sprintf(
		"arn:aws:iam::%s:role/%s",
		input.AWSAccountNumber,
		input.RoleName,
	)

	sessionName := input.SessionName
	if len(sessionName) == 0 {
		sessionName = generateSessionName(input.RoleName, input.MFADeviceID)
	}

	stsInput := &sts.AssumeRoleInput{
		DurationSeconds: aws.Int64(int64(input.Duration)),
		RoleArn:         aws.String(roleArn),
		RoleSessionName: aws.String(sessionName),
		SerialNumber:    aws.String(input.MFADeviceID),
		TokenCode:       aws.String(input.TokenCode),
	}

	return svc.AssumeRole(stsInput)
}

type webOutInput struct {
	AWSAccessKeyID     string `required:"true"`
	AWSSecretAccessKey string `required:"true"`
	AccountName        string `required:"true"`
	AWSAccountNumber   string `required:"true"`
	RoleName           string `required:"true"`
	MFADeviceID        string `required:"true"`
	TokenCode          string `required:"true"`
	SessionName        string
	Duration           int
}

func webOut(input webOutInput) (string, error) {
	assumeInput := assumeRoleInput{
		AWSAccessKeyID:     input.AWSAccessKeyID,
		AWSSecretAccessKey: input.AWSSecretAccessKey,
		AWSAccountNumber:   input.AWSAccountNumber,
		RoleName:           input.RoleName,
		MFADeviceID:        input.MFADeviceID,
		TokenCode:          input.TokenCode,
		SessionName:        input.SessionName,
		Duration:           input.Duration,
	}

	result, err := assumeRole(assumeInput)
	if err != nil {
		return "", err
	}
	tmpCredentials := struct {
		SessionID    string `json:"sessionId"`
		SessionKey   string `json:"sessionKey"`
		SessionToken string `json:"sessionToken"`
	}{
		SessionID:    aws.StringValue(result.Credentials.AccessKeyId),
		SessionKey:   aws.StringValue(result.Credentials.SecretAccessKey),
		SessionToken: aws.StringValue(result.Credentials.SessionToken),
	}

	credentialsJson, err := json.Marshal(&tmpCredentials)
	if err != nil {
		return "", err
	}

	federationRequestParams := fmt.Sprintf(
		"?Action=getSigninToken&SessionDuration=%d&Session=%s",
		input.Duration,
		url.QueryEscape(string(credentialsJson)),
	)

	tokenResp, err := http.Get("https://signin.aws.amazon.com/federation" + federationRequestParams)
	if err != nil {
		return "", err
	}

	tokenRespObj := struct {
		SigninToken string `json:"SigninToken"`
	}{}

	responseBody, _ := ioutil.ReadAll(tokenResp.Body)
	if err := json.Unmarshal(responseBody, &tokenRespObj); err != nil {
		return "", err
	}

	signinRequestParams := fmt.Sprintf(
		"?Action=login&Issuer=aws-session-cli&Destination=%s&SigninToken=%s",
		url.QueryEscape("https://console.aws.amazon.com/"),
		tokenRespObj.SigninToken,
	)

	return "https://signin.aws.amazon.com/federation" + signinRequestParams, nil
}

type credentialsOutInput struct {
	AWSAccessKeyID     string `required:"true"`
	AWSSecretAccessKey string `required:"true"`
	AccountName        string `required:"true"`
	AWSAccountNumber   string `required:"true"`
	RoleName           string `required:"true"`
	MFADeviceID        string `required:"true"`
	TokenCode          string `required:"true"`
	SessionName        string
	Duration           int
	Region             string `required:"true"`
	UserShell          string
}

func credentialOut(input credentialsOutInput) (string, error) {
	assumeInput := assumeRoleInput{
		AWSAccessKeyID:     input.AWSAccessKeyID,
		AWSSecretAccessKey: input.AWSSecretAccessKey,
		AWSAccountNumber:   input.AWSAccountNumber,
		RoleName:           input.RoleName,
		MFADeviceID:        input.MFADeviceID,
		TokenCode:          input.TokenCode,
		SessionName:        input.SessionName,
		Duration:           input.Duration,
	}

	result, err := assumeRole(assumeInput)
	if err != nil {
		return "", err
	}

	expiration := strconv.FormatInt(result.Credentials.Expiration.Unix(), 10)
	tmplVariables := EnvVariables{
		AccountName:     input.AccountName,
		AccountID:       input.AWSAccountNumber,
		Region:          input.Region,
		AccessKeyID:     aws.StringValue(result.Credentials.AccessKeyId),
		TokenExpiration: expiration,
		SecretAccessKey: aws.StringValue(result.Credentials.SecretAccessKey),
		SessionToken:    aws.StringValue(result.Credentials.SessionToken),
	}

	// Set UserShell
	userShell := input.UserShell
	if len(userShell) == 0 {
		userShell = detectShell()
	}
	switch userShell {
	case "powershell":
		tmplVariables.Delimiter = " = '"
		tmplVariables.Prefix = "$env:"
		tmplVariables.Suffix = "'\n"
	case "docker":
		tmplVariables.Delimiter = "="
		tmplVariables.Prefix = " -e "
		tmplVariables.Suffix = ""
	case "cmd":
		tmplVariables.Delimiter = "="
		tmplVariables.Prefix = "set "
		tmplVariables.Suffix = "\n"
	default:
		tmplVariables.Delimiter = "="
		tmplVariables.Prefix = "export "
		tmplVariables.Suffix = "\n"
	}

	var tmpl *template.Template
	var buffer bytes.Buffer
	tmpl = template.Must(template.New("envVariables").Parse(envTemplate()))
	if err := tmpl.Execute(&buffer, tmplVariables); err != nil {
		return "", err
	}

	return buffer.String(), nil
}

package main

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"
)

func PromptDefaultInputs() bool {
	promptConfirm := PromptEnvDetails()

	if len(migrationReq.SecretScope) == 0 {
		promptConfirm = true
		migrationReq.SecretScope = SelectInput("Scope for secrets & secret managers:", scopes, Project)
	}

	if len(migrationReq.ConnectorScope) == 0 {
		promptConfirm = true
		migrationReq.ConnectorScope = SelectInput("Scope for connectors:", scopes, Project)
	}

	if len(migrationReq.TemplateScope) == 0 {
		promptConfirm = true
		migrationReq.TemplateScope = SelectInput("Scope for templates:", scopes, Project)
	}

	return promptConfirm
}

func PromptEnvDetails() bool {
	promptConfirm := false

	if len(migrationReq.Environment) == 0 {
		promptConfirm = true
		migrationReq.Environment = SelectInput("Which environment?", []string{Dev, QA, Prod, Prod3}, Dev)
	}

	// Check if auth is provided. If not provided then request for one
	if len(migrationReq.Auth) == 0 {
		migrationReq.Auth = TextInput("The environment variable 'HARNESS_MIGRATOR_AUTH' is not set. What is the api key?")
	}

	if migrationReq.Environment == "Prod" && len(migrationReq.UrlCG) != 0 && len(migrationReq.UrlNG) != 0 {
		if !ParseNGUrl() {
			if len(migrationReq.Account) == 0 {
				promptConfirm = true
				migrationReq.Account = TextInput("Account that you wish to migrate:")
			}
			PromptOrgAndProject([]string{migrationReq.SecretScope, migrationReq.ConnectorScope, migrationReq.TemplateScope})
		}
		if !ParseCGUrl() {
			if len(migrationReq.AppId) == 0 {
				promptConfirm = true
				migrationReq.AppId = TextInput("AppId that you wish to migrate:")
			}
		}
	}

	if migrationReq.Environment == "Dev" || migrationReq.AllowInsecureReq {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	if len(migrationReq.Account) == 0 {
		promptConfirm = true
		migrationReq.Account = TextInput("Account that you wish to migrate:")
	}
	return promptConfirm
}

func PromptOrgAndProject(scope []string) bool {
	promptConfirm := false
	promptOrg := len(migrationReq.OrgIdentifier) == 0 && ContainsAny(scope, []string{Org, Project})
	promptProject := len(migrationReq.ProjectIdentifier) == 0 && ContainsAny(scope, []string{Project})

	if promptOrg {
		promptConfirm = true
		migrationReq.OrgIdentifier = TextInput("Which Org?")
	}
	if promptProject {
		promptConfirm = true
		migrationReq.ProjectIdentifier = TextInput("Which Project?")
	}
	return promptConfirm
}

func ParseNGUrl() bool {
	re := regexp.MustCompile(`https:\/\/.*\.harness\.io/ng/#/account/([a-zA-Z0-9-]+)/.*/orgs/([a-zA-Z0-9_]+)/projects/([a-zA-Z0-9_]+)/.*`)
	u, err := url.Parse(migrationReq.UrlNG)
	if err != nil {
		log.Warning(err)
		return false
	}

	if !re.MatchString(u.String()) {
		log.Warning("NG URL did not match the expected format, skipping...")
		return false
	} else {
		u, _ := url.Parse(migrationReq.UrlNG)
		fragment := u.Fragment

		migrationReq.Account = strings.Split(fragment, "/")[2]
		migrationReq.OrgIdentifier = strings.Split(fragment, "/")[5]
		migrationReq.ProjectIdentifier = strings.Split(fragment, "/")[7]
		migrationReq.ProjectName = strings.Split(fragment, "/")[7]
		return true
	}

}

func ParseCGUrl() bool {
	re := regexp.MustCompile(`https:\/\/.*\.harness\.io\/#\/account\/[a-zA-Z0-9-]+\/app\/[a-zA-Z0-9-_]+/.*`)
	u, err := url.Parse(migrationReq.UrlCG)
	if err != nil {
		log.Warning(err)
		return false
	}

	if !re.MatchString(u.String()) {
		log.Warning("CG URL did not match the expected format, skipping...")
		return false
	} else {
		u, _ := url.Parse(migrationReq.UrlCG)
		fragment := u.Fragment

		migrationReq.AppId = strings.Split(fragment, "/")[4]
		return true
	}
}

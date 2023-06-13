package request

import (
	"regexp"
)

const (
	// Runtime options: NoPHP, PHP71FPM, PHP74FPM, PHP81FPM

	// RuntimeNoPHP is the NoPHP runtime option.
	RuntimeNoPHP = "NoPHP"
	// RuntimePHP71FPM is the PHP71FPM runtime option.
	RuntimePHP71FPM = "PHP71FPM"
	// RuntimePHP74FPM is the PHP74FPM runtime option.
	RuntimePHP74FPM = "PHP74FPM"
	// RuntimePHP81FPM is the PHP81FPM runtime option.
	RuntimePHP81FPM = "PHP81FPM"

	// Database options: no, mysql

	// DatabaseNo is the no database option.
	DatabaseNo = "no"
	// DatabaseMySQL is the mysql database option.
	DatabaseMySQL = "mysql"

	// EmailRegex is the regex for the email validation.
	EmailRegex = `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`
)

var (
	// RuntimeOptions is the list of the available runtime options.
	RuntimeOptions = []string{RuntimeNoPHP, RuntimePHP71FPM, RuntimePHP74FPM, RuntimePHP81FPM}
	// DatabaseOptions is the list of the available database options.
	DatabaseOptions = []string{DatabaseNo, DatabaseMySQL}
)

// CreateProjectRequest represents the request body of the /project/create endpoint.
type CreateProjectRequest struct {
	Client     string `json:"project-client"`
	Command    string `json:"command"`
	Database   string `json:"project-database"`
	Name       string `json:"project-name"`
	OwnerEmail string `json:"project-owner-email"`
	Runtime    string `json:"project-runtime"`
	Staging    string `json:"env-staging"`
	Production string `json:"env-production"`
}

// Validate is for the data validation.
// It returns the array of the errors.
// If the data is fine, it returns empty arra.
func (r *CreateProjectRequest) Validate() map[string][]string {
	validationErrorMap := make(map[string][]string)
	// At least one environment has to be selected.
	if r.Staging == "false" && r.Production == "false" {
		var validationErrors []string
		validationErrors = append(validationErrors, "At least one environment has to be selected")
		validationErrorMap["env-staging"] = validationErrors
		validationErrorMap["env-production"] = validationErrors
	}
	// Name validation. It has to be 3-30 character long.
	lenName := len(r.Name)
	if lenName < 3 || lenName > 30 {
		var validationErrors []string
		validationErrors = append(validationErrors, "Invalid project name length: it has to be 3-30 char")
		validationErrorMap["project-name"] = validationErrors
	}
	// Client validation, it has to be 3-30 character long.
	lenClient := len(r.Client)
	if lenClient < 3 || lenClient > 30 {
		var validationErrors []string
		validationErrors = append(validationErrors, "Invalid project client length: it has to be 3-30 char")
		validationErrorMap["project-client"] = validationErrors
	}
	// Owner email validation.
	// It has to be a valid email address format.
	regexpEmail := regexp.MustCompile(EmailRegex)
	if regexpEmail.MatchString(r.OwnerEmail) == false {
		var validationErrors []string
		validationErrors = append(validationErrors, "Invalid project owner email: "+r.OwnerEmail)
		validationErrorMap["project-owner-email"] = validationErrors
	}
	// Runtime validation.
	// It has to be one of the available runtime options.
	var runtimeValid bool
	for _, runtime := range RuntimeOptions {
		if runtime == r.Runtime {
			runtimeValid = true
			break
		}
	}
	if runtimeValid == false {
		var validationErrors []string
		validationErrors = append(validationErrors, "Invalid project runtime")
		validationErrorMap["project-runtime"] = validationErrors
	}
	// Database validation.
	// It has to be one of the available database options.
	var databaseValid bool
	for _, database := range DatabaseOptions {
		if database == r.Database {
			databaseValid = true
			break
		}
	}
	if databaseValid == false {
		var validationErrors []string
		validationErrors = append(validationErrors, "Invalid project database")
		validationErrorMap["project-database"] = validationErrors
	}
	return validationErrorMap
}

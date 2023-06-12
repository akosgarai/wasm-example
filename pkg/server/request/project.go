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
	Name       string `json:"project-name"`
	Client     string `json:"project-client"`
	OwnerEmail string `json:"project-owner-email"`
	Runtime    string `json:"project-runtime"`
	Database   string `json:"project-database"`
}

// Validate is for the data validation.
// It returns the array of the errors.
// If the data is fine, it returns empty arra.
func (r *CreateProjectRequest) Validate() []string {
	var validationErrors []string
	// Name validation. It has to be 3-30 character long.
	lenName := len(r.Name)
	if lenName < 3 || lenName > 30 {
		validationErrors = append(validationErrors, "Invalid project name length: it has to be 3-30 char")
	}
	// Client validation, it has to be 3-30 character long.
	lenClient := len(r.Client)
	if lenClient < 3 || lenClient > 30 {
		validationErrors = append(validationErrors, "Invalid project client length: it has to be 3-30 char")
	}
	// Owner email validation.
	// It has to be a valid email address format.
	regexpEmail := regexp.MustCompile(EmailRegex)
	if regexpEmail.MatchString(r.OwnerEmail) == false {
		validationErrors = append(validationErrors, "Invalid project owner email: "+r.OwnerEmail)
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
		validationErrors = append(validationErrors, "Invalid project runtime")
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
		validationErrors = append(validationErrors, "Invalid project database")
	}
	return validationErrors
}

package request

import (
	"regexp"

	"github.com/akosgarai/wasm-example/pkg/server/models"
	"gorm.io/gorm"
)

const (
	// EmailRegex is the regex for the email validation.
	EmailRegex = `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`
)

// CreateProjectRequest represents the request body of the /project/create endpoint.
type CreateProjectRequest struct {
	Client      string               `json:"project-client"`
	Command     string               `json:"command"`
	Database    int                  `json:"project-database"`
	Name        string               `json:"project-name"`
	OwnerEmail  string               `json:"project-owner-email"`
	Runtime     int                  `json:"project-runtime"`
	Environment []ProjectEnvironment `json:"project-environment"`
	// The following fields are not part of the request body.
	ClientID         uint
	ProjectID        uint
	validationErrors map[string][]string
}

// ProjectEnvironment represents the environment of the project.
type ProjectEnvironment struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Validate is for the data validation.
// It returns the array of the errors.
// If the data is fine, it returns empty arra.
func (r *CreateProjectRequest) Validate(db *gorm.DB) map[string][]string {
	r.validateEnvironment(db)
	// Name validation. It has to be 3-30 character long.
	r.validateName(db)
	// Client validation, it has to be 3-30 character long.
	r.validateClient(db)
	// Owner email validation.
	r.validateOwnerEmail()
	// Runtime validation.
	r.validateRuntime(db)
	// Database validation.
	r.validateDatabase(db)
	return r.validationErrors
}

// project name validation. It has to be 3-30 character long.
func (r *CreateProjectRequest) validateName(db *gorm.DB) {
	lenName := len(r.Name)
	if lenName < 3 || lenName > 30 {
		r.addValidationError("project-name", "Invalid project name length: it has to be 3-30 char")
	}
	// setup project id. if the project is not in the database, it will be created.
	var project models.Project
	project.Name = r.Name
	if err := db.Where("name = ?", r.Name).FirstOrCreate(&project).Error; err != nil {
		r.addValidationError("project-name", "Invalid project name")
		return
	}
	r.ProjectID = project.ID
}

// client validation. It has to be 3-30 character long.
func (r *CreateProjectRequest) validateClient(db *gorm.DB) {
	lenClient := len(r.Client)
	if lenClient < 3 || lenClient > 30 {
		r.addValidationError("project-client", "Invalid project client length: it has to be 3-30 char")
		return
	}
	// setup client id. if the client is not in the database, it will be created.
	var client models.Client
	client.Name = r.Client
	if err := db.Where("name = ?", r.Client).FirstOrCreate(&client).Error; err != nil {
		r.addValidationError("project-client", "Invalid project client")
		return
	}
	r.ClientID = client.ID
}

// owner email validation. It has to be a valid email address format.
func (r *CreateProjectRequest) validateOwnerEmail() {
	regexpEmail := regexp.MustCompile(EmailRegex)
	if regexpEmail.MatchString(r.OwnerEmail) == false {
		r.addValidationError("project-owner-email", "Invalid project owner email: "+r.OwnerEmail)
	}
}

// runtime validation. It has to be one of the available runtime id options stored in the database.
func (r *CreateProjectRequest) validateRuntime(db *gorm.DB) {
	var runtime models.Runtime
	if err := db.Where("id = ?", r.Runtime).First(&runtime).Error; err != nil {
		r.addValidationError("project-runtime", "Invalid project runtime")
	}
}

// database validation. It has to be one of the available database id options stored in the database.
func (r *CreateProjectRequest) validateDatabase(db *gorm.DB) {
	var database models.Dbtype
	if err := db.Where("id = ?", r.Database).First(&database).Error; err != nil {
		r.addValidationError("project-database", "Invalid project database")
	}
}

// environment validation. It has to be at least one environment selected.
func (r *CreateProjectRequest) validateEnvironment(db *gorm.DB) {
	if len(r.Environment) == 0 {
		// At least one environment has to be selected.
		environments := []models.Environment{}
		db.Find(&environments)
		for _, env := range environments {
			r.addValidationError("env-"+env.Name, "At least one environment has to be selected")
		}
	}
	// environment validation
	for _, environment := range r.Environment {
		var env models.Environment
		if db.First(&env, environment.ID).RowsAffected == 0 {
			r.addValidationError("env-"+environment.Name, "Invalid environment id")
		}
	}
}

// Extend the error list with the new error.
func (r *CreateProjectRequest) addValidationError(key string, error string) {
	if r.validationErrors == nil {
		r.validationErrors = make(map[string][]string)
	}
	r.validationErrors[key] = append(r.validationErrors[key], error)
}

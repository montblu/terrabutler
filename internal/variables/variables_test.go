package variables

import (
	"testing"

	"github.com/l58193/terrabutler/tree/Rewrite-Go/terrabutler/internal/settings"
	"github.com/l58193/terrabutler/tree/Rewrite-Go/terrabutler/internal/utils"

	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func TestGeneratePassword(t *testing.T) {

	//Define password size
	size := 25
	password := generate_password(size)
	assert.EqualValues(t, size, len(password), "Failed, the password its not being generated at the required size.")

}

func TestGenerateVarFiles(t *testing.T) {

	//Name of the temporary environment
	env := "tempEnv"
	org := "org"
	sites := []string{"site1", "site2"}
	firebase_credentials := "DUMMY"
	mail_password := "DUMMY"

	//Define the settings file to get organization, sites, firebase_credentials and mail_password, all the settings defined to be used in the .j2 file
	settings.Conf.Set("general.organization", org)
	settings.Conf.Set("sites.ordered", sites)
	settings.Conf.Set("environments.temporary.secrets.firebase_credentials", firebase_credentials)
	settings.Conf.Set("environments.temporary.secrets.mail_password", mail_password)

	//Define the path templates and the template files (The env and a another one)
	variablesPath := "WRITE_PATH"
	utils.Paths["variables"] = variablesPath
	utils.Paths["templates"] = "PATH"
	envTemplate := utils.Paths["templates"] + "/env.j2"
	template := utils.Paths["templates"] + "/template.j2"

	templateText := `environment = "{{env}}"
sites_list    = [
{%+ for site in sites -%}
  "{{ site }}",
{%+ endfor -%}
]
firebaseCredentials = "{{firebase_credentials}}"
mailPassword = "{{mail_password}}"`

	generatedText := `environment = "` + env + `"
sites_list    = [
"site1",
"site2",
]
firebaseCredentials = "` + firebase_credentials + `"
mailPassword = "` + mail_password + `"`

	// Use the in-memory filesystem
	fs := afero.NewMemMapFs()

	afero.WriteFile(fs, envTemplate, []byte(templateText), 0644)
	afero.WriteFile(fs, template, []byte(templateText), 0644)

	//Call the function
	assert.NoError(t, Generate_var_files(env, fs), "Failed, the test wasn't supposed to fail, a valid template was provided.")

	envTemplateData, err := afero.ReadFile(fs, variablesPath+"/"+org+"-"+env+".tfvars")
	assert.NoError(t, err, "Failed, An error occurred while getting the data from the generated env template.")

	templateData, err := afero.ReadFile(fs, variablesPath+"/"+org+"-"+env+"-"+"template"+".tfvars")
	assert.NoError(t, err, "Failed, An error occurred while getting the data from the generated template.")

	assert.Equal(t, string(generatedText), string(envTemplateData), "Failed, the env template file was not written correctly.")
	assert.Equal(t, string(generatedText), string(templateData), "Failed, the template file was not written correctly.")
}

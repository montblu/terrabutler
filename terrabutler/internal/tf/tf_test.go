package tf

import (
	"terrabutler/internal/settings"
	"terrabutler/internal/utils"
	"testing"

	"github.com/stretchr/testify/assert"
)

//It wont be tested the execution of terraform commands, because its inviable to test a running executable.
//Creating mocks doesn't represent the "true output" of the execution of the command

// It only be tested the Building of the terraform commands and if its correct.
// If this passes the execution of terraform will also be correct, external factors, like not having terraform installed are reported as a error in the program.

// This test will also validate the NeedOptionsBuilder function
func TestArgsPrint(t *testing.T) {

	//Define the variables begin used
	site := "site"
	current_env = "env"
	org := "org"
	default_env := "default"
	backend_dir := "BACKEND_PATH"
	variables_dir := "VARIABLES_PATH"

	//Defining the global constants used
	settings.Conf.Set("general.organization", org)
	settings.Conf.Set("environments.default.name", default_env)
	utils.Paths["backends"] = backend_dir
	utils.Paths["variables"] = variables_dir

	//Command Init on site inception
	initArgsInception := ArgsPrint("init", "inception")
	//Valid Output
	initArgsInceptionOutput := "-backend-config " + backend_dir + "/" + org + "-" + default_env + "-inception.tfvars"

	//Command Init on another site
	initArgsSite := ArgsPrint("init", site)
	initArgsSiteOutput := "-backend-config " + backend_dir + "/" + org + "-" + current_env + "-" + site + ".tfvars"

	//Command Plan
	planArgs := ArgsPrint("plan", site)
	//Command Apply
	applyArgsSite := ArgsPrint("apply", site)

	// Plan and Apply should have the same output
	planAndApplyArgsOutput := "-var-file " + variables_dir + "/global.tfvars" +
		" -var-file " + variables_dir + "/" + org + "-" + current_env + ".tfvars" +
		" -var-file " + variables_dir + "/" + org + "-" + current_env + "-" + site + ".tfvars"

	assert.Equal(t, initArgsInceptionOutput, initArgsInception, "Failed building arguments for init, the argument created for the backend option for the site inception is invalid.")
	assert.Equal(t, initArgsSiteOutput, initArgsSite, "Failed building arguments for init, the argument created for the backend option for a generic site is invalid.")
	assert.Equal(t, planAndApplyArgsOutput, planArgs, "Failed building arguments for plan, the arguments created for the var option for a generic site is invalid.")
	assert.Equal(t, planAndApplyArgsOutput, applyArgsSite, "Failed building arguments for apply, the arguments created for the var option for a generic site is invalid.")

}

func TestCommandBuilder(t *testing.T) {
	//A valid structure for a terraform command
	validTerraformCommand1 := CommandBuilder("cmd", "site", []string{"arg1", "arg2"}, []string{"flag1", "flag2"}, "")
	validOutput1 := []string{"terraform", "cmd", "flag1", "flag2", "arg1", "arg2"}

	//A valid structure for a terraform command with a subcommand, where the subcommand should be split
	validTerraformCommand2 := CommandBuilder("cmd subcommand", "site", []string{"arg1", "arg2"}, []string{"flag1", "flag2"}, "")
	validOutput2 := []string{"terraform", "cmd", "subcommand", "flag1", "flag2", "arg1", "arg2"}

	assert.Equal(t, validOutput1, validTerraformCommand1, "Failed, the first command build is not valid.")
	assert.Equal(t, validOutput2, validTerraformCommand2, "Failed, the second command build is not valid.")

}

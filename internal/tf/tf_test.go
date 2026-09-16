package tf

import (
	"errors"
	"os"
	"sync"
	"testing"

	"github.com/montblu/terrabutler/internal/settings"
	"github.com/montblu/terrabutler/internal/utils"
	"github.com/spf13/afero"

	"github.com/stretchr/testify/assert"
)

// It wont be tested the execution of terraform commands, because its inviable to test a running executable.
// Creating mocks doesn't represent the "true output" of the execution of the command

// It only be tested the Building of the terraform commands and if its correct.
// If this passes the execution of terraform will also be correct, external factors, like not having terraform installed are reported as a error in the program.

// This test will also validate the NeedOptionsBuilder function
func TestArgsPrint(t *testing.T) {

	fs := afero.NewMemMapFs()

	// Define the variables begin used
	site := "site"
	utils.SetCurrentEnvForTest("env")
	org := "org"
	default_env := "default"
	backend_dir := "BACKEND_PATH"
	variables_dir := "VARIABLES_PATH"

	// Defining the global constants used
	_ = settings.Conf.Set("general.organization", org)
	_ = settings.Conf.Set("environments.default.name", default_env)
	utils.Paths["backends"] = backend_dir
	utils.Paths["variables"] = variables_dir

	// Command Init on site inception
	initArgsInception := ArgsPrint("init", "inception", fs)
	// Valid Output
	initArgsInceptionOutput := "-backend-config " + backend_dir + "/" + org + "-" + default_env + "-inception.tfvars"

	// Command Init on another site
	initArgsSite := ArgsPrint("init", site, fs)
	initArgsSiteOutput := "-backend-config " + backend_dir + "/" + org + "-" + "env" + "-" + site + ".tfvars"

	// Command Plan
	planArgs := ArgsPrint("plan", site, fs)
	// Command Apply
	applyArgsSite := ArgsPrint("apply", site, fs)

	// Plan and Apply should have the same output
	planAndApplyArgsOutput := "-var-file " + variables_dir + "/global.tfvars" +
		" -var-file " + variables_dir + "/" + org + "-" + "env" + ".tfvars" +
		" -var-file " + variables_dir + "/" + org + "-" + "env" + "-" + site + ".tfvars"

	assert.Equal(t, initArgsInceptionOutput, initArgsInception, "Failed building arguments for init, the argument created for the backend option for the site inception is invalid.")
	assert.Equal(t, initArgsSiteOutput, initArgsSite, "Failed building arguments for init, the argument created for the backend option for a generic site is invalid.")
	assert.Equal(t, planAndApplyArgsOutput, planArgs, "Failed building arguments for plan, the arguments created for the var option for a generic site is invalid.")
	assert.Equal(t, planAndApplyArgsOutput, applyArgsSite, "Failed building arguments for apply, the arguments created for the var option for a generic site is invalid.")

}

func TestCommandBuilder(t *testing.T) {

	fs := afero.NewMemMapFs()

	// A valid structure for a terraform command
	validTerraformCommand1 := CommandBuilder("cmd", "site", []string{"arg1", "arg2"}, []string{"flag1", "flag2"}, "", fs)
	validOutput1 := []string{"terraform", "cmd", "flag1", "flag2", "arg1", "arg2"}

	// A valid structure for a terraform command with a subcommand, where the subcommand should be split
	validTerraformCommand2 := CommandBuilder("cmd subcommand", "site", []string{"arg1", "arg2"}, []string{"flag1", "flag2"}, "", fs)
	validOutput2 := []string{"terraform", "cmd", "subcommand", "flag1", "flag2", "arg1", "arg2"}

	assert.Equal(t, validOutput1, validTerraformCommand1, "Failed, the first command build is not valid.")
	assert.Equal(t, validOutput2, validTerraformCommand2, "Failed, the second command build is not valid.")

}

func TestTerraformUserAgent(t *testing.T) {
	// userAgent is package-level global state; restore it so this test does not
	// leak into others.
	original := *userAgent.Load()
	defer func() { userAgent.Store(&original) }()

	// The default (used when the build version is unknown) carries the product
	def := apnUserAgent + " terrabutler"
	userAgent.Store(&def)
	assert.Contains(t, TerraformEnv(), "TF_APPEND_USER_AGENT="+apnUserAgent+" terrabutler",
		"TerraformEnv should append the default user agent when no version is set.")

	// SetUserAgent tags the user agent with the build version.
	SetUserAgent("1.2.3")
	assert.Equal(t, apnUserAgent+" terrabutler/1.2.3", *userAgent.Load(),
		"SetUserAgent should build a versioned user agent with the APN attribution tag.")

	env := TerraformEnv()
	assert.Contains(t, env, "TF_APPEND_USER_AGENT="+apnUserAgent+" terrabutler/1.2.3",
		"TerraformEnv should append the versioned TF_APPEND_USER_AGENT variable.")
	assert.Equal(t, "APN_1.1/pc_2me4418zjym9qcbpo4pbtq8rl$", apnUserAgent,
		"The APN attribution tag must match the partner-assigned value.")

	// The parent environment must still be inherited, not replaced.
	t.Setenv("TERRABUTLER_UA_INHERIT_CHECK", "present")
	assert.Contains(t, TerraformEnv(), "TERRABUTLER_UA_INHERIT_CHECK=present",
		"TerraformEnv should inherit the parent environment, not overwrite it.")

	// TerraformEnv must not mutate os.Environ(); it should only add one entry.
	assert.Equal(t, len(os.Environ())+1, len(TerraformEnv()),
		"TerraformEnv should add exactly one variable to the inherited environment.")
}

func TestInitAllSitesSuccess(t *testing.T) {
	var mu sync.Mutex
	callCount := 0
	calledSites := []string{}

	commandRunnerNoVisibleOutputVar = func(command, site string, args, options []string, needed_options string, fs afero.Fs) ([]byte, error) {
		mu.Lock()
		callCount++
		calledSites = append(calledSites, site)
		mu.Unlock()
		return nil, nil
	}

	settings.Conf.Set("sites.ordered", []string{"inception", "site-a", "site-b", "site-c"}) //nolint:errcheck

	fs := afero.NewMemMapFs()

	// Creating the environment file, with the new environment
	// This is done because if the environments files from the sites fails to be created
	// The values in fileSiteX will be ""
	// And this to test this, the utils.GetCurrentEnv() should return a value different from ""
	newEnv := "newEnv"
	utils.Paths["environment"] = "ROOT/site_inception/.terraform/environment"
	_ = afero.WriteFile(fs, utils.Paths["environment"], []byte(newEnv), 0644)

	// Environment that will be set in the file environment of the sites
	env := "oldEnv"
	// Environment that will be set in the file environment of the sites if those succeed
	newEnv = utils.GetCurrentEnv(fs)

	err := InitAllSites(env, fs)

	//Get each value for the sites to check if the environment was set correctly
	fileSiteA, _ := afero.ReadFile(fs, utils.Paths["root"]+"/site_site-a/.terraform/.terrabutler_env")
	fileSiteB, _ := afero.ReadFile(fs, utils.Paths["root"]+"/site_site-b/.terraform/.terrabutler_env")
	fileSiteC, _ := afero.ReadFile(fs, utils.Paths["root"]+"/site_site-c/.terraform/.terrabutler_env")

	assert.NoError(t, err, "Init should not fail")
	assert.Equal(t, 3, callCount, "Should init 3 sites (inception removed)")
	assert.Equal(t, 3, len(calledSites), "All 3 sites should be called")
	assert.NotContains(t, calledSites, "inception", "Inception should be filtered out")
	//Test the value in the files
	assert.Equal(t, newEnv, string(fileSiteA), "site-a should have a .terrabutler_environment file with "+newEnv)
	assert.Equal(t, newEnv, string(fileSiteB), "site-b should have a .terrabutler_environment file with "+newEnv)
	assert.Equal(t, newEnv, string(fileSiteC), "site-c should have .terrabutler_environment file with "+newEnv)
}

func TestInitAllSitesWithErrors(t *testing.T) {
	commandRunnerNoVisibleOutputVar = func(command, site string, args, options []string, needed_options string, fs afero.Fs) ([]byte, error) {
		if site == "site-b" {
			return nil, errors.New("backend unavailable")
		}
		return nil, nil
	}

	settings.Conf.Set("sites.ordered", []string{"site-a", "site-b", "site-c"}) //nolint:errcheck

	fs := afero.NewMemMapFs()

	newEnv := "newEnv"
	utils.Paths["environment"] = "ROOT/site_inception/.terraform/environment"
	_ = afero.WriteFile(fs, utils.Paths["environment"], []byte(newEnv), 0644)

	// Environment that will be set in the file environment of the sites if those fail
	oldEnv := "oldEnv"
	// Environment that will be set in the file environment of the sites if those succeed
	newEnv = utils.GetCurrentEnv(fs)

	err := InitAllSites(oldEnv, fs)

	//Get each value for the sites to check if the environment was set correctly
	fileSiteA, _ := afero.ReadFile(fs, utils.Paths["root"]+"/site_site-a/.terraform/.terrabutler_env")
	fileSiteB, _ := afero.ReadFile(fs, utils.Paths["root"]+"/site_site-b/.terraform/.terrabutler_env")
	fileSiteC, _ := afero.ReadFile(fs, utils.Paths["root"]+"/site_site-c/.terraform/.terrabutler_env")

	assert.Error(t, err, "Init should return error when some sites fail")
	assert.Contains(t, err.Error(), "1/3 sites failed", "Error should indicate failure count")
	//Test the value in the files
	assert.Equal(t, newEnv, string(fileSiteA), "site-a should have a .terrabutler_environment file with "+newEnv)
	assert.Equal(t, oldEnv, string(fileSiteB), "site-b should have a .terrabutler_environment file with "+oldEnv)
	assert.Equal(t, newEnv, string(fileSiteC), "site-c should have .terrabutler_environment file with "+newEnv)
}

func TestInitAllSitesEmptyList(t *testing.T) {
	// Only inception in the list — after filtering, nothing left
	settings.Conf.Set("sites.ordered", []string{"inception"}) //nolint:errcheck

	fs := afero.NewMemMapFs()
	env := "env"

	err := InitAllSites(env, fs)
	assert.NoError(t, err, "Should not error when no sites to init")
}

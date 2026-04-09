package main

import (
	"context"
	"fmt"
	"log"

	"github.com/hashicorp/go-version"
	"github.com/hashicorp/hc-install/product"
	"github.com/hashicorp/hc-install/releases"
	"github.com/hashicorp/terraform-exec/tfexec"
	"go.uber.org/zap"
)

// Initializes Terraform Globally
var tf = init_terraform()

// Initializes Terraform
func init_terraform() *tfexec.Terraform {
	//This downloads the terraform locally to be used
	installer := &releases.ExactVersion{
		Product: product.Terraform,
		Version: version.Must(version.NewVersion("1.10.5")),
	}

	execPath, err := installer.Install(context.Background())
	if err != nil {
		log.Fatalf("error installing Terraform: %s", err)
	}

	fmt.Println(execPath)

	//For better performance, and not always downloading the terraform when using terrabutler, using the locally installed terraform installed by mise.
	//It's needed to make a var env of this for the final stage
	workingDir := paths["inception"]
	tf, err := tfexec.NewTerraform(workingDir, execPath)
	if err != nil {
		log.Fatalf("error running NewTerraform: %s", err)
	}

	return tf
}

//Commands that required special options

func tf_init(site string, args []string, options []string, needed_options string) {

}

func tf_plan(site string, args []string, options []string, needed_options string) {

}

func tf_apply(site string, args []string, options []string, needed_options string) {

}

//All the other tf commands used

func tf_console(site string, options []string) {}

func tf_destroy(site string, options []string) {}

func tf_fmt(site string, options []string) {}

func tf_force_unlock(site string, options []string, lock_id string) {}

func tf_generate_options(site string, options []string, choice string) {}

func tf_import(site string, options []string, address string, id string) {}

func tf_output(site string, options []string) {}

func tf_providers_lock(site string, options []string, providers []string) {}

func tf_providers_mirror(site string, options []string, target_dir string) {}

func tf_providers_schema(site string, options []string) {}

func tf_refresh(site string, options []string) {}

func tf_show(site string, options []string, path string) {}

// Not included in the tfexec
func tf_state_list() {}

func tf_state_mv() {}

func tf_state_pull() {}

func tf_state_push() {}

// Not included in the tfexec
func tf_state_replace_provider() {}

func tf_state_rm() {}

// Not included in the tfexec
func tf_state_show() {}

func tf_taint(site string, options []string, address string) {}

func tf_untaint(site string, options []string, address string) {}

func tf_validate(site string, options []string) {

	/*json, err := tf.Validate(context.Background())*/
	/*if err != nil {
		logger.Error("An error has occurred", zap.Error(err))
	}*/

}

func tf_version(json_option bool) {

	version, _, err := tf.Version(context.Background(), false)
	if err != nil {
		logger.Error("An error has occurred", zap.Error(err))
	}
	//There exists a version plain text and Json Format, but its private for some reason?
	if json_option == true {
	}

	logger.Info("TerraForm Version : " + version.String())
}

// New commands to be used in all sites

func tf_destroy_all_sites() {}

func tf_apply_all_sites() {}

func tf_init_all_sites() {}

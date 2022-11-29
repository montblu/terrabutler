# Create a new project

Before proceeding make sure that you have followed the [Installation](installation.md). 

## Configure direnv

For direnv to work properly it needs to be hooked into the shell.
Check direnv docs to know more: https://direnv.net/docs/hook.html

Make sure that you have followed direnv documentation before proceeding. 


## Configure Terrabutler

### Create your project folder

```shell
mkdir <project_name>
```

### Download Template

Start by downloading the Terrabutler Template Project from the repository source below: 

[![Version-shield]](https://github.com/lucascanero/terrabutler-template/archive/refs/heads/example-template.zip)

Copy the files inside `./terrabutler-template/` to the root of your project folder.

```shell
$ cp -a /terrabutler-template/. /<project_name>/
```

### Create a new workspace

Before configuring terrabutler, inside `<project_name>/site_inception` folder, you will need to create a Terraform Workspace: 
For example, we are gonna call it "staging"

```shell
$ cd site_inception
$ terraform workspace new staging
```

### Change Variables

Run the script `./<project_name>/config_template.sh`, with the following arguments, located inside the project folder root:

```shell
$ ./config_template.sh -d <domain> -e <environment_name> -p <project_name>

USAGE:
   ./config_template [FLAG] [STRING]

FLAGS:
   -o <organization_name>   The name for your organization.  
                            Example: -o example

   -d <domain>              The domain of your organization. 
                            Example: -d example.com

   -e <environment_name>    The environment name of your organization. 
                            Example: -e staging
```


???+ danger
    This script only works with the template folder! Don't use it in another project folders!
    `./terrabutler-template/config_template.sh`

## Terraform

Start by installing Terraform with tfenv:

```shell
$ tfenv install 
```

### Initialize the Project
 
Perform an Terraform Initialization inside site_inception:

```shell
$ cd /site_inception/
$ terraform init
```
Perform an Terraform Apply with the `.tfvars` inside `/config/variables/`:

`./configs/variables/global.tfvars` </br>
`./configs/variables/<project_name>-<environment_name>.tfvars"` </br>
`./configs/variables/<project_name>-<environment_name>-inception.tfvars"`</br>

```shell
$ terraform apply -var-file="../configs/variables/global.tfvars" -var-file="../configs/variables/<project_name>-<environment_name>.tfvars" -var-file="../configs/variables/<project_name>-<environment_name>-inception.tfvars"
```

### Change local to remote backend

#### Uncomment the remote backend line
Remove the commented line as below, in the `terrabutler-template` in the path `./site-inception/terraform.tf` to change from local to remote.

```
backend "s3" {}
```

Perform an Terraform initialization with the backend config file,located in `/configs/backend/` to update the new changes:

```shell
$ terraform init -backend-config=<inception_backend_path>  
```

???+ tip
   Example:
   $ terraform init -backend-config="./configs/backends/<project_name>-<environment_name>-inception.tfvars"

#### Delete the local state

Delete the local state in `/site_inception`, as it is not necessary:

```shell
$ rm -rf terraform.tfstate.d/ terraform.tfstate terraform.tfstate.backup
```

### Terrabutler commands

#### Terrabutler Plan

Perform an Terrabutler Plan to ensure everything is ok, and plan any changes:

```shell
$ terrabutler tf -site inception plan
```
#### Terrabutler Apply

Perform an Terrabutler Apply to apply the planned changes:

```shell
$ terrabutler tf -site inception apply
```
#### Terrabutler Destroy

Perform an Terrabutler Destroy to delete or undo any changes made in the provider:

```shell 
$ terrabutler tf -site inception destroy
```


 [Version-shield]: https://img.shields.io/badge/terrabutler_Template-Download-%23121011.svg?style=for-the-badge&logo=github&colorA=273133&colorB=0093ee "Latest version"

 

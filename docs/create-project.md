# Create Project

Before proceeding make sure that you have followed the [Installation](installation.md). 

## Configure direnv

For direnv to work properly it needs to be hooked into the shell.

Make sure that you have hooked direnv into your bash shell before proceeding. 

If not, add the following line at the end of your `~/.bashrc` file:

```
eval "$(direnv hook bash)"
```

## Add direnv to your project

Create a `.envrc` file in the project root directory, containing the following lines:

```
export TERRABUTLER_ENABLE=TRUE
export TERRABUTLER_ROOT=$(pwd)
```

## Configure Terrabutler

Start by downloading the Terrabutler Template Project from the repository release below: 

[![Version-shield]](https://raw.githubusercontent.com/lucascanero/terrabutler-template/***********)

Inside the terrabutler-template.zip, downloaded before, in the `~/configs/` folder, copy the template file `settings.yml`, and edit the variables as below:


```
environments:
  default:
    domain: example.com
    name: staging
    profile_name: example-staging
    region: eu-central-1
  permanent:
    - staging
  temporary:
    secrets:
      firebase_credentials: DUMMY
      mail_password: DUMMY
general:
  organization: example-organization
  secrets_key_id: alias/secrets
sites:
  ordered:
    - inception
    - network
```

???+ tip
    You can get this template file in repository folder: 
    `~/terrabutler-template/configs/settings.yml`

## Terraform

Start installing Terraform with tfenv:

```shell
$ tfenv install 
```

### Change Variables

Change all project variables and settings in the template, starting by replacing `example` to your project name.

```shell
Project Variables File:
~/terrabutler-template/variables/example-staging.tfvars
```

### Initialize the Project
 
Perform an Terraform Initialization:

```shell
$ terraform init
```
Perform an Terraform Apply:

```shell
$ terraform apply
```

### Change local to remote backend

#### Edit the bucket configuration variables

Edit the variable files inside `/configs/backends/` as your bucket configuration in AWS, like the example below:
```
region         = "eu-central-1"
profile        = "example-development"
key            = "staging-network.tfstate"
bucket         = "mb-staging-site-network-tfstate"
dynamodb_table = "mb_staging_site_network_tfstatelock"

```

#### Uncommit the remote backend line
 Remove the commited line as below, in the `terrabutler-template` in the path `/site-inception/terraform.tf` to change from local to remote.

```
backend "s3" {}
```

Perform an Terraform initialization with the backend config to update the new changes:

```shell
terraform init -backend-config="/configs/backends/example-staging.tfvars"
```







 [Version-shield]: https://img.shields.io/badge/terrabutler_Template-Download-%23121011.svg?style=for-the-badge&logo=github&colorA=273133&colorB=0093ee "Latest version"

 

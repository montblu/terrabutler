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

First of all, you will need to create the configs directory for Terrabutler:

```shell
mkdir configs
```

Inside the `~/configs/`, copy the template file `settings.yml` below, and edit the variables:


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
    `~/terrabutler/template/configs/settings.yml`

## Terraform




 
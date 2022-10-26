# Basic Usage

In this section we quickly cover the basic commands of **Terrabutler**.
The usage of it can always be seen by using the help menu inside of every
command or subcommand.

Example:

``` shell
terrabutler tf -site inception apply --help
```

The command above shows all the arguments and options that can be used when
running that command.

## Usage 

``` shell
terrabutler [global options] command [subcommand] [arguments] [options]
```

## Global options

All global options can be placed at the command level.

* `--help`, `-help`, `-h`: Show help menu.
* `--version`, `-version`: Show version of **Terrabutler**.

## Commands

The commands are:

- `env`: Manage environments
- `init`: Initialize the manager
- `tf`: Manage terraform commands

### Command `env`

Subcommands:

- `delete`: "Delete an environment"
- `list`: "List environments"
- `new`: "Create a new environment"
- `reload`: "Reload the current environment"
- `select`: "Select a environment"
- `show`: "Show the name of the current environment"

Example:

``` shell
terrabutler env select staging
```

The command above change the current environment to `staging`.

### Command `tf`

???+ tip
    The `tf` subcommands are the Terraform commands

Subcommands:

- `apply`: "Create or update infrastructure"
- `console`: "Try Terraform expressions at an interactive command..."
- `destroy`: "Prepare your working directory for other commands"
- `fmt`: "Reformat your configuration in the standardstyle"
- `force-unlock`: "Release a stuck lock on the current workspace"
- `generate-options`: "Generate terraform options"
- `import`: "Associate existing infrastructure with a Terraform..."
- `init`: "Prepare your working directory for other commands"
- `output`: "Show output values from your root module"
- `plan`: "Show changes required by the current configuration"
- `providers`: "Show the providers required for this configuration"
- `refresh`: "Update the state to match remote systems"
- `show`: "Show the current state or a saved plan"
- `state`: "Advanced state management"
- `taint`: "Mark a resource instance as not fully functional"
- `untaint`: "Remove the 'tainted' state from a resource instance"
- `validate`: "Validate the configuration files"
- `version`: "Show the current Terraform version"

Example:

``` shell
terrabutler tf -site inception apply
```

The command above run a `terraform apply` command inside the `site inception` in
the current environment.

### Command `init`

Has no subcommands


#!/usr/bin/env python3

"""
Python wrapper used to manage environments and terraform commands
"""

import click
from colorama import Fore
from sys import exit
from os import path
from terrabutler.__init__ import (
    __name__,
    __version__
)
from terrabutler.env import (
    create_env,
    delete_env,
    get_current_env,
    set_current_env,
    get_available_envs
)
from terrabutler.tf import (
    terraform_args_print,
    terraform_command_runner
)
from terrabutler.settings import (
    get_settings,
    validate_settings
)
from terrabutler.inception import (
    inception_init,
    inception_init_needed
)
from terrabutler.utils import is_semantic_version


VERSION = (f"v{__version__}" if is_semantic_version(__version__)
           else __version__
           )


@click.group(context_settings=dict(help_option_names=['-h', '-help',
                                                      '--help']))
@click.version_option(version=VERSION, prog_name=__name__.capitalize(),
                      message='%(prog)s: %(version)s')
def main():
    validate_settings()


@main.group(name="env", help="Manage environments")
def env_cli():
    inception_init_needed()


@main.command(name="init", help="Initialize the manager")
def init_cli():
    inception_init()


@main.group(name="tf", help="Manage terraform commands")
@click.pass_context
@click.option('-site', metavar='SITE', required=True, help="Site where to run"
                                                           "terraform.")
def tf_cli(ctx, site):
    inception_init_needed()
    ctx.ensure_object(dict)
    ctx.obj['SITE'] = site

    sites = get_settings()["sites"]["ordered"]
    if "site_" in site:
        print(Fore.YELLOW + "There is no need to pass the 'site_' prefix.")
        site = site.replace("site_", "")
    if path.exists(f"site_{site}") and site not in sites:
        print(Fore.RED + f"The site '{site}' exists but is missing inside the"
              " config file.")
        exit(1)
    elif site not in sites:
        print(Fore.RED + f"The site '{site}' does not exist.")
        exit(1)


@env_cli.command(name="delete", help="Delete an environment")
@click.argument('NAME')
@click.option("-d", is_flag=True,
              help="Destroy all sites by inverse order.")
@click.option("-y", is_flag=True,
              help="Delete without asking for confirmation.")
@click.option("-s3", is_flag=True,
              help="Access S3 instead of parsing terraform output")
def env_delete_cli(name, y, d, s3):
    delete_env(name, y, d, s3)


@env_cli.command(name="list", help="List environments")
@click.option("-s3", is_flag=True,
              help="Access S3 instead of parsing terraform output")
def env_list_cli(s3):
    current_env = get_current_env()

    for env in get_available_envs(s3):
        if env == current_env:
            print(f"\u2192 {env}")
        else:
            print(env)


@env_cli.command(name="new", help="Create a new environment")
@click.argument('NAME')
@click.option("-y", is_flag=True,
              help="Delete without asking for confirmation.")
@click.option("-t", "-temp", is_flag=True,
              help="Create a temporary environment.")
@click.option("-a", "-apply", is_flag=True,
              help="Apply all terraform sites prior the creation" +
              " of the environment")
@click.option("-s3", is_flag=True,
              help="Access S3 instead of parsing terraform output")
def env_new_cli(name, y, t, a, s3):
    create_env(name, y, t, a, s3)


@env_cli.command(name="select", help="Select a environment")
@click.argument('NAME')
@click.option("-s3", is_flag=True,
              help="Access S3 instead of parsing terraform output")
def env_select_cli(name, s3):
    set_current_env(name, s3)


@env_cli.command(name="show", help="Show the name of the current environment")
def env_show_cli():
    print(get_current_env())


@tf_cli.command(name="apply", help="Create or update infrastructure")
@click.option("-auto-approve", is_flag=True,
              help="Skip interactive approval of plan before applying.")
@click.option("-destroy", is_flag=True,
              help="Select the 'destroy' planning mode, which creates a plan"
                   " to destroy all objects currently managed by this"
                   " Terraform configuration instead of the usual behavior.")
@click.option("-input", default=True,
              help="Ask for input for variables if not directly set.")
@click.option("-lock", default=True,
              help="Don't hold a state lock during backend migration. This is"
                   " dangerous if others might concurrently run commands"
                   " against the same workspace.")
@click.option("-lock-timeout", help="Duration to retry a state lock.")
@click.option("-no-color", is_flag=True,
              help="If specified, output won't contain any color.")
@click.option("-refresh-only", is_flag=True,
              help="Select the 'refresh only' planning mode, which checks"
                   " whether remote objects still match the outcome of the"
                   " most recent Terraform apply but does not propose any"
                   " actions to undo any changes made outside of Terraform.")
@click.option("-refresh", default=True,
              help="Skip checking for external changes to remote objects while"
              " creating the plan. This can potentially make planning faster,"
              " but at the expense of possibly planning against a stale record"
              " of the remote system state.")
@click.option("-target", multiple=True,
              help="Limit the planning operation to only the given module,"
                   " resource, or resource instance and all of its"
                   " dependencies. You can use this option multiple times to"
                   " include more than one object. This is for exceptional use"
                   " only.")
@click.option("-var", multiple=True,
              help="Set a value for one of the input variables in the"
              " root module of the configuration. Use this option"
              " more than once to set more than one variable.")
@click.pass_context
def tf_apply_cli(ctx, auto_approve, destroy, input, lock, lock_timeout,
                 no_color, refresh_only, refresh, target, var):
    args = []

    if auto_approve:
        args.append("-auto-approve")
    if destroy:
        args.append("-destroy")
    if input is False:
        args.append("-input=false")
    if lock is False:
        args.append("-lock=false")
    if lock_timeout:
        args.append(f"-lock-timeout={lock_timeout}")
    if no_color:
        args.append("-no-color")
    if refresh_only:
        args.append("-refresh-only")
    if refresh is False:
        args.append("-refresh=false")
    if target:
        for name in target:
            args.append(f"-target={name}")
    if var:
        for name in var:
            args.append(f"-var='{name}'")

    terraform_command_runner("apply", args, "var", ctx.obj['SITE'])


@tf_cli.command(name="console", help="Try Terraform expressions at an "
                                     "interactive command prompt")
@click.option("-state", help="Legacy option for the local backend only."
              " See the local backend's documentation for more information.")
@click.option("-var", multiple=True,
              help="Set a variable in the Terraform configuration. "
                   "This flag can be set multiple times.")
@click.pass_context
def tf_console_cli(ctx, state, var):
    args = []

    if state:
        args.append(f"-state={state}")
    if var:
        for name in var:
            args.append(f"-var='{name}'")

    terraform_command_runner("console", args, "var", ctx.obj['SITE'])


@tf_cli.command(name="destroy", help="Prepare your working directory for other"
                                     " commands")
@click.option("-auto-approve", is_flag=True,
              help="Skip interactive approval of plan before applying.")
@click.option("-input", default=True,
              help="Ask for input for variables if not directly set.")
@click.option("-lock", default=True,
              help="Don't hold a state lock during backend migration. This is"
                   " dangerous if others might concurrently run commands"
                   " against the same workspace.")
@click.option("-lock-timeout", help="Duration to retry a state lock.")
@click.option("-no-color", is_flag=True,
              help="If specified, output won't contain any color.")
@click.option("-refresh-only", is_flag=True,
              help="Select the 'refresh only' planning mode, which checks"
                   " whether remote objects still match the outcome of the"
                   " most recent Terraform apply but does not propose any"
                   " actions to undo any changes made outside of Terraform.")
@click.option("-refresh", default=True,
              help="Skip checking for external changes to remote objects while"
              " creating the plan. This can potentially make planning faster,"
              " but at the expense of possibly planning against a stale record"
              " of the remote system state.")
@click.option("-target", multiple=True,
              help="Limit the planning operation to only the given module,"
                   " resource, or resource instance and all of its"
                   " dependencies. You can use this option multiple times to"
                   " include more than one object. This is for exceptional use"
                   " only.")
@click.option("-var", multiple=True,
              help="Set a value for one of the input variables in the"
              " root module of the configuration. Use this option"
              " more than once to set more than one variable.")
@click.pass_context
def tf_destroy_cli(ctx, auto_approve, input, lock, lock_timeout, no_color,
                   refresh_only, refresh, target, var):
    args = []

    if auto_approve:
        args.append("-auto-approve")
    if input is False:
        args.append("-input=false")
    if lock is False:
        args.append("-lock=false")
    if lock_timeout:
        args.append(f"-lock-timeout={lock_timeout}")
    if no_color:
        args.append("-no-color")
    if refresh_only:
        args.append("-refresh-only")
    if refresh is False:
        args.append("-refresh=false")
    if target:
        for name in target:
            args.append(f"-target={name}")
    if var:
        for name in var:
            args.append(f"-var='{name}'")

    terraform_command_runner("destroy", args, "var", ctx.obj['SITE'])


@tf_cli.command(name="fmt", help="Reformat your configuration in the standard"
                                 "style")
@click.option("-diff", is_flag=True,
              help="Display diffs of formatting changes")
@click.option("-no-color", is_flag=True,
              help="If specified, output won't contain any color.")
@click.option("-recursive", is_flag=True,
              help="Also process files in subdirectories. By default, only the"
                   " given directory (or current directory) is processed.")
@click.pass_context
def tf_fmt_cli(ctx, diff, no_color, recursive):
    args = []

    if diff:
        args.append("-diff")
    if no_color:
        args.append("-no-color")
    if recursive:
        args.append("-recursive")

    terraform_command_runner("fmt", args, "none", ctx.obj['SITE'])


@tf_cli.command(name="force-unlock", help="Release a stuck lock on the current"
                                          " workspace")
@click.argument("LOCK_ID")
@click.option("-force", is_flag=True,
              help="Don't ask for input for unlock confirmation.")
@click.pass_context
def tf_force_unlock_cli(ctx, lock_id, force):
    args = []

    if force:
        args.append("-force")
    args.append(lock_id)

    terraform_command_runner("force-unlock", args, "", ctx.obj['SITE'])


@tf_cli.command(name="generate-arguments", help="Generate terraform arguments")
@click.argument("command", type=click.Choice(["init", "plan", "apply"]))
@click.pass_context
def tf_generate_arguments_cli(ctx, command):
    print(Fore.GREEN + "Needed args:" + Fore.RESET)
    print(terraform_args_print(command, ctx.obj['SITE']))


@tf_cli.command(name="import", help="Associate existing infrastructure with a"
                                    " Terraform resource")
@click.argument("ADDR")
@click.argument("ID")
@click.option("-allow-missing-config", is_flag=True,
              help="Allow import when no resource configuration block exists.")
@click.option("-input", default=True,
              help="Ask for input for variables if not directly set.")
@click.option("-lock", default=True,
              help="Don't hold a state lock during the operation. This is "
                   "dangerous if others might concurrently run commands "
                   "against the same workspace.")
@click.option("-no-color", is_flag=True,
              help="If specified, output won't contain any color.")
@click.option("-var", multiple=True,
              help="Set a variable in the Terraform configuration. "
                   "This flag can be set multiple times.")
@click.option("-ignore-remote-version",
              help="A rare option used for the remote backend only. See the"
                   " remote backend documentation for more information.")
@click.pass_context
def tf_import_cli(ctx, addr, id, allow_missing_config, input, lock, no_color,
                  var, ignore_remote_version):
    print(Fore.RED + "Function not implemented yet!")


@tf_cli.command(name="init", help="Prepare your working directory for other"
                                  " commands")
@click.option("-backend", default=True,
              help="Disable backend or Terraform Cloud initialization for this"
                   " configuration and use what what was previously"
                   " initialized instead.")
@click.option("-force-copy", is_flag=True,
              help="Allow import when no resource configuration block exists.")
@click.option("-get", default=True,
              help="Disable downloading modules for this configuration.")
@click.option("-input", default=True,
              help="Disable interactive prompts. Note that some actions may"
                   " require interactive prompts and will error if input is"
                   " disabled.")
@click.option("-lock", default=True,
              help="Don't hold a state lock during backend migration. This is"
                   " dangerous if others might concurrently run commands"
                   " against the same workspace.")
@click.option("-no-color", is_flag=True,
              help="If specified, output won't contain any color.")
@click.option("-reconfigure", is_flag=True,
              help="Reconfigure a backend, ignoring any saved configuration.")
@click.option("-migrate-state", is_flag=True,
              help="Reconfigure a backend, and attempt to migrate any existing"
                   " state.")
@click.option("-upgrade", is_flag=True,
              help="Install the latest module and provider versions allowed"
              " within configured constraints, overriding the default behavior"
              " of selecting exactly the version recorded in the dependency"
              " lockfile.")
@click.option("-lockfile",
              help="Set a dependency lockfile mode. Currently only 'readonly'"
              " is valid.")
@click.option("-ignore-remote-version", is_flag=True,
              help="A rare option used for Terraform Cloud and the remote"
              " backend only. Set this to ignore checking that the local and"
              " remote Terraform versions use compatible state representations"
              ", making an operation proceed even when there is a potential"
              " mismatch. See the documentation on configuring Terraform with"
              " Terraform Cloud for more information.")
@click.pass_context
def tf_init_cli(ctx, backend, force_copy, get, input, lock, no_color,
                reconfigure, migrate_state, upgrade, lockfile,
                ignore_remote_version):

    args = []

    if backend is False:
        args.append("-backend=false")
    if force_copy:
        args.append("-force-copy")
    if get is False:
        args.append("-get=false")
    if input is False:
        args.append("-input=false")
    if lock is False:
        args.append("-lock=false")
    if no_color:
        args.append("-no-color")
    if reconfigure:
        args.append("-reconfigure")
    if migrate_state:
        args.append("-migrate-state")
    if upgrade:
        args.append("-upgrade")
    if lockfile:
        args.append(f"-lockfile={lockfile}")
    if ignore_remote_version:
        args.append("-ignore-remote-version")

    terraform_command_runner("init", args, "backend", ctx.obj['SITE'])


@tf_cli.command(name="output", help="Show output values from your root module")
@click.option("-no-color", is_flag=True,
              help="If specified, output won't contain any color.")
@click.option("-json", is_flag=True,
              help="If specified, machine readable output will be printed in"
              " JSON format.")
@click.option("-raw", is_flag=True,
              help="For value types that can be automatically converted to a"
              " string, will print the raw string directly, rather than a"
              " human-oriented representation of the value.")
@click.pass_context
def tf_output_cli(ctx, no_color, json, raw):
    args = []

    if no_color:
        args.append("-no-color")
    if json:
        args.append("-json")
    if raw:
        args.append("-raw")

    terraform_command_runner("output", args, "", ctx.obj['SITE'])


@tf_cli.command(name="plan", help="Show changes required by the current"
                                  " configuration")
@click.option("-destroy", is_flag=True,
              help="Select the 'destroy' planning mode, which creates a plan"
                   " to destroy all objects currently managed by this"
                   " Terraform configuration instead of the usual behavior.")
@click.option("-input", default=True,
              help="Ask for input for variables if not directly set.")
@click.option("-lock", default=True,
              help="Don't hold a state lock during backend migration. This is"
                   " dangerous if others might concurrently run commands"
                   " against the same workspace.")
@click.option("-lock-timeout", help="Duration to retry a state lock.")
@click.option("-no-color", is_flag=True,
              help="If specified, output won't contain any color.")
@click.option("-refresh-only", is_flag=True,
              help="Select the 'refresh only' planning mode, which checks"
                   " whether remote objects still match the outcome of the"
                   " most recent Terraform apply but does not propose any"
                   " actions to undo any changes made outside of Terraform.")
@click.option("-refresh", default=True,
              help="Skip checking for external changes to remote objects while"
              " creating the plan. This can potentially make planning faster,"
              " but at the expense of possibly planning against a stale record"
              " of the remote system state.")
@click.option("-target", multiple=True,
              help="Limit the planning operation to only the given module,"
                   " resource, or resource instance and all of its"
                   " dependencies. You can use this option multiple times to"
                   " include more than one object. This is for exceptional use"
                   " only.")
@click.option("-var", multiple=True,
              help="Set a value for one of the input variables in the"
              " root module of the configuration. Use this option"
              " more than once to set more than one variable.")
@click.pass_context
def tf_plan_cli(ctx, destroy, input, lock, lock_timeout, no_color,
                refresh_only, refresh, target, var):
    args = []

    if destroy:
        args.append("-destroy")
    if input is False:
        args.append("-input=false")
    if lock is False:
        args.append("-lock=false")
    if lock_timeout:
        args.append(f"-lock-timeout={lock_timeout}")
    if no_color:
        args.append("-no-color")
    if refresh_only:
        args.append("-refresh-only")
    if refresh is False:
        args.append("-refresh=false")
    if target:
        for t in target:
            args.append(f"-target={t}")
    if var:
        for v in var:
            args.append(f"-var={v}")

    terraform_command_runner("plan", args, "var", ctx.obj['SITE'])


@tf_cli.command(name="providers", help="Show the providers required for this"
                                       " configuration")
@click.pass_context
def tf_providers_cli():
    print(Fore.RED + "Function not implemented yet!")


@tf_cli.command(name="refresh", help="Update the state to match remote"
                                     " systems")
@click.option("-input", default=True,
              help="Ask for input for variables if not directly set.")
@click.option("-lock", default=True,
              help="Don't hold a state lock during the operation. This is "
                   "dangerous if others might concurrently run commands "
                   "against the same workspace.")
@click.option("-no-color", is_flag=True,
              help="If specified, output won't contain any color.")
@click.option("-target", multiple=True,
              help="Resource to target. Operation will be limited to this"
              " resource and its dependencies. This flag can be used"
              " multiple times.")
@click.option("-var", multiple=True,
              help="Set a variable in the Terraform configuration. "
                   "This flag can be set multiple times.")
@click.pass_context
def tf_refresh_cli(ctx, input, lock, no_color, target, var):
    print(Fore.RED + "Function not implemented yet!")


@tf_cli.command(name="show", help="Show the current state or a saved plan")
@click.option("-no-color", is_flag=True,
              help="If specified, output won't contain any color.")
@click.option("-json", is_flag=True,
              help="If specified, machine readable output will be printed in"
              " JSON format.")
@click.pass_context
def tf_show_cli(ctx, no_color, json):
    print(Fore.RED + "Function not implemented yet!")


@tf_cli.command(name="state", help="Advanced state management")
def tf_state_cli():
    print(Fore.RED + "Function not implemented yet!")


@tf_cli.command(name="taint", help="Mark a resource instance as not fully"
                                   " functional")
@click.option("-allow-missing", is_flag=True,
              help="If specified, the command will succeed (exit code 0) even"
                   " if the resource is missing.")
@click.option("-lock", default=True,
              help="Don't hold a state lock during the operation. This is "
                   "dangerous if others might concurrently run commands "
                   "against the same workspace.")
@click.option("-lock-timeout", help="Duration to retry a state lock.")
@click.option("-ignore-remote-version", is_flag=True,
              help="A rare option used for the remote backend only. See the"
                   " remote backend documentation for more information.")
@click.pass_context
def tf_taint_cli(ctx, allow_missing, lock, lock_timeout,
                 ignore_remote_version):
    print(Fore.RED + "Function not implemented yet!")


@tf_cli.command(name="untaint", help="Remove the 'tainted' state from a"
                                     " resource instance")
@click.option("-allow-missing", is_flag=True,
              help="If specified, the command will succeed (exit code 0) even"
                   " if the resource is missing.")
@click.option("-lock", default=True,
              help="Don't hold a state lock during the operation. This is "
                   "dangerous if others might concurrently run commands "
                   "against the same workspace.")
@click.option("-lock-timeout", help="Duration to retry a state lock.")
@click.option("-ignore-remote-version", is_flag=True,
              help="A rare option used for the remote backend only. See the"
                   " remote backend documentation for more information.")
@click.pass_context
def tf_untaint_cli(ctx, allow_missing, lock, lock_timeout,
                   ignore_remote_version):
    print(Fore.RED + "Function not implemented yet!")


@tf_cli.command(name="version", help="Show the current Terraform version")
@click.option("-json", help="Output the version information as a JSON object.")
@click.pass_context
def tf_version_cli(ctx, json):
    args = []

    if json:
        args.append("-json")

    terraform_command_runner("version", args, "", ctx.obj['SITE'])

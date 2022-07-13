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
    options = []

    if auto_approve:
        options.append("-auto-approve")
    if destroy:
        options.append("-destroy")
    if input is False:
        options.append("-input=false")
    if lock is False:
        options.append("-lock=false")
    if lock_timeout:
        options.append(f"-lock-timeout={lock_timeout}")
    if no_color:
        options.append("-no-color")
    if refresh_only:
        options.append("-refresh-only")
    if refresh is False:
        options.append("-refresh=false")
    if target:
        for name in target:
            options.append(f"-target={name}")
    if var:
        for name in var:
            options.append(f"-var='{name}'")

    terraform_command_runner("apply", ctx.obj['SITE'], options=options,
                             needed_options="var")


@tf_cli.command(name="console", help="Try Terraform expressions at an "
                                     "interactive command prompt")
@click.option("-state", help="Legacy option for the local backend only."
              " See the local backend's documentation for more information.")
@click.option("-var", multiple=True,
              help="Set a variable in the Terraform configuration. "
                   "This flag can be set multiple times.")
@click.pass_context
def tf_console_cli(ctx, state, var):
    options = []

    if state:
        options.append(f"-state={state}")
    if var:
        for name in var:
            options.append(f"-var='{name}'")

    terraform_command_runner("console", ctx.obj['SITE'], options=options,
                             needed_options="var")


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
    options = []

    if auto_approve:
        options.append("-auto-approve")
    if input is False:
        options.append("-input=false")
    if lock is False:
        options.append("-lock=false")
    if lock_timeout:
        options.append(f"-lock-timeout={lock_timeout}")
    if no_color:
        options.append("-no-color")
    if refresh_only:
        options.append("-refresh-only")
    if refresh is False:
        options.append("-refresh=false")
    if target:
        for name in target:
            options.append(f"-target={name}")
    if var:
        for name in var:
            options.append(f"-var='{name}'")

    terraform_command_runner("destroy", ctx.obj['SITE'], options=options,
                             needed_options="var")


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
    options = []

    if diff:
        options.append("-diff")
    if no_color:
        options.append("-no-color")
    if recursive:
        options.append("-recursive")

    terraform_command_runner("fmt", ctx.obj['SITE'], options=options)


@tf_cli.command(name="force-unlock", help="Release a stuck lock on the current"
                                          " workspace")
@click.argument("LOCK_ID")
@click.option("-force", is_flag=True,
              help="Don't ask for input for unlock confirmation.")
@click.pass_context
def tf_force_unlock_cli(ctx, lock_id, force):
    args, options = []

    args.append(lock_id)
    if force:
        options.append("-force")

    terraform_command_runner("force-unlock", ctx.obj['SITE'], args=args,
                             options=options)


@tf_cli.command(name="generate-options", help="Generate terraform options")
@click.argument("command", type=click.Choice(["init", "plan", "apply"]))
@click.pass_context
def tf_generate_options_cli(ctx, command):
    print(Fore.GREEN + "Options:" + Fore.RESET)
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
    args, options = [], []

    args.append(addr)
    args.append(id)
    if allow_missing_config:
        options.append("-allow-missing-config")
    if input is False:
        options.append("-input=false")
    if input is False:
        options.append("-input=false")
    if no_color:
        options.append("-no-color")
    if var:
        for v in var:
            options.append(f"-var={v}")
    if ignore_remote_version:
        options.append("-ignore-remote-version")

    terraform_command_runner("import", ctx.obj['SITE'], args=args,
                             options=options, needed_options="var")


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

    options = []

    if backend is False:
        options.append("-backend=false")
    if force_copy:
        options.append("-force-copy")
    if get is False:
        options.append("-get=false")
    if input is False:
        options.append("-input=false")
    if lock is False:
        options.append("-lock=false")
    if no_color:
        options.append("-no-color")
    if reconfigure:
        options.append("-reconfigure")
    if migrate_state:
        options.append("-migrate-state")
    if upgrade:
        options.append("-upgrade")
    if lockfile:
        options.append(f"-lockfile={lockfile}")
    if ignore_remote_version:
        options.append("-ignore-remote-version")

    terraform_command_runner("init", ctx.obj['SITE'], options=options,
                             needed_options="backend")


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
    options = []

    if no_color:
        options.append("-no-color")
    if json:
        options.append("-json")
    if raw:
        options.append("-raw")

    terraform_command_runner("output", ctx.obj['SITE'], options=options)


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
    options = []

    if destroy:
        options.append("-destroy")
    if input is False:
        options.append("-input=false")
    if lock is False:
        options.append("-lock=false")
    if lock_timeout:
        options.append(f"-lock-timeout={lock_timeout}")
    if no_color:
        options.append("-no-color")
    if refresh_only:
        options.append("-refresh-only")
    if refresh is False:
        options.append("-refresh=false")
    if target:
        for t in target:
            options.append(f"-target={t}")
    if var:
        for v in var:
            options.append(f"-var={v}")

    terraform_command_runner("plan", ctx.obj['SITE'], options=options,
                             needed_options="var")


@tf_cli.group(name="providers", help="Show the providers required for this"
                                     " configuration",
              invoke_without_command=True)
@click.pass_context
def tf_providers_cli(ctx):
    if ctx.invoked_subcommand is None:
        terraform_command_runner("providers", ctx.obj['SITE'])


@tf_providers_cli.command(name="lock", help="Write out dependency locks for"
                                            " the configured providers")
@click.argument("providers", nargs=-1, required=True)
@click.option("-fs-mirror", help="Consult the given filesystem mirror"
                                 " directory instead of the origin registry"
                                 " for each of the given providers.")
@click.option("-net-mirror", help="Consult the given network mirror"
                                  " (given as a base URL) instead of the"
                                  " origin registry for each of the given"
                                  " providers.")
@click.option("-platform", help="Choose a target platform to request package"
                                " checksums for.")
@click.pass_context
def tf_providers_lock_cli(ctx, providers, fs_mirror, net_mirror, platform):
    args, options = []

    args += providers
    if fs_mirror:
        options.append(f"-fs-mirror={fs_mirror}")
    if net_mirror:
        options.append(f"-net-mirror={net_mirror}")
    if platform:
        options.append(f"-platform={platform}")

    terraform_command_runner("providers lock", ctx.obj['SITE'], args=args,
                             options=options)


@tf_providers_cli.command(name="mirror", help="Save local copies of all"
                                              " required provider plugins")
@click.argument("target-dir", required=True)
@click.option("-platform", help="Choose a target platform to request package"
                                " checksums for.")
@click.pass_context
def tf_providers_mirror_cli(ctx, target_dir, platform):
    args, options = []

    args.append(target_dir)
    if platform:
        options.append(f"-platform={platform}")

    terraform_command_runner("providers lock", ctx.obj['SITE'], args=args,
                             options=options)


@tf_providers_cli.command(name="schema", help="Show schemas for the providers"
                                              " used in the configuration")
@click.option("-json", help="Prints out a json representation of the schemas"
                            " for all providers used in the current"
                            " configuration.", required=True, is_flag=True)
@click.pass_context
def tf_providers_schema_cli(ctx, json):
    options = []

    if json:
        options.append("-json")

    terraform_command_runner("providers schema", ctx.obj['SITE'],
                             options=options)


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
    options = []

    if input is False:
        options.append("-input=false")
    if lock is False:
        options.append("-lock=false")
    if no_color:
        options.append("-no-color")
    if target:
        for t in target:
            options.append(f"-target={t}")
    if var:
        for v in var:
            options.append(f"-var={v}")

    terraform_command_runner("refresh", ctx.obj['SITE'], options=options,
                             needed_options="var")


@tf_cli.command(name="show", help="Show the current state or a saved plan")
@click.argument("PATH", required=False)
@click.option("-no-color", is_flag=True,
              help="If specified, output won't contain any color.")
@click.option("-json", is_flag=True,
              help="If specified, machine readable output will be printed in"
              " JSON format.")
@click.pass_context
def tf_show_cli(ctx, path, no_color, json):
    args, options = []

    args.append(path)
    if no_color:
        options.append("-no-color")
    if json:
        options.append("-json")

    terraform_command_runner("show", ctx.obj['SITE'], args=args,
                             options=options)


@tf_cli.group(name="state", help="Advanced state management")
@click.pass_context
def tf_state_cli(ctx):
    pass


@tf_state_cli.command(name="list", help="List resources in the state")
@click.argument("address", nargs=-1, required=False)
@click.option("-state", help="Path to a Terraform state file to use to look up"
                             " Terraform-managed resources. By default,"
                             " Terraform will consult the state of the"
                             " currently-selected workspace.")
@click.option("-id", help="Filters the results to include only instances"
                          " whoseresource types have an attribute named 'id'"
                          " whose value equals the given id string.")
@click.pass_context
def tf_state_list_cli(ctx, address, state, id):
    args, options = []

    if len(address) > 0:
        args += address
    if state:
        options.append(f"-state={state}")
    if id:
        options.append(f"-id={id}")

    terraform_command_runner("state list", ctx.obj['SITE'], args=args,
                             options=options)


@tf_state_cli.command(name="mv", help="Move an item in the state")
@click.argument("source", required=True)
@click.argument("destination", required=True)
@click.option("-dry-run", is_flag=True, help="If set, prints out what would've"
                                             " been moved but doesn't actually"
                                             " move anything.")
@click.option("-lock", is_flag=True, default=True,
              help="Don't hold a state lock during the operation. This is"
                   " dangerous if others might concurrently run commands"
                   " against the same workspace.")
@click.option("-lock-timeout", help="Duration to retry a state lock.")
@click.option("-ignore-remote-version", is_flag=True,
              help="A rare option used for the remote backend only."
              "See the remote backend documentation for more information.")
@click.pass_context
def tf_state_mv_cli(ctx, source, destination, dry_run, lock, lock_timeout,
                    ignore_remote_version):
    args, options = []

    args.append(source)
    args.append(destination)
    if dry_run:
        options.append("-dry-run")
    if lock is False:
        options.append("-lock=false")
    if lock_timeout:
        options.append(f"-lock-timeout={lock_timeout}")
    if ignore_remote_version:
        options.append("-ignore-remote-version")

    terraform_command_runner("state mv", ctx.obj['SITE'], args=args,
                             options=options)


@tf_state_cli.command(name="pull", help="Pull current state and output to"
                                        " stdout")
@click.pass_context
def tf_state_pull_cli(ctx):
    terraform_command_runner("state pull", ctx.obj['SITE'])


@tf_state_cli.command(name="push", help="Update remote state from a local"
                                        " state file")
@click.argument("path", required=True)
@click.option("-force", is_flag=True,
              help="Write the state even if lineages don't match or the remote"
                   " serial is higher.")
@click.option("-lock", is_flag=True, default=True,
              help="Don't hold a state lock during the operation. This is"
                   " dangerous if others might concurrently run commands"
                   " against the same workspace.")
@click.option("-lock-timeout", help="Duration to retry a state lock.")
@click.pass_context
def tf_state_push_cli(ctx, path, force, lock, lock_timeout):
    args, options = []

    args.append(path)
    if force:
        options.append("-force")
    if lock is False:
        options.append("-lock=false")
    if lock_timeout:
        options.append(f"-lock-timeout={lock_timeout}")

    terraform_command_runner("state push", ctx.obj['SITE'], args=args,
                             options=options)


@tf_state_cli.command(name="replace-provider",
                      help="Replace provider for resources in the Terraform"
                           " state.")
@click.argument("from_provider_fqdn", required=True)
@click.argument("to_provider_fqdn", required=True)
@click.option("-auto-approve", is_flag=True,
              help="Skip interactive approval of plan before applying.")
@click.option("-lock", is_flag=True, default=True,
              help="Don't hold a state lock during the operation. This is"
                   " dangerous if others might concurrently run commands"
                   " against the same workspace.")
@click.option("-lock-timeout", help="Duration to retry a state lock.")
@click.option("-ignore-remote-version", is_flag=True,
              help="A rare option used for the remote backend only. See the"
                   " remote backend documentation for more information.")
@click.pass_context
def tf_state_replace_cli(ctx, from_provider_fqdn, to_provider_fqdn,
                         auto_approve, lock, lock_timeout,
                         ignore_remote_version):
    args, options = []

    args.append(from_provider_fqdn)
    args.append(to_provider_fqdn)
    if auto_approve:
        options.append("-auto-approve")
    if lock is False:
        options.append("-lock=false")
    if lock_timeout:
        options.append(f"-lock-timeout={lock_timeout}")
    if ignore_remote_version:
        options.append("-ignore-remote-version")

    terraform_command_runner("state replace-provider", ctx.obj['SITE'],
                             args=args, options=options)


@tf_state_cli.command(name="rm",
                      help="Remove instances from the state")
@click.argument("address", nargs=-1, required=True)
@click.option("-dry-run", is_flag=True, help="If set, prints out what would've"
                                             " been moved but doesn't actually"
                                             " move anything.")
@click.option("-backup", help="Path where Terraform should write the backup"
                              " state.")
@click.option("-lock", is_flag=True, default=True,
              help="Don't hold a state lock during the operation. This is"
                   " dangerous if others might concurrently run commands"
                   " against the same workspace.")
@click.option("-lock-timeout", help="Duration to retry a state lock.")
@click.option("-state", help="Path to the state file to update. Defaults to"
                             " the current workspace state.")
@click.option("-ignore-remote-version", is_flag=True,
              help="A rare option used for the remote backend only. See the"
                   " remote backend documentation for more information.")
@click.pass_context
def tf_state_rm_cli(ctx, address, dry_run, backup, lock, lock_timeout, state,
                    ignore_remote_version):
    args, options = []

    args += address
    if dry_run:
        options.append("-dry-run")
    if backup:
        options.append(f"-backup={backup}")
    if lock is False:
        options.append("-lock=false")
    if lock_timeout:
        options.append(f"-lock-timeout={lock_timeout}")
    if state:
        options.append(f"-state={state}")
    if ignore_remote_version:
        options.append("-ignore-remote-version")

    terraform_command_runner("state rm", ctx.obj['SITE'], args=args,
                             options=options)


@tf_state_cli.command(name="show",
                      help="Show a resource in the state")
@click.argument("address", required=True)
@click.option("-state", help="Path to the state file to update. Defaults to"
                             " the current workspace state.")
@click.pass_context
def tf_state_show_cli(ctx, address, state):
    args, options = []

    args.append(address)
    if state:
        options.append(f"-state={state}")

    terraform_command_runner("state show", ctx.obj['SITE'], args=args,
                             options=options)


@tf_cli.command(name="taint", help="Mark a resource instance as not fully"
                                   " functional")
@click.argument("address", required=True)
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
def tf_taint_cli(ctx, address, allow_missing, lock, lock_timeout,
                 ignore_remote_version):
    args, options = []

    args.append(address)
    if allow_missing:
        options.append("-allow-missing")
    if lock is False:
        options.append("-lock=false")
    if lock_timeout:
        options.append(f"-lock-timeout={lock_timeout}")
    if ignore_remote_version:
        options.append("-ignore-remote-version")

    terraform_command_runner("taint", ctx.obj['SITE'], args=args,
                             options=options)


@tf_cli.command(name="untaint", help="Remove the 'tainted' state from a"
                                     " resource instance")
@click.argument("address", required=True)
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
def tf_untaint_cli(ctx, address, allow_missing, lock, lock_timeout,
                   ignore_remote_version):
    args, options = []

    args.append(address)
    if allow_missing:
        options.append("-allow-missing")
    if lock is False:
        options.append("-lock=false")
    if lock_timeout:
        options.append(f"-lock-timeout={lock_timeout}")
    if ignore_remote_version:
        options.append("-ignore-remote-version")

    terraform_command_runner("untaint", ctx.obj['SITE'], args=args,
                             options=options)


@tf_cli.command(name="version", help="Show the current Terraform version")
@click.option("-json", help="Output the version information as a JSON object.")
@click.pass_context
def tf_version_cli(ctx, json):
    options = []

    if json:
        options.append("-json")

    terraform_command_runner("version", ctx.obj['SITE'], options=options)


@tf_cli.command(name="validate", help="Validate the configuration files")
@click.option("-no-color", is_flag=True,
              help="If specified, output won't contain any color.")
@click.option("-json", is_flag=True,
              help="Output the version information as a JSON object.")
@click.pass_context
def tf_validate_cli(ctx, no_color, json):
    options = []

    if no_color:
        options.append("-no-color")
    if json:
        options.append("-json")

    terraform_command_runner("validate", ctx.obj['SITE'], options=options)

import os
import signal
import subprocess
from colorama import Fore
from sys import exit
from terrabutler.settings import get_settings
from terrabutler.utils import paths

# Values from Config
org = get_settings()["general"]["organization"]


def setup_tfenv(site):
    """
    Use desired version of terraform
    """
    process = subprocess.run(args=['tfenv', 'install'],
                             cwd=site,
                             stdout=subprocess.DEVNULL,
                             stderr=subprocess.DEVNULL)
    if process.returncode != 0:
        print('Error: tfenv failed to install the terraform version')
        exit(1)


def terraform_args_print(command, site):
    """
    Print Args
    """

    if command == "init":
        needed_options = "backend"
    elif command == "plan" or command == "apply":
        needed_options = "var"

    options = terraform_needed_options_builder(needed_options, site)
    return " ".join(options)


def terraform_needed_options_builder(needed_options, site):
    """
    Create array of needed options for backend or var files
    """
    from terrabutler.env import get_current_env
    env = get_current_env()

    if needed_options == "backend":
        backend_dir = paths["backends"]

        if site == "inception":  # Inception backend does only exist in dev
            return ["-backend-config",
                    f"{backend_dir}/{org}-dev-inception.tfvars"]
        else:
            return ["-backend-config",
                    f"{backend_dir}/{org}-{env}-{site}.tfvars"]

    elif needed_options == "var":
        variables_dir = paths["variables"]

        return ["-var-file", f"{variables_dir}/global.tfvars",
                "-var-file", f"{variables_dir}/{org}-{env}.tfvars",
                "-var-file", f"{variables_dir}/{org}-{env}-{site}.tfvars"
                ]

    else:  # If needed_options is empty, return empty array
        return []


def terraform_command_builder(command, site, args=[], options=[],
                              needed_options=""):
    """
    Create the command to run terraform
    """
    aux = ["terraform", command]  # Start base command

    if needed_options == "backend" or needed_options == "var":
        aux += terraform_needed_options_builder(needed_options, site)

    aux += options  # Add options passed by user
    aux += args  # Add args passed by user

    return aux


def terraform_command_runner(command, site, args=[], options=[],
                             needed_options=""):
    """
    Run tfenv and run the terraform command
    """
    from terrabutler.env import get_current_env
    site_dir = f"{paths['root']}/site_{site}"
    env = get_current_env()

    setup_tfenv(site_dir)

    command = terraform_command_builder(command, site, args=args,
                                        options=options,
                                        needed_options=needed_options)
    try:
        p = subprocess.Popen(args=command, cwd=site_dir)
        p.wait()

        if p.returncode != 0:
            print(Fore.RED + "There was an error while running the terraform"
                  " command.")
            exit(1)
    except KeyboardInterrupt:
        p.send_signal(signal.SIGINT)
        p.wait()
        exit(p.returncode)
    except subprocess.CalledProcessError:
        print(Fore.RED + f"There was an error while doing the {command}"
              f" command inside the '{site}' site in '{env}' environment.")
        exit(1)


def terraform_destroy_all_sites():
    """
    Destroy all sites by looping through all sites in reverse order
    """
    sites = list(reversed(get_settings()["sites"]["ordered"]))
    for site in sites:
        terraform_command_runner("destroy", [], "var", site)


def terraform_apply_all_sites():
    """
    Destroy all sites by looping through all sites
    """
    sites = list(get_settings()["sites"]["ordered"])
    for site in sites:
        if site != "inception":
            terraform_command_runner("init", ["-reconfigure"], "backend", site)
        terraform_command_runner("apply", [], "var", site)


def terraform_init_all_sites():
    """
    Init all sites by looping through all sites (inception doesn't need a init)
    """
    sites = list(get_settings()["sites"]["ordered"])
    if "inception" in sites:
        sites.remove("inception")
    for site in sites:
        print(Fore.YELLOW + f"Initializing {site} site")
        terraform_command_runner("init", site, options=["-reconfigure"],
                                 needed_options="backend")

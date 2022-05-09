import os
import signal
import subprocess
from colorama import Fore
from terrabuttler.settings import get_settings

# Values from Config
backend_dir = os.path.realpath(get_settings()["locations"]["backend_dir"])
environment_file = os.path.realpath(get_settings()
                                    ["locations"]["environment_file"])
inception_dir = os.path.realpath(get_settings()["locations"]["inception_dir"])
templates_dir = os.path.realpath(get_settings()["locations"]["templates_dir"])
variables_dir = os.path.realpath(get_settings()["locations"]["variables_dir"])


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
        needed_args = "backend"
    elif command == "plan" or command == "apply":
        needed_args = "var"

    args = terraform_args_builder(needed_args, site, backend_dir,
                                  variables_dir)
    return " ".join(args)


def terraform_args_builder(needed_args, site, backend_dir, var_dir):
    """
    Create array of needed arguments for backend or var files
    """
    from terrabuttler.env import get_current_env
    env = get_current_env()

    if needed_args == "backend":
        if site == "inception":
            return ["-backend-config",
                    f"{backend_dir}/pl-dev-inception.tfvars"]
        else:
            return ["-backend-config",
                    f"{backend_dir}/pl-{env}-{site}.tfvars"]

    elif needed_args == "var":
        return ["-var-file", f"{variables_dir}/global.tfvars",
                "-var-file", f"{variables_dir}/pl-{env}.tfvars",
                "-var-file", f"{variables_dir}/pl-{env}-{site}.tfvars"
                ]

    return []


def terraform_command_builder(command, args, needed_args, site,
                              backend_dir, var_dir):
    """
    Create the command to run terraform
    """
    base_command = ["terraform", command]

    base_command += args
    base_command += terraform_args_builder(needed_args, site, backend_dir,
                                           var_dir)

    return base_command


def terraform_command_runner(command, args, needed_args, site):
    """
    Run tfenv and run the terraform command
    """
    from terrabuttler.env import get_current_env
    site_dir = os.path.realpath(f"site_{site}")
    env = get_current_env()

    setup_tfenv(site_dir)

    command = terraform_command_builder(command, args, needed_args, site,
                                        backend_dir, variables_dir)
    try:
        p = subprocess.Popen(args=command, cwd=site_dir)
        p.wait()
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
        terraform_command_runner("init", ["-reconfigure"], "backend", site)

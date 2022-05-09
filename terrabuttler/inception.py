from colorama import Fore
from os import path
import subprocess


def inception_init_check():
    site_dir = path.realpath("site_inception")

    if (path.exists(f"{site_dir}/.terraform") and
            path.exists(f"{site_dir}/.terraform/environment")):
        return True
    return False


def inception_init_needed():
    if not inception_init_check():
        print(Fore.RED + "The initialization hasn't been made yet.\n\n"
              "Please execute the following command to initialize it:\n"
              "./run.py init")
        exit(1)


def inception_init():
    from terrabuttler.env import reload_direnv
    from terrabuttler.settings import get_settings
    site_dir = path.realpath(get_settings()["locations"]["inception_dir"])
    backend_dir = path.realpath(get_settings()["locations"]["backend_dir"])

    if not inception_init_check():
        try:
            subprocess.run(args=["terraform", "init", "-backend-config",
                                 f"{backend_dir}/pl-dev-inception.tfvars"],
                           cwd=site_dir, stdout=subprocess.DEVNULL,
                           stderr=subprocess.STDOUT)
        except subprocess.CalledProcessError:
            print(Fore.RED + "There was an error while doing the initializing")
            exit(1)

        try:
            with open(f"{site_dir}/.terraform/environment", "w") as f:
                f.write("dev")
        except FileNotFoundError:
            print(Fore.RED + "The file that manages the environments could not"
                  " be created.")
            exit(1)

        reload_direnv()
        print(Fore.GREEN + "The initialization was successfull!")

    else:
        print(Fore.YELLOW + "The initialization was already done.")

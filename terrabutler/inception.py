from terrabutler.utils import paths
from colorama import Fore
from os import path
from sys import exit
import subprocess


def inception_init_check():
    dir = paths["inception"]

    if (path.exists(f"{dir}/.terraform") and
            path.exists(f"{dir}/.terraform/environment")):
        return True
    return False


def inception_init_needed():
    if not inception_init_check():
        print(Fore.RED + "The initialization hasn't been made yet.\n\n"
              "Please execute the following command to initialize it:\n"
              "./run.py init")
        exit(1)


def inception_init():
    from terrabutler.settings import get_settings
    org = get_settings()["general"]["organization"]
    default_env_name = get_settings()["environments"]["default"]["name"]
    inception_dir = paths["inception"]
    backend_dir = paths["backends"]

    if not inception_init_check():
        try:
            subprocess.run(args=["terraform", "init", "-backend-config",
                                 f"{backend_dir}/{org}-{default_env_name}"
                                 "-inception.tfvars"],
                           cwd=inception_dir, stdout=subprocess.DEVNULL,
                           stderr=subprocess.STDOUT)
        except subprocess.CalledProcessError:
            print(Fore.RED + "There was an error while doing the initializing")
            exit(1)

        try:
            with open(f"{inception_dir}/.terraform/environment", "w") as f:
                f.write(default_env_name)
        except FileNotFoundError:
            print(Fore.RED + "The file that manages the environments could not"
                  " be created.")
            exit(1)

        print(Fore.GREEN + "The initialization was successfull!")

    else:
        print(Fore.YELLOW + "The initialization was already done.")

from colorama import Fore
from os import getenv

ROOT_PATH = getenv("TERRABUTLER_ROOT")
if getenv("TERRABUTLER_ENABLE") != "true":
    print(Fore.YELLOW + "Terrabutler is not currently enabled on this"
          " folder. Please set 'TERRABUTLER_ENABLE' in your environment"
          " to true to enable it." + Fore.RESET)
    exit(1)

paths = {
    "backends": ROOT_PATH + "/configs/backends",
    "environment": ROOT_PATH + "/site_inception/.terraform/environment",
    "inception": ROOT_PATH + "/site_inception",
    "root": ROOT_PATH,
    "settings": ROOT_PATH + "/configs/settings.yml",
    "templates": ROOT_PATH + "/configs/templates",
    "variables": ROOT_PATH + "/configs/variables"
}

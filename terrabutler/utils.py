from terrabutler.requirements import check_requirements
from os import getenv

check_requirements()
ROOT_PATH = getenv("TERRABUTLER_ROOT")

paths = {
    "backends": ROOT_PATH + "/configs/backends",
    "environment": ROOT_PATH + "/site_inception/.terraform/environment",
    "inception": ROOT_PATH + "/site_inception",
    "root": ROOT_PATH,
    "settings": ROOT_PATH + "/configs/settings.yml",
    "templates": ROOT_PATH + "/configs/templates",
    "variables": ROOT_PATH + "/configs/variables"
}

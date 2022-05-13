from terrabutler.requirements import check_requirements
from sys import exit
from semantic_version import Version
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


def is_semantic_version(version):
    """
    Check if the version corresponds to the semantic versioning.
    """
    try:
        Version(version)
    except ValueError:
        return False
    except Exception as e:
        print(f"There was an error while parsing version: {e}")
        exit(1)
    return True

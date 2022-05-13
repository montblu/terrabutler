from terrabutler.utils import paths
from colorama import Fore
from os import path
from schema import Schema, SchemaError
import yaml

PATH = paths["settings"]
SCHEMA = Schema({
    "general": {
        "organization": str,
        "secrets_key_id": str
    },
    "sites": {
        "ordered": list
    },
    "environments": {
        "default": {
            "domain": str,
            "name": str,
            "profile_name": str,
            "region": str
        },
        "permanent": list,
        "temporary": {
            "secrets": {
                "firebase_credentials": str,
                "mail_password": str
            }
        }
    }
})


def get_settings():
    """
    Returns the settings object
    """
    with open(PATH) as f:
        return yaml.safe_load(f)


def validate_settings():
    """
    Validade settings file
    """
    try:
        with open(PATH, "r") as f:
            configuration = yaml.safe_load(f)
    except FileNotFoundError:
        print(f"File {PATH} not found.  Aborting")
        exit(1)
    except OSError:
        print(f"OS error occurred trying to open {PATH}")
        exit(1)
    except Exception as err:
        print(f"Unexpected error when reading {PATH}: {err}")
        exit(1)

    try:
        SCHEMA.validate(configuration)
    except SchemaError as se:
        print("Your settings file is not using the needed values:"
              f" {se}")
        exit(1)


def write_settings(yaml_file):
    """
    Write settings file
    """
    with open(PATH, 'w') as f:
        yaml.safe_dump(yaml_file, f, default_flow_style=False)

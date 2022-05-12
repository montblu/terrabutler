from terrabutler import paths
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


def check_settings():
    """
    Check if the settings file exists
    """
    if not path.exists(PATH):
        print(Fore.YELLOW + "The settings file does not exist."
              "\n\nPlease create a 'settings.yml' file inside the 'configs'"
              " folder.")
        exit(1)


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
    with open(PATH) as f:
        configuration = yaml.safe_load(f)

    try:
        SCHEMA.validate(configuration)
    except SchemaError as se:
        raise se


def write_settings(yaml_file):
    """
    Write settings file
    """
    with open(PATH, 'w') as f:
        yaml.safe_dump(yaml_file, f, default_flow_style=False)

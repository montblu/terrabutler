from os import listdir
import boto3
from base64 import b64encode
from jinja2 import (
    Environment,
    FileSystemLoader
)
from secrets import choice
from string import (
    ascii_letters,
    digits
)
from terrabutler.settings import get_settings
from terrabutler.utils import paths


REGION = get_settings()["environments"]["default"]["region"]
ORG = get_settings()["general"]["organization"]
KEY_ID = get_settings()["general"]["secrets_key_id"]


def generate_var_files(env):
    """
    Create a variables files for a given environment
    """
    templates = listdir(paths["templates"])
    file_loader = FileSystemLoader(paths["templates"])
    environment = Environment(loader=file_loader)
    sites = list(get_settings()["sites"]["ordered"])
    firebase_credentials = (get_settings()["environments"]["temporary"]
                            ["secrets"]["firebase_credentials"])
    mail_password = (get_settings()["environments"]["temporary"]
                     ["secrets"]["mail_password"])
    if "inception" in sites:
        sites.remove("inception")
    for template in templates:
        temp = environment.get_template(template)
        output = temp.render(env=env,
                             generate_encrypted_password=(
                               generate_encrypted_password
                             ),
                             sites=sorted(sites),
                             mail_password=mail_password,
                             firebase_credentials=firebase_credentials)
        name = template.replace(".j2", "")
        if name == 'env':
            with open(f"{paths['variables']}/{ORG}-{env}.tfvars", "w") as fh:
                fh.write(output)
        else:
            with open(f"{paths['variables']}/{ORG}-{env}-{name}"
                      ".tfvars", "w") as fh:
                fh.write(output)


def generate_password(size):
    """
    Generate a password with a given size
    """
    characters = ascii_letters + digits
    password = "".join(choice(characters) for i in range(size))
    return password


def encrypt_password(password):
    """
    Encrypt password with AWS KMS
    """
    environment = boto3.session.Session(profile_name=f"{ORG}-dev",
                                        region_name=REGION)
    kms = environment.client("kms")
    encrypted = kms.encrypt(KeyId=KEY_ID, Plaintext=password)
    password_encrypted = encrypted[u'CiphertextBlob']
    password_encoded = b64encode(password_encrypted)
    return password_encoded.decode("utf-8")


def generate_encrypted_password(size):
    """
    Generate random password with a desired size and encrypt it with AWS KMS
    """
    password = generate_password(size)
    return encrypt_password(password)

from os import (
    listdir,
    path
)
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
from terrabuttler.settings import get_settings


REGION = get_settings()["environments"]["default"]["region"]
TEMPLATES_DIR = path.realpath(get_settings()["locations"]["templates_dir"])
VARIABLES_DIR = path.realpath(get_settings()["locations"]["variables_dir"])
ORG = get_settings()["general"]["organization"]
KEY_ID = get_settings()["general"]["secrets_key_id"]


def generate_var_files(env):
    """
    Create a variables files for a given environment
    """
    templates = listdir(TEMPLATES_DIR)
    file_loader = FileSystemLoader(TEMPLATES_DIR)
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
            with open(f"{VARIABLES_DIR}/{ORG}-{env}.tfvars", "w") as fh:
                fh.write(output)
        else:
            with open(f"{VARIABLES_DIR}/{ORG}-{env}-{name}.tfvars", "w")as fh:
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
    environment = boto3.session.Session(profile_name="pl-dev",
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

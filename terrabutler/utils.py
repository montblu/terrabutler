from os import getenv

ROOT_PATH = getenv("TERRABUTLER_ROOT")

paths = {
    "root": ROOT_PATH,
    "backends": ROOT_PATH + "/configs/backends",
    "templates": ROOT_PATH + "/configs/templates",
    "variables": ROOT_PATH + "/configs/variables",
    "environment": ROOT_PATH + "/site_inception/.terraform/environment",
    "inception": ROOT_PATH + "/site_inception",
}

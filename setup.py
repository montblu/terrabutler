import os
from setuptools import setup


# Utility function to read the README file.
# Used for the long_description.  It's nice, because now 1) we have a top level
# README file and 2) it's easier to type in the README file than to put a raw
# string in below ...
def read(fname):
    return open(os.path.join(os.path.dirname(__file__), fname)).read()


setup(
    name="terrabutler",
    version="1.0.0",
    author="AnguloSólido",
    description=("A tool to manage Terraform projects easier."),
    long_description=read("README.md"),
    license="GPL-3.0",
    keywords="terraform manager",
    url="https://github.com/angulo-solido/terrabutler",
    scripts=["bin/terrabutler"],
    packages=["terrabutler"],
    install_requires=[
        "boto3>=1.20.0",
        "botocore>=1.20.0",
        "click>=8.0.0",
        "jmespath>=0.10.0",
        "mccabe>=0.6.0",
        "pycodestyle>=2.8.0",
        "pyflakes>=2.4.0",
        "python-dateutil>=2.8.0",
        "s3transfer>=0.5.0",
        "six>=1.16.0",
        "urllib3>=1.26.0"
    ],
    python_requires=">= 3.6",
    classifiers=[
        "Development Status :: 4 - Beta",
        "Intended Audience :: Developers",
        "Intended Audience :: System Administrators",
        "License :: OSI Approved :: GNU General Public License v3 (GPLv3)",
        "Natural Language :: English",
        "Programming Language :: Python",
        "Programming Language :: Python :: 3",
        "Programming Language :: Python :: 3.6",
        "Programming Language :: Python :: 3.7",
        "Programming Language :: Python :: 3.8",
        "Programming Language :: Python :: 3.9",
        "Programming Language :: Python :: 3.10",
        "Topic :: Utilities"
    ],
    project_urls={
        "Source": "https://github.com/angulo-solido/terrabutler",
        "Changelog": "https://github.com/angulo-solido/terrabutler/releases",
    },
)

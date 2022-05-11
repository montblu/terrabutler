# Terrabutler

## Requirements

The tools required to run the application are:
- [direnv](https://direnv.net/)
- [tfenv](https://github.com/tfutils/tfenv)

### Installation of direnv

Follow the officials docs to install **direnv** [here](https://direnv.net/docs/installation.html)

### Installation of tfenv

Follow the officials docs to install **tfenv** [here](https://github.com/tfutils/tfenv#installation)

## Install

1. Download the installer 
`curl -o terrabutler.tar.gz https://terrabutler-public.s3.amazonaws.com/releases/terrabutler-linux-x86_64-latest.tar.gz`

2.Create folder to extract the installer
`mkdir -p terrabutler`

3. Extract the installer
`tar -xf terrabutler.tar.gz -C terrabutler`

3. Run the installer
`sudo ./terrabutler/install`

4. Remove the installer
`rm terrabutler*`
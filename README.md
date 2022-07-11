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
```
curl -o /tmp/terrabutler.tar.gz https://terrabutler-public.s3.amazonaws.com/releases/terrabutler-linux-x86_64-latest.tar.gz
```

2. Extract the installer
```
tar -xf /tmp/terrabutler.tar.gz -C /tmp
```

3. Run the installer
```
sudo /tmp/terrabutler/install
```
or without sudo:
```
/tmp/terrabutler/install -i ~/.local/share/terrabutler -b ~/.local/bin
```

4. Remove the archive and the installer
```
rm -rf /tmp/terrabutler*
```

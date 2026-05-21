# Installation

Before proceeding make sure that you have the [requirements](requirements.md).

## Installing from the binaries

### Downloading the binaries

All the binaries versions are available inside the [releases pages](https://github.com/l58193/terrabutler/releases/)

To download the latest binaries run the following command:

``` shell
https://github.com/l58193/terrabutler/releases/download/v1.0.0/terrabutler_Darwin_arm64.tar.gz
wget -qO- https://github.com/l58193/terrabutler/releases/download/<VERSION>/terrabutler_<OS>_<ARCH>_<VERSION>.tar.gz | tar -zxvf - terrabutler
```

Where `<VERSION>` is the version of the release.

For example, to download **Terrabutler v1.0.0** for Linux x64, just run:

``` shell
wget -qO- https://github.com/l58193/terrabutler/releases/download/v1.0.0/terrabutler_Linux_x86_64.tar.gz | tar -zxvf - terrabutler
```

To download **Terrabutler v1.0.0** for MacOS arm64, just run:

``` shell
wget -qO- https://github.com/l58193/terrabutler/releases/download/v1.0.0/terrabutler_Darwin_arm64.tar.gz | tar -zxvf - terrabutler
```

> [!WARNING]
    In case you have installed downloaded the archive via Safari you won't be able to use Terrabutler as the file will be marked as "quarantined" and you won't be able to use it as stated in this [issue](https://github.com/borgbackup/borg/issues/5622#issuecomment-774617595). 

### Install the binaries

To install the binaries into your system simply run the installer script inside the `terrabutler` folder:

``` shell
sudo terrabutler/install
```

All the binaries will be placed inside the `/usr/local/share/terrabutler` folder and the bin inside the `/usr/local/bin` folder


> [!TIP]
    If you wanna set the location where terrabutler will be installed you can define it by passing arguments when running the install script.
    This arguments can be seen by running:
    
    ``` shell
    terrabutler/install -h
    ```
    
> [!TIP]
    Example of installing for local user only **(no need to run the install script as sudo)**:
    
    ``` shell
    terrabutler/install -i ~/.local/share/terrabutler -b ~/.local/bin
    ```

## Building from source

To build the package from source it need to have [GO](https://go.dev/doc/install) installed, at least version 1.26.0.


### Building locally the project

Download the source of Terrabutler.

Inside the terrabutler folder, do:
 ``` shell
 go mod tidy
 ``` 
 And run then:
 ```shell
 go build
 ```

Now a executable named terrabutler will be created and ready to be executed.

You can also run ```go install``` install the terrabutler in your system.

### Using go install

If the environment variable PATH is setup correctly during the installation of [GO](https://go.dev/doc/install), you can build terrabutler from the source and it becomes a executable in your system.

You can do:
``` shell 
go install github.com/l58193/terrabutler@Rewrite-Go
```

And now ```terrabutler``` is now installed in your system.

### Common problems with go install 

> [!WARNING]
After using ```go install``` and ```terrabutler``` it's not a valid command, verify if your PATH value for Go matches with your Go installation path for binaries.

Normally it should be ```GOPATH/bin```, where you can consult the GOPATH value with ```go env GOPATH```.


## Check if the installation was successful

You should be able to run:

``` shell
terrabutler -h
```

And the output should be the help menu:


``` shell
NAME:
   terrabutler - The utility that helps keeping your IaC in one piece

USAGE:
   terrabutler [OPTIONS] COMMAND [ARGS]...

COMMANDS:
   version  Show version and exit
   env      Manage environments
   init     Initialize the manager
   tf       Manage terraform commands

GLOBAL OPTIONS:
  --help, -H, -h       Show this message and exit
```

If the output is not similar to the one above, then **something went wrong during the installation**.

## How to update to a newer version?

### Installed with the binaries

If you have Terrabutler already installed and want to update to a newer version just do all the installation process again. When you will be running the installer script it should prompt if you want to upgrade the Terrabutler, as we can see below:

```shell
Found preexisting Terrabutler installation: /usr/share/terrabutler.
Do you want to replace it? [y/N]
```

Just press y and press enter and Terrabutler should be updated to the version that you downloaded.

> [!CAUTION]
In case you have installed the Terrabutler earlier with different locations, you will need to pass them, otherwise the install script won't prompt you to update it.

### Installed by building locally the project

Download the newer version of Terrabutler and repeat the [process](#building-locally-the-project).

### Installed using go install

If used go install, it's only need to run the [command again](#using-go-install).

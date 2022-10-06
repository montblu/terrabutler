# Installation

Before proceeding make sure that you have the [requirements](requirements.md).

## Download the binaries and extract them

All the binaries versions are available inside the [releases pages](https://github.com/angulo-solido/terrabutler/releases)

To download the latest binaries run the following command:

``` shell
wget -qO- https://terrabutler-public.s3.amazonaws.com/releases/terrabutler-linux-x86_64-latest.tar.gz | tar -zxvf - terrabutler
```

???+ tip
    If you wanna get a specific version just change the `<VERSION>` in the command below. All the **Terrabutler releases** are
    available in the [repository releases pages](https://github.com/angulo-solido/terrabutler/releases)

    ``` shell
    wget -qO- https://terrabutler-public.s3.amazonaws.com/releases/terrabutler-linux-x86_64-<VERSION>.tar.gz | tar -zxvf - terrabutler
    ```

    For example, to download **Terrabutler v0.1.0**, just run:

    ```
    wget -qO- https://terrabutler-public.s3.amazonaws.com/releases/terrabutler-linux-x86_64-v0.1.0.tar.gz | tar -zxvf - terrabutler
    ```

## Install the binaries

To install the binaries into your system simply run the installer script inside the `terrabutler` folder:

``` shell
sudo terrabutler/install
```

All the binaries will be placed inside the `/usr/share/terrabutler` folder and the bin inside the `/usr/bin` folder

???+ tip
    If you wanna set the location where terrabutler will be installed you can define it by passing arguments when running the install script.
    This arguments can be seen by running:
    
    ``` shell
    terrabutler/install -h
    ```

    Example of installing for local user only **(no need to run the install script as sudo)**:
    
    ``` shell
    terrabutler/install -i ~/.local/share/terrabutler -b ~/.local/bin
    ```

## Check if the installation was successful

You should be able to run:

``` shell
terrabutler --version
```

and the output should be:


``` shell
Terrabutler: v0.1.0
```

If the output is not similar to the one above, then **something went wrong during the installation**.

## How to update to a newer version?

If you have **Terrabutler** already installed and want to update to a newer version just do all the installation process again.
When you will be running the installer script it should prompt if you want to upgrade the **Terrabutler**, as we can see below:

``` shell
Found preexisting Terrabutler installation: /usr/share/terrabutler.
Do you want to replace it? [y/N]
```

Just press `y` and press *enter* and **Terrabutler** should be updated to the version that you downloaded.

???+ danger
    In case you have installed the **Terrabutler** earlier with different locations, you will need to pass them, otherwise the
    install script **won't prompt** you to update it.

# Installation

## Download the binaries and extract them

All the binaries versions are available inside the [releases pages](https://github.com/angulo-solido/terrabutler/releases)

To download the latest binaries run the following command:

``` shell
wget -qO- https://terrabutler-public.s3.amazonaws.com/releases/terrabutler-linux-x86_64-latest.tar.gz | tar -zxvf - terrabutler
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
    
    ```
    terrabutler/install -h
    ```

    Example of installing for local user only **(no need to run the install script as sudo)**:
    
    ```
    terrabutler/install -i ~/.local/share/terrabutler -b ~/.local/bin
    ```

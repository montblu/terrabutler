# Installation

Before proceeding make sure that you have the [requirements](requirements.md).

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

After using ```go install``` and ```terrabutler``` it's not a valid command, verify if your PATH value for Go matches with your Go installation path for binaries.

Normally it should be GOPATH/bin, where you can consult the GOPATH value with ```go env GOPATH```.


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

### Installed by building locally the project

Download the newer version of Terrabutler and repeat the [process](#building-locally-the-project).

### Installed using go install

If used go install, it's only need to run the [command again](#using-go-install).

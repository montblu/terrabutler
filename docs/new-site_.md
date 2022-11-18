# Add a new site_

## Create a new directory 

In the project folder root, create a new diretory called `site_<site-name>`:

```shell
$ mkdir site_<site-name>
```
???+ tip
  Replace `<site-name>`  with your new site name!

## Configure the new `site_`


### Add the new site to `site_inception` variables file

Add `<site-name>` to the following file, in `inception_projects = []` line:

```
./configs/variables/<project>-<env>-inception.tfvars
```

### Add the new site to Terrabutler settings

Add `<site-name>` to the following file, in `sites: ordered: - <site-name>` line:

```
./configs/settings.yml
```

### Add a new variable file to the Variables Folder

Create a new file called `<project-name>-<env>-<site-name>.tfvars` to the following directory:

```
./configs/variables/<project-name>-<env>-<site-name>.tfvars
```

### Perform an apply in `site_inception`

Run the following command, to update the new configuration on `site_inception`:

```shell
$ terrabutler tf -site inception apply
```

## Add files to the new `site_`

### Add terraform files

Run the following commands to copy the following Terraform files to the site_`<site-name>`:

```shell
$ cd site_<site-name>
```

```shell
$ cp ../site_inception/terraform.tf .
```

```shell
$ cp ../site_inception/provider.tf .
```

### Create Symbolic Links

Run the following commands to create the Symbolic Links:

```shell
$ cd site_<site-name>
```

```shell
$ ln -s ../globals/data.tf ./data_global.tf
```
```shell
$ ln -s ../globals/locals.tf ./locals_globals.tf
```
```shell
$ $ ln -s ../globals/variables.tf ./variables_globals.tf
```
### Perform an init in `site_`

Run the following command:

```shell
$ terrabutler tf -site data init 
```

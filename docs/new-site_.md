# New site_

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
../configs/variables/<project>-<env>-inception.tfvars
```

### Add the new site to Terrabutler settings

Add `<site-name>` to the following file, in `sites: ordered: - <site-name>` line:

```
../configs/settings.yml
```

### Add a new variable file to the Variables Folder

Create a new file called `<project-name>-<env>-<site-name>.tfvars` to the following directory:

```
../configs/variables/<project-name>-<env>-<site-name>.tfvars
```

### Perform an apply in `site_inception`

Run the following command, to update the new configuration on `site_inception`:

```shell
$ terrabutler tf -site inception apply
```

## Add files to the new `site_`

### Add terraform files

Create the Terraform files inside the site_`<site-name>`

```
- site_<site-name>
  - provider.tf
  - terraform.tf
  - variables-global.tf (Symbolic Link) -> site_inception
```

### Perform an init in `site_`

Run the following command, with the `site_<site-name>` backend config:

```shell
$ terrabutler tf -site data init -backend-config="./configs/backends/<project_name>-<environment_name>-<site-name>.tfvars"
```

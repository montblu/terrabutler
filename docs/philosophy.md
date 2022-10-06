# Philosophy

Before settling for **Terrabutler**, it's a good idea to understand the philosophy behind the project,
in order to make sure it aligns with your goals. This page explains the problem and the solution.

## Problem

All our IaC projects are backed by Terraform. Instead of having a folder with every resource, we have
created different folders that have resources from the same category, and we call them `sites`. For
example inside the `network site` we will have the creation of `VPC`, `SG` and `WAF`, so this site is
where we have all the resources related to network. By splitting the code into various `sites` we have
smaller plans and it's more easy to manage the code. By having all the code divided into various `sites`
it's a bit more difficult to manage all the terraform variable files. And we wanted all of this for each
**environment**, so it needs to be divided into different environments, for example production and development.
All of the `sites` state needs to be have a remote state and with a lock function, to prevent usage at the same
time.

## Solution

Create a IaC with the `inception site` that will be responsible to manage the project environments via Terraform
workspaces tool and where the backends for each site will be created. Use **Terrabutler** to manage all the IAC
between all the sites.
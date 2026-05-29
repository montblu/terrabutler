<div align="center">

<img src="assets/logo.png" align="center"/>

# Terrabutler

**The utility that helps keeping your IaC in one piece**

![GitHub release (latest by date)](https://img.shields.io/github/v/release/montblu/terrabutler?color=8956c4&label=Latest%20Version&logo=Github&style=for-the-badge)
![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/montblu/terrabutler/release.yml?color=8956c4&logo=Github&style=for-the-badge)
![GitHub](https://img.shields.io/github/license/montblu/terrabutler?color=8956c4&logo=Github&style=for-the-badge)
![GitHub Repo stars](https://img.shields.io/github/stars/montblu/terrabutler?color=8956c4&label=Repo%20Stars&style=for-the-badge)

---

## What is Terrabutler?

Terrabutler is a **wrapper** written in [Python](https://www.python.org/) that helps maintaining IaC (Infrastructure as code) projects
using [Terraform](https://www.terraform.io/) by managing the **environments**.

## What is this?

This a rewrite of Terrabutler written in [GO](https://go.dev/), where it aims to implement all it's functionalities, support tests unit, and with better performance and scalability.

In this current version it supports all the principal functionalities except the connection to the Amazon S3 used in some parts of the old Terrabutler.
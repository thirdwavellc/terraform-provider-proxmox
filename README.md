# Terraform Provider Proxmox

This project defines a Terraform Provider plugin for manipulating resources in
the Proxmox VE open-source virtualization platform.

## Work In Progress

This repo is currently a work in progress, with limited functionality provided
at this point.

## Setup

We are using [dep](https://github.com/golang/dep) to manage go dependencies.
Once you have dep installed, to install the project's dependencies:

```bash
$ dep ensure
```

This will install them under `vendor`.

To build the project:

```bash
$ make
```

Currently this only builds for your local OS, however this may change in the
future to match how the build is done in
[go-proxmox](https://github.com/thirdwavellc/go-proxmox).

To install the plugin, follow the [official
docs](https://www.terraform.io/docs/configuration/providers.html#third-party-plugins).

The next time you run `terraform init` it should load the plugin.

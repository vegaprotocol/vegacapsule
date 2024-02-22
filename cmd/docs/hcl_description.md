# Capsule configuration docs

Capsule is a tool that allows you to run a custom network simulation locally on single a machine. It is an incredibly useful tool for anybody who wants to try experimenting on or with a Vega network without using a real network.

The configuration for Capsule is used to generate and bootstrap commands, and can be customised to fit your personal need. Under the hood, Capsule uses this configuration to generate a new network and store all its files in a single directory. You can then use [Nomad](https://www.nomadproject.io/) to deploy all generated services from the generation step. Nomad is built into the Capsule binaries.

The configuration uses the [HCL](https://github.com/hashicorp/hcl) language syntax, which is also used, for example, by [Terraform](https://www.terraform.io/).

This document explains all possible configuration options in Capsule.

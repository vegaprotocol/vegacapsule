#Capsule configuration docs

Capsule it's not a real network deployment tool - it is rather a tool that allows to run a custom network simulation locally on single a machine.
This means that it is a incredibly useful tool for anybody who wants to try Vega network without using a real network.

Capsule configuration is used by vegacapsule CLI network to generate and bootstrap commands and can be customised to personal need.
Under the hood Capsule will use this configuration to generate a new network a stores all it's files in a single directory and then
it uses [Nomad](https://www.nomadproject.io/) to deploy all generated services from the generation step.

The configuration uses the [HCL](https://github.com/hashicorp/hcl) language syntax, which is also used, for example, by [Terraform](https://www.terraform.io/).

This document explains all possible configuration options in Capsule.
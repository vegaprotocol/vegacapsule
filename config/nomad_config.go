package config

/*
description: |

	Allows the user to configure a [Nomad job](https://developer.hashicorp.com/nomad/docs/job-specification) definition to be run on Capsule.

example:

	type: hcl
	value: |
			nomad_job "clef" {
				job_template = "/path-to/nomad-job.tmpl"
			}
*/
type NomadConfig struct {
	/*
		description: |
			Name of the Nomad job.
		example:
			type: hcl
			value: |
				nomad_job "service-1" {
					...
				}
	*/
	Name string `hcl:"name,label"`

	/*
		description: |
			[Go template](templates.md) of a Nomad job template.

			The [nomad.PreGenerateTemplateCtx](templates.md#nomadpregeneratetemplatectx) can be used in the template. Example [example](jobs/clef.tmpl).
		optional_if: job_template_file
		note: |
				It is recommended that you use `job_template_file` param instead.
				If both `job_template` and `job_template_file` are defined, then `job_template`
				overrides `job_template_file`.
		examples:
			- type: hcl
			  value: |
						job_template = <<EOH
							...
						EOH

	*/
	JobTemplate *string `hcl:"job_template,optional"`

	/*
		description: |
			Same as `job_template` but it allows the user to link the Nomad job template as an external file.
		examples:
			- type: hcl
			  value: |
						job_template_file = "/your_path/nomad-job.tmpl"

	*/
	JobTemplateFile *string `hcl:"job_template_file,optional"`
}

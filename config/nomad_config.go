package config

type NomadConfig struct {
	Name            string  `hcl:"name,label"`
	JobTemplate     *string `hcl:"job_template,optional"`
	JobTemplateFile *string `hcl:"job_template_file,optional"`
}

package generator

import (
	"context"
	"fmt"

	"code.vegaprotocol.io/vegacapsule/config"
	"code.vegaprotocol.io/vegacapsule/generator/nomad"
)

func (g *Generator) startPreGenerateJobs(n config.NodeConfig, index int) ([]string, error) {
	templates, err := g.templatePreGenerateJobs(n.PreGenerate, index)
	if err != nil {
		return nil, fmt.Errorf("failed to template pre generate jobs for node set %q-%q: %w", n.Name, index, err)
	}
	preGenJobIDs, err := g.startNomadJobs(templates)
	if err != nil {
		return nil, fmt.Errorf("failed to start pre generate jobs for node set %s-%d: %w", n.Name, index, err)
	}

	return preGenJobIDs, nil
}

func (g *Generator) templatePreGenerateJobs(preGenConf *config.PreGenerate, index int) ([]string, error) {
	if preGenConf == nil {
		return []string{}, nil
	}

	jobTemplates := make([]string, 0, len(preGenConf.Nomad))
	for _, nc := range preGenConf.Nomad {
		if nc.JobTemplate == nil {
			continue
		}

		template, err := nomad.GeneratePreGenerateTemplate(*nc.JobTemplate, nomad.PreGenerateTemplateCtx{
			Name:  nc.Name,
			Index: index,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to template nomad job for pre generate %q: %w", nc.Name, err)
		}

		jobTemplates = append(jobTemplates, template.String())
	}

	return jobTemplates, nil
}

func (g *Generator) startNomadJobs(rawNomadJobs []string) ([]string, error) {
	if len(rawNomadJobs) == 0 {
		return rawNomadJobs, nil
	}

	jobs, err := g.jobRunner.RunRawNomadJobs(context.Background(), rawNomadJobs)
	if err != nil {
		return nil, fmt.Errorf("failed to run node set pre generate job: %w", err)
	}

	jobIDs := make([]string, 0, len(jobs))
	for _, j := range jobs {
		jobIDs = append(jobIDs, *j.ID)
	}

	return jobIDs, nil
}

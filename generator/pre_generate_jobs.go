package generator

import (
	"context"
	"fmt"
	"time"

	"code.vegaprotocol.io/vegacapsule/config"
	"code.vegaprotocol.io/vegacapsule/generator/nomad"
	"code.vegaprotocol.io/vegacapsule/types"
)

func (g *Generator) startPreGenerateJobs(n config.NodeConfig, index int) ([]types.NomadJob, error) {
	templates, err := g.templatePreGenerateJobs(n.PreGenerate, index)
	if err != nil {
		return nil, fmt.Errorf("failed to template pre generate jobs for node set %q-%q: %w", n.Name, index, err)
	}
	preGenJobs, err := g.startNomadJobs(templates)
	if err != nil {
		return nil, fmt.Errorf("failed to start pre generate jobs for node set %s-%d: %w", n.Name, index, err)
	}

	return preGenJobs, nil
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

func (g *Generator) startNomadJobs(rawNomadJobs []string) ([]types.NomadJob, error) {
	if len(rawNomadJobs) == 0 {
		return nil, nil
	}

	jobs, err := g.jobRunner.RunRawNomadJobs(context.Background(), rawNomadJobs)
	if err != nil {
		return nil, fmt.Errorf("failed to run node set pre generate job: %w", err)
	}

	jobIDs := make([]types.NomadJob, 0, len(jobs))
	for _, j := range jobs {
		jobIDs = append(jobIDs, types.NomadJob{
			ID:          *j.NomadJob.ID,
			NomadJobRaw: j.RawJob,
		})
	}

	return jobIDs, nil
}

func (g *Generator) stopNomadJobs() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*20)
	defer cancel()

	return g.jobRunner.StopNetwork(ctx, nil, false)
}

package types

type JobIDMap map[string]bool

func (jm JobIDMap) ToSlice() []string {
	slice := make([]string, 0, len(jm))
	for id := range jm {
		slice = append(slice, id)
	}
	return slice
}

type NetworkJobs struct {
	NodesSetsJobIDs JobIDMap
	ExtraJobIDs     JobIDMap
	FaucetJobID     string
	WalletJobID     string
}

func (nj *NetworkJobs) MergeNetworkJobs(toMerge NetworkJobs) {
	if toMerge.FaucetJobID != "" {
		nj.FaucetJobID = toMerge.FaucetJobID
	}

	if toMerge.WalletJobID != "" {
		nj.WalletJobID = toMerge.WalletJobID
	}

	for id := range toMerge.NodesSetsJobIDs {
		if nj.NodesSetsJobIDs == nil {
			nj.NodesSetsJobIDs = JobIDMap{}
		}

		nj.NodesSetsJobIDs[id] = true
	}

	for id := range toMerge.ExtraJobIDs {
		if nj.ExtraJobIDs == nil {
			nj.NodesSetsJobIDs = JobIDMap{}
		}

		nj.ExtraJobIDs[id] = true
	}
}

func (nj *NetworkJobs) Exists(jobID string) bool {
	if _, ok := nj.NodesSetsJobIDs[jobID]; ok {
		return true
	}
	if _, ok := nj.ExtraJobIDs[jobID]; ok {
		return true
	}
	if nj.FaucetJobID == jobID {
		return true
	}
	if nj.WalletJobID == jobID {
		return true
	}

	return false
}

func (nj *NetworkJobs) AddExtraJobIDs(ids []string) {
	if nj.ExtraJobIDs == nil {
		nj.ExtraJobIDs = JobIDMap{}
	}

	for _, id := range ids {
		nj.ExtraJobIDs[id] = true
	}
}

func (nj *NetworkJobs) RemoveRunningJobsIDs(ids []string) {
	for _, id := range ids {
		if nj.NodesSetsJobIDs != nil {
			delete(nj.NodesSetsJobIDs, id)
		}
		if nj.ExtraJobIDs != nil {
			delete(nj.ExtraJobIDs, id)
		}

		if nj.FaucetJobID == id {
			nj.FaucetJobID = ""
		}
		if nj.WalletJobID == id {
			nj.WalletJobID = ""
		}
	}
}

func (nj *NetworkJobs) ToSlice() []string {
	out := append(nj.NodesSetsJobIDs.ToSlice(), nj.ExtraJobIDs.ToSlice()...)

	if nj.FaucetJobID != "" {
		out = append(out, nj.FaucetJobID)
	}

	if nj.WalletJobID != "" {
		out = append(out, nj.WalletJobID)
	}

	return out
}

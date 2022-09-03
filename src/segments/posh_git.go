package segments

import "encoding/json"

const (
	poshGitEnv = "POSH_GIT_STATUS"
)

type poshGit struct {
	StashCount   int            `json:"StashCount"`
	AheadBy      int            `json:"AheadBy"`
	Index        *poshGitStatus `json:"Index"`
	RepoName     string         `json:"RepoName"`
	HasWorking   bool           `json:"HasWorking"`
	Branch       string         `json:"Branch"`
	HasIndex     bool           `json:"HasIndex"`
	GitDir       string         `json:"GitDir"`
	BehindBy     int            `json:"BehindBy"`
	HasUntracked bool           `json:"HasUntracked"`
	Working      *poshGitStatus `json:"Working"`
	Upstream     string         `json:"Upstream"`
}

type poshGitStatus struct {
	Added    []string `json:"Added"`
	Modified []string `json:"Modified"`
	Deleted  []string `json:"Deleted"`
	Unmerged []string `json:"Unmerged"`
}

func (s *GitStatus) parsePoshGitStatus(p *poshGitStatus) {
	if p == nil {
		return
	}
	s.Added = len(p.Added)
	s.Deleted = len(p.Deleted)
	s.Modified = len(p.Modified)
	s.Unmerged = len(p.Unmerged)
}

func (g *Git) hasPoshGitStatus() bool {
	envStatus := g.env.Getenv(poshGitEnv)
	if len(envStatus) == 0 {
		return false
	}
	var posh poshGit
	err := json.Unmarshal([]byte(envStatus), &posh)
	if err != nil {
		return false
	}
	g.setDir(posh.GitDir)
	g.Working = &GitStatus{}
	g.Working.parsePoshGitStatus(posh.Working)
	g.Staging = &GitStatus{}
	g.Staging.parsePoshGitStatus(posh.Index)
	g.HEAD = posh.Branch
	g.StashCount = posh.StashCount
	g.Ahead = posh.AheadBy
	g.Behind = posh.BehindBy
	g.UpstreamGone = len(posh.Upstream) == 0
	g.Upstream = posh.Upstream
	g.setBranchStatus()
	if len(g.Upstream) != 0 && g.props.GetBool(FetchUpstreamIcon, false) {
		g.UpstreamIcon = g.getUpstreamIcon()
	}
	return true
}

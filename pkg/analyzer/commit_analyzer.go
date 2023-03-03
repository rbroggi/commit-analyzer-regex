package analyzer

import (
	"regexp"
	"strings"

	"github.com/go-semantic-release/semantic-release/v2/pkg/semrel"
)

var (
	CAVERSION              = "dev"
	commitPattern          = regexp.MustCompile(`^(\w*)(?:\((.*)\))?(!)?: (.*)$`)
	breakingPattern        = regexp.MustCompile("BREAKING CHANGES?")
	mentionedIssuesPattern = regexp.MustCompile(`#(\d+)`)
	mentionedUsersPattern  = regexp.MustCompile(`(?i)@([a-z\d]([a-z\d]|-[a-z\d])+)`)
)

func extractMentions(re *regexp.Regexp, s string) string {
	ret := make([]string, 0)
	for _, m := range re.FindAllStringSubmatch(s, -1) {
		ret = append(ret, m[1])
	}
	return strings.Join(ret, ",")
}

type DefaultCommitAnalyzer struct {
	MinorRegexp *regexp.Regexp
	PatchRegexp *regexp.Regexp
}

func (da *DefaultCommitAnalyzer) Init(config map[string]string) error {
	minorPattern := config["minor"]
	if minorPattern != "" {
		minorRegexp, err := regexp.Compile(minorPattern)
		if err != nil {
			return err
		}
		da.MinorRegexp = minorRegexp
	} else {
		da.MinorRegexp = regexp.MustCompile("feat")
	}

	patchPattern := config["patch"]
	if patchPattern != "" {
		patchRegexp, err := regexp.Compile(patchPattern)
		if err != nil {
			return err
		}
		da.PatchRegexp = patchRegexp
	} else {
		da.PatchRegexp = regexp.MustCompile("fix")
	}

	return nil
}

func (da *DefaultCommitAnalyzer) Name() string {
	return "commit-analyzer-regex"
}

func (da *DefaultCommitAnalyzer) Version() string {
	return CAVERSION
}

func (da *DefaultCommitAnalyzer) analyzeSingleCommit(rawCommit *semrel.RawCommit) *semrel.Commit {
	c := &semrel.Commit{
		SHA:         rawCommit.SHA,
		Raw:         strings.Split(rawCommit.RawMessage, "\n"),
		Change:      &semrel.Change{},
		Annotations: rawCommit.Annotations,
	}
	c.Annotations["mentioned_issues"] = extractMentions(mentionedIssuesPattern, rawCommit.RawMessage)
	c.Annotations["mentioned_users"] = extractMentions(mentionedUsersPattern, rawCommit.RawMessage)

	found := commitPattern.FindAllStringSubmatch(c.Raw[0], -1)
	if len(found) < 1 {
		return c
	}
	c.Type = strings.ToLower(found[0][1])
	c.Scope = found[0][2]
	breakingChange := found[0][3]
	c.Message = found[0][4]

	isMajorChange := breakingPattern.MatchString(rawCommit.RawMessage)
	isMinorChange := da.MinorRegexp.MatchString(c.Type)
	isPatchChange := da.PatchRegexp.MatchString(c.Type)

	if len(breakingChange) > 0 {
		isMajorChange = true
		isMinorChange = false
		isPatchChange = false
	}

	c.Change = &semrel.Change{
		Major: isMajorChange,
		Minor: isMinorChange,
		Patch: isPatchChange,
	}
	return c
}

func (da *DefaultCommitAnalyzer) Analyze(rawCommits []*semrel.RawCommit) []*semrel.Commit {
	ret := make([]*semrel.Commit, len(rawCommits))
	for i, c := range rawCommits {
		ret[i] = da.analyzeSingleCommit(c)
	}
	return ret
}

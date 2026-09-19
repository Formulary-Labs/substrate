// Package evidence provides URL classification for compliance artifact evidence links.
// It identifies which system an evidence URL points at and flags PENDING placeholders.
package evidence

import "strings"

// linkPatterns maps link types to substrings that identify them.
// More-specific patterns appear before broader ones so the first match wins.
var linkPatterns = []struct {
	linkType string
	keywords []string
}{
	{"Confluence", []string{"atlassian.net/wiki", "confluence"}},
	{"Jira", []string{"atlassian.net/browse", "atlassian.net/jira"}},
	{"GitHub", []string{"github.com"}},
	{"GitLab", []string{"gitlab"}},
	{"Google Doc", []string{"docs.google.com"}},
	{"Red Hat Docs", []string{"access.redhat.com", "docs.redhat.com"}},
	{"Red Hat", []string{"redhat.com"}},
	{"Miro", []string{"miro.com"}},
	{"ServiceNow", []string{"service-now.com"}},
	{"Splunk", []string{"splunkcloud.com"}},
	{"Grafana", []string{"grafana"}},
}

// ClassifyLink returns (linkType, isPending) for an evidence URL or placeholder.
// linkType is one of: "Confluence", "Jira", "GitHub", "GitLab", "Google Doc",
// "Red Hat Docs", "Red Hat", "Miro", "ServiceNow", "Splunk", "Grafana", "Other",
// or "PENDING" when isPending is true.
// isPending is true when url is empty or starts with "PENDING" (case-insensitive).
func ClassifyLink(url string) (linkType string, isPending bool) {
	s := strings.TrimSpace(url)
	if s == "" || strings.HasPrefix(strings.ToUpper(s), "PENDING") {
		return "PENDING", true
	}
	lower := strings.ToLower(s)
	for _, p := range linkPatterns {
		for _, kw := range p.keywords {
			if strings.Contains(lower, kw) {
				return p.linkType, false
			}
		}
	}
	return "Other", false
}

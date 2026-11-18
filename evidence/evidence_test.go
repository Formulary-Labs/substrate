package evidence_test

import (
	"testing"

	"github.com/Formulary-Labs/substrate/evidence"
)

func TestClassifyLink(t *testing.T) {
	tests := []struct {
		input       string
		wantType    string
		wantPending bool
	}{
		{"", "PENDING", true},
		{"PENDING", "PENDING", true},
		{"PENDING — evidence not yet collected", "PENDING", true},
		{"pending", "PENDING", true},
		{"https://github.com/org/repo/issues/1", "GitHub", false},
		{"https://example.atlassian.net/wiki/spaces/SEC/pages/1", "Confluence", false},
		{"https://example.atlassian.net/browse/SEC-123", "Jira", false},
		{"https://gitlab.com/org/repo/-/issues/1", "GitLab", false},
		{"https://docs.google.com/document/d/abc", "Google Doc", false},
		{"https://miro.com/board/abc", "Miro", false},
		{"https://company.service-now.com/tickets/123", "ServiceNow", false},
		{"https://company.splunkcloud.com/en-US/app", "Splunk", false},
		{"https://grafana.example.com/d/abc", "Grafana", false},
		{"https://example.com/some-document.pdf", "Other", false},
		// An employer or vendor domain with no entry above must classify as Other,
		// not as a special case: the vocabulary is publisher-neutral.
		{"https://blog.example-vendor.com/post", "Other", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			gotType, gotPending := evidence.ClassifyLink(tt.input)
			if gotType != tt.wantType {
				t.Errorf("ClassifyLink(%q) linkType = %q, want %q", tt.input, gotType, tt.wantType)
			}
			if gotPending != tt.wantPending {
				t.Errorf("ClassifyLink(%q) isPending = %v, want %v", tt.input, gotPending, tt.wantPending)
			}
		})
	}
}

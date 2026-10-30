package artifact

import (
	"fmt"
	"strings"
	"unicode"

	gemara "github.com/gemaraproj/go-gemara"
)

// UnifiedControl is one control, independent of whether the file was gemara or OSCAL.
type UnifiedControl struct {
	ID           string
	Title        string
	Group        string
	Objective    string
	Requirements []UnifiedRequirement
}

// UnifiedRequirement is one assessment requirement or OSCAL part other than the statement.
type UnifiedRequirement struct {
	ID   string
	Text string
}

// UnifiedGroup is a control family.
type UnifiedGroup struct {
	ID    string
	Title string
}

// UnifiedCatalog is a control catalog loaded from either interchange language.
// Gemara is set for gemara ControlCatalog files and is the original value,
// not a copy. OSCAL is set for native OSCAL catalogs.
type UnifiedCatalog struct {
	Format   Format
	ID       string
	Title    string
	Groups   []UnifiedGroup
	Controls []UnifiedControl
	Gemara   *ControlCatalog
	OSCAL    *OSCALCatalog
}

// LoadCatalog loads a gemara ControlCatalog or an OSCAL catalog.
// Other OSCAL models return an error; use LoadOSCAL for those.
func LoadCatalog(path string) (*UnifiedCatalog, error) {
	format, err := DetectFormat(path)
	if err != nil {
		return nil, err
	}
	switch format {
	case FormatGemara:
		cat, err := LoadControlCatalog(path)
		if err != nil {
			return nil, err
		}
		return unifiedFromGemara(cat), nil
	case FormatOSCAL:
		cat, err := LoadOSCALCatalog(path)
		if err != nil {
			return nil, err
		}
		return unifiedFromOSCAL(cat), nil
	default:
		return nil, fmt.Errorf("loading catalog %q: unrecognized format (expected a gemara ControlCatalog or an OSCAL catalog)", path)
	}
}

// AsControlCatalog returns the gemara catalog tools already understand.
// For an OSCAL source this is a projection, not the original file.
func (c *UnifiedCatalog) AsControlCatalog() *ControlCatalog {
	if c == nil {
		return nil
	}
	if c.Gemara != nil {
		return c.Gemara
	}
	out := &ControlCatalog{Title: c.Title}
	out.Metadata.Id = c.ID
	out.Metadata.Description = c.Title
	for _, g := range c.Groups {
		out.Groups = append(out.Groups, gemara.Group{Id: g.ID, Title: g.Title})
	}
	for _, ctrl := range c.Controls {
		projected := gemara.Control{
			Id:        ctrl.ID,
			Title:     ctrl.Title,
			Group:     ctrl.Group,
			Objective: ctrl.Objective,
		}
		for _, req := range ctrl.Requirements {
			projected.AssessmentRequirements = append(projected.AssessmentRequirements, gemara.AssessmentRequirement{
				Id:   req.ID,
				Text: req.Text,
			})
		}
		out.Controls = append(out.Controls, projected)
	}
	return out
}

func unifiedFromGemara(cat *ControlCatalog) *UnifiedCatalog {
	if cat == nil {
		return &UnifiedCatalog{Format: FormatGemara}
	}
	out := &UnifiedCatalog{
		Format: FormatGemara,
		ID:     cat.Metadata.Id,
		Title:  cat.Title,
		Gemara: cat,
	}
	for _, g := range cat.Groups {
		out.Groups = append(out.Groups, UnifiedGroup{ID: g.Id, Title: g.Title})
	}
	for _, ctrl := range cat.Controls {
		uc := UnifiedControl{
			ID:        ctrl.Id,
			Title:     ctrl.Title,
			Group:     ctrl.Group,
			Objective: ctrl.Objective,
		}
		for _, req := range ctrl.AssessmentRequirements {
			uc.Requirements = append(uc.Requirements, UnifiedRequirement{ID: req.Id, Text: req.Text})
		}
		out.Controls = append(out.Controls, uc)
	}
	return out
}

func unifiedFromOSCAL(cat *OSCALCatalog) *UnifiedCatalog {
	out := &UnifiedCatalog{
		Format: FormatOSCAL,
		OSCAL:  cat,
	}
	if cat == nil {
		return out
	}
	out.ID = oscalFrameworkID(cat)
	out.Title = cat.Metadata.Title
	if cat.Groups != nil {
		flattenOSCALGroups(out, *cat.Groups, "")
	}
	appendOSCALControls(out, cat.Controls, "")
	return out
}

func flattenOSCALGroups(out *UnifiedCatalog, groups []OSCALGroup, parent string) {
	for _, g := range groups {
		id := g.ID
		if id == "" {
			id = parent
		}
		if g.ID != "" || g.Title != "" {
			out.Groups = append(out.Groups, UnifiedGroup{ID: g.ID, Title: g.Title})
		}
		appendOSCALControls(out, g.Controls, id)
		if g.Groups != nil {
			flattenOSCALGroups(out, *g.Groups, id)
		}
	}
}

func appendOSCALControls(out *UnifiedCatalog, controls *[]OSCALControl, groupID string) {
	if controls == nil {
		return
	}
	for _, ctrl := range *controls {
		out.Controls = append(out.Controls, unifiedControl(ctrl, groupID))
		appendOSCALControls(out, ctrl.Controls, groupID)
	}
}

func unifiedControl(ctrl OSCALControl, groupID string) UnifiedControl {
	uc := UnifiedControl{
		ID:        ctrl.ID,
		Title:     ctrl.Title,
		Group:     groupID,
		Objective: oscalStatement(ctrl),
	}
	if ctrl.Parts != nil {
		for _, part := range *ctrl.Parts {
			if part.Name == "statement" {
				continue
			}
			if strings.TrimSpace(part.Prose) == "" {
				continue
			}
			id := part.ID
			if id == "" {
				id = part.Name
			}
			uc.Requirements = append(uc.Requirements, UnifiedRequirement{ID: id, Text: part.Prose})
		}
	}
	return uc
}

func oscalStatement(ctrl OSCALControl) string {
	if ctrl.Parts == nil {
		return ""
	}
	var fallback string
	for _, part := range *ctrl.Parts {
		if strings.TrimSpace(part.Prose) == "" {
			continue
		}
		if part.Name == "statement" {
			return part.Prose
		}
		if fallback == "" {
			fallback = part.Prose
		}
	}
	return fallback
}

func oscalFrameworkID(cat *OSCALCatalog) string {
	if cat.Metadata.DocumentIds != nil {
		for _, id := range *cat.Metadata.DocumentIds {
			if strings.TrimSpace(id.Identifier) != "" {
				return id.Identifier
			}
		}
	}
	if slug := slugTitle(cat.Metadata.Title); slug != "" {
		return slug
	}
	return cat.UUID
}

func slugTitle(title string) string {
	var b strings.Builder
	prevDash := false
	for _, r := range strings.ToLower(strings.TrimSpace(title)) {
		switch {
		case unicode.IsLetter(r) || unicode.IsDigit(r):
			b.WriteRune(r)
			prevDash = false
		default:
			if b.Len() > 0 && !prevDash {
				b.WriteByte('-')
				prevDash = true
			}
		}
	}
	return strings.Trim(b.String(), "-")
}

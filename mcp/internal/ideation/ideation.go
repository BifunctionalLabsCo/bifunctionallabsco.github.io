package ideation

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type Project struct {
	Name             string   `json:"name"`
	Purpose          string   `json:"purpose"`
	Audience         string   `json:"audience"`
	ProductType      string   `json:"productType"`
	Success          string   `json:"success"`
	InheritedTraits  []string `json:"inheritedTraits,omitempty"`
	DistinctTraits   []string `json:"distinctTraits,omitempty"`
	InitialWorkflows []string `json:"initialWorkflows,omitempty"`
	Constraints      []string `json:"constraints,omitempty"`
	Tone             string   `json:"tone,omitempty"`
}

type Intake struct {
	Complete  bool     `json:"complete"`
	Questions []string `json:"questions"`
	NextStep  string   `json:"nextStep"`
}

func SocraticIntake(p Project) Intake {
	questions := make([]string, 0, 8)
	if strings.TrimSpace(p.Purpose) == "" {
		questions = append(questions, "What concrete change should this project create for its users?")
	}
	if strings.TrimSpace(p.Audience) == "" {
		questions = append(questions, "Who is the primary user, and what do they already understand?")
	}
	if strings.TrimSpace(p.ProductType) == "" {
		questions = append(questions, "What are we making: an application, service, site, library, or experiment?")
	}
	if strings.TrimSpace(p.Success) == "" {
		questions = append(questions, "What observable outcome would make the first release successful?")
	}
	if len(p.InheritedTraits) == 0 {
		questions = append(questions, "Which parts of Bifunctional should this project unmistakably inherit?")
	}
	if len(p.DistinctTraits) == 0 {
		questions = append(questions, "Where should this project deliberately differ from the parent brand?")
	}
	if len(p.InitialWorkflows) == 0 {
		questions = append(questions, "What are the first three user journeys, screens, or service interactions?")
	}
	if len(p.Constraints) == 0 {
		questions = append(questions, "Which technical, legal, content, timeline, or asset constraints already exist?")
	}
	next := "Answer the questions, then call run_socratic_intake again."
	if len(questions) == 0 {
		next = "The intake is complete. Generate the project brief and design documents."
	}
	return Intake{Complete: len(questions) == 0, Questions: questions, NextStep: next}
}

func Validate(p Project) error {
	missing := make([]string, 0, 5)
	for name, value := range map[string]string{"name": p.Name, "purpose": p.Purpose, "audience": p.Audience, "productType": p.ProductType, "success": p.Success} {
		if strings.TrimSpace(value) == "" {
			missing = append(missing, name)
		}
	}
	if len(missing) > 0 {
		return fmt.Errorf("project is incomplete; missing %s", strings.Join(missing, ", "))
	}
	return nil
}

func ProjectBrief(p Project) (string, error) {
	if err := Validate(p); err != nil {
		return "", err
	}
	return frontmatter(p.Name+" project brief", "brief") + fmt.Sprintf(`# %s

## Purpose

%s

## Audience

%s

## Product shape

%s

## Definition of success

%s

## Initial workflows

%s

## Constraints

%s

## Non-goals

- Features that do not directly support the stated outcome.
- Brand novelty without a product reason.
- Operational complexity before demonstrated need.
`, p.Name, p.Purpose, p.Audience, p.ProductType, p.Success, bullets(p.InitialWorkflows), bullets(p.Constraints)), nil
}

func BrandDirection(p Project) (string, error) {
	if err := Validate(p); err != nil {
		return "", err
	}
	inherited := append([]string{"clarity over noise", "boutique over bloated", "structured but warm", "technical but human"}, p.InheritedTraits...)
	return frontmatter(p.Name+" brand direction", "design") + fmt.Sprintf(`# %s brand direction

## Narrative

%s exists to %s. It should feel credible to %s without borrowing the visual language of a generic software product.

## Inherited from Bifunctional

%s

## Project-specific character

%s

## Voice

%s

## Guardrails

- Family resemblance, not a clone of the parent website.
- Editorial clarity before visual decoration.
- Distinctive choices must support the product purpose.
- Public assets must contain no private client or internal operating knowledge.
`, p.Name, p.Name, strings.TrimSuffix(p.Purpose, "."), p.Audience, bullets(inherited), bullets(p.DistinctTraits), defaultValue(p.Tone, "Direct, observant, precise, and human. Avoid inflated agency language.")), nil
}

func DesignSystem(p Project) (string, error) {
	if err := Validate(p); err != nil {
		return "", err
	}
	return frontmatter(p.Name+" design system", "design") + fmt.Sprintf(`# %s design system

## Design principles

- Make hierarchy visible before adding interaction.
- Use typography as the primary structural material.
- Keep density proportional to the task.
- Reserve motion for transitions that explain state.
- Meet WCAG 2.2 AA contrast, focus, keyboard, and reduced-motion expectations.

## Foundations

### Typography

Choose one expressive display family and one highly legible text family. Define a fluid type scale, a readable measure, and explicit numeric styles before component work.

### Color

Start from a restrained Bifunctional family resemblance, then derive semantic roles for surface, text, border, accent, success, warning, and danger. Validate every foreground/background pair.

### Spacing and layout

Use a 4px base unit, named spacing tokens, a content grid, and deliberate breakpoints derived from content rather than devices.

### Motion

Define duration and easing tokens. Support prefers-reduced-motion and never make animation the only carrier of meaning.

## Initial component scope

%s

## Product adaptations

%s

## Definition of done

- Tokens are machine-readable and documented.
- Components include default, hover, focus, active, disabled, loading, empty, and error states where relevant.
- Desktop and mobile layouts are verified.
- Accessibility checks are automated where possible and manually sampled.
`, p.Name, bullets(p.InitialWorkflows), bullets(p.DistinctTraits)), nil
}

func frontmatter(title, kind string) string {
	return fmt.Sprintf("---\ntitle: %s\nowner: bifunctional\nstatus: draft\nupdated: %s\ntags: [%s, okf]\n---\n\n", title, time.Now().UTC().Format("2006-01-02"), kind)
}

func bullets(values []string) string {
	if len(values) == 0 {
		return "- To be decided."
	}
	var b strings.Builder
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			fmt.Fprintf(&b, "- %s\n", strings.TrimSpace(value))
		}
	}
	return strings.TrimSuffix(b.String(), "\n")
}

func defaultValue(value, fallback string) string {
	if strings.TrimSpace(value) == "" {
		return fallback
	}
	return value
}

var ErrExists = errors.New("destination already exists")

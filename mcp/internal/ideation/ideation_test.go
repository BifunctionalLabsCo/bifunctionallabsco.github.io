package ideation

import (
	"strings"
	"testing"
)

func TestSocraticIntakeAndGenerators(t *testing.T) {
	intake := SocraticIntake(Project{Name: "Only a name"})
	if intake.Complete || len(intake.Questions) != 8 {
		t.Fatalf("unexpected intake: %#v", intake)
	}
	p := completeProject()
	if intake := SocraticIntake(p); !intake.Complete || len(intake.Questions) != 0 {
		t.Fatalf("complete project was not complete: %#v", intake)
	}
	for name, generate := range map[string]func(Project) (string, error){"brief": ProjectBrief, "brand": BrandDirection, "design": DesignSystem} {
		doc, err := generate(p)
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		if !strings.HasPrefix(doc, "---\n") || !strings.Contains(doc, "updated:") || !strings.Contains(doc, "okf") {
			t.Errorf("%s did not produce OKF document", name)
		}
	}
}

func TestValidateRejectsIncompleteProject(t *testing.T) {
	if err := Validate(Project{Name: "Incomplete"}); err == nil {
		t.Fatal("expected incomplete project to fail")
	}
}

func completeProject() Project {
	return Project{Name: "Atlas", Purpose: "Make knowledge retrievable.", Audience: "Operators", ProductType: "Service", Success: "Answers include sources", InheritedTraits: []string{"clear"}, DistinctTraits: []string{"dense"}, InitialWorkflows: []string{"Search"}, Constraints: []string{"Private by default"}}
}

package okf

import "testing"

func TestParseAndValidate(t *testing.T) {
	input := "---\ntitle: Test\nowner: bifunctional\nstatus: active\nupdated: 2026-08-07\ntags: [docs, okf]\n---\n\n# Body\n"
	meta, body, err := Parse(input)
	if err != nil {
		t.Fatal(err)
	}
	if meta.Title != "Test" || len(meta.Tags) != 2 || body != "# Body\n" {
		t.Fatalf("unexpected parse: %#v %q", meta, body)
	}
	if warnings := Validate(meta); len(warnings) != 0 {
		t.Fatalf("unexpected warnings: %v", warnings)
	}
}

func TestParseRejectsInvalidDate(t *testing.T) {
	_, _, err := Parse("---\nupdated: yesterday\n---\n")
	if err == nil {
		t.Fatal("expected invalid date to fail")
	}
}

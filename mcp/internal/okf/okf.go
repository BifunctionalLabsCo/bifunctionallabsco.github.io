package okf

import (
	"bufio"
	"fmt"
	"strings"
	"time"
)

type Metadata struct {
	Title   string   `json:"title,omitempty"`
	Owner   string   `json:"owner,omitempty"`
	Status  string   `json:"status,omitempty"`
	Updated string   `json:"updated,omitempty"`
	Tags    []string `json:"tags,omitempty"`
}

type Document struct {
	Path       string   `json:"path"`
	Metadata   Metadata `json:"metadata"`
	Content    string   `json:"content"`
	Redactions int      `json:"redactions"`
}

func Parse(content string) (Metadata, string, error) {
	if !strings.HasPrefix(content, "---\n") && !strings.HasPrefix(content, "---\r\n") {
		return Metadata{}, content, nil
	}
	scanner := bufio.NewScanner(strings.NewReader(content))
	scanner.Scan()
	meta := Metadata{}
	closed := false
	lineCount := 1
	for scanner.Scan() {
		lineCount++
		line := strings.TrimSpace(scanner.Text())
		if line == "---" {
			closed = true
			break
		}
		key, value, ok := strings.Cut(line, ":")
		if !ok {
			return Metadata{}, "", fmt.Errorf("invalid frontmatter line %d", lineCount)
		}
		value = strings.Trim(strings.TrimSpace(value), `"'`)
		switch strings.TrimSpace(key) {
		case "title":
			meta.Title = value
		case "owner":
			meta.Owner = value
		case "status":
			meta.Status = value
		case "updated":
			if value != "" {
				if _, err := time.Parse("2006-01-02", value); err != nil {
					return Metadata{}, "", fmt.Errorf("updated must be YYYY-MM-DD: %w", err)
				}
			}
			meta.Updated = value
		case "tags":
			value = strings.Trim(value, "[]")
			for _, tag := range strings.Split(value, ",") {
				if tag = strings.TrimSpace(tag); tag != "" {
					meta.Tags = append(meta.Tags, tag)
				}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return Metadata{}, "", err
	}
	if !closed {
		return Metadata{}, "", fmt.Errorf("unclosed frontmatter")
	}
	lines := strings.Split(content, "\n")
	return meta, strings.TrimLeft(strings.Join(lines[lineCount:], "\n"), "\r\n"), nil
}

func Validate(meta Metadata) []string {
	var warnings []string
	if meta.Title == "" {
		warnings = append(warnings, "missing title")
	}
	if meta.Owner == "" {
		warnings = append(warnings, "missing owner")
	}
	if meta.Status == "" {
		warnings = append(warnings, "missing status")
	}
	if meta.Updated == "" {
		warnings = append(warnings, "missing updated")
	}
	return warnings
}

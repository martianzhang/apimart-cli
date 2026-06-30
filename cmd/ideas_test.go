package cmd

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"
)

// --- resolveIdeasKeywords ---

func TestResolveIdeasKeywords_args(t *testing.T) {
	got, err := resolveIdeasKeywords([]string{"cinematic", "portrait"})
	if err != nil {
		t.Fatalf("resolveIdeasKeywords() returned error: %v", err)
	}
	if got != "cinematic portrait" {
		t.Errorf("resolveIdeasKeywords() = %q, want %q", got, "cinematic portrait")
	}
}

func TestResolveIdeasKeywords_singleArg(t *testing.T) {
	got, err := resolveIdeasKeywords([]string{"portrait"})
	if err != nil {
		t.Fatalf("resolveIdeasKeywords() returned error: %v", err)
	}
	if got != "portrait" {
		t.Errorf("resolveIdeasKeywords() = %q, want %q", got, "portrait")
	}
}

func TestResolveIdeasKeywords_noArgs(t *testing.T) {
	got, err := resolveIdeasKeywords(nil)
	if err != nil {
		t.Fatalf("resolveIdeasKeywords() returned error: %v", err)
	}
	if got != "" {
		t.Errorf("resolveIdeasKeywords() = %q, want empty", got)
	}
}

// --- scoreEntry ---

func TestScoreEntry_matchTitle(t *testing.T) {
	e := IdeaEntry{Title: "Cinematic Portrait", Prompt: "A photo"}
	score := scoreEntry(e, []string{"portrait"})
	if score < 2 {
		t.Errorf("scoreEntry() = %d, want >= 2 (title match bonus)", score)
	}
}

func TestScoreEntry_matchPrompt(t *testing.T) {
	e := IdeaEntry{Title: "Something", Prompt: "cinematic lighting portrait"}
	score := scoreEntry(e, []string{"portrait"})
	if score < 1 {
		t.Errorf("scoreEntry() = %d, want >= 1", score)
	}
}

func TestScoreEntry_noMatch(t *testing.T) {
	e := IdeaEntry{Title: "Cat", Prompt: "A cute cat"}
	score := scoreEntry(e, []string{"portrait"})
	if score != 0 {
		t.Errorf("scoreEntry() = %d, want 0", score)
	}
}

func TestScoreEntry_multiKeyword(t *testing.T) {
	e := IdeaEntry{Title: "Cat Portrait", Prompt: "A cat portrait photo"}
	score := scoreEntry(e, []string{"cat", "portrait"})
	if score < 4 {
		t.Errorf("scoreEntry() = %d, want >= 4", score)
	}
}

func TestScoreEntry_caseInsensitive(t *testing.T) {
	e := IdeaEntry{Title: "CINEMATIC PORTRAIT", Prompt: "test"}
	score := scoreEntry(e, []string{"cinematic"})
	if score < 2 {
		t.Errorf("scoreEntry() should be case-insensitive, got %d", score)
	}
}

func TestScoreEntry_zhMatch(t *testing.T) {
	e := IdeaEntry{TitleZh: "电影感肖像", PromptZh: "一张电影感肖像照片"}
	score := scoreEntry(e, []string{"电影"})
	if score < 1 {
		t.Errorf("scoreEntry() should match zh text, got %d", score)
	}
}

// --- searchIdeas ---

func TestSearchIdeas_emptyQuery(t *testing.T) {
	entries := []IdeaEntry{{Title: "Test", Prompt: "test"}}
	results := searchIdeas(entries, "")
	if results != nil {
		t.Errorf("searchIdeas('') = %v, want nil", results)
	}
}

func TestSearchIdeas_sortedByScore(t *testing.T) {
	entries := []IdeaEntry{
		{Title: "Cat Portrait", Prompt: "a cat"},
		{Title: "Something Else", Prompt: "portrait photo"},
	}
	results := searchIdeas(entries, "portrait")
	if len(results) != 2 {
		t.Fatalf("searchIdeas() = %d results, want 2", len(results))
	}
	if results[0].score < results[1].score {
		t.Errorf("title match should rank higher, got %d < %d", results[0].score, results[1].score)
	}
}

func TestSearchIdeas_noMatches(t *testing.T) {
	entries := []IdeaEntry{{Title: "Cat", Prompt: "meow"}}
	results := searchIdeas(entries, "portrait")
	if len(results) != 0 {
		t.Errorf("searchIdeas() = %d results, want 0", len(results))
	}
}

// --- outputMarkdown ---

func captureStdout(fn func()) string {
	r, w, _ := os.Pipe()
	old := os.Stdout
	os.Stdout = w
	fn()
	w.Close()
	os.Stdout = old
	var buf bytes.Buffer
	buf.ReadFrom(r)
	return buf.String()
}

func TestOutputMarkdown_multipleResults(t *testing.T) {
	results := []searchResult{
		{entry: IdeaEntry{Title: "Test One", Prompt: "prompt one", Author: "Alice", License: "MIT"}, score: 3},
		{entry: IdeaEntry{Title: "Test Two", Prompt: "prompt two", Author: "Bob", SourceURL: "https://example.com"}, score: 1},
	}

	output := captureStdout(func() {
		if err := outputMarkdown(results, "test", 2, nil); err != nil {
			t.Errorf("outputMarkdown() returned error: %v", err)
		}
	})

	checks := []struct {
		name string
		want string
	}{
		{"title", "# Ideas: test"},
		{"count", "找到 2 个结果"},
		{"first heading", "## Test One"},
		{"second heading", "## Test Two"},
		{"first prompt", "```\nprompt one\n```"},
		{"second prompt", "```\nprompt two\n```"},
		{"author", "作者: Alice"},
		{"license", "MIT"},
		{"source link", "[来源]"},
		{"separator", "---"},
	}
	for _, c := range checks {
		t.Run(c.name, func(t *testing.T) {
			if !strings.Contains(output, c.want) {
				t.Errorf("output missing %q", c.want)
			}
		})
	}
}

func TestOutputMarkdown_singleResultNoSeparator(t *testing.T) {
	results := []searchResult{
		{entry: IdeaEntry{Title: "Only One", Prompt: "single"}, score: 1},
	}
	output := captureStdout(func() {
		if err := outputMarkdown(results, "test", 1, nil); err != nil {
			t.Errorf("outputMarkdown() returned error: %v", err)
		}
	})
	if !strings.Contains(output, "## Only One") {
		t.Errorf("output missing title")
	}
	if strings.Contains(output, "---") {
		t.Errorf("single result should not have separator")
	}
}

func TestOutputMarkdown_zhPrompt(t *testing.T) {
	results := []searchResult{
		{entry: IdeaEntry{Title: "ZH Test", Prompt: "english prompt", PromptZh: "中文提示词", Lang: "zh"}, score: 1},
	}
	output := captureStdout(func() {
		if err := outputMarkdown(results, "test", 1, nil); err != nil {
			t.Errorf("outputMarkdown() returned error: %v", err)
		}
	})
	if !strings.Contains(output, "中文提示词") {
		t.Errorf("zh entry should show zh prompt, got:\n%s", output)
	}
}

func TestOutputMarkdown_images(t *testing.T) {
	results := []searchResult{
		{entry: IdeaEntry{Title: "With Img", Prompt: "test", ImageURLs: []string{"https://example.com/img.jpg"}}, score: 1},
	}
	output := captureStdout(func() {
		if err := outputMarkdown(results, "test", 1, nil); err != nil {
			t.Errorf("outputMarkdown() returned error: %v", err)
		}
	})
	if !strings.Contains(output, "![参考图 1]") {
		t.Errorf("output missing image reference")
	}
}

func TestOutputMarkdown_emptyTitle(t *testing.T) {
	results := []searchResult{
		{entry: IdeaEntry{Prompt: "just a prompt"}, score: 1},
	}
	output := captureStdout(func() {
		if err := outputMarkdown(results, "test", 1, nil); err != nil {
			t.Errorf("outputMarkdown() returned error: %v", err)
		}
	})
	if !strings.Contains(output, "## Result 1") {
		t.Errorf("empty title should fallback to 'Result 1'")
	}
}

// --- outputJSON ---

func TestOutputJSON(t *testing.T) {
	results := []searchResult{
		{entry: IdeaEntry{Title: "JSON Test", Prompt: "test prompt"}, score: 1},
	}
	output := captureStdout(func() {
		if err := outputJSON(results, 1); err != nil {
			t.Errorf("outputJSON() returned error: %v", err)
		}
	})
	var parsed struct {
		Total   int         `json:"total"`
		Results []IdeaEntry `json:"results"`
	}
	if err := json.Unmarshal([]byte(output), &parsed); err != nil {
		t.Fatalf("invalid JSON: %v\n%s", err, output)
	}
	if parsed.Total != 1 {
		t.Errorf("total = %d, want 1", parsed.Total)
	}
	if len(parsed.Results) != 1 {
		t.Errorf("results = %d, want 1", len(parsed.Results))
	}
	if parsed.Results[0].Title != "JSON Test" {
		t.Errorf("title = %q", parsed.Results[0].Title)
	}
}

// --- localImagePath ---

func TestLocalImagePath_empty(t *testing.T) {
	if got := localImagePath(""); got != "" {
		t.Errorf("localImagePath('') = %q, want empty", got)
	}
}

func TestLocalImagePath_fullURL(t *testing.T) {
	got := localImagePath("https://example.com/path/to/img.jpg")
	if !strings.HasSuffix(got, "img.jpg") {
		t.Errorf("localImagePath() = %q, should end with img.jpg", got)
	}
}

// --- default constants ---

func TestDefaultConstants(t *testing.T) {
	if ideasDefaultLimit != 8 {
		t.Errorf("ideasDefaultLimit = %d, want 8", ideasDefaultLimit)
	}
	// ideasDefaultPageSize was removed; page-size functionality removed
}

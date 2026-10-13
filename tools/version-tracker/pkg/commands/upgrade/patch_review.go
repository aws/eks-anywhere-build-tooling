package upgrade

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

const maxPatchReviewLineRunes = 240

type patchReviewSummary struct {
	Route      string
	PatchCount int
	Summary    string
	Validation string
}

var (
	patchReviewMarkup     = strings.NewReplacer("*", "", "__", "", "~~", "", "`", "")
	patchReviewListMarker = regexp.MustCompile(`^([-*+]|[0-9]+[.)])\s+`)
	patchReviewNoise      = regexp.MustCompile(`(?i)(^(validation|tests?|checks?|build|logs?|output|files?|summary|human review)\s*:|\b(go test|go build|make test|git apply)\b|(^|[;.!?]\s+)((all|the)\s+)?(tests?|checks?|builds?|validation)(\s+(and|&)\s+(tests?|checks?|builds?|validation))*\s+((has|have|was|were|is|are)\s+)?(passed|passing|succeeded|successful|failed|skipped|not run)\b|(^|[;.!?]\s+)(validation|tests?|checks?|build|logs?|output|files?)\s*:)`)
	patchReviewDump       = regexp.MustCompile(`(?i)(^(diff --git|@@|traceback|error:|fatal:|\{|[^ ]+\.(go|mod|sum|patch|yaml|json|md)\s*:|[0-9]{4}-[0-9]{2}-[0-9]{2}[T ])|\{\s*["']|\b(aws_access_key_id|aws_secret_access_key|aws_session_token|accesskeyid|secretaccesskey|sessiontoken|github_token|password|secret|token)["']?\s*[:=]|(AKIA|ASIA)[A-Z0-9]{16}|-----BEGIN|level=(debug|info|warn|error))`)
	patchReviewURL        = regexp.MustCompile(`(?i)https?://\S+`)
)

// formatPatchReviewComment uses model text only for the change and review focus.
// Validation comes from the harness; every output line is bounded plain text.
func formatPatchReviewComment(review patchReviewSummary) string {
	route := "AI-assisted"
	deterministic := strings.EqualFold(strings.TrimSpace(review.Route), "deterministic")
	if deterministic {
		route = "deterministic"
	}
	change, focus := patchReviewSemanticLines(review.Summary)
	if change == "" {
		if deterministic {
			change = "Regenerated Autoscaler provider/dependency patches for target release."
		} else {
			change = "Repaired patch intent for target release."
		}
	}
	if focus == "" {
		switch {
		case patchReviewAlreadyUpstream(change):
			focus = "Confirm upstream behavior covers the removed patch's intent."
		case deterministic:
			focus = "CAPI-only provider selection and dependency updates."
		default:
			focus = "Confirm original intent is retained."
		}
	}
	validation := patchReviewPlainText(review.Validation)
	if validation == "" {
		validation = "Not reported by harness."
	}
	count := max(0, review.PatchCount)
	patches := "patches"
	if count == 1 {
		patches = "patch"
	}
	return strings.Join([]string{
		boundedPatchReviewLine(fmt.Sprintf("%s repair (%d %s): %s", route, count, patches, change)),
		boundedPatchReviewLine("Review: " + focus),
		boundedPatchReviewLine("Validation: " + validation),
		"Human review required.",
	}, "\n")
}

func patchReviewSemanticLines(summary string) (change, focus string) {
	var changes, focuses, fallbacks []string
	seen := make(map[string]bool)
	appendUnique := func(items *[]string, category, text string) {
		key := category + strings.ToLower(strings.TrimRight(text, "."))
		if !seen[key] {
			*items = append(*items, text)
			seen[key] = true
		}
	}
	inFence := false
	// Treat all common line separators consistently before choosing semantic lines.
	summary = strings.Map(func(r rune) rune {
		switch r {
		case '\r', '\v', '\f', '\u0085', '\u2028', '\u2029':
			return '\n'
		default:
			return r
		}
	}, summary)
	for _, raw := range strings.Split(summary, "\n") {
		line := strings.TrimLeft(strings.TrimSpace(raw), "> ")
		if strings.HasPrefix(line, "```") || strings.HasPrefix(line, "~~~") {
			inFence = !inFence
			continue
		}
		if strings.HasPrefix(line, "#") {
			continue
		}
		marker := patchReviewListMarker.FindStringIndex(line)
		bullet := marker != nil
		if bullet {
			line = strings.TrimSpace(line[marker[1]:])
		}
		// Normalize controls before detecting URLs so invisible characters
		// cannot hide a link that becomes active in the published plain text.
		line = patchReviewPlainText(line)
		line = patchReviewPlainText(patchReviewURL.ReplaceAllString(line, ""))
		label, text, labeled := strings.Cut(line, ":")
		if labeled && (strings.EqualFold(label, "Change") || strings.EqualFold(label, "Review")) {
			text = strings.TrimSpace(text)
			if !patchReviewUsefulText(text) {
				continue
			}
			if strings.EqualFold(label, "Change") {
				appendUnique(&changes, "change:", text)
			} else {
				appendUnique(&focuses, "review:", text)
			}
			continue
		}
		if inFence || bullet || !patchReviewUsefulText(line) {
			continue
		}
		// Older summaries sometimes place the review focus after a semicolon.
		if before, after, ok := strings.Cut(line, "; review "); ok {
			line = before
			if patchReviewUsefulText(after) {
				appendUnique(&focuses, "review:", after)
			}
		}
		appendUnique(&fallbacks, "fallback:", line)
	}
	haveLabeledChanges := len(changes) > 0
	for _, fallback := range fallbacks {
		if !haveLabeledChanges || patchReviewAlreadyUpstream(fallback) {
			appendUnique(&changes, "change:", fallback)
		}
	}
	for i, text := range changes {
		if strings.EqualFold(text, "already upstream") {
			changes[i] = "Removed obsolete patch; its intent is already upstream."
		}
	}
	var specificFocuses []string
	for _, text := range focuses {
		if !patchReviewGenericFocus(text) {
			specificFocuses = append(specificFocuses, text)
		}
	}
	if len(specificFocuses) > 0 {
		focuses = specificFocuses
	}
	return joinPatchReviewTopics(changes), joinPatchReviewTopics(focuses)
}

func joinPatchReviewTopics(topics []string) string {
	if len(topics) == 0 {
		return ""
	}
	// Reserve room for the route/count prefix and show several substantive
	// changes, rather than letting one long model summary hide the others.
	const budget = 180
	suffix := ""
	if len(topics) > 3 {
		suffix = fmt.Sprintf(" (+%d more)", len(topics)-3)
		topics = topics[:3]
	}
	text := strings.Join(topics, " ") + suffix
	if len([]rune(text)) <= budget {
		return text
	}
	perTopic := (budget - len([]rune(suffix)) - len(topics) + 1) / len(topics)
	short := make([]string, len(topics))
	for i, topic := range topics {
		short[i] = boundedPatchReviewText(topic, perTopic)
	}
	return strings.Join(short, " ") + suffix
}

func patchReviewUsefulText(text string) bool {
	if text == "" || strings.HasSuffix(text, ":") || patchReviewNoise.MatchString(text) || patchReviewDump.MatchString(text) {
		return false
	}
	// Bare paths, numbered inventories, and shell/log records are not summaries.
	first := strings.Fields(text)[0]
	if strings.Contains(first, "/") || strings.HasPrefix(first, "[") ||
		strings.HasPrefix(first, "$") || strings.HasSuffix(first, ".") && unicode.IsDigit([]rune(first)[0]) {
		return false
	}
	switch strings.TrimRight(strings.ToLower(text), ".") {
	case "done", "summary", "changes", "repair summary", "all checks passed",
		"updated imports", "adjusted whitespace", "formatting only":
		return false
	}
	return true
}

func patchReviewGenericFocus(text string) bool {
	switch strings.TrimRight(strings.ToLower(text), ".") {
	case "original intent retained", "original intent preserved", "confirm original intent is retained",
		"none", "no risks", "no non-obvious risks":
		return true
	}
	return false
}

func patchReviewAlreadyUpstream(text string) bool {
	return strings.Contains(strings.ToLower(text), "already upstream")
}

func patchReviewPlainText(text string) string {
	text = patchReviewMarkup.Replace(text)
	text = strings.ReplaceAll(text, "alreadyupstream", "already upstream")
	text = strings.Map(func(r rune) rune {
		if unicode.IsSpace(r) {
			return ' '
		}
		if unicode.IsControl(r) || unicode.Is(unicode.Cf, r) {
			return -1
		}
		return r
	}, text)
	return strings.Join(strings.Fields(strings.Trim(text, "*_ ")), " ")
}

func boundedPatchReviewLine(line string) string {
	return boundedPatchReviewText(line, maxPatchReviewLineRunes)
}

func boundedPatchReviewText(line string, limit int) string {
	runes := []rune(line)
	if len(runes) <= limit {
		return line
	}
	head := string(runes[:limit-1])
	if space := strings.LastIndexByte(head, ' '); space > len(head)/2 {
		head = head[:space]
	}
	return strings.TrimSpace(head) + "…"
}

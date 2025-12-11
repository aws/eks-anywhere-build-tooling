package types

// FixPatchesOptions represents the options that can be passed to the `fix-patches` command.
type FixPatchesOptions struct {
	ProjectName         string
	PRNumber            int
	Auto                bool
	MaxAttempts         int
	Model               string
	Region              string
	ComplexityThreshold int
	CommentOnPR         bool
	Push                bool
	Branch              string
}

// PatchContext contains all context needed for LLM to fix a patch.
type PatchContext struct {
	ProjectName       string
	OriginalPatch     string
	PatchAuthor       string // From: name <email>
	PatchDate         string // Date: timestamp
	PatchSubject      string // Subject: [PATCH] description
	FailedHunks       []FailedHunk
	CurrentFileState  map[string]string
	PatchIntent       string
	BuildError        string
	PreviousAttempts  []string
	TokenCount        int
	ApplicationResult *PatchApplicationResult // Information about how patch was applied
	AllFileContexts   map[string]string       // filename -> current content around changed lines (for ALL files in patch)
}

// FailedHunk represents a single failed hunk from a .rej file.
type FailedHunk struct {
	FilePath        string
	HunkIndex       int
	OriginalLines   []string
	Context         string
	LineNumber      int      // Approximate line number where conflict occurred
	ExpectedContext []string // What the patch expects to find (context lines from .rej)
	ActualContext   []string // What's actually in the current file
	Differences     []string // Specific differences between expected and actual
}

// PatchApplicationResult contains information about how a patch was applied.
type PatchApplicationResult struct {
	OffsetFiles     map[string]int    // filename -> line offset (e.g., "go.sum" -> 2 means applied at line+2)
	GitOutput       string            // Full git apply output for reference
	PristineContent map[string]string // filename -> content BEFORE git apply (pristine state)
}

// PatchFix represents the LLM-generated patch fix.
type PatchFix struct {
	Patch      string
	TokensUsed int
	Cost       float64
}

// PatchFailureEvent represents the EventBridge event published when patches fail.
type PatchFailureEvent struct {
	ProjectName    string   `json:"project"`
	PRNumber       int      `json:"pr_number"`
	PRBranch       string   `json:"branch"`
	PRURL          string   `json:"pr_url"`
	FailedPatches  []string `json:"failed_patches"`
	FailedFiles    []string `json:"failed_files"`
	PatchesApplied int      `json:"patches_applied"`
	TotalPatches   int      `json:"total_patches"`
	Timestamp      string   `json:"timestamp"`
}

// PatchFixError represents a categorized error from patch fixing.
type PatchFixError struct {
	Code    string
	Message string
	Details map[string]interface{}
}

// Error codes for patch fixing failures.
const (
	ErrorComplexityTooHigh   = "COMPLEXITY_TOO_HIGH"
	ErrorBuildFailed         = "BUILD_FAILED"
	ErrorSemanticDrift       = "SEMANTIC_DRIFT"
	ErrorMaxAttemptsExceeded = "MAX_ATTEMPTS_EXCEEDED"
	ErrorContextTooLarge     = "CONTEXT_TOO_LARGE"
	ErrorBedrockAPI          = "BEDROCK_API_ERROR"
	ErrorGitOperation        = "GIT_OPERATION_ERROR"
)

func (e *PatchFixError) Error() string {
	return e.Message
}

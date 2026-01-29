package cmd

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/commands/fixpatches"
	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/types"
)

var fixPatchesOptions = &types.FixPatchesOptions{}

// fixPatchesCmd is the command used to fix failed patches using LLM assistance.
var fixPatchesCmd = &cobra.Command{
	Use:   "fix-patches --project <project name> --pr <PR number>",
	Short: "Fix failed patches using LLM assistance",
	Long:  "Use this command to automatically fix failed patch applications using AWS Bedrock LLM. The command analyzes .rej files, generates fixes, and validates them through compilation.",
	Run: func(cmd *cobra.Command, args []string) {
		err := fixpatches.Run(fixPatchesOptions)
		if err != nil {
			log.Fatalf("Error fixing patches: %v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(fixPatchesCmd)

	// Required flags
	fixPatchesCmd.Flags().StringVar(&fixPatchesOptions.ProjectName, "project", "", "Specify the project name (e.g., aquasecurity/trivy)")
	fixPatchesCmd.Flags().IntVar(&fixPatchesOptions.PRNumber, "pr", 0, "Specify the PR number")

	// Optional flags
	fixPatchesCmd.Flags().BoolVar(&fixPatchesOptions.Auto, "auto", false, "Auto-detect project and PR from current git branch")
	fixPatchesCmd.Flags().IntVar(&fixPatchesOptions.MaxAttempts, "max-attempts", 5, "Maximum number of fix attempts")
	fixPatchesCmd.Flags().StringVar(&fixPatchesOptions.Model, "model", "anthropic.claude-sonnet-4-5-20250929-v1:0", "Bedrock model ID to use (Claude Sonnet 4.5 - 200K context, 200K tokens/min)")
	fixPatchesCmd.Flags().StringVar(&fixPatchesOptions.Region, "region", "us-west-2", "AWS region for Bedrock API")
	fixPatchesCmd.Flags().IntVar(&fixPatchesOptions.ComplexityThreshold, "complexity-threshold", 10, "Skip patches exceeding this complexity score")
	fixPatchesCmd.Flags().BoolVar(&fixPatchesOptions.CommentOnPR, "comment-on-pr", false, "Post success/failure comment on the PR")
	fixPatchesCmd.Flags().BoolVar(&fixPatchesOptions.Push, "push", false, "Commit and push fixed patches to the PR branch")
	fixPatchesCmd.Flags().StringVar(&fixPatchesOptions.Branch, "branch", "", "Branch name to push to (required if --push is used)")

	// Mark required flags (unless --auto is used)
	fixPatchesCmd.MarkFlagsRequiredTogether("project", "pr")
}

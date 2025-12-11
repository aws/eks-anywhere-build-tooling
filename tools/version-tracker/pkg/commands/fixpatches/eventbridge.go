package fixpatches

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge"
	"github.com/aws/aws-sdk-go-v2/service/eventbridge/types"

	"github.com/aws/eks-anywhere-build-tooling/tools/version-tracker/pkg/util/logger"
)

// PatchFailureEvent represents a patch failure event to be sent to EventBridge.
type PatchFailureEvent struct {
	Project       string   `json:"project"`
	PRNumber      int      `json:"pr_number"`
	Branch        string   `json:"branch"`
	FailedPatches []string `json:"failed_patches"`
	Reason        string   `json:"reason"`
	RepoOwner     string   `json:"repo_owner"`
	RepoName      string   `json:"repo_name"`
}

// PublishPatchFailureEvent publishes a patch failure event to EventBridge.
// This triggers the Lambda function to automatically fix patches.
func PublishPatchFailureEvent(event PatchFailureEvent) error {
	// Check if auto-patch-fix is enabled
	enableAutoPatchFix := os.Getenv("ENABLE_AUTO_PATCH_FIX")
	if enableAutoPatchFix != "true" {
		logger.Info("Auto patch fix is disabled, skipping EventBridge event", "ENABLE_AUTO_PATCH_FIX", enableAutoPatchFix)
		return nil
	}

	eventBusName := os.Getenv("PATCH_FIXER_EVENT_BUS")
	if eventBusName == "" {
		eventBusName = "default" // Use default event bus if not specified
	}

	logger.Info("Publishing patch failure event to EventBridge",
		"project", event.Project,
		"pr", event.PRNumber,
		"event_bus", eventBusName)

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		return fmt.Errorf("loading AWS config: %v", err)
	}

	client := eventbridge.NewFromConfig(cfg)

	// Marshal event to JSON
	eventDetail, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshaling event: %v", err)
	}

	// Put event to EventBridge
	input := &eventbridge.PutEventsInput{
		Entries: []types.PutEventsRequestEntry{
			{
				Source:       aws.String("eks-anywhere.version-tracker"),
				DetailType:   aws.String("PatchFailure"),
				Detail:       aws.String(string(eventDetail)),
				EventBusName: aws.String(eventBusName),
			},
		},
	}

	output, err := client.PutEvents(ctx, input)
	if err != nil {
		return fmt.Errorf("putting event to EventBridge: %v", err)
	}

	// Check for failed entries
	if output.FailedEntryCount > 0 {
		return fmt.Errorf("failed to publish %d events", output.FailedEntryCount)
	}

	logger.Info("Successfully published patch failure event to EventBridge",
		"project", event.Project,
		"pr", event.PRNumber)

	return nil
}

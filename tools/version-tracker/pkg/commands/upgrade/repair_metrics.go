package upgrade

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	patchRepairMetricsEnv = "PATCH_FIXER_METRICS_ENABLED"
	patchRepairNamespace  = "EKSAnywhere/VersionTracker"

	metricPatchConflictDetected = "PatchConflictDetected"
	metricPatchRepairAttempted  = "PatchRepairAttempted"
	metricPatchRepairSucceeded  = "PatchRepairSucceeded"
	metricPatchRepairFailed     = "PatchRepairFailed"
	metricPatchRepairSkipped    = "PatchRepairSkipped"
	metricPatchRepairDuration   = "PatchRepairDuration"
	metricPatchRepairPRCreated  = "PatchRepairPullRequestCreated"
)

var patchRepairMetricsWriter io.Writer = os.Stdout

func emitPatchRepairMetric(
	name string,
	project string,
	branch string,
	route string,
	result string,
	reason string,
	duration time.Duration,
) {
	enabled, err := strconv.ParseBool(strings.TrimSpace(os.Getenv(patchRepairMetricsEnv)))
	if err != nil || !enabled {
		return
	}
	if branch == "" {
		branch = "unknown"
	}
	if route == "" {
		route = "unknown"
	}

	metrics := []map[string]string{{"Name": name, "Unit": "Count"}}
	if duration > 0 {
		metrics = append(metrics, map[string]string{"Name": metricPatchRepairDuration, "Unit": "Milliseconds"})
	}
	payload := map[string]any{
		"_aws": map[string]any{
			"Timestamp": time.Now().UnixMilli(),
			"CloudWatchMetrics": []map[string]any{{
				"Namespace":  patchRepairNamespace,
				"Dimensions": [][]string{{"Branch", "Route"}},
				"Metrics":    metrics,
			}},
		},
		"Branch":  branch,
		"Route":   route,
		"Project": project,
		"BuildID": os.Getenv("CODEBUILD_BUILD_ID"),
		"Result":  result,
		"Reason":  boundedMetricReason(reason),
		name:      1,
	}
	if duration > 0 {
		payload[metricPatchRepairDuration] = duration.Milliseconds()
	}

	encoded, err := json.Marshal(payload)
	if err != nil {
		return
	}
	_, _ = fmt.Fprintln(patchRepairMetricsWriter, string(encoded))
}

func boundedMetricReason(reason string) string {
	const maxBytes = 1024
	if len(reason) <= maxBytes {
		return reason
	}
	return reason[len(reason)-maxBytes:]
}

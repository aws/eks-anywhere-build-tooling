package upgrade

import (
	"bytes"
	"encoding/json"
	"testing"
	"time"
)

func TestEmitPatchRepairMetric(t *testing.T) {
	t.Setenv(patchRepairMetricsEnv, "true")
	t.Setenv("CODEBUILD_BUILD_ID", "build-id")
	var output bytes.Buffer
	originalWriter := patchRepairMetricsWriter
	patchRepairMetricsWriter = &output
	t.Cleanup(func() {
		patchRepairMetricsWriter = originalWriter
	})

	emitPatchRepairMetric(
		metricPatchRepairSucceeded,
		"kubernetes/autoscaler",
		"main",
		"deterministic",
		"validated",
		"",
		1500*time.Millisecond,
	)

	var payload map[string]any
	if err := json.Unmarshal(output.Bytes(), &payload); err != nil {
		t.Fatal(err)
	}
	if payload["Branch"] != "main" || payload["Route"] != "deterministic" {
		t.Fatalf("dimensions = %#v", payload)
	}
	if payload["Project"] != "kubernetes/autoscaler" || payload["BuildID"] != "build-id" {
		t.Fatalf("properties = %#v", payload)
	}
	if payload[metricPatchRepairSucceeded] != float64(1) {
		t.Fatalf("success metric = %#v", payload[metricPatchRepairSucceeded])
	}
	if payload[metricPatchRepairDuration] != float64(1500) {
		t.Fatalf("duration metric = %#v", payload[metricPatchRepairDuration])
	}
	awsMetadata, ok := payload["_aws"].(map[string]any)
	if !ok {
		t.Fatalf("_aws = %#v", payload["_aws"])
	}
	cloudWatchMetrics := awsMetadata["CloudWatchMetrics"].([]any)
	metricDirective := cloudWatchMetrics[0].(map[string]any)
	if metricDirective["Namespace"] != patchRepairNamespace {
		t.Fatalf("namespace = %#v", metricDirective["Namespace"])
	}
	metricDefinitions := metricDirective["Metrics"].([]any)
	if len(metricDefinitions) != 2 {
		t.Fatalf("metric definitions = %#v, want count and duration", metricDefinitions)
	}
}

func TestEmitPatchRepairMetricDisabled(t *testing.T) {
	t.Setenv(patchRepairMetricsEnv, "false")
	var output bytes.Buffer
	originalWriter := patchRepairMetricsWriter
	patchRepairMetricsWriter = &output
	t.Cleanup(func() {
		patchRepairMetricsWriter = originalWriter
	})

	emitPatchRepairMetric(metricPatchConflictDetected, "example/project", "main", "pending", "detected", "", 0)

	if output.Len() != 0 {
		t.Fatalf("metric output = %q, want empty", output.String())
	}
}

func TestEmitPatchRepairPullRequestCreatedMetric(t *testing.T) {
	t.Setenv(patchRepairMetricsEnv, "true")
	var output bytes.Buffer
	originalWriter := patchRepairMetricsWriter
	patchRepairMetricsWriter = &output
	t.Cleanup(func() {
		patchRepairMetricsWriter = originalWriter
	})

	emitPatchRepairMetric(
		metricPatchRepairPRCreated,
		"example/project",
		"main",
		"strands",
		"created",
		"",
		0,
	)

	var payload map[string]any
	if err := json.Unmarshal(output.Bytes(), &payload); err != nil {
		t.Fatal(err)
	}
	if payload[metricPatchRepairPRCreated] != float64(1) {
		t.Fatalf("pull request metric = %#v", payload[metricPatchRepairPRCreated])
	}
	if payload["Route"] != "strands" || payload["Result"] != "created" {
		t.Fatalf("properties = %#v", payload)
	}
}

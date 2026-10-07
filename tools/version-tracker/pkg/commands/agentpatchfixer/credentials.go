package agentpatchfixer

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type agentCredentials struct {
	AccessKeyID     string    `json:"AccessKeyId"`
	SecretAccessKey string    `json:"SecretAccessKey"`
	SessionToken    string    `json:"SessionToken"`
	Expiration      time.Time `json:"Expiration"`
}

// assumeAgentRole keeps the build role's credentials in the Go harness. AWS CLI
// is already required by the upgrade build; its response stays in memory.
func assumeAgentRole(ctx context.Context) (agentCredentials, error) {
	roleARN := strings.TrimSpace(os.Getenv(AgentRoleARNEnv))
	if roleARN == "" {
		return agentCredentials{}, fmt.Errorf("%s is required for agent patch repair", AgentRoleARNEnv)
	}
	duration := defaultTimeout
	if deadline, ok := ctx.Deadline(); ok {
		duration = time.Until(deadline)
	}
	if err := ctx.Err(); err != nil {
		return agentCredentials{}, err
	}
	// STS sessions must last at least 15 minutes. Leave a minute for startup and
	// reject attempts that exceed the one-hour role-chaining limit.
	sessionSeconds := max(900, int(duration.Seconds())+61)
	if sessionSeconds > 3600 {
		return agentCredentials{}, fmt.Errorf("agent timeout exceeds the one-hour role session limit")
	}
	assumeCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(assumeCtx, "aws", "sts", "assume-role",
		"--role-arn", roleARN,
		"--role-session-name", "patch-fixer-agent",
		"--duration-seconds", strconv.Itoa(sessionSeconds),
		"--output", "json",
		"--no-cli-pager",
		"--no-cli-auto-prompt",
	)
	output, err := cmd.Output()
	if assumeCtx.Err() != nil {
		return agentCredentials{}, assumeCtx.Err()
	}
	if err != nil {
		// Never include captured output: a successful STS response contains secrets.
		return agentCredentials{}, fmt.Errorf("assuming patch fixer agent role: %w", err)
	}
	var response struct {
		Credentials agentCredentials `json:"Credentials"`
	}
	if err := json.Unmarshal(output, &response); err != nil {
		return agentCredentials{}, fmt.Errorf("parsing patch fixer agent credentials: %w", err)
	}
	credentials := response.Credentials
	if credentials.AccessKeyID == "" || credentials.SecretAccessKey == "" || credentials.SessionToken == "" {
		return agentCredentials{}, fmt.Errorf("STS returned incomplete patch fixer agent credentials")
	}
	if time.Until(credentials.Expiration) < duration {
		return agentCredentials{}, fmt.Errorf("patch fixer agent credentials expire before the repair deadline")
	}
	return credentials, nil
}

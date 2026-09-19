package agentpatchfixer

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/document"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
)

type fakeBedrockClient struct {
	input  *bedrockruntime.ConverseInput
	output *bedrockruntime.ConverseOutput
}

func (f *fakeBedrockClient) Converse(
	_ context.Context,
	input *bedrockruntime.ConverseInput,
	_ ...func(*bedrockruntime.Options),
) (*bedrockruntime.ConverseOutput, error) {
	f.input = input
	return f.output, nil
}

func TestBedrockModelConvertsToolConversation(t *testing.T) {
	fake := &fakeBedrockClient{output: &bedrockruntime.ConverseOutput{
		Output: &types.ConverseOutputMemberMessage{Value: types.Message{
			Role: types.ConversationRoleAssistant,
			Content: []types.ContentBlock{
				&types.ContentBlockMemberToolUse{Value: types.ToolUseBlock{
					ToolUseId: aws.String("tool-1"),
					Name:      aws.String("read_file"),
					Input: document.NewLazyDocument(map[string]any{
						"path":       "main.go",
						"start_line": 1,
					}),
				}},
			},
		}},
		StopReason: types.StopReasonToolUse,
		Usage: &types.TokenUsage{
			InputTokens:  aws.Int32(12),
			OutputTokens: aws.Int32(4),
			TotalTokens:  aws.Int32(16),
		},
	}}
	model := &bedrockModel{client: fake, modelID: "test-model"}
	response, err := model.Converse(context.Background(), modelRequest{
		System: "system",
		Messages: []modelMessage{{
			Role:    "user",
			Content: []contentBlock{{Text: "repair"}},
		}},
		Tools:     patchTools(),
		MaxTokens: 16384,
	})
	if err != nil {
		t.Fatal(err)
	}
	if aws.ToString(fake.input.ModelId) != "test-model" {
		t.Fatalf("model ID = %q", aws.ToString(fake.input.ModelId))
	}
	if len(fake.input.ToolConfig.Tools) != len(patchTools()) {
		t.Fatalf("tools = %d", len(fake.input.ToolConfig.Tools))
	}
	if response.StopReason != "tool_use" || response.InputTokens != 12 || response.OutputTokens != 4 {
		t.Fatalf("response = %#v", response)
	}
	if len(response.Message.Content) != 1 || response.Message.Content[0].ToolUse == nil {
		t.Fatalf("content = %#v", response.Message.Content)
	}
	use := response.Message.Content[0].ToolUse
	if use.Name != "read_file" || use.Input["path"] != "main.go" {
		t.Fatalf("tool use = %#v", use)
	}

	replayed, err := toBedrockMessage(response.Message)
	if err != nil {
		t.Fatal(err)
	}
	if replayed.Content[0] != fake.output.Output.(*types.ConverseOutputMemberMessage).Value.Content[0] {
		t.Fatal("assistant content was not preserved for the next Converse turn")
	}
}

func TestToBedrockMessageMarksToolErrors(t *testing.T) {
	message, err := toBedrockMessage(modelMessage{
		Role: "user",
		Content: []contentBlock{{ToolResult: &toolResult{
			ToolUseID: "tool-1",
			Content:   "blocked",
			IsError:   true,
		}}},
	})
	if err != nil {
		t.Fatal(err)
	}
	result, ok := message.Content[0].(*types.ContentBlockMemberToolResult)
	if !ok {
		t.Fatalf("content type = %T", message.Content[0])
	}
	if result.Value.Status != types.ToolResultStatusError {
		t.Fatalf("status = %q", result.Value.Status)
	}
}

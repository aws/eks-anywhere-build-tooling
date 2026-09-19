package agentpatchfixer

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/document"
	"github.com/aws/aws-sdk-go-v2/service/bedrockruntime/types"
)

type bedrockConverseAPI interface {
	Converse(context.Context, *bedrockruntime.ConverseInput, ...func(*bedrockruntime.Options)) (*bedrockruntime.ConverseOutput, error)
}

type bedrockModel struct {
	client  bedrockConverseAPI
	modelID string
}

func newBedrockModel(ctx context.Context, modelID, region string) (model, error) {
	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(region),
		config.WithRetryMaxAttempts(3),
		config.WithRetryMode(aws.RetryModeAdaptive),
	)
	if err != nil {
		return nil, fmt.Errorf("loading AWS configuration: %v", err)
	}
	return &bedrockModel{
		client:  bedrockruntime.NewFromConfig(cfg),
		modelID: modelID,
	}, nil
}

func (m *bedrockModel) Converse(ctx context.Context, request modelRequest) (modelResponse, error) {
	messages := make([]types.Message, 0, len(request.Messages))
	for _, message := range request.Messages {
		converted, err := toBedrockMessage(message)
		if err != nil {
			return modelResponse{}, err
		}
		messages = append(messages, converted)
	}
	tools := make([]types.Tool, 0, len(request.Tools))
	for _, tool := range request.Tools {
		tools = append(tools, &types.ToolMemberToolSpec{Value: types.ToolSpecification{
			Name:        aws.String(tool.Name),
			Description: aws.String(tool.Description),
			InputSchema: &types.ToolInputSchemaMemberJson{
				Value: document.NewLazyDocument(tool.InputSchema),
			},
		}})
	}

	output, err := m.client.Converse(ctx, &bedrockruntime.ConverseInput{
		ModelId:  aws.String(m.modelID),
		Messages: messages,
		System: []types.SystemContentBlock{
			&types.SystemContentBlockMemberText{Value: request.System},
		},
		InferenceConfig: &types.InferenceConfiguration{
			MaxTokens: aws.Int32(request.MaxTokens),
		},
		ToolConfig: &types.ToolConfiguration{Tools: tools},
	})
	if err != nil {
		return modelResponse{}, err
	}
	messageOutput, ok := output.Output.(*types.ConverseOutputMemberMessage)
	if !ok {
		return modelResponse{}, fmt.Errorf("Bedrock returned an unsupported Converse output %T", output.Output)
	}
	message, err := fromBedrockMessage(messageOutput.Value)
	if err != nil {
		return modelResponse{}, err
	}

	response := modelResponse{
		Message:    message,
		StopReason: string(output.StopReason),
	}
	if output.Usage != nil {
		response.InputTokens = int(aws.ToInt32(output.Usage.InputTokens))
		response.OutputTokens = int(aws.ToInt32(output.Usage.OutputTokens))
	}
	return response, nil
}

func toBedrockMessage(message modelMessage) (types.Message, error) {
	content := make([]types.ContentBlock, 0, len(message.Content))
	for _, block := range message.Content {
		if raw, ok := block.raw.(types.ContentBlock); ok {
			content = append(content, raw)
			continue
		}
		switch {
		case block.Text != "":
			content = append(content, &types.ContentBlockMemberText{Value: block.Text})
		case block.ToolUse != nil:
			content = append(content, &types.ContentBlockMemberToolUse{Value: types.ToolUseBlock{
				ToolUseId: aws.String(block.ToolUse.ID),
				Name:      aws.String(block.ToolUse.Name),
				Input:     document.NewLazyDocument(block.ToolUse.Input),
			}})
		case block.ToolResult != nil:
			status := types.ToolResultStatusSuccess
			if block.ToolResult.IsError {
				status = types.ToolResultStatusError
			}
			content = append(content, &types.ContentBlockMemberToolResult{Value: types.ToolResultBlock{
				ToolUseId: aws.String(block.ToolResult.ToolUseID),
				Status:    status,
				Content: []types.ToolResultContentBlock{
					&types.ToolResultContentBlockMemberText{Value: block.ToolResult.Content},
				},
			}})
		default:
			return types.Message{}, fmt.Errorf("message contains an empty content block")
		}
	}
	role := types.ConversationRoleUser
	if message.Role == "assistant" {
		role = types.ConversationRoleAssistant
	}
	return types.Message{Role: role, Content: content}, nil
}

func fromBedrockMessage(message types.Message) (modelMessage, error) {
	result := modelMessage{Role: string(message.Role)}
	for _, raw := range message.Content {
		block := contentBlock{raw: raw}
		switch value := raw.(type) {
		case *types.ContentBlockMemberText:
			block.Text = value.Value
		case *types.ContentBlockMemberToolUse:
			input := make(map[string]any)
			if err := value.Value.Input.UnmarshalSmithyDocument(&input); err != nil {
				content, marshalErr := value.Value.Input.MarshalSmithyDocument()
				if marshalErr != nil {
					return modelMessage{}, fmt.Errorf("decoding tool input for %s: %v", aws.ToString(value.Value.Name), err)
				}
				if jsonErr := json.Unmarshal(content, &input); jsonErr != nil {
					return modelMessage{}, fmt.Errorf("decoding tool input for %s: %v", aws.ToString(value.Value.Name), jsonErr)
				}
			}
			block.ToolUse = &toolUse{
				ID:    aws.ToString(value.Value.ToolUseId),
				Name:  aws.ToString(value.Value.Name),
				Input: input,
			}
		}
		result.Content = append(result.Content, block)
	}
	return result, nil
}

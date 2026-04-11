package openrouter

import (
	"context"
	"fmt"

	_ "embed"

	"github.com/revrost/go-openrouter"
)

const systemPrompt = `
You are an advanced quantitative crypto trader.

Your task is to analyze structured market data
and output trading intelligence.

Rules:

- Do NOT hallucinate levels
- Use only provided data
- Be deterministic
- Use numbers
- Always produce structured output
- Always estimate probabilities

Focus on:

- Trend structure
- Liquidity zones
- Momentum
- Volatility
- Risk management

Return output in telegram markdown v2 format which could be:
`

//go:embed format_options.txt
var formatOptions string

func (s *Service) Analyze(ctx context.Context, openrouterModel, req string) (string, error) {
	resp, err := s.client.CreateChatCompletion(ctx, openrouter.ChatCompletionRequest{
		Model: openrouterModel,
		Messages: []openrouter.ChatCompletionMessage{
			openrouter.SystemMessage(
				fmt.Sprintf("%s\n\n%s", systemPrompt, formatOptions)),
			openrouter.UserMessage(req),
		},
	})
	if err != nil {
		return "", err
	}

	return resp.Choices[0].Message.Content.Text, nil
}

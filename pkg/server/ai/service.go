package ai

import "github.com/ca-risken/common/pkg/logging"

type AIService struct {
	logger logging.Logger
}

func NewAIService(l logging.Logger) *AIService {
	return &AIService{logger: l}
}

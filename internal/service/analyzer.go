package service

import (
	"fmt"

	"auditor.z9fr.xyz/server/internal/lib"
)

type AnalyzerServiceImpl struct {
	log lib.Logger
	env *lib.Env
}

func NewAnalyzerServiceImpl(
	log lib.Logger,
	env *lib.Env,
) *AnalyzerServiceImpl {
	return &AnalyzerServiceImpl{
		log: log,
		env: env,
	}
}

func (s *AnalyzerServiceImpl) InitiateAnalyzer() (string, error) {
	fmt.Print("Analyze start ")
	return "ok", nil
}

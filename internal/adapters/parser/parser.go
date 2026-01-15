package parser

import (
	"context"
)

type Parser interface {
	Parse(ctx context.Context, filename string, content []byte) ([]Symbol, error)
}

type Symbol struct {
	Name string
	Type string
	Code string
}

type GenericParser struct {
}

func NewGenericParser() *GenericParser {
	return &GenericParser{}
}

func (p *GenericParser) Parse(ctx context.Context, filename string, content []byte) ([]Symbol, error) {
	lang := DetectLanguage(filename)
	return ExtractSymbols(lang, content)
}

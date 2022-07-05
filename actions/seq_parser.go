package actions

import (
	"fmt"
	"os"
	"regexp"
)

var envNamePattern = `^[a-zA-Z_][a-zA-Z0-9_]*$`

type Pattern struct {
	start rune
	end   rune
}

type Parser interface {
	Transform(input byte) ([]byte, error)
	Flush() []byte
}

type sequenceParser struct {
	begin      string
	beginIndex int
	end        string
	endIndex   int
	envName    string
	regexEnv   *regexp.Regexp
}

func NewSeqParser(pattern []string) Parser {
	r, _ := regexp.Compile(envNamePattern)

	begin, end := extractPattern(pattern)

	return &sequenceParser{
		begin:      begin,
		beginIndex: 0,
		end:        end,
		endIndex:   0,
		regexEnv:   r,
	}
}

func (p *sequenceParser) Transform(input byte) ([]byte, error) {
	switch true {
	case p.isMatchedBegin(input):
		p.beginIndex += 1
		return nil, nil
	case p.isMatchedFileName(input):
		p.envName += string(input)
		return nil, nil
	case p.isMatchedEnd(input):
		p.endIndex += 1
		if p.endIndex == len(p.end) {
			output := p.getEnvValue()
			p.Reset()
			return output, nil
		}
		return nil, nil
	default:
		output := p.begin[:p.beginIndex] + p.envName + p.end[:p.endIndex]
		p.Reset()
		if p.isMatchedBegin(input) {
			p.beginIndex += 1
			return []byte(output), nil
		} else {
			return []byte(output + string(input)), nil
		}
	}
}

func (p *sequenceParser) isMatchedBegin(input byte) bool {
	return !p.isEndBegin() && input == p.begin[p.beginIndex]
}

func (p *sequenceParser) isMatchedFileName(input byte) bool {
	if p.isEndBegin() && p.endIndex == 0 {
		currentName := p.envName + string(input)
		return p.regexEnv.MatchString(currentName)
	}

	return false
}

func (p *sequenceParser) isMatchedEnd(input byte) bool {
	if !p.isEndBegin() || len(p.envName) == 0 || p.isEndEnd() {
		return false
	}

	return p.end[p.endIndex] == input
}

func (p *sequenceParser) isEndBegin() bool {
	return p.beginIndex == len(p.begin)
}

func (p *sequenceParser) isEndEnd() bool {
	return p.endIndex == len(p.end)
}

func (p *sequenceParser) Reset() {
	p.beginIndex = 0
	p.endIndex = 0
	p.envName = ""
}

func (p *sequenceParser) Flush() []byte {
	output := p.begin[:p.beginIndex] + p.envName + p.end[:p.endIndex]
	p.Reset()
	return []byte(output)
}

func (p *sequenceParser) getEnvValue() []byte {
	return []byte(os.Getenv(p.envName))
}

func extractPattern(pattern []string) (string, string) {
	if len(pattern) != 2 {
		pattern = []string{"{{", "}}"}
	}
	begin := fmt.Sprintf("%senv \"", pattern[0])
	end := fmt.Sprintf("\"%s", pattern[1])

	return begin, end
}

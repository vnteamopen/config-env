package actions

import (
	"fmt"
	"os"
	"regexp"
)

const (
	envNamePattern      = `^[a-zA-Z_][a-zA-Z0-9_]*$`
	defaultValuePattern = `^[^"]+$`
)

type Pattern struct {
	start rune
	end   rune
}

type Parser interface {
	Transform(input byte) ([]byte, error)
	Flush() []byte
}

type sequenceParser struct {
	begin                  string
	beginIndex             int
	end                    string
	endIndex               int
	beginDefaultValue      string
	beginDefaultValueIndex int
	envName                string
	defaultValue           string
	regexEnv               *regexp.Regexp
	regexDefaultValue      *regexp.Regexp
}

func NewSeqParser(pattern []string) Parser {
	rEnv, _ := regexp.Compile(envNamePattern)
	rDefaultValue, _ := regexp.Compile(defaultValuePattern)

	begin, end := extractPattern(pattern)

	return &sequenceParser{
		begin:                  begin,
		beginIndex:             0,
		end:                    end,
		endIndex:               0,
		beginDefaultValue:      "\" \"",
		beginDefaultValueIndex: 0,
		regexEnv:               rEnv,
		regexDefaultValue:      rDefaultValue,
	}
}

func (p *sequenceParser) Transform(input byte) ([]byte, error) {
	switch true {
	case p.isMatchedBegin(input):
		p.beginIndex += 1
		return nil, nil
	case p.isMatchedValue(input):
		p.envName += string(input)
		return nil, nil
	case p.isMatchedBeginDefaultValue(input):
		p.beginDefaultValueIndex += 1
		return nil, nil
	case p.isMatchedDefaultValue(input):
		p.defaultValue += string(input)
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
		output := p.Flush()
		if p.isMatchedBegin(input) {
			p.beginIndex += 1
			return []byte(output), nil
		} else {
			return append(output, byte(input)), nil
		}
	}
}

func (p *sequenceParser) isMatchedBegin(input byte) bool {
	return !p.isEndBegin() && input == p.begin[p.beginIndex]
}

func (p *sequenceParser) isMatchedValue(input byte) bool {
	if p.isEndBegin() && p.beginDefaultValueIndex == 0 && p.endIndex == 0 {
		currentName := p.envName + string(input)
		return p.regexEnv.MatchString(currentName)
	}

	return false
}

func (p *sequenceParser) isMatchedBeginDefaultValue(input byte) bool {
	if len(p.envName) > 0 && p.endIndex == 0 && !p.isEndBeginDefaultValue() {
		return input == p.beginDefaultValue[p.beginDefaultValueIndex]
	}

	return false
}

func (p *sequenceParser) isMatchedDefaultValue(input byte) bool {
	if p.isEndBeginDefaultValue() && p.endIndex == 0 {
		currentName := p.defaultValue + string(input)
		return p.regexDefaultValue.MatchString(currentName)
	}

	return false
}

func (p *sequenceParser) isMatchedEnd(input byte) bool {
	incorrectDefaultValue := p.beginDefaultValue[0:p.beginDefaultValueIndex] + p.defaultValue
	lenMatching := getMin(len(incorrectDefaultValue), len(p.end))
	if lenMatching > 0 && incorrectDefaultValue[0:lenMatching] == p.end[0:lenMatching] {
		p.beginDefaultValueIndex = 0
		p.defaultValue = ""
		p.endIndex = lenMatching
	}

	if !p.isEndBegin() || len(p.envName) == 0 || p.isEndEnd() {
		return false
	}

	return p.end[p.endIndex] == input
}

func getMin(n1 int, n2 int) int {
	if n1 > n2 {
		return n2
	}
	return n1
}

func (p *sequenceParser) isEndBegin() bool {
	return p.beginIndex == len(p.begin)
}

func (p *sequenceParser) isEndBeginDefaultValue() bool {
	return p.beginDefaultValueIndex == len(p.beginDefaultValue)
}

func (p *sequenceParser) isEndEnd() bool {
	return p.endIndex == len(p.end)
}

func (p *sequenceParser) Reset() {
	p.beginIndex = 0
	p.endIndex = 0
	p.beginDefaultValueIndex = 0
	p.envName = ""
	p.defaultValue = ""
}

func (p *sequenceParser) Flush() []byte {
	output := p.begin[:p.beginIndex] +
		p.envName +
		p.beginDefaultValue[:p.beginDefaultValueIndex] + p.defaultValue +
		p.end[:p.endIndex]
	p.Reset()
	return []byte(output)
}

func (p *sequenceParser) getEnvValue() []byte {
	value := []byte(os.Getenv(p.envName))
	if len(value) == 0 {
		value = []byte(p.defaultValue)
	}
	return value
}

func extractPattern(pattern []string) (string, string) {
	if len(pattern) != 2 {
		pattern = []string{"{{", "}}"}
	}
	begin := fmt.Sprintf("%senv \"", pattern[0])
	end := fmt.Sprintf("\"%s", pattern[1])

	return begin, end
}

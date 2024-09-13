package parser

import (
	model "Audio2TextService/internal/models/json/yandexstt/get"
	"encoding/json"
	"strings"

	"github.com/rs/zerolog"
)

type Parser struct {
	Logger *zerolog.Logger
}

func New(logger *zerolog.Logger) *Parser {
	return &Parser{Logger: logger}
}

// Разбивает текст на массив предложений заданной длины
func (p *Parser) splitLine(line string, length int) []string {

	p.Logger.Info().Msg("Splitting line")
	defer p.Logger.Info().Msg("Finished splitting line")

	lines := make([]string, 0, length/len(line)+1)

	words := strings.Split(line, " ")

	currentLength := 0

	str := ""

	for _, word := range words {
		str += word + " "
		currentLength += len(word)
		if len(str) >= length {
			lines = append(lines, strings.TrimSpace(str))
			currentLength = 0
			str = ""
		}
	}

	return lines
}

// Подготавливает массив предложений к парсингу
func (p *Parser) prepareLines(lines []string, uniqPhraseSplitter string, maxLength int) []string {

	p.Logger.Info().Msg("Preparing lines for parsing")
	defer p.Logger.Info().Msg("Finished preparing lines for parsing")

	fullText := strings.Join(lines, " ")
	textLength := len(fullText)
	arrayLength := textLength/maxLength + 1

	newLines := make([]string, 0, arrayLength)

	// Разбиваем текст на блоки заданной длины
	for _, line := range lines {
		if len(line) <= maxLength {
			newLines = append(newLines, uniqPhraseSplitter+" "+line)
		} else {
			parts := p.splitLine(line, maxLength)
			parts[0] = uniqPhraseSplitter + " " + parts[0]
			newLines = append(newLines, parts...)
		}
	}

	return newLines
}

func (p *Parser) Parse(lines [][]byte, uniqPhraseSplitter string, maxLength int) ([]string, error) {

	p.Logger.Info().Msg("Parsing raw data to text")
	defer p.Logger.Info().Msg("Finished parsing raw data to text")

	trueLines := make([][]byte, 0, len(lines))

	for _, line := range lines {
		if strings.Contains(string(line), "\"final\":{\"alternatives\":") {
			trueLines = append(trueLines, line)
		}
	}

	GetResponses := make([]model.Response, len(trueLines))
	for i, line := range trueLines {
		err := json.Unmarshal(line, &(GetResponses[i]))
		if err != nil {
			p.Logger.Error().Msg(err.Error())
			return nil, err
		}
	}

	currentChannelTag := GetResponses[0].Result.ChannelTag
	speachParts := make([]string, 0, len(GetResponses))
	var currentChannelText string = ""

	for i, resp := range GetResponses {
		if resp.Result.ChannelTag != currentChannelTag {
			currentChannelTag = resp.Result.ChannelTag
			speachParts = append(speachParts, strings.TrimSpace(currentChannelText))
			currentChannelText = ""
		}
		currentChannelText += " "
		currentChannelText += resp.Result.Final.Alternatives[0].Text
		if i == (len(GetResponses) - 1) {
			speachParts = append(speachParts, strings.TrimSpace(currentChannelText))
		}
	}

	return p.prepareLines(speachParts, uniqPhraseSplitter, maxLength), nil
}

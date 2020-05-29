package dialogflow

import (
	"strings"

	"github.com/niko-cb/covid19datascraper/server/model"
)

// SMessage: the S stands for Symptom
type SMessage struct {
	Content *SContent
}

// SContent: the S stands for Symptom
type SContent struct {
	Header string
	Common string
	Rare   string
	Severe string
}

func (c *SContent) String() string {
	return Header + c.Common + c.Rare + c.Severe
}

func NewSMessage(s *model.Symptoms) *SMessage {
	return &SMessage{
		&SContent{
			Header: Header,
			Common: CommonPrefix + makeString(s.Common),
			Rare:   RarePrefix + makeString(s.Rare),
			Severe: SeverePrefix + makeString(s.Severe),
		},
	}
}

func makeString(s []string) string {
	return strings.Join(s, ", ")
}

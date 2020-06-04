package models

import (
	"strings"

	"github.com/niko-cb/covid19datascraper/app/models/helpers/dialogflow"
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
	return dialogflow.Header + c.Common + c.Rare + c.Severe
}

func NewSMessage(s *Symptoms) *SMessage {
	return &SMessage{
		&SContent{
			Header: dialogflow.Header,
			Common: dialogflow.CommonPrefix + makeString(s.Common),
			Rare:   dialogflow.RarePrefix + makeString(s.Rare),
			Severe: dialogflow.SeverePrefix + makeString(s.Severe),
		},
	}
}

func makeString(s []string) string {
	return strings.Join(s, ", ")
}

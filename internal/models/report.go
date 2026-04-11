package models

type Tokens struct {
	InPortfolio       string
	GainersAndLoosers string
}

type Report struct {
	MarketCap string
	Tokens    Tokens
	AISummary string
}

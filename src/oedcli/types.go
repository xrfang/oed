package oed

type GrammaticalFeature struct {
	Text string
	Type string
}

type Example struct {
	Text string
}

type Sense struct {
	Definitions []string
	Domains     []string
	Examples    []Example
}

type Entry struct {
	GrammaticalFeatures []GrammaticalFeature
	Senses              []Sense
}

type Pronunciation struct {
	AudioFile        string
	Dialects         []string
	PhoneticSpelling string
}

type LexicalEntry struct {
	Entries         []Entry
	LexicalCategory string
	Pronunciations  []Pronunciation
	Text            string
}

type QueryResult struct {
	LexicalEntries []LexicalEntry
	Word           string
}

type QueryReply struct {
	Result QueryResult
}

package oed

type note struct {
	Text string
	Type string
}

type Example struct {
	Text string
}

type ThesaurusLink struct {
	EntryID string `json:"entry_id"`
	SenseID string `json:"sense_id"`
}

type Sense struct {
	Definitions    []string
	Domains        []string
	Examples       []Example
	Notes          []note
	SubSenses      []Sense
	ThesaurusLinks []ThesaurusLink
}

type Entry struct {
	Senses []Sense
}

type Pronunciation struct {
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
}

type QueryReply struct {
	Results []QueryResult
}

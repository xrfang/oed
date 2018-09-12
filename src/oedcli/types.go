package oed

type note struct {
	Text string
	Type string
}

type textEntry struct {
	Text string
}

type ThesaurusLink struct {
	EntryID string `json:"entry_id"`
	SenseID string `json:"sense_id"`
}

type Sense struct {
	ID             string
	Definitions    []string
	Domains        []string
	Examples       []textEntry
	Notes          []note
	Registers      []string
	SubSenses      []Sense
	ThesaurusLinks []ThesaurusLink
	Synonyms       []textEntry
	Antonyms       []textEntry
}

type Entry struct {
	GrammaticalFeatures []note
	Senses              []Sense
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

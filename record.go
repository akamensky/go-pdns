package pdns

// RRsets structure with JSON API metadata
type RRsets struct {
	Sets []*RRset `json:"rrsets"`
}

// RRset struct
type RRset struct {
	Name       string     `json:"name"`
	Type       recordType `json:"type"`
	TTL        uint       `json:"ttl"`
	Records    []*Record  `json:"records"`
	ChangeType string     `json:"changetype"`
}

// Record structure with JSON API metadata
type Record struct {
	Content  string `json:"content"`
	Disabled bool   `json:"disabled"`
	SetPTR   bool   `json:"set-ptr"`
}

type recordType string

const (
	RecordTypeA     recordType = "A"
	RecordTypeAAAA  recordType = "AAAA"
	RecordTypeCNAME recordType = "CNAME"
	RecordTypeTXT   recordType = "TXT"
	RecordTypeNS    recordType = "NS"
	RecordTypeALIAS recordType = "ALIAS"
	RecordTypeCAA   recordType = "CAA"
	RecordTypeMX    recordType = "MX"
	RecordTypeSOA   recordType = "SOA"
	RecordTypePTR   recordType = "PTR"
)

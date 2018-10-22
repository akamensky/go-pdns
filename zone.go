package pdns

import "encoding/json"

// Zone struct
type Zone struct {
	// Internal use fields
	client *Client `json:"-"`

	// Required fields for create
	Name        string   `json:"name"`
	Kind        zoneKind `json:"kind"`
	Nameservers []string `json:"nameservers"`

	// Not required fields
	ID     string  `json:"id,omitempty"`
	URL    string  `json:"url,omitempty"`
	Serial int     `json:"serial,omitempty"`
	RRsets []RRset `json:"rrsets,omitempty"`
}

type zoneKind string

const (
	ZoneKindNative zoneKind = "Native"
	ZoneKindMaster zoneKind = "Master"
	ZoneKindSlave  zoneKind = "Slave"
)

func (z *Zone) String() string {
	b, _ := json.MarshalIndent(z, "", "\t")
	return string(b)
}

func (z *Zone) Delete() error {
	result := new(interface{})
	f := new(failure)

	_, err := z.client.getSling().Delete(z.URL).Receive(result, f)
	if err != nil {
		return err
	}

	if err := f.getError(); err != nil {
		return err
	}

	return nil
}

func (z *Zone) AddRecord(name string, rtype recordType, ttl uint, content []string) error {
	return z.ChangeRecord(name, rtype, ttl, content)
}

func (z *Zone) ChangeRecord(name string, rtype recordType, ttl uint, content []string) error {
	rrset := new(RRset)
	rrset.Name = name
	rrset.Type = rtype
	rrset.TTL = ttl
	rrset.ChangeType = "REPLACE"

	for _, c := range content {
		r := &Record{Content: c, Disabled: false, SetPTR: false}
		rrset.Records = append(rrset.Records, r)
	}

	return z.patchRRSet(rrset)
}

func (z *Zone) DeleteRecord(name string, rtype recordType) error {
	rrset := new(RRset)
	rrset.Name = name
	rrset.Type = rtype
	rrset.ChangeType = "DELETE"

	return z.patchRRSet(rrset)
}

func (z *Zone) patchRRSet(rrset *RRset) error {
	// Make sure name is FQDN
	//rrset.Name = rrset.Name

	rrsets := new(RRsets)
	rrsets.Sets = append(rrsets.Sets, rrset)

	result := new(interface{})
	f := new(failure)

	_, err := z.client.getSling().Patch(z.URL).BodyJSON(rrsets).Receive(result, f)
	if err != nil {
		return err
	}

	if err := f.getError(); err != nil {
		return err
	}

	return nil
}

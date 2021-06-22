package parties

type PartyAttributes map[string][]string

func (a PartyAttributes) Get(key string) string {
	if values, ok := a[key]; ok {
		if len(values) > 0 {
			return values[0]
		}
	}
	return ""
}

func (a PartyAttributes) Set(key, value string) {
	a[key] = []string{value}
}

func (a PartyAttributes) Add(key, value string) {
	a[key] = append(a[key], value)
}

type Party struct {
	ID         string          `json:"id" bson:"id"`
	PartyTypes []string        `json:"partyTypes" bson:"partyTypes"`
	Attributes PartyAttributes `json:"attributes" bson:"attributes"`
}

func (p *Party) HasPartyType(partyType string) bool {
	for _, p := range p.PartyTypes {
		if p == partyType {
			return true
		}
	}
	return false
}

func (p *Party) AddPartyType(partyType string) {
	if p.HasPartyType(partyType) {
		return
	}
	p.PartyTypes = append(p.PartyTypes, partyType)
}

func (p *Party) String() string {
	return p.ID
}

type PartyList struct {
	Items []*Party `json:"items" bson:"items"`
}

func (p *Party) GetAttribute(name string) []string {
	if v, ok := p.Attributes[name]; ok {
		return v
	}
	return nil
}

func (p *Party) HasAttribute(name string) bool {
	_, ok := p.Attributes[name]
	return ok
}

func (p *Party) FindAttributeValue(name string) interface{} {
	if v, ok := p.Attributes[name]; ok {
		return v
	}
	return nil
}

func (pl *PartyList) FindByID(id string) *Party {
	for _, v := range pl.Items {
		if id == v.ID {
			return v
		}
	}
	return nil
}

package eui

type Registration struct {
	Prefix string
	Name   string
}

type Registrations []Registration

func (r Registrations) Len() int           { return len(r) }
func (r Registrations) Less(i, j int) bool { return r[i].Prefix < r[j].Prefix }
func (r Registrations) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }

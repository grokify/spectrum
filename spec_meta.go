package spectrum

type SpecMeta struct {
	Counts SpecMetaCounts
	Names  SpecMetaNames
}

func NewSpecMeta() *SpecMeta {
	return &SpecMeta{
		Counts: SpecMetaCounts{},
		Names:  SpecMetaNames{}}
}

func (m *SpecMeta) Inflate() {
	m.Counts = m.Names.Counts()
}

type SpecMetaCounts struct {
	Endpoints int
	Paths     int
	Models    int
}

type SpecMetaNames struct {
	Endpoints []string
	Models    []string
	Paths     []string
}

func (n SpecMetaNames) Counts() SpecMetaCounts {
	return SpecMetaCounts{
		Endpoints: len(n.Endpoints),
		Models:    len(n.Models),
		Paths:     len(n.Paths),
	}
}

type SpecMetaSet struct {
	Data map[string]SpecMeta
}

func NewSpecMetaSet() *SpecMetaSet {
	return &SpecMetaSet{Data: map[string]SpecMeta{}}
}

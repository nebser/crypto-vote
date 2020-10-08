package party

type Party struct {
	Name    string `json:"name"`
	Address string `json:"address"`
	Balance int    `json:"balance"`
}

type Parties []Party

func (p Parties) Len() int {
	return len(p)
}

func (p Parties) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func (p Parties) Less(i, j int) bool {
	switch {
	case p[i].Balance < p[j].Balance:
		return true
	case p[i].Balance > p[j].Balance:
		return false
	default:
		return p[i].Name <= p[j].Name
	}
}

type GetPartyFn func(string) (*Party, error)

type GetPartiesFn func() (Parties, error)

type SavePartyFn func(Party) error

package party

type Party struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

type Parties []Party

type GetPartyFn func(string) (*Party, error)

type GetPartiesFn func() (Parties, error)

type SavePartyFn func(Party) error

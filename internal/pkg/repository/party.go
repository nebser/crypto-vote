package repository

import (
	"encoding/json"

	"github.com/boltdb/bolt"
	_party "github.com/nebser/crypto-vote/internal/pkg/party"
	"github.com/pkg/errors"
)

type party struct {
	Name    string `json:"name"`
	Address string `json:"address"`
}

func partiesBucket() []byte {
	return []byte("parties")
}

func newParty(p _party.Party) party {
	return party{
		Address: p.Address,
		Name:    p.Name,
	}
}

func (p party) toParty() _party.Party {
	return _party.Party{
		Address: p.Address,
		Name:    p.Name,
	}
}

func SaveParty(db *bolt.DB) _party.SavePartyFn {
	return func(party _party.Party) error {
		return db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket(partiesBucket())
			if b == nil {
				created, err := tx.CreateBucket(partiesBucket())
				if err != nil {
					return errors.Wrapf(err, "Failed to create bucket %s", partiesBucket())
				}
				b = created
			}
			raw, err := json.Marshal(newParty(party))
			if err != nil {
				return errors.Wrap(err, "Failed to serialize party")
			}
			if err := b.Put([]byte(party.Address), raw); err != nil {
				return errors.Wrapf(err, "Failed to save party %#v", party)
			}
			return nil
		})
	}
}

func GetParty(db *bolt.DB) _party.GetPartyFn {
	return func(address string) (*_party.Party, error) {
		var result *_party.Party
		err := db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket(partiesBucket())
			if b == nil {
				return nil
			}
			raw := b.Get([]byte(address))
			if raw == nil {
				return nil
			}
			var partyRaw party
			if err := json.Unmarshal(raw, &partyRaw); err != nil {
				return errors.Wrap(err, "Failed to unmarshal parties")
			}
			party := partyRaw.toParty()
			result = &party
			return nil
		})
		return result, err
	}
}

func GetParties(db *bolt.DB) _party.GetPartiesFn {
	return func() (_party.Parties, error) {
		result := _party.Parties{}
		err := db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket(partiesBucket())
			if b == nil {
				return nil
			}
			c := b.Cursor()
			for addressRaw, partyRaw := c.First(); addressRaw != nil; addressRaw, partyRaw = c.Next() {
				var dbParty party
				if err := json.Unmarshal(partyRaw, &dbParty); err != nil {
					return errors.Wrapf(err, "Failed to unmarshal db party %s", partyRaw)
				}
				result = append(result, dbParty.toParty())
			}
			return nil
		})
		return result, err
	}
}

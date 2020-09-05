package repository

import "github.com/boltdb/bolt"

func transactionsArray(transactions ...func(*bolt.Tx) error) func(*bolt.Tx) error {
	return func(tx *bolt.Tx) error {
		for _, t := range transactions {
			if err := t(tx); err != nil {
				return err
			}
		}
		return nil
	}
}

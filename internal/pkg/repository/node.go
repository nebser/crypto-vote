package repository

import (
	"encoding/json"

	"github.com/boltdb/bolt"
	"github.com/nebser/crypto-vote/internal/pkg/blockchain"
	"github.com/pkg/errors"
)

func nodesBucket() []byte {
	return []byte("nodes")
}

func SaveNode(db *bolt.DB) blockchain.SaveNodeFn {
	return func(node blockchain.Node) error {
		err := db.Update(func(tx *bolt.Tx) error {
			b := tx.Bucket(nodesBucket())
			if b == nil {
				created, err := tx.CreateBucket(nodesBucket())
				if err != nil {
					return errors.Wrap(err, "Failed to create nodes bucket")
				}
				b = created
			}
			rawNode, err := json.Marshal(node)
			if err != nil {
				return errors.Wrapf(err, "Failed to marshal node %#v", node)
			}
			if err := b.Put([]byte(node.ID), rawNode); err != nil {
				return errors.Wrapf(err, "Failed to put node %#v", node)
			}
			return nil
		})
		return err
	}
}

func GetNode(db *bolt.DB) blockchain.GetNodeFn {
	return func(id string) (*blockchain.Node, error) {
		var node *blockchain.Node
		err := db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket(nodesBucket())
			if b == nil {
				return nil
			}
			raw := b.Get([]byte(id))
			if len(raw) == 0 {
				return nil
			}
			if err := json.Unmarshal(raw, node); err != nil {
				return errors.Wrapf(err, "Failed to unmarshal %s into node while getting node %s", raw, id)
			}
			return nil
		})
		return node, err
	}
}

func GetNodes(db *bolt.DB) blockchain.GetNodesFn {
	return func() (blockchain.Nodes, error) {
		var nodes blockchain.Nodes
		assemble := func(k, v []byte) error {
			var node blockchain.Node
			if err := json.Unmarshal(v, &node); err != nil {
				return errors.Wrapf(err, "Failed unmarshal value %s into node", v)
			}
			nodes = append(nodes, node)
			return nil
		}
		err := db.View(func(tx *bolt.Tx) error {
			b := tx.Bucket(nodesBucket())
			if b == nil {
				return nil
			}
			if err := b.ForEach(assemble); err != nil {
				return errors.Wrap(err, "Failed to retrieve all nodes")
			}
			return nil
		})
		return nodes, err
	}
}

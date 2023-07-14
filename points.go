package main

import (
	"encoding/binary"
	"fmt"

	"go.etcd.io/bbolt"
)

func setPoints(id string, amount int) error {
	if err := POINTS_DB.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("points"))
		if err != nil {
			return fmt.Errorf("error creating bucket: %w", err)
		}

		n := binary.BigEndian.AppendUint64([]byte{}, uint64(amount))
		if err := b.Put([]byte(id), n); err != nil {
			return fmt.Errorf("error setting points: %w", err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("error in transaction: %w", err)
	}

	return nil
}

func getPoints(id string) (int, error) {
	var points int
	if err := POINTS_DB.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("points"))
		if err != nil {
			return fmt.Errorf("error creating bucket: %w", err)
		}

		n := b.Get([]byte(id))
		if n != nil {
			points = int(binary.BigEndian.Uint64(b.Get([]byte(id))))
		}

		return nil
	}); err != nil {
		return points, fmt.Errorf("error in transaction: %w", err)
	}

	return points, nil
}

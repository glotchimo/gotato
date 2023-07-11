package main

import (
	"encoding/binary"
	"fmt"

	"go.etcd.io/bbolt"
)

func reward(id string) (int, error) {
	var points int
	if err := POINTS_DB.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("points"))
		if err != nil {
			return fmt.Errorf("error creating bucket: %w", err)
		}

		var oldPoints uint64
		oldPointsRaw := b.Get([]byte(id))
		if oldPointsRaw != nil {
			oldPoints = binary.BigEndian.Uint64(b.Get([]byte(id)))
			points = int(oldPoints) + REWARD
		}

		newPoints := binary.BigEndian.AppendUint64([]byte{}, oldPoints+uint64(REWARD))
		if err := b.Put([]byte(id), newPoints); err != nil {
			return fmt.Errorf("error putting new score: %w", err)
		}

		return nil
	}); err != nil {
		return points, fmt.Errorf("error in transaction: %w", err)
	}

	return points, nil
}

func getPoints(id string) (int, error) {
	var points int
	if err := POINTS_DB.Update(func(tx *bbolt.Tx) error {
		b, err := tx.CreateBucketIfNotExists([]byte("points"))
		if err != nil {
			return fmt.Errorf("error creating bucket: %w", err)
		}

		oldScoreRaw := b.Get([]byte(id))
		if oldScoreRaw != nil {
			points = int(binary.BigEndian.Uint64(b.Get([]byte(id))))
		}

		return nil
	}); err != nil {
		return points, fmt.Errorf("error in transaction: %w", err)
	}

	return points, nil
}

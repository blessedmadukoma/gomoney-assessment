package seeds

import (
	"context"
	"log"
	"time"

	"github.com/blessedmadukoma/gomoney-assessment/data"
)

func (s Seed) UsersSeeder() {
	ctx := context.Background()
	for _, user := range data.Users {
		// meta, _ := json.Marshal(&c.SerializedMeta)

		now := time.Now()
		user.CreatedAt = now
		user.UpdatedAt = now

		_, err := s.collections["users"].InsertOne(ctx, user)

		if err != nil {
			log.Fatalf("cannot seed users table: %v", err)
		}
	}
}

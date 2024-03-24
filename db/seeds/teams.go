package seeds

import (
	"context"
	"log"
	"time"

	"github.com/blessedmadukoma/gomoney-assessment/data"
)

func (s Seed) TeamsSeeder() {
	ctx := context.Background()
	for _, team := range data.Teams {
		// meta, _ := json.Marshal(&c.SerializedMeta)

		now := time.Now()
		team.CreatedAt = now
		team.UpdatedAt = now

		_, err := s.collections["teams"].InsertOne(ctx, team)

		if err != nil {
			log.Fatalf("cannot seed teams table: %v", err)
		}
	}
}

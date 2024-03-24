package seeds

import (
	"context"
	"log"
	"time"

	"github.com/blessedmadukoma/gomoney-assessment/data"
)

func (s Seed) FixturesSeeder() {
	ctx := context.Background()
	for _, fixture := range data.Fixtures {
		// meta, _ := json.Marshal(&c.SerializedMeta)

		now := time.Now()
		fixture.CreatedAt = now
		fixture.UpdatedAt = now

		_, err := s.collections["fixtures"].InsertOne(ctx, fixture)

		if err != nil {
			log.Fatalf("cannot seed fixtures table: %v", err)
		}
	}
}

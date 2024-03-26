package seeds

import (
	"log"
	"time"

	"github.com/blessedmadukoma/gomoney-assessment/data"
	"github.com/gin-gonic/gin"
)

func (s Seed) FixturesSeeder() {
	// ctx := context.Background()
	ctx := &gin.Context{}
	for _, fixture := range data.Fixtures {
		// meta, _ := json.Marshal(&c.SerializedMeta)

		now := time.Now()
		fixture.CreatedAt = now
		fixture.UpdatedAt = now
		fixture.Fixture = srv.GenerateFixture(ctx, fixture.Home, fixture.Away)
		fixture.Link = srv.GenerateFixtureLink(ctx, fixture.Home, fixture.Away)

		_, err := s.collections["fixtures"].InsertOne(ctx, fixture)

		if err != nil {
			log.Fatalf("cannot seed fixtures table: %v", err)
		}
	}
}

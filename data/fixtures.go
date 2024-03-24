package data

import db "github.com/blessedmadukoma/gomoney-assessment/db/models"

var Fixtures = []db.CreateFixturesParams{
	{
		Home:   "Aston Villa FC",
		Away:   "Arsenal FC",
		Status: "pending",
	},
	{
		Home:   "Chelsea FC",
		Away:   "Manchester United",
		Status: "completed",
	},
	{
		Home:   "Arsenal FC",
		Away:   "Liverpool FC",
		Status: "completed",
	},
	{
		Home:   "Everton",
		Away:   "Watford",
		Status: "pending",
	},
	{
		Home:   "Chelsea FC",
		Away:   "Manchester City",
		Status: "completed",
	},
	{
		Home:   "Manchester United",
		Away:   "Liverpool FC",
		Status: "completed",
	},
	{
		Home:   "Arsenal FC",
		Away:   "Newcastle FC",
		Status: "pending",
	},
	{
		Home:   "chelsea",
		Away:   "Burnley",
		Status: "completed",
	},
	{
		Home:   "liverpool",
		Away:   "Brentford",
		Status: "pending",
	},
}

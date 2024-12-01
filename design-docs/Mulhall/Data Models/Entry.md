Represents a single week's picks for a single Contestant, including both their selected (locked-in) and suggested picks.

- ID `primary key`
- Schedule `Schedule.ID`
- Selected `Team.ID`
- Suggested `[]Team.ID`
## Example
```
{
	"id": <unique auto-generated ID>,
	"schedule": "scheduleId",
	"selected": "team1Id",
	"suggested": [
		"team2Id",
		"team3Id",
		// etc.
	]
}
```
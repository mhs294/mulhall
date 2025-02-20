Represents the slate of matchups for a single week from which Contestants will Select a Pick.
Pre-defined static content, but is subject to change and the system must be able to account for this (e.g. - playoffs, game cancellations)

- ID `primary key`
- Year `int`
- Week `int`
	- 1-18 for NFL regular season
	- 19 = Wild Card weekend
	- 20 = Divisional Round weekend
	- 21 = Conference Championship weekend
	- 22 = Super Bowl weekend
- Opens `datetime`
	- Date/Time when the schedule opens for Picks
- Closes `datetime`
	- Date/Time when the schedule closes for Picks
- Matchups `[]Matchup`
## Example
```
{
	"id": <unique auto-generated ID>,
	"year": 2024,
	"week": 1,
	"start": "2025-09-02T08:00:00.000Z",
	"end": "2025-09-09T07:59:59.999Z"
	"opens": "2025-09-02T08:00:00.000Z",
	"closes": "2025-09-04T23:45:00.000Z",
	"matchups": [
		"matchupId1",
		"matchupId2",
		// etc.
	]
}
```
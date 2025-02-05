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
- Matchups `[]Matchup.ID`
	-
## Example
```
{
	"id": <unique auto-generated ID>,
	"year": 2024,
	"week": 1,
	"matchups": [
		"matchupId1",
		"matchupId2",
		// etc.
	]
}
```
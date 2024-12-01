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
- Matchups `[]Matchup`
	- Matchup `object`
		- Away Team `Team.ID`
		- Home Team `Team.ID`
		- Date/Time `datetime`
## Example
```
{
	"id": <unique auto-generated ID>,
	"year": 2024,
	"week": 1,
	"matchups": [
		{
			"awayTeam": "team1Id",
			"homeTeam": "team2Id",
			"dateTime": "2024-09-06T00:15:00.000Z"
		},
		// etc.
	]
}
```
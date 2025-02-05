Represents a single week's picks for a single Contestant, including both their selected (locked-in) and suggested picks.

- ID `primary key`
- Schedule `Schedule.ID`
	- The Schedule that contains the Matchups available for Picks with this Entry
- SelectedPick `Pick`
	- The Pick that will be counted as the Contestant's submission for the Entry
- SuggestedPick `[]Pick`
	- Suggested Picks that can be elevated to Selected Pick by a Team OWNER
## Example
```
{
	"id": <unique auto-generated ID>,
	"schedule": "scheduleId",
	"selectedPick": {
		"matchup": "matchupId1",
		"team": "teamId1"
	},
	"suggestedPicks": [
		{
			"matchup": "matchupId2",
			"team": "teamId2"
		},
		{
			"matchup": "matchupId3",
			"team": "teamId3"
		},
		// etc.
	]
}
```
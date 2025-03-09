Represents a single week's picks for a single Contestant, including both their selected (locked-in) and suggested picks.

- ID `primary key`
- Contestant `Contestant.ID`
- Schedule `Schedule.ID`
	- The Schedule that contains the Matchups available for Picks with this Entry
- SelectedPick `Pick`
	- The Pick that will be counted as the Contestant's submission for the Entry
	- Will always be managed by the system to ensure it contains exactly 0 or 1 key/value pairs
- SuggestedPick `[]Pick`
	- Suggested Picks that can be elevated to Selected Pick by a Team OWNER
## Example
```
{
	"id": <unique auto-generated ID>,
	"schedule": "scheduleId",
	"selectedPick": {
		"matchupId1": "teamId1"
	},
	"suggestedPicks": {
		"matchupId2": "teamId2",
		"matchupId3": "teamId3"
	}
}
```
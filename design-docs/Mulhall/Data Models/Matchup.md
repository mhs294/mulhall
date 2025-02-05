 Represents a single game between two Teams at a specific date/time.

- ID `primary key`
- Away Team `Team.ID`
- Home Team `Team.ID`
- Date/Time `datetime`
## Example
```
{
	"id": <unique auto-generated ID>,
	"awayTeam": "teamId1",
	"homeTeam": "teamId2",
	"dateTime": "2025-02-09T23:55:00.000Z"
}
```
Represents a record of a single event that occurs as a result of a user interacting with the system. This allows for strict record-keeping of all changes made in case something requires investigation or troubleshooting.

- ID `primary key`
- User `User.ID`
- Action `string, enumerated`
- Resource `string, enumerated`
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
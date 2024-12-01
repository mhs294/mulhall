Represents a single entrant within the Pool, which can be managed by one or more Users. Each Contestant must have at least one User with the OWNER Role associated to it.

- ID `primary key`
- Name `string`
- Users `[]User.ID`
- Status `string, enumerated`
	- See [this](obsidian://open?vault=Mulhall&file=Brainstorming%2FContestant%20Statuses.canvas) for more info

## Example
```
{
	"id": <unique auto-generated ID>,
	"name": "Swamy Says",
	"users": [
		"userId1",
		// etc.
	],
	"status": "ACTIVE"
}
```
Represents a single entrant within the Pool, which can be managed by one or more Users. Each Contestant must have at least one User with the OWNER Role associated to it.

- ID `primary key`
- Name `string`
- Authorized Users `[]struct`
	- User ID `User.ID`
	- Role `string, enumerated`
		- See [this](obsidian://open?vault=Mulhall&file=Brainstorming%2FRoles%20%2B%20Functions.canvas) for more info
- Status `string, enumerated`
	- See [this](obsidian://open?vault=Mulhall&file=Brainstorming%2FContestant%20Statuses.canvas) for more info

## Example
```
{
	"id": <unique auto-generated ID>,
	"name": "Swamy Says",
	"authorizedUusers": [
		{
			"userId": "userId1",
			"role": "OWNER"
		},
		// etc.
	],
	"status": "ACTIVE"
}
```
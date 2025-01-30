Represents a Pool, which defines the rules for the elimination game as well as the participating Contestants.

- ID `primary key`
- Contestants `[]Contestant.ID`
- Active `bool`
	- (logical delete)

## Example
```
{
	"id": <unique auto-generated ID>,
	"contestants": [
		"contestantId1",
		"contestantId2",
		// etc.
	]
	"active": true
}
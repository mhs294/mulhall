Represents a Pool, which defines the rules for the elimination game as well as the participating Contestants.

- ID `primary key`
- Name `string`
- Contestants `map[Contestant.ID]struct{}`
	- (acts as a Set of `Contestant.ID` values)
- Active `bool`
	- (logical delete)
- Completed `bool`
	- Indicates whether the pool has met its completion condition or whether it is still ongoing (might need to make a Pool-specific status enum later)

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
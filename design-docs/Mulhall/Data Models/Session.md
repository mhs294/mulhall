Represents a single session for a single logged in User. The session is used to authenticate the User for rendering site views as well as any actions performed via API calls.

- ID `primary key`
- Expiration `date/time`

## Example
```
{
	"id": <unique auto-generated ID>,
	"email": 
	"expiration": "2024-12-25T00:00:00.000Z"
}
```
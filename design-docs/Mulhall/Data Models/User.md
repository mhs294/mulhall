Represents an individual user of the application, who is an authorized member of one or more Pool Contestants.
Users are created upon invitation and then Authenticated is set to true when the invitation is accepted using the email link and a password is established.

- ID `primary key`
- Email `string`
- Salt `string`
	- Unique, randomly generated salt for the user's current password, which is used to produce the resulting Hash from the user's original password string.
- Hash `string`
	- Hash value used to compare against salted and hashed password input when a user attempts to login.
- Role `string, enumerated`
	- See [this](obsidian://open?vault=Mulhall&file=Brainstorming%2FRoles%20%2B%20Functions.canvas) for more info
- Active `bool`
	- Indicates whether this user is active and allowed to login. This will be set to true initially after the user accepts their invitation from an email invite link and creates their password.

## Example
```
{
	"id": <unique auto-generated ID>,
	"email": "bob@gmail.com",
	"salt": "hd723hk$23@kkjd",
	"hash": "af7cbbc1eacf780e70344af1a4b16698",
	"role": "OWNER",
	"active": true
}
```
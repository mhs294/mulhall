	Represents an invitation for a [User](User) to join the Pool, either by an Administrator, or by a Contestant Owner. This data model is used only as a temporary mechanism to authenticate and confirm invitations via email links.

- ID `primary key`
- Email `string`
- Contestant `Contestant.ID`
	- The Contestant the invited User will be associated with once their Invite has been accepted
- Role `string, enumerated`
	- The Role for the new User on the associated Contestant once the Invite has been accepted
	- See [this](obsidian://open?vault=Mulhall&file=Brainstorming%2FRoles%20%2B%20Functions.canvas) for more info
- Inviting User `User.ID`
	- The User responsible for creating and sending the invitation to the new User
- Token `string`
- Expiration `date/time`
- Accepted `bool`
	- Indicates whether the email invitation link corresponding to the Invite has been accepted by the intended user.
## Example
```
{
	"id": <unique auto-generated ID>,
	"email": "bob@gmail.com",
	"invitingUser": "userId1",
	"token": "19z2tzroufu0mdwyigs0qc92u",
	"expiration": "2024-12-25T00:00:00.000Z"
	"accepted": false
}
```

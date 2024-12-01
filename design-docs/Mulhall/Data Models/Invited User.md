Represents a [User](User) that has been invited to join the Pool by an administrator.
This data model is used only as a temporary mechanism to authenticate and confirm invitations via email links.

- ID `primary key`
- Email `string`
- Confirmation Key `string`
## Example
```
{
	"id": <unique auto-generated ID>,
	"email": "bob@gmail.com",
	"confirmationKey: "19z2tzroufu0mdwyigs0qc92u"
}
```

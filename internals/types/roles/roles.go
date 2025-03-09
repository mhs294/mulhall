package roles

type Role string

// The enumerated Roles available to an Authorized User of a Pool Contestant.
const (
	// OWNER is authorized full control over a Contestant, including managing its picks and authorized Users.
	OWNER Role = "Owner"
	// MANAGER is authorized to manage picks for a Contestant, including Suggested and Selected picks.
	MANAGER Role = "Manager"
	// VIEWER is authorized to view picks and Suggest picks for a Contestant, but may not Select a pick.
	VIEWER Role = "Viewer"
)

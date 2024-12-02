package roles

type Role string

// The enumerated Roles available to an Authorized User of a Pool Contestant.
const (
	// OWNER is authorized full control over a Contestant, including managing its Picks and authorized Users.
	OWNER Role = "Owner"
	// MANAGER is authorized to manage Picks for a Contestant, including Suggested and Selected Picks.
	MANAGER Role = "Manager"
	// VIEWER is authorized to view Picks and Suggest Picks for a Contestant, but may not Select a Pick.
	VIEWER Role = "Viewer"
)

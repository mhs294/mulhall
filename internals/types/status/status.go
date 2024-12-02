package status

type Status string

// The enumerates Statuses for a Pool Contestant.
const (
	// The Contestant is still eligible to win the Pool.
	ACTIVE Status = "Active"
	// The Contestant has been eliminated from the Pool, but is eligible for reinstatement.
	ELIMINATED Status = "Eliminated"
	// The Contestant has been eliminated from the Pool and is not eligible for reinstatement.
	DISQUALIFIED Status = "Disqualified"
)

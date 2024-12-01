package types

// Team contains the basic identifying details of an NFL team.
type Team struct {
	ID        string `json:"id"`
	Shorthand string `json:"shorthand"`
	Location  string `json:"location"`
	Name      string `json:"name"`
}

package wikitionaryparser

type Language struct {
	Name     string            `json:"name"`
	Sections map[string]string `json:"sections"`
	Details  []Detail          `json:"details"`
}

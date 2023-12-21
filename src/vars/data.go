package vars

type GenminiResponse struct {
	Candidates     []Candidates   `json:"candidates"`
	PromptFeedback PromptFeedback `json:"promptFeedback"`
}
type Parts struct {
	Text string `json:"text"`
}
type Content struct {
	Parts []Parts `json:"parts"`
	Role  string  `json:"role"`
}
type SafetyRatings struct {
	Category    string `json:"category"`
	Probability string `json:"probability"`
}
type Candidates struct {
	Content       Content         `json:"content"`
	FinishReason  string          `json:"finishReason"`
	Index         int             `json:"index"`
	SafetyRatings []SafetyRatings `json:"safetyRatings"`
}
type PromptFeedback struct {
	SafetyRatings []SafetyRatings `json:"safetyRatings"`
}

type Contents struct {
	Parts Parts `json:"parts"`
}

type GenminiRequest struct {
	Contents []Contents `json:"contents"`
}

package models

type Integration struct {
	Model
	SlackToken string `json:"slack_token"`
	SlackWorkspace string `json:"slack_workspace"`
	ZendeskSubdomain int32 `json:"slac_subdomain"`
}

package models

type Integration struct {
	Model
	SlackToken string `json:"slack_token"`
	SlackTeamId string `json:"slack_team_id"`
	SlackWorkspace string `json:"slack_workspace"`
	ZendeskSubdomain string `json:"zendesk_subdomain"`
	ZendeskInstancePushId string `json:"zendesk_instance_push_id"`
	ZendeskToken string `json:"zendesk_token"`
}

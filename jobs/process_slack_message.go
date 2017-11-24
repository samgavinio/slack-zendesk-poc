package jobs

import (
	"fmt"
	"time"
	"encoding/json"
	"bytes"
	"net/http"

	"github.com/zendesk/slack-poc/operation"
	"github.com/zendesk/slack-poc/models"
)

type (
	externalResource struct {
		ExternalId string `json:"external_id"`
		Message string `json:"message"`
		CreatedAt string `json:"created_at"`
		Author author `json:"author"`
		AllowChannelback bool `json:"allow_channelback"`
	}
	author struct {
		ExternalId string `json:"external_id"`
		Name string `json:"name"`
	}
	postBody struct {
		ExternalResources []externalResource `json:"external_resources"`
		InstancePushId string `json:"instance_push_id"`
	}
)

func ProcessSlackMessage(queue string, args ...interface{}) error {
	fmt.Println("You are working on a job.")
	payload, ok := args[0].(map[string]interface{})
	if !ok {
		return fmt.Errorf("Invalid parameters %v to insert worker. Expected map, was %T.", args, args[0])
	}

	event, ok := payload["event"].(map[string]interface{})
	if !ok {
		return fmt.Errorf("Invalid parameters to insert worker. Expected map, was %T.", event)
	}

	teamId := payload["team_id"].(string)
	message := event["text"].(string)
	userId := event["user"].(string)
	eventId := event["event_ts"].(string)
	channel := event["channel"].(string)

	var integration models.Integration
	query := operation.DB.Where(&models.Integration{
		SlackTeamId: teamId,
	})
	if query.First(&integration).RecordNotFound() {
		return fmt.Errorf("Integration for slack team %s not found", teamId)
	}

	person := author{
		ExternalId: userId,
		Name: "Somebody from Slack",
	}
	now := time.Now().UTC()
	resource := externalResource{
		ExternalId: fmt.Sprintf("%s-%s", channel, eventId),
		Message: message,
		CreatedAt: now.Format("2006-01-02T15:04:05Z"),
		Author: person,
		AllowChannelback: true,
	}
	body := postBody{
		ExternalResources: []externalResource{resource},
		InstancePushId: integration.ZendeskInstancePushId,
	}

	response, err := pushChannels(body, integration)
	if err != nil {
		return err
	}

	fmt.Println(response)

	return nil
}

func pushChannels(body postBody, integration models.Integration) (response *http.Response, err error) {
	endpoint := fmt.Sprintf("https://%s.zendesk.com/api/v2/any_channel/push.json", integration.ZendeskSubdomain)
	jsonPayload, err := json.Marshal(body)
    if err != nil {
		return nil, err
    }
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer([]byte(jsonPayload)))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", integration.ZendeskToken))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return resp, nil
}

# Zendesk-Slack POC

:warning: The code in this repo is for exploration purpose, not production ready and likely to have serious security flaws.

A proof-of-concept Slack app built in Go to send messages to/from Slack. Built on top of Slack's event API and Zendesk's Channel Framework

#### Install dependencies.

This project uses [glide](https://github.com/Masterminds/glide#install) to manage dependencies.

```
    glide install
```

##### Configure the app by editing `config/parameters.json`

```
    {
        "slack_verification_token": "see-app-credentialzes-in-slack",
        "slack_app_client_id": "see-app-credentials-in-slack",
        "slack_app_client_secret": "see-app-credentials-in-slack",
        "database_host": "localhost",
        "database_port": 3306,
        "database_username": "root",
        "database_password": null,
        "database_name": "slack_zendesk"
    }
```

##### Migrate the database schema by executing.

Create your database before hand.

```
    go run migrate.go
```

#### Start the app by executing:

```golang
    go run server.go --queues=zendesk
```

#### Start the goworker by executing

```golang
    go run worker.go --queues=zendesk
```

#### Package the zat app in the /zat directory

```
    zat package
```

1. Upload the package as a private app
2. Go to the `channels integration` section in Zendesk and configure the integration
![Installation Flow](https://i.imgur.com/PRvhL4K.gif)

#### What does it support

- Slack OAuth, initiated from Channels Framework app
- Add bot user to channel, then send message to create a new ticket in Zendesk. (Does not perform any kind of message filtering)
- Replies from Zendesk are automatically sent to Slack

#### Slack App Set-up

![OAuth Setup](https://i.imgur.com/4aZoW8o.png)
![Bot Events](https://i.imgur.com/G15GGUo.png)
![Event Subscription](https://i.imgur.com/ecy5Tob.png)
{
	"blocks": [
		{
			"type": "section",
			"text": {
				"type": "mrkdwn",
				"text": "Hey there! 👋 I'm the *Artifact Publisher Bot*. I'm here to help you create and manage artifact deployments using Slack."
			}
		},
		{
			"type": "header",
			"text": {
				"type": "plain_text",
				"text": "Workflow",
				"emoji": true
			}
		},
		{
			"type": "divider"
		},
		{
			"type": "section",
			"text": {
				"type": "mrkdwn",
				"text": "💡 When you merge new code into the `development` branch, an artifact with an updated tag will be created and deployed into the `dev` environment. I will inform you about it."
			}
		},
		{
			"type": "section",
			"text": {
				"type": "mrkdwn",
				"text": "To control the bot, you can use the following commands:"
			}
		},
		{
			"type": "section",
			"text": {
				"type": "mrkdwn",
				"text": "1️⃣ Utilize the *`/artifact promote ...`* command. Enter `/artifact promote <artifact> [qa|staging|prod]`. This command facilitates the deployment of the artifact in the designated environment. In case you don't specify a preferred environment, the artifact will be deployed sequentially to the next environment in the order specified."
			}
		},
		{
			"type": "section",
			"text": {
				"type": "mrkdwn",
				"text": "2️⃣ The *`/artifact rollback <qa|staging|prod>`* command will assist you in reverting to the previous version of the artifact that was installed in the specified environment."
			}
		},
		{
			"type": "section",
			"text": {
				"type": "mrkdwn",
				"text": "3️⃣ The *`/artifact list`* command provides a list of artifacts deployed in the environments `dev`, `qa`, `staging`, `prod`."
			}
		},
		{
			"type": "section",
			"text": {
				"type": "mrkdwn",
				"text": "4️⃣ The *`/artifact diff <src_env> <dst_env>`* command will show the differences between artifacts in the specified environments."
			}
		},
		{
			"type": "section",
			"text": {
				"type": "mrkdwn",
				"text": "➕ To start tracking artifact deployment tasks, *add me to a channel* and I'll introduce myself. I'm usually added to a team or project channel. Type `/invite @ArtifactBot` from the channel or pick a channel on the right."
			},
			"accessory": {
				"type": "conversations_select",
				"placeholder": {
					"type": "plain_text",
					"text": "Select a channel...",
					"emoji": true
				}
			}
		},
		{
			"type": "divider"
		},
		{
			"type": "context",
			"elements": [
				{
					"type": "mrkdwn",
					"text": "❓Get help at any time with `/artifact help` or type *help* in a DM with me"
				}
			]
		}
	]
}
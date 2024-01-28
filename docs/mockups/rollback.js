{
	"blocks": [
		{
			"type": "section",
			"text": {
				"type": "mrkdwn",
				"text": "*Artifact rollback:*\n\n*<fakeLink.toEmployeeProfile.com|app_name:1.0.2>*"
			}
		},
		{
			"type": "section",
			"fields": [
				{
					"type": "mrkdwn",
					"text": "*Environment:*\nProduction"
				},
				{
					"type": "mrkdwn",
					"text": "*Trigged by:*\n <fakeLink.toEmployeeProfile.com|@Username>"
				},
				{
					"type": "mrkdwn",
					"text": "*Timestamp:*\n2024-01-28 12:45:23 UTC"
				},
				{
					"type": "mrkdwn",
					"text": "*Reason:*\nRolback unstable artifact."
				},
				{
					"type": "mrkdwn",
					"text": "*State:*\nRollback to the previous\n artifact `app_name:1.0.1`"
				},
				{
					"type": "mrkdwn",
					"text": "*Command:*\n `/artifact rollback 1.0.2 prod`"
				}
			]
		}
	]
}
# Run command

Cron plugin for running Slack command

## Configuration

```yaml
command: "/who"
channel: "my-team-name"
token: "xoxp-xxx"
text: "xxx-xxxxxx"
```

* `command`: slash command to be executed. Leading backslash is required.
* `channel`: ID of the public channel to execute the command in
* `token`: legacy token of the user need to click on the button
* `text`: additional parameters provided to the slash command

# Use invite not at

For user who have the bad habit of invite user using message like "@\<username\>", the bot answers with a thread and a link to the documentation

## Configuration

```json
{
    "message": " /invite is your friend https://get.slack.help/hc/en-us/articles/201980108-Invite-members-to-a-channel",
    "reactions_good": ["thumbsup", "white_check_mark"],
    "reactions_bad": ["x", "no_good", "thumbsdown"]
}
```

* `message`: message send to the user (in a thread) who invite user with @username way instead of /invite @username
* `reactions_good`: reactions added to the response
* `reactions_bad`: reactions added to the user message

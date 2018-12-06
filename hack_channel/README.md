# Hack channel

Re-enable @channel / @everyone when is disable for simple user

## Configuration

```json
{
    "message": "Make %s great again!",
    "channel_keyword": ["@channel"], 
    "here_keyword": ["@everyone", "@here"], 
    "on_public": false,
    "enabled_channel": [""],
    "disabled_channel": [""]
}
```

* `message`: message send by the bot
* `channel_keyword`: list of words that triggering @channel
* `here_keyword`: list of words that triggering @everyone
* `on_public`: if bot should mention channel or everyone on public channel
* `enabled_channel`: list of channel id for this feature is enable
* `disabled_channel`: list of channel id for this feature is disable (override enabled)

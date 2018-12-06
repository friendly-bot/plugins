# Hack channel

Re-enable @channel / @everyone when is disable for simple user

## Configuration

```json
{
    "channel_keyword": "@channel", 
    "everyone_keyword": "@everyone", 
    "on_public": false,
    "enabled_channel": [""],
    "disabled_channel": [""]
}
```

* `channel_keyword`: word that triggering @channel
* `everyone_keyword`: word that triggering @everyone
* `on_public`: if bot should mention channel or everyone on public channel
* `enabled_channel`: list of channel id for this feature is enable
* `disabled_channel`: list of channel id for this feature is disable (override enabled)

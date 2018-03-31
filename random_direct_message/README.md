# Random direct message

Send random message to a random active user (who is not a bot or a guest) at specific time (schedule by cron style)

## Configuration

```json
{
    "probability": 95,
    "chance": 100,
    "messages": ["I like your style !", "You should be proud of you :clap:"],
    "talk_after": 48
}
```

* `probability`: used with `chance`, roll a dice if probability < [0:chance[ then run this cron
* `chance`: see `probability`
* `messages`: array of message they can be send to one user, reaction are available
* `talk_after`: number of hour to wait before re-send message to the same user 

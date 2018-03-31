# Auto attachment action

Auto click on the first action button in one attachment for someone.

## Configuration

```json
{
    "text_attachment": "Click first for win",
    "team": "my-team-name",
    "token": "xoxp-xxx"
}
```

* `text_attachment`: text contains in attachment for match (use strings.Contains)
* `team`: team name (can be found in your slack domain https://\<xxx\>.slack.com)
* `token`: legacy token of the user need to click on the button

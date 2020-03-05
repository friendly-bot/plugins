# Keyword reaction

Add reaction based on regex

## Configuration

```yaml
reactions:
  cookie: "(?i)(^| )cookies?($| )"
```

* `reactions`: is a map\[reaction\]regex, key is the reaction to add when the regex match with the message

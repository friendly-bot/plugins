# Keyword reaction

Add reaction based on keywords. Bot only match full word

## Configuration

```json
{
    "reactions": {
        "cookie": ["cookie", "cookies"]
    },
    "keywords": {
        "happy birthday": ["birthday", "clap", "tada"]
    }
}
```

* `reactions`: is a map\[reaction\][]keyword, key is the reaction to add when one keyword in array is found
* `keywords`: is a map\[keyword\][]reaction, if key is found, add all reactions in the array 

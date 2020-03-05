# Random reaction

Bot can add reaction on message based on probability. You can configure which reaction use and what probability this one can appears

## Configuration
```yaml
chance: 1000
reactions:
  heart": 20
  yellow_heart: 1
```

* `chance`: roll a dice, if reactions\[reaction\] < [0;chance[; add reaction
* `reactions`: is a map\[reaction\]probability, key is the reaction and value is used with `chance` 

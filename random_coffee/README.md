# Random Coffee

Create random groups for a coffee time

## Configuration

```yaml
header: "Random Coffee time!"
footer: ""
prefix: "•"
group_of: 2
max_number_of_group: 10
channel: ""
separator: ", "
```

* `header`: Message send to the channel, with the list of the groups (before list)
* `footer`: Message send to the channel, with the list of the groups (after list)
* `group_of`: Number of users in one group
* `max_number_of_group`: Number of groups to create each run
* `channel`: Channel to use for select users for the random coffee
* `separator`: String to use when the users selected for the coffee is displayed

# go-get-cli
CLI tooling for Go Get to make it a bit more enjoyable.

## usage 

### Fetch the index 
Fetch the latest index. This is recommended to be done upon installing the CLI tool as the initial fetch can take a bit. This can be called again whenever you want a fresh version of the index.

```
go-get-cli fetch
```

### notes and planning 

#### Features
- go-get-cli (open ui to select sub commands)
- go-get-cli search (search packages via overall index)
- go-get-cli discover (uses awesome-go parsed data to provide context based packages)
  - flag `-c` categories to sort via categories 
- go-get-cli remove (find packages in local dir and allows removal features)
- MORE?

### Roadmap:
- [ ] Parse and create index from index.go/index (can use goroutines to get each section by the time (need to go from now back to 2019 or sth))
- [ ] Parse and create context based index from awesome-go 
- [ ] Think about a good and fast way of storing these (likely just going to be writing to a local file)
      - [ ] TODO: How can we handle updates? maybe store last update time in write file and on load work our way back from there (ie. LastWrite=2024-01-01:00:00:000, and recursively work back until we reach there again.)
- [ ] Build out UI view for option selection 
      - [ ] TODO: Decide between bespoke UI libs or just attempting to make a modular one that can service all components.
- [ ] Search 
- [ ] Discover 
- [ ] Remove 

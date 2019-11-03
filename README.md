# Battleship
Battleship is a simple two player game written in Golang, where players take turns guessing locations of enemy ships and try to sink them.

Compiling is done with
```
go build
```

To run tests, do
```
go test
```

Games can be played with a human player, ai, or some combination of the two. By default, ai boards are shown. To hide them, call battleship with the flag --no-show-ai
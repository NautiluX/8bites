all:
	go build -o 8bites main.go

run: all
	./8bites

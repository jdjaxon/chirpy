BIN=out

run: build
	./$(BIN)

build:
	go build -o $(BIN)

clean:
	$(RM) $(BIN)

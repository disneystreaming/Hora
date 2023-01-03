NAME = hora
VERSION = 0.1.0
# Included as relative path because of compile step, if not Go will look in GOROOT
SRC_DIR = ./src
BIN_DIR = bin

# constructs bin directory for Go executable to reside in
$(BIN_DIR):
	mkdir $(BIN_DIR)

###########################
####### Go Targets ########
###########################

.PHONY: swag compile run test

# Builds out swagger documentation and outputs the specs to the docs/ directory
swag:
	go install github.com/swaggo/swag/cmd/swag@v1.8.8
	swag init -d $(SRC_DIR) -g main/server.go -o $(SRC_DIR)/main/docs

# Compiles build of the src/main Go application and outputs binary to the bin/ directory
# This is dependent on the swagger documentation and bin/ directory being present
compile: $(BIN_DIR) swag
	go build -ldflags "-X main.HoraVersion=$(VERSION)" -o $(BIN_DIR)/$(NAME) $(SRC_DIR)/main

# Runs the compiled version of the Go application located in the bin/ directory
run:
	./$(BIN_DIR)/$(NAME)

# Runs all go unit tests
test:
	go test ./...

###########################
## Miscellaneous Targets ##
###########################

.PHONY: clean

# Cleans up artifacts
clean:
	rm -rf $(BIN_DIR)
	rm -rf $(SRC_DIR)/main/docs
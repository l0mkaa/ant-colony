PROJECTNAME = Ant Colony
BIN = $(CURDIR)/bin

SOURCEGLFW = $(CURDIR)/cmd/glfw
BINGL= $(BIN)/glfw

all: help

## build-glfw
build-glfw:
	go build -v -o $(BINGL)/glfw $(SOURCEGLFW)/main.go
## run-glfw
run_glfw:
	go run $(SOURCEGLFW)/main.go

## clean
clean:
	rm -r $(BIN)

## help
.PHONY: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo
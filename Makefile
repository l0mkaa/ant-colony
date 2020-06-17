PROJECTNAME = Ant Colony
BIN = $(CURDIR)/bin

SOURCEGLFW = $(CURDIR)/cmd/glfw
BINGL= $(BIN)/glfw

SOURCEWA = $(CURDIR)/cmd/web-assembly
BINWA = $(BIN)/web-assembly

all: help

## build-glfw
build-glfw:
	go build -v -o $(BINGL)/glfw $(SOURCEGLFW)/main.go
## run-glfw
run-glfw:
	go run $(SOURCEGLFW)/main.go

## build_web-assembly
build-web-assembly:
	GOOS=js GOARCH=wasm go build -v -o $(BINWA)/main.wasm $(SOURCEWA)/main.go
	cp $(SOURCEWA)/index.html $(BINWA)
	cp $(SOURCEWA)/wasm_exec.js $(BINWA)
## run_web-assembly
run-web-assembly: build-web-assembly
	xdg-open http://localhost:8080/
	goexec 'http.ListenAndServe(":8080", http.FileServer(http.Dir("$(BINWA)")))'

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
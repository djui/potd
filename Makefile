all:
	@go get -d && go build -o "potd"

clean:
	@go clean -i

init:
	@echo "Run: $$ source gvp && gpm"

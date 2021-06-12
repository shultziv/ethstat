build_cmd:
	go build -ldflags "-s -w" -o $(CURDIR)/build/ethstatcmd $(CURDIR)/cmd/ethstatcmd/main.go
PROJECT=gotrash

CMD_PATH=~/bin


install:
	# GOTRASH_PATH should be an empty folder, serve as the data dir of gotrash
	go build -o ${CMD_PATH}/gotrash -ldflags "-X github.com/Troublor/go-trash/cmd.GOTRASH_PATH=${GOTRASH_PATH}"
	@echo gotrash command installed at ${CMD_PATH}/gotrash

test:
	go test -v ./...
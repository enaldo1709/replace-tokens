run:
	go run src/main.go $(arg1) $(arg2) $(arg3) $(arg4)

build:
	go build -o replacetokens src/main.go 

install:
	go build -o replacetokens src/main.go 
	chmod a+x replacetokens
	mv -f replacetokens ${GOBIN}
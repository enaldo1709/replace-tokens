run:
	go run src/main.go $(arg1) $(arg2) $(arg3) $(arg4)

build:
	go build -o replacetokens src/main.go
	chmod a+x replacetokens
	
build-windows:
	go build -o replacetokens.exe src/main.go 

install:
	go build -o replacetokens src/main.go 
	chmod a+x replacetokens
	mv -f replacetokens ${GOBIN}

install-to-system:
	go build -o replacetokens src/main.go 
	chmod a+x replacetokens
	sudo mv -f replacetokens /bin
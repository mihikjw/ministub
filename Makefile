APPLICATION_NAME = ministub

all:
	@$(MAKE) bootstrap build success || $(MAKE) failure

bootstrap: 
	sh scripts/bootstrap.sh

build:
	sh scripts/build.sh

install:
	cp ./bin/ministub /usr/local/bin

clean:
	sh scripts/clean.sh

test:
	sh scripts/test.sh

success:
	printf "\n\e[1;32mBuild Successful\e[0m\n"

failure:
	printf "\n\e[1;31mBuild Failure\e[0m\n"
	exit 1
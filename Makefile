.PHONY: build clean deploy gomodgen

build: gomodgen
	export GO111MODULE=on
	./gobuild.sh

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh

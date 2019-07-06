.PHONY: build clean deploy gomodgen

build: gomodgen
	./gobuild.sh

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh

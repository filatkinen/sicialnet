BIN := "./build/cliurl"
FILE = $(URL)


build-img:
	docker-compose -f docker-compose.yml build

run-img: build-img
	docker-compose -f docker-compose.yml up -d


up: run-img

down:
	docker-compose -f docker-compose.yml down \
		 --rmi local \
		--volumes \
		--remove-orphans \
		--timeout 5; \




.PHONY: up


.PHONY: migrate-up migrate-create consume-rss start-server

migrate-create:
	echo \# create migration name="$(name)"
	go run cmd/main.go migrate-create $(name)

migrate-up:
	go run cmd/main.go migrate-up

consume-rss:
	go run cmd/main.go consume-rss


start-server:
	go run cmd/main.go serve
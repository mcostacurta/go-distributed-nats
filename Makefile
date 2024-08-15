up:
	docker compose up

down:
	docker compose down

sub:
	go run consumer/main.go

pub:
	echo ${MSG}
	go run publisher/main.go -msg="${MSG}"
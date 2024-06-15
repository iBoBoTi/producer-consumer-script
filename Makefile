start-producer:
	@cd ./producer &&  go run producer.go
	
start-consumer:
	go run consumer/consumer.go

docker-up:
	docker-compose  --env-file .env up -d
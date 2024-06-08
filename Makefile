build-and-push:
	@swag init -g ./cmd/app/main.go
	@docker build -t maksud1/fitness .
	@docker push maksud1/fitness
	@ssh root@158.160.62.249 "docker rm -f back || 1"
	@ssh root@158.160.62.249 "docker pull maksud1/fitness"
	@ssh root@158.160.62.249 "HOST=158.160.62.249:8000 docker run -d -p 8000:8000 --name=back --rm maksud1/fitness"
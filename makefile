.PHONY=model/build
model/build:
	@echo "Build model"
	@cd ml_model && make model/all-steps

.PHONY=services/build
services/build:
	@echo "Build services"
	@docker compose build

.PHONY=services/up
services/up:
	@echo "Run services"
	@docker compose up -d
	@cd ml_facade && make db/migrations/up && bash threshold_setup.sh

.PHONY=services/stop
services/stop:
	@echo "Stop services"
	@docker compose stop

.PHONY=services/down
services/down:
	@echo "Remove services"
	@docker compose down

.PHONY=setup-all
setup-all: model/build services/build services/up

.PHONY=run/send-data

rabbitmq ?= true
requests ?= 10000

run/send-data:
	@echo "Send data"
	@cd scripts/send_data && go run . --rabbitmq=$(rabbitmq) --requests=$(requests)

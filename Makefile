.PHONY: docker-b docker-f docker-up docker-down docker-rb docker-rf run-b run-f

# d = docker
# Windows users must run docker desktop before running these commands
db:
	@docker compose up --build b

df:
	@docker compose up --build f

dup:
	@docker compose up --build

ddown:
	@docker compose down

drb:
	@docker compose build b

drf:
	@docker compose build f


b: #backend
	@go run -C backend ./cmd/server

f: #frontend
	@cd frontend && bun run dev

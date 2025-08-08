.PHONY: docker-b docker-f docker-up docker-down docker-rb docker-rf run

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
	@cd frontend && bun run net

testdb: #test database
	@go run -C backend ./cmd/test

run:
	@$(MAKE) -j2 b f

stop:
	@xargs kill < .dev.pids 2>/dev/null || true
	@rm -f .dev.pids

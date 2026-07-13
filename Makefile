BINARY := $(notdir $(CURDIR))
APP := $(notdir $(CURDIR))
# Ports used by the dev servers (frontend, backend, and PocketBase-style API)
PORTS := 3000 3001


init:
	fastmod --hidden myapp $(APP) --glob '!Makefile'
	find . -depth \( -type f -o -type d \) -name '*myapp*' | while read -r f; do \
		mv -- "$$f" "$$(dirname "$$f")/$$(basename "$$f" | sed 's/myapp/$(APP))/g')"; \
	done

.PHONY: frontend-deps
frontend-deps:
	cd frontend && pnpm install

.PHONY: build-frontend
build-frontend: frontend-deps
	cd frontend && pnpm run build

.PHONY: build
build: build-frontend
	go build -o $(BINARY) ./cmd/$(BINARY)

.PHONY: kill-ports
kill-ports:
	@for port in $(PORTS); do \
		pid=$$(lsof -ti tcp:$$port); \
		if [ -n "$$pid" ]; then \
			echo "Killing process on port $$port (pid $$pid)"; \
			kill -9 $$pid; \
		fi \
	done


.PHONY: server
server: kill-ports
	#./kithara migrate up --dir=pb_data
	./$(BINARY) superuser upsert admin@mail.internal password --dir=pb_data
	./$(BINARY) serve

# --------------
.PHONY: clean
	rm -fr ./tmp/ # air

# port: 3001
.PHONY: dev-front
dev-front: clean
	npx concurrently -n "frontend,backend" -c "blue,green" "cd frontend && pnpm dev" "./$(BINARY) serve"

# port: 3000
.PHONY: dev-back
dev-back: clean
	npx concurrently -n "frontend,backend" -c "blue,green" "cd frontend && pnpm watch" "air"


.PHONY: test
test:
	cd frontend && pnpm test
	go test ./...


format:
	cd frontend && pnpm exec prettier --write "src/**/*.{js,jsx,css}"

migrate-collections:
	ls -1 migrations/*.go | sort | head -n -1 | xargs rm -f
	yes | go run ./cmd/kithara migrate collections # 開発初期限定
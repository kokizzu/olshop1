
setup:
	go get -u -v github.com/kokizzu/gotro@latest
	go install github.com/air-verse/air@latest
	go install github.com/fatih/gomodifytags@latest
	go install github.com/kokizzu/replacer@latest
	go install github.com/akbarfa49/farify@latest
	go install golang.org/x/tools/cmd/goimports@latest
	#curl -fsSL https://get.pnpm.io/install.sh | bash -
	curl -fsSL https://bun.sh/install | bash
	cd svelte ; pnpm i
	# alt: task/Taskfile

local-postgres:
	docker exec -it olshop1-cockroach1-1 cockroach sql --insecure

modtidy:
	sudo chmod -R a+rwx cockroach1 && go mod tidy

views:
	# generate views and routes
	./gen-views.sh

startdeps:
	docker compose up

startbe:
	air

startfe:
	cd svelte; npm run dev
task1:
	go fmt ./...

task2:
	go vet ./...

task3:
	golangci-lint run


bench_handler:
	go test -run ^$$ -bench=. ./internal/handlers/ -count 2 -benchmem
task5:
	go test -cover ./...

task6:
	errcheck ./...

task7:
	staticcheck ./...

task9:
	go mod tidy
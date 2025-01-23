.PHONY: swagger

swagger:
	docker exec -it api sh -c "swag init -d cmd --parseDependency --parseInternal --parseDepth 4 -o swagger/docs"

test_coverage:
	-go test -v ./... -coverprofile=coverage.out
	go tool cover -html coverage.out -o cover.html
	@open cover.html
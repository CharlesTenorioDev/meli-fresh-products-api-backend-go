.PHONY: swagger

swagger:
	docker exec -it api sh -c "swag init -d cmd --parseDependency --parseInternal --parseDepth 4 -o swagger/docs"

test_coverage:
	-go test -v ./... -coverprofile=coverage.out
	go tool cover -html coverage.out -o cover.html
	@open cover.html

coverage_buyer_service_test:
	-go test -v ./internal/service -run TestBuyer -coverprofile=coverage.out
	@cat coverage.out | (head -n 1 coverage.out && grep "github.com/meli-fresh-products-api-backend-t1/internal/service/buyer" coverage.out) > buyer_coverage.out
	@rm coverage.out
	go tool cover -html buyer_coverage.out -o cover.html
	@open cover.html

coverage_buyer_handler_test:
	-go test -v ./internal/handler -run "^TestHandler.*UnitTest$$" -coverprofile=coverage.out
	@cat coverage.out | (head -n 1 coverage.out && grep "github.com/meli-fresh-products-api-backend-t1/internal/handler/buyer" coverage.out) > buyer_coverage.out
	@rm coverage.out
	go tool cover -html buyer_coverage.out -o cover.html
	@open cover.html

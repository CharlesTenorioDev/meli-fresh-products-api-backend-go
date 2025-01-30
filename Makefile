.PHONY: swagger test_coverage

swagger:
	docker exec -it api sh -c "swag init -d cmd --parseDependency --parseInternal --parseDepth 4 -o swagger/docs"

test_coverage:
	-go test ./... -coverprofile=coverage.out
	go tool cover -html coverage.out -o cover.html
	@open cover.html

unit_test_employees_handler:
	@# Excludes lines that end with [no tests to run], [no test files] and no tests to run
	go test -v ./... -run "^TestHandler.*EmployeeUnitTest$$" | grep -v -E "\[no tests to run\]$$|\[no test files\]$$|no tests to run$$" 

unit_test_employees_svc:
	@# Excludes lines that end with [no tests to run], [no test files] and no tests to run
	go test -v ./... -run "^Test.*EmployeeUnitTestService$$" | grep -v -E "\[no tests to run\]$$|\[no test files\]$$|no tests to run$$"

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

coverage_employee_unit_test:
	-go test -v ./... -run "EmployeeUnitTest$$" -coverprofile=coverage.out
	@cat coverage.out | (head -n 1 coverage.out && grep -E "^(github.com/meli-fresh-products-api-backend-t1/internal/service/employee.go|github.com/meli-fresh-products-api-backend-t1/internal/handler/employee)" coverage.out) > employee_coverage.out
	@rm coverage.out
	go tool cover -html employee_coverage.out -o cover.html
	@open cover.html

coverage_section_unit_test:
	-go test -v ./... -run "SectionUnitTest$$" -coverprofile=coverage.out
	@cat coverage.out | (head -n 1 coverage.out && grep -E "^(github.com/meli-fresh-products-api-backend-t1/internal/service/section.go|github.com/meli-fresh-products-api-backend-t1/internal/handler/section)" coverage.out) > section_coverage.out
	@rm coverage.out
	go tool cover -html section_coverage.out -o cover.html
	@open cover.html

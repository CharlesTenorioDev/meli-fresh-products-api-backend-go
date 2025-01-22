unit_test_employees:
	@# Excludes lines that end with [no tests to run], [no test files] and no tests to run
	go test -v ./... -run "^TestHandler.*EmployeeUnitTest$$" | grep -v -E "\[no tests to run\]$$|\[no test files\]$$|no tests to run$$" 

test_coverage:
	-go test -v ./... -coverprofile=coverage.out
	go tool cover -html coverage.out -o cover.html
	@open cover.html

swagger:
	docker exec -it api sh -c "swag init -d cmd --parseDependency --parseInternal --parseDepth 4 -o swagger/docs"

.PHONY: swagger test_coverage

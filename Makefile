unit_test_employees:
	@# Excludes lines that end with [no tests to run], [no test files] and no tests to run
	go test -v ./... -run "^TestHandler.*EmployeeUnitTest$$" | grep -v -E "\[no tests to run\]$$|\[no test files\]$$|no tests to run$$" 

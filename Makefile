.PHONY: test coverage coverage-html coverage-report coverage-full clean-coverage

test:
	@echo "Running all tests..."
	@go test -v ./tests/...

coverage:
	@echo "Running all tests with coverage..."
	@mkdir -p ./tests/coverage
	@go test -v -coverprofile=./tests/coverage/coverage.out -covermode=atomic -coverpkg=./... ./tests/...
	@go tool cover -func=./tests/coverage/coverage.out

coverage-html:
	@echo "Generating HTML coverage report..."
	@mkdir -p ./tests/coverage
	@go test -v -coverprofile=./tests/coverage/coverage.out -covermode=atomic -coverpkg=./... ./tests/...
	@go tool cover -html=./tests/coverage/coverage.out -o ./tests/coverage/coverage.html
	@echo "Coverage report generated: ./tests/coverage/coverage.html"

coverage-report: coverage-html
	@echo ""
	@echo "=== Coverage Summary ==="
	@go tool cover -func=./tests/coverage/coverage.out | tail -1

coverage-full:
	@echo "Running tests and analyzing complete codebase coverage..."
	@mkdir -p ./tests/coverage
	@go test -v -coverprofile=./tests/coverage/coverage.out -covermode=atomic -coverpkg=./... ./tests/... > /dev/null 2>&1 || true
	@echo ""
	@echo "=== Files with Test Coverage ==="
	@go tool cover -func=./tests/coverage/coverage.out
	@echo ""
	@echo "=== Files with NO Test Coverage (0%) ==="
	@find . -name "*.go" -not -path "./tests/*" -not -path "./.git/*" -not -name "*_test.go" | while read file; do \
		pkg=$$(echo $$file | sed 's|^\./||' | sed 's|\.go$$||'); \
		if ! grep -q "$$pkg" ./tests/coverage/coverage.out 2>/dev/null; then \
			echo "  $$file"; \
		fi; \
	done
	@echo ""
	@echo "Note: The coverage percentage above only includes files that were executed by tests."
	@echo "Files listed under 'NO Test Coverage' are not included in the percentage calculation."

clean-coverage:
	@echo "Cleaning coverage files..."
	@rm -rf ./tests/coverage
	@echo "Coverage files removed"

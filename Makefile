export GOOS := linux
export GOARCH := amd64
export CGO_ENABLED := 0

# Creates the 'bootstrap' file (binary of the code)
build:
	go build -tags lambda.norpc -o bootstrap main.go

# Deploys to AWS (same as npm run deploy:dev)
deploy_dev:
	serverless deploy --stage dev

# Deletes the deployment (same as npm run delete:lambda)
delete_dev:
	serverless remove --stage dev

# Deletes the current 'bootstrap' file (compiled), removes the .serverless folder, and creates a new 'bootstrap' binary
update_lambda:
	powershell -Command "& { If (Test-Path 'bootstrap') { Remove-Item 'bootstrap' }; If (Test-Path '.serverless') { Remove-Item '.serverless' -Recurse }; If (Test-Path 'tmp') { Remove-Item 'tmp' -Recurse }; make build }"

# 'Should' run locally, DOES NOT WORK
start:
	powershell "$$env:ENVIRONMENT = 'dev'; air"

# Creates the .serverless folder, but without deploying
package:
	serverless package --stage dev

.PHONY: build deploy_dev delete_dev update_lambda run_lambda package

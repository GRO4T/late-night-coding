# QUIZ
Implementation of quiz exercise from https://github.com/gophercises/quiz/tree/master
## Build and run
```bash
cd cmd/quiz
go build -o ../../bin/quiz
cd ../../
bin/quiz -f problems.csv -t 30
```
## Tests
```
go test ./...
```
With coverage
```
# run tests and create a coverprofile
go test ./... -coverprofile=cover.out
# open the interactive UI to check the Coverage Repor
go tool cover -html=cover.out
```
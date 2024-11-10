dev:
	./bin/wgo -file=.html -file=.go -file=.css -xfile=./static/css/output.css ./bin/tailwindcss-linux-x64 -i static/css/input.css -o static/css/styles.css :: go run cmd/gofind/main.go 

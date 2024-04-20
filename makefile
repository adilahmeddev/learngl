build: 
	@go build -o ./out/learngl ./... 

run: build
	./out/learngl



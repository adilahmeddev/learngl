ifdef OS
	Outfile = learngl.exe
else
	Outfile = learngl
endif

build: 
	@go build -o ./out/${Outfile} ./... 

run: build
	@./out/${Outfile}



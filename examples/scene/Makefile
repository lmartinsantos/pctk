.PHONY: all
all: resources

.PHONY: resources
resources: resources.idx resources.dat

resources.idx resources.dat: 
	go run ../../cmd/pctk pack resources

.PHONY: run
run: resources
	go run main.go

.PHONY: clean
clean:
	rm -f resources.idx resources.dat
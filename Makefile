bump:
	bump_version patch filedb.go 

push: 
	git push origin --tags

release: test bump push
	git push --tags

test:
	go test ./... -cover -bench=. -test.benchtime=3s;

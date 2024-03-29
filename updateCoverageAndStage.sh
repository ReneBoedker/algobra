#!/bin/bash

for subpkg in $(go list ./... | sed "s#github.com/ReneBoedker/algobra#.#g")
do
	## Run the coverage test
	go test -coverprofile=$subpkg/coverage.out -coverpkg $subpkg/... $subpkg/... > /dev/null
	cover=$(go tool cover -func=$subpkg/coverage.out | tail -n 1 | sed "s/^total:[^0-9]*\([0-9\.]\+\)%.*/\1/")

	echo $subpkg: $cover
	
	if [[ $(echo $cover | tail -n1) == FAIL* ]]
	then
		## Test or build failed; abort commit
		exit 1
	fi

	if [[ $cover == "?"* ]] || [ ! -f $subpkg/README.md ]
	then
		## If no tests are defined or README does not exist, skip
		continue
	fi

	## Determine badge colour based on the result
	if (( $(echo "${cover} > 90" | bc -l)))
	then
		colour="brightgreen"
	elif (( $(echo "${cover} > 80" | bc -l)))
	then
		colour="green"
	elif (( $(echo "${cover} > 70" | bc -l)))
	then
		colour="yellowgreen"
	elif (( $(echo "${cover} > 60" | bc -l)))
	then
		colour="yellow"
	elif (( $(echo "${cover} > 50" | bc -l)))
	then
		colour="orange"
	else
		colour="red"
	fi

	## Update badges in README files
	sed -i "s@\(!\[coverage-badge\]\)([^)]*)@\1(https://img.shields.io/badge/coverage-${cover}%25-${colour}?cacheSeconds=86400\&style=flat)@" ./$subpkg/README.md

	## Add the changed files to the commit
	git add ./$subpkg/README.md
done

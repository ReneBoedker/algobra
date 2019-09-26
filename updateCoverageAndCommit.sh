#!/bin/bash

for subpkg in basic bivariate errors primefield univariate
do
	## Run the coverage test
	cover=$(go test -cover algobra/$subpkg | sed "s/.*coverage: \([0-9\.]\+\)%.*/\1/")

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

	## Add the files to the commit
	git add ./$subpkg/README.md
done

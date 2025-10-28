#!/bin/bash
repository=
descriptionsDirectory=
i=0
start=4090
cd ./descriptions
rm -r generated
mkdir generated
while [ $i -lt $1 ]
do
	sum=$[$i+$start]
	python3 ./materialize_descriptions_temperature.py $sum
i=$[$i+1]
done
i=0
cd ../
cd ./vsensor-temperature/
while [ $i -lt $1 ]
do
	sum=$[$i+$start]
	go run . --server.port=$sum --server.repository=$repository --server.description=$descriptionsDirectory &
	echo "running on "$sum
i=$[$i+1]
done

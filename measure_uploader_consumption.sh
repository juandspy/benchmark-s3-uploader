#!/bin/bash

RESULTS_CSV_PATH="results.csv"

if test -f "$RESULTS_CSV_PATH"; then
    echo "$RESULTS_CSV_PATH exists."
else
    echo "$RESULTS_CSV_PATH doesn't exist. Creating stucture. . ."
    echo "total_csvs,memory_consumption" > $RESULTS_CSV_PATH
fi

total_csvs=0
tmp_output=$(mktemp)

function measure_usage() {
    echo "total csvs is $total_csvs"
    export TOTAL_CSVS="$total_csvs"
    go build .
    /usr/bin/time -f "%P %M" "./m" 2> "$tmp_output"

    memory_consumption=$(cat $tmp_output | tail -n 1 | awk 'NF>1{print $NF}')
    echo "memory consumption is $memory_consumption KB"

    echo "$total_csvs,$memory_consumption" >> $RESULTS_CSV_PATH
}

for total_csvs in {0..100..10}
do
    measure_usage
done


# Pretty print the CSV
cat $RESULTS_CSV_PATH | column -t -s,

python ploter.py

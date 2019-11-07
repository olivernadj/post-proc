#!/usr/bin/env bash

################################################################
# Script name      : fakeload.sh                               #
# Description      : Makes randomised apache bench requests    #
# Original Author  : Oliver Nadj <mr.oliver.nadj@gmail.com>    #
################################################################


# validate the arguments
if [[ -z "$1" ]]; then HOST="http://localhost:8080/"; else HOST=$1; fi

generateActions () {
    if [[ -z "$3" ]]; then SLEEP_RAND=5; else SLEEP_RAND=$3; fi
    (
        while true   # Endless loop.
        do
            sleep $(( RANDOM % SLEEP_RAND + 1 ))
            echo "-X POST $1 -H \"Source-Type: $2\" -H \"Content-Type: application/json\" -d \"{ \"action\": \"action\", \"state\": \"new\"}\""
            curl -X POST "$1" -H "Source-Type: $2" -H "Content-Type: application/json" -d "{ \"action\": \"action\", \"state\": \"new\"}"
        done
    ) &
}


generateActions "$HOST""v1/action" "server" 3
generateActions "$HOST""v1/action" "client" 2
generateActions "$HOST""v1/action" "payment" 5

trap ctrl_c INT

function ctrl_c() {
        echo "Bye"
        pkill -P $$
        exit 1
}

while [[ 1 ]]
do
    sleep 1
done


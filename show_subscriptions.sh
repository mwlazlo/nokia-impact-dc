#!/bin/bash
# Register application with Nokia IMPACT and setup subscriptions

set -e

source config.sh

CURLBIN=$(which curl)
PARAMS=(--header "Authorization: Basic ${IMPACT_AUTH}")
PARAMS+=(-sS)
PARAMS+=(--header "Content-Type: application/json" --header "Accept: application/json")

[ -z "${CURLBIN}" ] && {
    echo "Could not find curl binary"
    exit 1
}

curl_cmd() {
    tmp=$(mktemp)
    ${CURLBIN} "${PARAMS[@]}" "$@" > $tmp
    cat $tmp | python -m json.tool
    if [ $? ] ; then
        cat $tmp
        rm $tmp
        exit 1
    fi
    rm $tmp
}

extract_subscription_id() {
    tee /dev/stderr | 
        grep subscriptionId | 
        sed 's/.*subscriptionId": *"//; s/".*//;'
}

list_subscription_ids() {
    for type in resources lifecycleEvents; do
        curl_cmd https://impact.idc.nokia.com/m2m/subscriptions?type=$type |
             extract_subscription_id
    done
}

for type in resources lifecycleEvents; do
    echo "Current subscriptions: ${type}"
    curl_cmd https://impact.idc.nokia.com/m2m/subscriptions?type=${type}
done

# vim: set ts=4 sw=4 et

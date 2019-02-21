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
    ${CURLBIN} "${PARAMS[@]}" "$@" | python -m json.tool
}

json_app_register() {
    cat <<EOF
{
    "headers": {
        "authorization": "Basic ${UBIIK_AUTH}"
    },
    "url": "${UBIIK_BASE_URL}/callback"
}
EOF
}

json_lifecycle_subscription_request() {
    cat <<EOF
{
 'events': ['deregistration','update','registration','expiration'],
 'deletionPolicy': 0,
 'groupName': '${IMPACT_GROUP}',
 'subscriptionType': 'lifecycleEvents'
}
EOF
}

json_resource_subscription_request() {
    cat <<EOF
{
 'deletionPolicy': 0,
 'groupName': '${IMPACT_GROUP}',
 'subscriptionType': 'resources',
 'resources': [
EOF
    first=1
    for r in $(cat resources.txt); do
        if [ $first = 1 ]; then
            first=0
        else 
            echo -n ,
        fi
        echo "{'resourcePath': '$r'}"
    done
    cat <<EOF
 ]
}
EOF
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

delete_subscription() {
    echo "Deleting subscription $1" 
    curl_cmd -X DELETE https://impact.idc.nokia.com/m2m/subscriptions/$1
}


echo "Clear existing subscriptions"
for sid in $(list_subscription_ids); do 
    delete_subscription $sid
done

echo "Register application URL"
curl_cmd -X PUT -d "$(json_app_register)" \
    https://impact.idc.nokia.com/m2m/applications/registration


echo "Register for lifecycle events"
lc_subsId=$(curl_cmd -X POST -d "$(json_lifecycle_subscription_request)" \
    https://impact.idc.nokia.com/m2m/subscriptions?type=lifecycleEvents |
        extract_subscription_id)

echo "Updating web service with current subscription ID ${lc_subsId}"
./ubiik-backend.sh setLifecycleSubscriptionId "${lc_subsId}"

echo "Register for resource events"
rs_subsId=$(curl_cmd -X POST -d "$(json_resource_subscription_request)" \
    https://impact.idc.nokia.com/m2m/subscriptions?type=resources |
        extract_subscription_id)
./ubiik-backend.sh setResourceSubscriptionId "${rs_subsId}"

for type in resources lifecycleEvents; do
    echo "Current subscriptions: ${type}"
    curl_cmd https://impact.idc.nokia.com/m2m/subscriptions?type=${type}
done

# vim: set ts=4 sw=4 et

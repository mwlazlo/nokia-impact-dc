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

curl() {
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

echo "Register application URL"
curl -X PUT -d "$(json_app_register)" https://impact.idc.nokia.com/m2m/applications/registration

echo "Clear existing subscriptions"
curl -X DELETE https://impact.idc.nokia.com/m2m/subscriptions

echo "Validate existing subscriptions"
curl https://impact.idc.nokia.com/m2m/subscriptions?type=lifecycleEvents    

echo "Register for lifecycle events"
subscriptionIdBuf=$(mktemp)
trap "rm -f ${subscriptionIdBuf}" EXIT
curl -X POST -d "$(json_lifecycle_subscription_request)" https://impact.idc.nokia.com/m2m/subscriptions?type=lifecycleEvents > ${subscriptionIdBuf}
echo MDFLKMSDFLKM
cat ${subscriptionIdBuf}
subscriptionId=$(cat ${subscriptionIdBuf} | python -c "import sys, json; print json.load(sys.stdin)['subscriptionId']")

echo "Updating web service with current subscription ID $(cat ${subscriptionIdBuf})"
./ubiik-backend.sh setLifecycleSubscriptionId "${subscriptionId}"

for type in resources lifecycleEvents; do
    echo "Current subscriptions: ${type}"
    curl https://impact.idc.nokia.com/m2m/subscriptions?type=${type}
done

# vim: set ts=4 sw=4 et

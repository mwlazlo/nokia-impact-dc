#!/bin/bash
# Register application with Nokia IMPACT and setup subscriptions

set -e

source config.sh

CURLBIN=$(which curl)
PARAMS=(--header "Authorization: Basic ${UBIIK_AUTH}")
PARAMS+=(-sS)
PARAMS+=(--header "Content-Type: application/json" --header "Accept: application/json")

[ -z "${CURLBIN}" ] && {
    echo "Could not find curl binary"
    exit 1
}

curl() {
    ${CURLBIN} "${PARAMS[@]}" "$@" 
}

json_update_subscription() {
    id=$1
    cat <<EOF
{
    "id": "$id"
}
EOF
}

die() {
    echo $*
    exit 1
}

set_lifecycle_id() {
    id=$1
    [ -z "${id}" ] && die "This option requires an additional argument"
    curl -X POST -d "$(json_update_subscription $id)" ${UBIIK_BASE_URL}/setLifecycleSubscription
}

set_resource_id() {
    id=$1
    [ -z "${id}" ] && die "This option requires an additional argument"
    curl -X POST -d "$(json_update_subscription $id)" ${UBIIK_BASE_URL}/setResourceSubscription
}

case $1 in 
    setLifecycleSubscriptionId)
        set_lifecycle_id $2
    ;;
    setResourceSubscriptionId)
        set_resource_id $2
    ;;

    *)
    echo "Usage: $0 setLifecycleSubscriptionId | setResourceSubscriptionId"
    exit 1
    ;;
esac

# vim: set ts=4 sw=4 et

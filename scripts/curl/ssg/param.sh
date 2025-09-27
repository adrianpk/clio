#!/bin/bash

# Description: Tests the Param API endpoints (CRUD cycle).
# Usage: ./param.sh

# --- Configuration ---
source "$(dirname "$0")/_config.sh"
RESOURCE="params"
RANDOM_SUFFIX=$((RANDOM % 10000))
PARAM_NAME="test.param.$RANDOM_SUFFIX"
PARAM_REF_KEY="ssg.test.param.$RANDOM_SUFFIX"

# --- Main Functions ---
create_param() {
    echo "--- 1. POST /$RESOURCE (Create New Param) ---" >&2 # Write to stderr
    PAYLOAD=$(cat <<EOF
{
    "name": "$PARAM_NAME",
    "description": "A test parameter created via API.",
    "value": "initial_value",
    "ref_key": "$PARAM_REF_KEY"
}
EOF
)
    response=$(curl -s -X POST "$BASE_URL/$RESOURCE" -H "Content-Type: application/json" -d "$PAYLOAD")
    echo "$response"
}

get_param() {
    local id=$1
    echo "--- 2. GET /$RESOURCE/{id} (Verify Creation) ---"
    curl -s -X GET "$BASE_URL/$RESOURCE/$id" | jq .
}

get_param_by_name() {
    local name=$1
    echo "--- 3. GET /$RESOURCE/name/{name} (Get by Name) ---"
    curl -s -X GET "$BASE_URL/$RESOURCE/name/$name" | jq .
}

get_param_by_ref_key() {
    local ref_key=$1
    echo "--- 4. GET /$RESOURCE/refkey/{ref_key} (Get by RefKey) ---"
    curl -s -X GET "$BASE_URL/$RESOURCE/refkey/$ref_key" | jq .
}

list_params() {
    echo "--- 5. GET /$RESOURCE (List Params) ---"
    curl -s -X GET "$BASE_URL/$RESOURCE" | jq .
}

update_param() {
    local id=$1
    echo "--- 6. PUT /$RESOURCE/{id} (Update Param) ---"
    PAYLOAD=$(cat <<EOF
{
    "name": "$PARAM_NAME",
    "description": "An updated test parameter.",
    "value": "updated_value",
    "ref_key": "$PARAM_REF_KEY"
}
EOF
)
    curl -s -X PUT "$BASE_URL/$RESOURCE/$id" -H "Content-Type: application/json" -d "$PAYLOAD" | jq .
}

delete_param() {
    local id=$1
    echo "--- 7. DELETE /$RESOURCE/{id} (Delete Param) ---"
    curl -s -X DELETE "$BASE_URL/$RESOURCE/$id" | jq .
}

# --- Main Execution ---

echo "--- Running Full Param CRUD Test Cycle ---"

response_data=$(create_param)
PARAM_ID=$(echo "$response_data" | jq -r '.data.param.id')

if [ -z "$PARAM_ID" ] || [ "$PARAM_ID" == "null" ]; then
    echo "Failed to create param or capture ID. Aborting."
    exit 1
fi
echo "Captured Param ID: $PARAM_ID"
sleep 1

get_param "$PARAM_ID"
sleep 1

get_param_by_name "$PARAM_NAME"
sleep 1

get_param_by_ref_key "$PARAM_REF_KEY"
sleep 1

list_params
sleep 1

update_param "$PARAM_ID"
sleep 1

delete_param "$PARAM_ID"
echo "--- Param CRUD Test Cycle Finished ---"

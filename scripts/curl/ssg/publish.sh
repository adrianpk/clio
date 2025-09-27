#!/bin/bash

echo "Publishing site..."

# Optional: Pass a commit message in the request body
# COMMIT_MESSAGE="Site update from Clio CLI"

API_URL="http://localhost:8081/api/v1/ssg/publish"
echo "Calling API: ${API_URL}"

if [ -n "$COMMIT_MESSAGE" ]; then
  echo "With commit message: ${COMMIT_MESSAGE}"
  curl -X POST "${API_URL}" \
    -H "Content-Type: application/json" \
    -d "{\"commit_message\": \"${COMMIT_MESSAGE}\"}"
else
  echo "Without specific commit message."
  curl -X POST "${API_URL}" \
    -H "Content-Type: application/json" \
    -d "{}"
fi

echo "" # Newline for better readability
echo "Site publish request sent. Check above for API response."
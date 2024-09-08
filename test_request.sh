#!/bin/bash

# 1. Valid Request with Unique ID
echo "Test 1: Valid Request with Unique ID"
curl -X GET "http://localhost:8080/api/verve/accept?id=1"
echo -e "\n"

# 2. Valid Request with External GET Endpoint
echo "Test 2: Valid Request with External GET Endpoint"
curl -X GET "http://localhost:8080/api/verve/accept?id=2&endpoint=http://httpbin.org/get"
echo -e "\n"

# 3. Invalid Request (Missing ID)
echo "Test 3: Invalid Request (Missing ID)"
curl -X GET "http://localhost:8080/api/verve/accept"
echo -e "\n"

# 4. Duplicate Request (Same ID Sent Twice)
echo "Test 4: Duplicate Request (First Request)"
curl -X GET "http://localhost:8080/api/verve/accept?id=3"
echo -e "\n"

echo "Test 4: Duplicate Request (Second Request)"
curl -X GET "http://localhost:8080/api/verve/accept?id=3"
echo -e "\n"

# 5. Valid Request with External POST Endpoint (Extension 1)
echo "Test 5: Valid Request with External POST Endpoint"
curl -X POST -H "Content-Type: application/json" -d '{"count":10}' "http://localhost:8080/api/verve/accept?id=4&endpoint=http://httpbin.org/post"
echo -e "\n"

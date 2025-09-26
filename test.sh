#!/bin/bash

# Test script for the Git server
echo "Testing Git Server..."
echo "===================="

# Test 1: List repositories
echo "1. Listing repositories:"
curl -s http://localhost:8080/repos | jq .
echo ""

# Test 2: Create a new repository
echo "2. Creating a new repository:"
curl -X POST "http://localhost:8080/repos?name=test-$(date +%s)"
echo ""
echo ""

# Test 3: List repositories again
echo "3. Listing repositories after creation:"
curl -s http://localhost:8080/repos | jq .
echo ""

# Test 4: Test Git protocol endpoint
echo "4. Testing Git protocol endpoint:"
curl -s "http://localhost:8080/test-repo/info/refs?service=git-upload-pack"
echo ""
echo ""

# Test 5: Add a file to repository
echo "5. Adding a file to repository:"
curl -X POST "http://localhost:8080/test-repo/files" \
  -d "filename=test.txt" \
  -d "content=Hello, Git Server!"
echo ""
echo ""

echo "All tests completed!"
echo "Visit http://localhost:8080 to see the web interface"

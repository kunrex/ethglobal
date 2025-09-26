#!/bin/bash

echo "🚀 Lighthouse (Filecoin) API Test"
echo "=================================="
echo ""

# Check if server is running
if ! curl -s http://localhost:8080/repos > /dev/null; then
    echo "❌ Server is not running. Please start it with: ./git-server"
    exit 1
fi

echo "✅ Server is running on http://localhost:8080"
echo ""

# Check if Lighthouse API key is set
if [ -z "$LIGHTHOUSE_API_KEY" ]; then
    echo "⚠️  LIGHTHOUSE_API_KEY environment variable not set."
    echo "   Set it with: export LIGHTHOUSE_API_KEY=your_api_key_here"
    echo "   Get your API key from: https://files.lighthouse.storage/"
    echo ""
    echo "   Testing without API key (endpoints will show warnings)..."
    echo ""
fi

# Test 1: Get Lighthouse help
echo "📋 1. Getting Lighthouse API help:"
curl -s http://localhost:8080/lighthouse/help | jq .
echo ""

# Test 2: Upload text content
echo "📝 2. Uploading text content to Filecoin:"
curl -X POST -H "Content-Type: application/json" \
  -d '{"content":"Hello from Git Server! This is a test upload to Filecoin network.","filename":"test.txt"}' \
  http://localhost:8080/lighthouse/upload-text
echo ""
echo ""

# Test 3: Create a test file and upload it
echo "📁 3. Creating and uploading a test file:"
echo "This is a test file for Lighthouse upload." > test-file.txt
curl -X POST -F 'file=@test-file.txt' http://localhost:8080/lighthouse/upload
echo ""
echo ""

# Test 4: List recent uploads
echo "📋 4. Listing recent uploads:"
curl -s http://localhost:8080/lighthouse/uploads | jq .
echo ""

# Clean up test file
rm -f test-file.txt

echo "🎉 Lighthouse API tests completed!"
echo ""
echo "Note: If you see errors, make sure to:"
echo "  1. Set LIGHTHOUSE_API_KEY environment variable"
echo "  2. Get your API key from https://files.lighthouse.storage/"
echo "  3. Ensure you have sufficient credits in your Lighthouse account"
echo ""
echo "Available Lighthouse endpoints:"
echo "  • Upload file:     curl -X POST -F 'file=@example.txt' http://localhost:8080/lighthouse/upload"
echo "  • Upload text:     curl -X POST -H 'Content-Type: application/json' -d '{\"content\":\"Hello\",\"filename\":\"hello.txt\"}' http://localhost:8080/lighthouse/upload-text"
echo "  • Download file:   curl 'http://localhost:8080/lighthouse/download?cid=QmXXX...' -o downloaded_file"
echo "  • Get file info:   curl 'http://localhost:8080/lighthouse/file-info?cid=QmXXX...'"
echo "  • Get help:        curl http://localhost:8080/lighthouse/help"

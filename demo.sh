#!/bin/bash

echo "🚀 Git Server Demo"
echo "=================="
echo ""

# Check if server is running
if ! curl -s http://localhost:8080/repos > /dev/null; then
    echo "❌ Server is not running. Please start it with: ./ethglobal"
    exit 1
fi

echo "✅ Server is running on http://localhost:8080"
echo ""

# Demo 1: List repositories
echo "📋 1. Listing existing repositories:"
curl -s http://localhost:8080/repos | jq .
echo ""

# Demo 2: Create a new repository
echo "🆕 2. Creating a new repository:"
REPO_NAME="demo-$(date +%s)"
curl -X POST "http://localhost:8080/repos?name=$REPO_NAME"
echo ""
echo ""

# Demo 3: List repositories again
echo "📋 3. Listing repositories after creation:"
curl -s http://localhost:8080/repos | jq .
echo ""

# Demo 4: Test Git protocol endpoint
echo "🔗 4. Testing Git protocol endpoint:"
curl -s "http://localhost:8080/$REPO_NAME/info/refs?service=git-upload-pack" | head -c 50
echo "..."
echo ""

# Demo 5: Clone the repository
echo "📥 5. Cloning the repository with Git:"
git clone "http://localhost:8080/$REPO_NAME" "cloned-$REPO_NAME"
echo ""

# Demo 6: Show cloned directory
echo "📁 6. Contents of cloned repository:"
ls -la "cloned-$REPO_NAME"
echo ""

# Demo 7: Web interface
echo "🌐 7. Web interface available at:"
echo "   http://localhost:8080"
echo ""

echo "🎉 Demo completed successfully!"
echo ""
echo "Available commands:"
echo "  • List repos:    curl http://localhost:8080/repos"
echo "  • Create repo:   curl -X POST 'http://localhost:8080/repos?name=repo-name'"
echo "  • Clone repo:    git clone http://localhost:8080/repo-name"
echo "  • Web interface: open http://localhost:8080"

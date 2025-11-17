#!/bin/bash

# Test runner script for Bifrost SDK

set -e

case "$1" in
  setup)
    echo "🔧 Running setup checks..."
    go test -v ./sdk -run "^TestSetup"
    ;;
  unit)
    echo "🧪 Running unit tests..."
    go test -v ./developpement_tests
    ;;
  integration)
    echo "🔗 Running integration tests..."
    go test -v ./developpement_tests
    ;;
  all)
    echo "🚀 Running all tests..."
    echo ""
    echo "1️⃣ Setup checks..."
    go test -v ./sdk -run "^TestSetup"
    echo ""
    echo "2️⃣ Unit tests..."
    go test -v ./developpement_tests
    echo ""
    echo "3️⃣ Integration tests..."
    go test -v ./developpement_tests
    echo ""
    echo "✅ All tests completed!"
    ;;
  *)
    echo "Usage: $0 {setup|unit|integration|all}"
    echo ""
    echo "  setup       - Run setup/configuration checks"
    echo "  unit        - Run fast unit tests"
    echo "  integration - Run integration tests (requires external services)"
    echo "  all         - Run all tests"
    exit 1
    ;;
esac

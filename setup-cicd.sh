#!/bin/bash

# CI/CD Pipeline Setup Script
# This script helps configure the GitHub Actions CI/CD pipeline for Render deployment

set -e

echo "🚀 CI/CD Pipeline Setup for Render Deployment"
echo "=============================================="

# Check if we're in a git repository
if ! git rev-parse --is-inside-work-tree > /dev/null 2>&1; then
    echo "❌ Not in a git repository"
    echo "Please run this script from your project root"
    exit 1
fi

echo "✅ Git repository detected"

# Check if GitHub CLI is installed
if command -v gh &> /dev/null; then
    echo "✅ GitHub CLI is available"
    
    # Check if user is authenticated
    if gh auth status > /dev/null 2>&1; then
        echo "✅ GitHub CLI is authenticated"
    else
        echo "⚠️  GitHub CLI not authenticated"
        echo "Run: gh auth login"
    fi
else
    echo "⚠️  GitHub CLI not found"
    echo "Install from: https://cli.github.com/"
fi

echo ""
echo "📁 Pipeline Files Created:"
echo "   .github/workflows/ci-cd.yml - Main CI/CD pipeline"
echo "   .golangci.yml - Linting configuration"
echo "   Dockerfile - Production-ready container image"
echo "   render.yaml - Render deployment configuration"
echo "   tests/load-test.js - k6 performance tests"

echo ""
echo "🔧 Required GitHub Secrets:"
echo ""
echo "Repository Secrets (GitHub → Settings → Secrets):"
echo "   RENDER_API_KEY=rnd_xxxxxxxxxxxxx"
echo "   RENDER_STAGING_SERVICE_ID=srv-xxxxxxxxxxxxx"
echo "   RENDER_PRODUCTION_SERVICE_ID=srv-xxxxxxxxxxxxx"
echo "   STAGING_HEALTH_URL=https://your-staging-app.onrender.com"
echo "   PRODUCTION_HEALTH_URL=https://your-production-app.onrender.com"
echo ""
echo "Optional Secrets:"
echo "   SLACK_WEBHOOK=https://hooks.slack.com/services/xxx/xxx/xxx"

echo ""
echo "🌐 Render Service Setup:"
echo ""
echo "1. Create Render Account: https://render.com"
echo "2. Create MongoDB Database:"
echo "   • Go to Dashboard → New → PostgreSQL/MongoDB"
echo "   • Name: demo-go-mongodb"
echo "   • Plan: Starter (Free)"
echo ""
echo "3. Create Redis Cache:"
echo "   • Go to Dashboard → New → Redis"
echo "   • Name: demo-go-redis"
echo "   • Plan: Starter (Free)"
echo ""
echo "4. Create Web Service:"
echo "   • Go to Dashboard → New → Web Service"
echo "   • Connect GitHub repository"
echo "   • Name: demo-go-api"
echo "   • Environment: Docker"
echo "   • Plan: Starter (Free)"
echo "   • Auto-Deploy: No (handled by GitHub Actions)"

echo ""
echo "🔑 Getting Render API Key:"
echo "   1. Go to Render Dashboard → Account Settings"
echo "   2. Click on 'API Keys' tab"
echo "   3. Generate new API key"
echo "   4. Copy the key (starts with 'rnd_')"

echo ""
echo "🎯 Pipeline Triggers:"
echo "   • Push to 'develop' → Deploy to Staging"
echo "   • Push to 'main' → Deploy to Production"
echo "   • Pull Request → Run Tests Only"

echo ""
echo "📊 Pipeline Features:"
echo "   ✅ Automated Testing (Unit + Integration)"
echo "   ✅ Code Linting (golangci-lint)"
echo "   ✅ Security Scanning (Gosec + Trivy)"
echo "   ✅ Docker Multi-architecture Build"
echo "   ✅ Performance Testing (k6)"
echo "   ✅ Health Checks"
echo "   ✅ Slack Notifications"
echo "   ✅ Artifact Management"

echo ""
echo "🧪 Testing the Pipeline:"
echo ""
echo "1. Commit and push changes:"
echo "   git add ."
echo "   git commit -m 'Add CI/CD pipeline'"
echo "   git push origin main"
echo ""
echo "2. Check GitHub Actions:"
echo "   • Go to GitHub → Actions tab"
echo "   • Monitor pipeline execution"
echo ""
echo "3. Verify deployment:"
echo "   • Check Render dashboard for deployment status"
echo "   • Test health endpoint: curl https://your-app.onrender.com/health"

echo ""
echo "🔍 Monitoring and Debugging:"
echo ""
echo "GitHub Actions Logs:"
echo "   gh run list"
echo "   gh run view <run-id>"
echo ""
echo "Render Logs:"
echo "   render logs --service-id srv-xxxxxxxxxxxxx --follow"
echo ""
echo "Local Testing:"
echo "   docker build -t demo-go-api ."
echo "   docker run -p 8080:8080 demo-go-api"
echo "   k6 run tests/load-test.js"

echo ""
echo "📚 Documentation:"
echo "   • Pipeline Guide: CICD_PIPELINE.md"
echo "   • Render Docs: https://render.com/docs"
echo "   • GitHub Actions: https://docs.github.com/actions"

echo ""
echo "⚠️  Important Notes:"
echo "   • Free tier has limited resources and may sleep"
echo "   • Production deployments require manual approval"
echo "   • Keep API keys and secrets secure"
echo "   • Monitor costs for production workloads"

echo ""
echo "🎉 CI/CD Pipeline setup complete!"
echo "Configure the required secrets and push to trigger your first deployment."

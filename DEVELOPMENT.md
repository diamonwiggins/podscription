# Development Setup

## Prerequisites
- Go 1.24+, Node.js 16+, OpenAI API key
- Task runner: `brew install go-task/tap/go-task`

## Quick Start
```bash
git clone https://github.com/diamonwiggins/podscription.git
cd podscription
export OPENAI_API_KEY="your-api-key-here"

# Terminal 1
task api

# Terminal 2  
task web
```

Visit: http://localhost:3000

## Tasks
```bash
task api              # Start Go backend
task web              # Start React frontend
task api:test         # Run tests
task web:lint         # Lint frontend
```

## Troubleshooting
- **Port 8080 in use**: `lsof -i :8080`
- **Missing API key**: `echo $OPENAI_API_KEY`
- **Node issues**: `rm -rf web/node_modules && cd web && npm install`
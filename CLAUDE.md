# Claude Code Skills Repository

Staging repo for claude-code skills, commands, and agents. Deployable content lives in `claude-files/`.

## Directory Structure

```
claude-files/           # Deployable to ~/.claude
├── CLAUDE.md           # Project standards (for target codebases)
├── skills/             # /skill-name triggers
│   └── <skill-name>/
│       └── SKILL.md
└── commands/           # Slash commands
ccd/                    # Deployment tool (Go)
deploy                  # Compiled binary
config.yaml             # Deployment config
```

## Authoring Skills
Use your skills-development agent to author skills in compliance for claude-code.

# Claude Code Stuff Deployment

We have a tool writen in golang living in ccd folder that we use to manage the deployments of the claude-files to the target folders (eg. ~/.claude)
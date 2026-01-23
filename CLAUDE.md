# Claude Code Skills Repository

Staging repo for claude-code skills, commands, and agents. Deployable content lives in `claude-files/`.

## Workflow
- When working on skills (anything in folder `./claude-files/skills`), always use the skill "Skill Development"


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

## Skill Structure

Each skill follows this structure:

```
skill-name/
├── SKILL.md              # Required - YAML frontmatter + instructions
├── assets/               # Optional - Files used in OUTPUT (templates, images, boilerplate)
│   └── TEMPLATE.md       # Example: document templates that get filled/copied
├── references/           # Optional - Documentation loaded INTO CONTEXT as needed
│   ├── patterns.md       # Example: detailed patterns documentation
│   └── advanced.md       # Example: advanced techniques
├── examples/             # Optional - Working code examples
│   └── example.sh
└── scripts/              # Optional - Executable utilities
    └── validate.sh
```

### Directory Purpose

| Directory | Purpose | Loaded Into Context? |
|-----------|---------|---------------------|
| `assets/` | Files used in output (templates, images, boilerplate) | No - copied/used in output |
| `references/` | Documentation Claude reads when needed | Yes - on demand |
| `examples/` | Working code examples users can copy | Yes - on demand |
| `scripts/` | Executable utilities for automation | No - executed directly |

### SKILL.md Requirements

```yaml
---
name: skill-name
version: 1.0.0
description: Third-person description with trigger phrases. This skill should be used when the user asks to "trigger phrase 1", "trigger phrase 2", or needs guidance on [specific topic].
---

# Skill Title

[Lean content - 1,500-2,000 words ideal, <3,000 max]
[Reference assets/references/examples/scripts as needed]
```

## Authoring Skills
Use your skill-development skill to author skills in compliance for claude-code.

# Claude Code Stuff Deployment

We have a tool writen in golang living in ccd folder that we use to manage the deployments of the claude-files to the target folders (eg. ~/.claude)
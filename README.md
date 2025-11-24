# Procedural Testing Framework (proctest)

A spike project exploring automated testing for technical documentation procedures.

## The Problem

Technical documentation often contains step-by-step procedures that readers follow to accomplish tasks. These procedures include executable elements like code snippets, shell commands, and UI interactions. Currently, there's no systematic way to verify that these procedures actually work as written.

**Challenges**:
- Procedures can break when software updates change APIs, UIs, or behavior
- Manual testing is time-consuming and doesn't scale across large documentation sets
- Writers discover broken procedures only when users report issues
- No visibility into which procedures have been tested and verified

## The Solution

`proctest` is a testing framework designed specifically for technical writers. It:

1. **Parses documentation** (ReStructuredText initially, MDX in the future) to extract procedures and testable actions
2. **Executes testable actions** (code blocks, shell commands, CLI commands, UI interactions, API calls)
3. **Validates results** to ensure procedures work as documented
4. **Reports failures** with clear, actionable error messages
5. **Integrates with CI** to catch regressions automatically

**Key Design Principles**:
- ✅ **Zero-config by default** - Works out of the box for 90% of cases
- ✅ **Writer-friendly** - Simple CLI, clear error messages, minimal technical overhead
- ✅ **Progressive disclosure** - Advanced features available when needed
- ✅ **Convention over configuration** - Smart defaults based on common patterns

## Documentation

### For Technical Writers

**→ [Usage Guide](usage.md)** - Learn how to use `proctest` to test your documentation

This guide covers:
- Quick start and installation
- Local development workflow
- Running tests and interpreting results
- CI integration
- Monorepo conventions and shared configurations
- Common patterns and troubleshooting

### For Developers

**→ [Technical Specification](technical-specification.md)** - Implementation details and architecture

This specification covers:
- System architecture and component design
- Data models and interfaces
- Parser implementation (RST, tabs, variants)
- Executor design for different testable action types
- Implementation plan (3 phases, 11 milestones)
- Testing strategy and risk assessment

### Supporting Documents

- **[Requirements](requirements.md)** - Original requirements and use cases that informed the design
- **[Architecture Options](architecture-options.md)** - Five architectural approaches evaluated, with rationale for selecting Option 5 (Convention + Configuration)

## Project Status

**Status**: 🔬 **Spike / Exploration Phase**

This is a spike project to explore feasibility and design. The documents in this repository represent:
- ✅ Comprehensive requirements analysis
- ✅ Architectural evaluation and decision-making
- ✅ Detailed technical specification
- ✅ Usage documentation and examples
- ⏳ No implementation yet

**Next Steps**:
1. Review and validate the approach with stakeholders
2. Prioritize features for MVP (Phase 1)
3. Begin implementation following the technical specification
4. Iterate based on real-world usage and feedback

## Quick Example

```bash
# Install
npm install -g @mongodb/proctest

# Test a single procedure
proctest test content/atlas/source/tutorial/getting-started.txt

# Test all procedures in a directory
proctest test content/atlas/source/

# Debug: See what proctest detects
proctest parse content/atlas/source/tutorial/getting-started.txt

# CI: Run verified procedures from registry
proctest test --registry code-example-tests/procedures/test-registry.json --reporter junit
```

## Key Features

### Phase 1 (MVP)
- ✅ RST parsing with procedure detection
- ✅ Code block execution (JavaScript, Python, PHP, etc.)
- ✅ Shell command execution
- ✅ Placeholder resolution from `.env` and `snooty.toml`
- ✅ Prerequisite detection and validation
- ✅ Automatic cleanup of test resources
- ✅ Human-readable and JUnit XML output

### Phase 2
- ✅ CLI command execution (mongosh, atlas-cli)
- ✅ Download testable actions
- ✅ Sub-procedure support (nested steps)
- ✅ Enhanced error reporting

### Phase 3
- ✅ UI testing with navigation mappings
- ✅ Atlas Admin API testing
- ✅ URL validation
- ✅ Screenshot capture on failures

## Technology Stack

- **Language**: TypeScript
- **Runtime**: Node.js 24+ LTS (native `.env` support)
- **Testing**: Jest with ts-jest
- **CLI Framework**: Commander.js
- **UI Automation**: Playwright (Phase 3)

## Contributing

This is currently a spike project. Once we move to implementation:

1. Follow the technical specification for architecture and design decisions
2. Write tests for all new functionality
3. Update documentation to reflect changes
4. Maintain the writer-first, zero-config philosophy

## License

TBD - MongoDB internal project

---

**Questions?** Review the [Usage Guide](usage.md) for writers or the [Technical Specification](technical-specification.md) for developers.


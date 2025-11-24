# Procedural Testing Framework (proctest)

A spike project exploring automated testing for technical documentation procedures.

## The Problem

Technical documentation often contains step-by-step procedures that readers follow to accomplish tasks. These procedures include executable elements like code snippets, shell commands, and UI interactions. While some tools exist for testing documentation (like [Doc Detective](https://doc-detective.com/)), they require writers to either annotate documentation files or maintain separate test files that can drift out of sync with the actual page content.

**Key Challenges**:
- Procedures can break when software updates change APIs, UIs, or behavior
- Manual testing is time-consuming and doesn't scale across large documentation sets
- Writers discover broken procedures only when users report issues
- **Existing tools require annotations or separate test files** that can become out-of-step with documentation
- **No automated test discovery** - writers must manually identify and configure each test
- **Point-in-time snapshots** - tests represent what the documentation *was*, not what it *is*

## The Solution

`proctest` is a testing framework designed specifically for technical writers. It tests procedures **as written** in the living documentation, with no annotations or separate test files required.

**How it works**:

1. **Automatic discovery** - Parses documentation (ReStructuredText initially, MDX in the future) to extract procedures and testable actions
2. **Tests living content** - Executes what's actually on the page right now, not a point-in-time snapshot
3. **No drift** - Tests always reflect current documentation because they *are* the documentation
4. **Validates results** - Ensures procedures work as documented
5. **Clear reporting** - Provides actionable error messages for writers
6. **CI integration** - Catches regressions automatically

**Key Differentiators**:
- ✅ **No annotations required** - Tests the documentation as-is
- ✅ **No separate test files** - Procedures and tests can't drift out of sync
- ✅ **Automatic test discovery** - Finds procedures without manual configuration
- ✅ **Zero-config by default** - Works out of the box for 90% of cases
- ✅ **Writer-friendly** - Simple CLI, clear error messages, minimal technical overhead
- ✅ **Progressive disclosure** - Advanced features available when needed

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
npm install -g @grove-platform/proctest

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

## Why This Spike?

### Comparison with Doc Detective

[Doc Detective](https://doc-detective.com) is an existing open-source documentation testing framework that shares similar goals. We're spiking on `proctest` to explore an approach that better fits MongoDB's specific needs and documentation structure.

| Feature | Doc Detective | proctest (This Spike) |
|---------|--------------|----------------------|
| **Markup Support** | Markdown, JSON test specs | **ReStructuredText (RST)**, MDX (future) |
| **Test Discovery** | Limited (scans for test specs) | **Automatic procedure detection** from documentation |
| **Test Definition** | Separate test files or inline annotations | **No separate files or annotations** - tests living documentation |
| **Drift Risk** | ⚠️ High - tests can become out-of-sync with docs | ✅ **Zero** - tests *are* the documentation |
| **Code Execution** | `runCode` action | Code, Shell, CLI (mongosh, atlas-cli) |
| **UI Testing** | `click`, `find`, `type` actions | **Navigation mappings** for generic UI instructions |
| **API Testing** | Generic `httpRequest` | **MongoDB Atlas Admin API** specific support |
| **Prerequisites** | Manual configuration | **Automatic detection and validation** |
| **Placeholders** | `loadVariables` from `.env` | **Multi-source resolution** (.env, snooty.toml, config) |
| **Variants/Tabs** | Not supported | **Automatic variant expansion** (tabs, composable tutorials) |
| **Cleanup** | Manual | **Automatic resource tracking and cleanup** |
| **Configuration** | Required `.doc-detective.json` | **Zero-config by default**, optional for edge cases |
| **Writer Workflow** | Write docs → Write separate tests → Maintain both | **Write docs → Tests automatically discovered** |

### Key Differentiators

**1. Tests Living Documentation**
- Doc Detective requires either inline annotations or separate test specification files
- `proctest` tests the documentation **as written** - no annotations, no separate files
- Eliminates drift between documentation and tests

**2. RST Support**
- Doc Detective supports Markdown but not ReStructuredText
- `proctest` is designed specifically for RST (MongoDB's current format)
- Lightweight, targeted RST parsing (not heavyweight Snooty parser)

**3. Automatic Test Discovery**
- Doc Detective requires explicit test specifications
- `proctest` automatically discovers procedures in documentation
- Writers don't need to configure what to test - it's automatic

**4. MongoDB-Specific Features**
- Atlas Admin API testing with authentication
- mongosh and atlas-cli command execution
- Integration with `snooty.toml` for source constants
- Support for MongoDB documentation patterns (tabs, composable tutorials)

**5. Zero-Config Philosophy**
- Doc Detective requires configuration files for most use cases
- `proctest` works out-of-the-box for 90% of cases
- Configuration only needed for edge cases

### Why Spike Instead of Adopt?

While Doc Detective is a solid framework, we're exploring `proctest` because:

1. **Format mismatch** - Our docs are in RST, not Markdown
2. **Drift prevention** - We want to test living documentation, not point-in-time snapshots
3. **Writer experience** - Zero-config, automatic discovery reduces friction
4. **MongoDB-specific needs** - Atlas API, mongosh, atlas-cli, snooty.toml integration
5. **Maintenance burden** - Separate test files create additional maintenance overhead

This spike validates whether a purpose-built solution better serves our needs than adapting an existing tool.

---

**Questions?** Review the [Usage Guide](usage.md) for writers or the [Technical Specification](technical-specification.md) for developers.


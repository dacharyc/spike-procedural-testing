# Architecture Design Options

## Executive Summary

**Recommendation**: **Option 5: Hybrid + Plugin Ready** ⭐

This architecture provides the optimal balance for MongoDB's procedural testing framework:
- ✅ **Zero config for 90% of users** - Writers can start immediately
- ✅ **Progressive complexity** - Add configuration only when needed
- ✅ **Future-proof architecture** - Plugin-ready without plugin overhead
- ✅ **3-4 week PoC timeline** - Minimal upfront investment (1-2 extra days for interface design)
- ✅ **No refactoring needed** - Stable architecture for 24+ months

### Quick Comparison

| What You Get | Option 1: Convention | Option 3: Hybrid | Option 5: Hybrid + Plugin Ready ⭐ |
|--------------|---------------------|------------------|-------------------------------------|
| **Zero config works** | ✅ Yes | ✅ Yes | ✅ Yes |
| **PoC timeline** | 1-2 weeks | 3-4 weeks | 3-4 weeks |
| **Config for complex cases** | ⚠️ Limited | ✅ Yes | ✅ Yes |
| **Plugin support (future)** | ❌ Requires refactor | ⚠️ Requires refactor | ✅ Built-in |
| **Refactoring risk** | 🔴 High | 🟡 Moderate | 🟢 Minimal |
| **Long-term stability** | 🟡 12 months | 🟢 18 months | 🟢 24+ months |

### Why Option 5?

**Near-term** (PoC): Identical user experience to Option 1 (zero config), with only 1-2 extra days to define clean interfaces.

**Mid-term** (Production): Identical capabilities to Option 3 (progressive configuration), with better code organization.

**Long-term** (Scale): Plugin-ready architecture enables community contributions and team customization without refactoring.

**Investment**: Small upfront cost (1-2 days) prevents months of refactoring later.

---

## Overview

This document presents five architectural approaches for the procedural testing framework, each optimized for different priorities while maintaining the core goal: enabling technical writers to validate that documentation procedures work as written.

## Key Design Considerations

Before exploring specific architectures, let's establish the key considerations that will guide our design:

### User-Centric Priorities

1. **Minimal Learning Curve**: Technical writers should be able to use the tool with minimal training
2. **Clear Error Messages**: When tests fail, writers need actionable information to fix the issue
3. **Predictable Behavior**: The tool should behave consistently and intuitively
4. **Low Maintenance Burden**: Writers shouldn't need to maintain complex configurations

### Technical Priorities

1. **Extensibility**: Easy to add new languages, directives, or execution strategies
2. **Debuggability**: When something goes wrong, it should be easy to diagnose
3. **Isolation**: Tests should not interfere with each other or the writer's environment
4. **Accuracy**: The tool should execute procedures exactly as a developer would follow them

### Documentation Quality Priorities

1. **Real-World Validation**: Tests should validate procedures work in real environments
2. **Comprehensive Coverage**: Support for code, CLI, API, and eventually UI interactions
3. **Version Awareness**: Handle different versions of tools, languages, and APIs
4. **Failure Transparency**: Make it obvious when documentation is broken

---

## Option 1: Convention-Over-Configuration (Simplest Writer Usage)

### Philosophy

"Zero configuration for 90% of use cases, with escape hatches for the other 10%"

### Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                     CLI Entry Point                          │
│                  (procedure-test <file>)                     │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                  Convention Engine                           │
│  • Auto-detect .env files                                    │
│  • Auto-detect snooty.toml                                   │
│  • Auto-discover test files in conventional locations        │
│  • Apply smart defaults for everything                       │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                    RST Parser                                │
│  • Parse procedures, steps, code blocks                      │
│  • Resolve includes and transclusions                        │
│  • Extract metadata (language, options)                      │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              Smart Placeholder Resolver                      │
│  • Use placeholder-consistency.md patterns                   │
│  • Fuzzy match environment variables                         │
│  • Auto-detect common patterns                               │
│  • Fail with suggestions when unresolvable                   │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                 Execution Orchestrator                       │
│  • Sequential test execution                                 │
│  • Per-step state management                                 │
│  • Automatic cleanup detection                               │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              Language-Specific Executors                     │
│  ┌──────────┬──────────┬──────────┬──────────┐              │
│  │JavaScript│  Python  │   PHP    │  Shell   │              │
│  └──────────┴──────────┴──────────┴──────────┘              │
│  • Auto-detect installed tools                               │
│  • Create isolated execution contexts                        │
│  • Accumulate code within steps                              │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                  Reporter                                    │
│  • Human-friendly output for writers                         │
│  • File/line numbers for failures                            │
│  • Suggested fixes when possible                             │
│  • Optional verbose mode                                     │
└─────────────────────────────────────────────────────────────┘
```

### Key Features

**Convention-Based Discovery**:
- Looks for `.env` in current directory, parent directories, or `~/.procedure-test/.env`
- Finds `snooty.toml` by walking up directory tree
- Discovers test files in `procedures/` directory or via simple file list

**Smart Defaults**:
- Automatically maps common placeholders to environment variables
- Uses fuzzy matching: `<connection-string>`, `<connectionString>`, `<connection string>` all map to `CONNECTION_STRING`
- Detects language from code block or file extension
- Assumes sequential execution, isolated state per procedure

**Minimal Configuration**:
Writers only need to:
1. Create a `.env` file with their credentials
2. Run `procedure-test <file>` or `procedure-test --all`

### Escape Hatches (Optional Configuration)

For the 10% of complex cases, writers can create an optional `.procedure-test.json`:

```json
{
  "placeholders": {
    "<my-weird-placeholder>": "MY_ENV_VAR"
  },
  "executors": {
    "javascript": {
      "runtime": "node",
      "version": ">=18"
    }
  },
  "cleanup": {
    "databases": true,
    "collections": true,
    "files": true
  }
}
```

### Pros

✅ **Extremely low barrier to entry** - Writers can start testing immediately
✅ **Minimal maintenance** - No configuration files to keep in sync
✅ **Self-documenting** - Conventions are discoverable and predictable
✅ **Fast iteration** - Writers can test changes without configuration updates
✅ **Reduced errors** - Fewer configuration files means fewer places for mistakes

### Cons

❌ **Less explicit** - Behavior might seem "magical" to some users
❌ **Harder to customize** - Complex scenarios require understanding escape hatches
❌ **Convention lock-in** - Changing conventions could break existing workflows
❌ **Debugging complexity** - When conventions fail, it's harder to understand why

### Best For

- Teams that value simplicity over flexibility
- Writers with minimal technical background
- Projects with consistent structure and patterns
- Rapid prototyping and iteration

---

## Option 2: Configuration-First (Flexible & Customizable)

### Philosophy

"Explicit configuration provides clarity, control, and customization"

### Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                     CLI Entry Point                          │
│            (procedure-test --config test.config.js)          │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                Configuration Loader                          │
│  • Load and validate configuration file                      │
│  • Merge with defaults                                       │
│  • Validate required settings                                │
│  • Support multiple config formats (JSON, JS, YAML)          │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                    RST Parser                                │
│  • Configurable directive handlers                           │
│  • Pluggable transclusion resolvers                          │
│  • Custom metadata extractors                                │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│            Configurable Placeholder Resolver                 │
│  • User-defined placeholder mappings                         │
│  • Custom resolver functions                                 │
│  • Fallback to smart defaults                                │
│  • Validation rules per placeholder                          │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                 Execution Orchestrator                       │
│  • Configurable execution strategy                           │
│  • Custom hooks (before/after step, procedure)               │
│  • Pluggable cleanup handlers                                │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              Language Executor Registry                      │
│  • Pluggable executor implementations                        │
│  • Custom executor configuration                             │
│  • Runtime version management                                │
│  • Execution environment customization                       │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              Configurable Reporter                           │
│  • Multiple output formats (JSON, TAP, JUnit, human)         │
│  • Custom formatters                                         │
│  • Configurable verbosity levels                             │
└─────────────────────────────────────────────────────────────┘
```

### Key Features

**Explicit Configuration File** (`procedure-test.config.js`):

```javascript
module.exports = {
  // Test discovery
  testFiles: [
    'testdata/atlas/source/**/*.txt',
    'testdata/drivers/source/**/*.txt'
  ],
  exclude: ['**/includes/**'],

  // Environment
  envFiles: ['.env', '.env.local'],
  snootyConfig: 'testdata/atlas/snooty.toml',

  // Placeholder resolution
  placeholders: {
    resolvers: [
      // Custom resolver function
      (placeholder, context) => {
        if (placeholder.match(/<cluster.*>/i)) {
          return process.env.ATLAS_CLUSTER_NAME;
        }
      },
      // Built-in fuzzy resolver
      'fuzzy-env-match',
      // Explicit mappings
      {
        '<connection-string>': 'MONGODB_URI',
        '<username>': 'DB_USERNAME',
        '<password>': 'DB_PASSWORD'
      }
    ],
    onUnresolved: 'fail' // or 'warn', 'skip'
  },

  // Execution
  executors: {
    javascript: {
      runtime: 'node',
      version: '>=18',
      timeout: 30000,
      env: { NODE_ENV: 'test' }
    },
    python: {
      runtime: 'python3',
      virtualenv: true,
      requirements: 'auto-detect'
    },
    shell: {
      shell: '/bin/bash',
      env: { TERM: 'xterm' }
    }
  },

  // Hooks for custom behavior
  hooks: {
    beforeAll: async (context) => {
      // Setup test database
    },
    beforeEach: async (procedure, context) => {
      // Per-procedure setup
    },
    afterEach: async (procedure, result, context) => {
      // Cleanup
    },
    afterAll: async (results, context) => {
      // Teardown
    }
  },

  // Cleanup
  cleanup: {
    databases: {
      enabled: true,
      pattern: /^test_/,
      onFailure: 'warn' // or 'block', 'ignore'
    },
    collections: {
      enabled: true,
      pattern: /^temp_/
    },
    files: {
      enabled: true,
      paths: ['./temp/**']
    }
  },

  // Reporting
  reporters: [
    'human', // Built-in human-friendly reporter
    ['json', { outputFile: 'test-results.json' }],
    ['custom', { formatter: './my-formatter.js' }]
  ],

  // Advanced
  stateManagement: {
    persistAcrossSteps: false,
    isolationLevel: 'procedure' // or 'step', 'file'
  },

  parallel: false, // Future: enable parallel execution

  verbose: false
};
```

### Simplified Configuration for Common Cases

For writers who don't need customization, provide a minimal config:

```javascript
module.exports = {
  testFiles: ['procedures/**/*.txt'],
  envFiles: ['.env']
};
```

### Pros

✅ **Maximum flexibility** - Every aspect is configurable
✅ **Explicit behavior** - No surprises, everything is documented in config
✅ **Powerful customization** - Hooks and custom functions for complex scenarios
✅ **Team-specific workflows** - Different teams can configure differently
✅ **Version control friendly** - Configuration is code, can be reviewed and versioned
✅ **Easier debugging** - Explicit configuration makes behavior predictable

### Cons

❌ **Higher learning curve** - Writers need to understand configuration options
❌ **More maintenance** - Configuration files need to be kept up-to-date
❌ **Potential for errors** - Misconfiguration can cause confusing failures
❌ **Overwhelming for simple cases** - Too many options for basic usage
❌ **Documentation burden** - Need comprehensive docs for all config options

### Best For

- Teams with diverse testing needs
- Projects requiring custom workflows
- Writers comfortable with configuration files
- Organizations with dedicated DevOps/tooling support
- Complex multi-project documentation repositories

---

## Option 3: Hybrid Approach (Smart Defaults + Progressive Disclosure)

### Philosophy

"Make simple things simple, and complex things possible"

### Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                     CLI Entry Point                          │
│         (procedure-test <file> [--config <file>])            │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              Smart Configuration Manager                     │
│  • Load optional config file if present                      │
│  • Apply convention-based defaults                           │
│  • Merge config with conventions (config wins)               │
│  • Validate and provide helpful error messages               │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                    RST Parser                                │
│  • Standard directive handlers (extensible via config)       │
│  • Auto-detect transclusion patterns                         │
│  • Metadata extraction with smart defaults                   │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│            Layered Placeholder Resolver                      │
│  1. Explicit config mappings (if provided)                   │
│  2. Fuzzy environment variable matching                      │
│  3. Smart pattern detection                                  │
│  4. Fail with helpful suggestions                            │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                 Execution Orchestrator                       │
│  • Convention-based execution (overridable)                  │
│  • Optional hooks (only if configured)                       │
│  • Automatic cleanup (configurable)                          │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              Language Executor Manager                       │
│  • Auto-detect and use installed tools                       │
│  • Allow version/runtime overrides via config                │
│  • Provide sensible defaults for each language               │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              Adaptive Reporter                               │
│  • Human-friendly by default                                 │
│  • Additional formats via config                             │
│  • Auto-adjust verbosity based on context                    │
└─────────────────────────────────────────────────────────────┘
```

### Key Features

**Zero Config for Simple Cases**:
```bash
# Just works with conventions
procedure-test testdata/atlas/source/connect-to-database-deployment.txt
```

**Progressive Configuration**:

**Level 1 - Minimal** (`.procedure-test.json`):
```json
{
  "testFiles": ["procedures/**/*.txt"]
}
```

**Level 2 - Common Customizations**:
```json
{
  "testFiles": ["procedures/**/*.txt"],
  "placeholders": {
    "<my-cluster>": "ATLAS_CLUSTER"
  },
  "cleanup": {
    "databases": true
  }
}
```

**Level 3 - Advanced** (`.procedure-test.js`):
```javascript
module.exports = {
  testFiles: ['procedures/**/*.txt'],

  // Only configure what you need to override
  placeholders: {
    // Custom resolver for complex cases
    resolver: (placeholder, context) => {
      if (placeholder.includes('cluster')) {
        return process.env.ATLAS_CLUSTER_NAME;
      }
      // Fall back to default fuzzy matching
      return null;
    }
  },

  // Optional hooks - only if you need them
  hooks: {
    afterEach: async (procedure, result) => {
      if (result.failed) {
        await captureDebugInfo(procedure);
      }
    }
  }
};
```

### Layered Placeholder Resolution

The resolver tries strategies in order:

1. **Explicit config mappings** (if config file exists)
2. **Fuzzy environment variable matching** (built-in intelligence)
3. **Pattern-based detection** (using placeholder-consistency.md)
4. **Fail with suggestions**:
   ```
   ❌ Could not resolve placeholder: <my-cluster-name>

   Suggestions:
   - Add to .env: MY_CLUSTER_NAME=...
   - Add to config: "placeholders": { "<my-cluster-name>": "ENV_VAR" }
   - Similar env vars found: ATLAS_CLUSTER, CLUSTER_NAME
   ```

### Smart Defaults with Override Points

**Default Behavior**:
- Auto-discover `.env` files
- Auto-detect `snooty.toml`
- Sequential execution
- Automatic cleanup of test databases
- Human-friendly output

**Override Points** (only configure if needed):
```javascript
{
  // Override discovery
  "envFiles": [".env.test", ".env"],

  // Override execution
  "executors": {
    "javascript": { "timeout": 60000 }
  },

  // Override cleanup
  "cleanup": {
    "databases": { "pattern": /^mytest_/ }
  },

  // Override reporting
  "reporters": ["human", "json"]
}
```

### CLI Flags for Quick Overrides

```bash
# Use conventions
procedure-test myfile.txt

# Override specific behavior
procedure-test myfile.txt --verbose
procedure-test myfile.txt --no-cleanup
procedure-test myfile.txt --env .env.staging
procedure-test myfile.txt --config custom.config.js

# Generate starter config
procedure-test --init
```

### Guided Configuration Generation

```bash
$ procedure-test --init

? Where are your test files? (procedures/**/*.txt)
? Do you have a .env file? (Y/n) y
? Do you need custom placeholder mappings? (y/N) n
? Do you want automatic database cleanup? (Y/n) y

✅ Created .procedure-test.json with your settings
```

### Pros

✅ **Low barrier to entry** - Works immediately with zero config
✅ **Grows with complexity** - Add configuration only when needed
✅ **Discoverable** - CLI flags and `--init` help users learn
✅ **Flexible** - Supports both simple and complex scenarios
✅ **Maintainable** - Minimal config for most cases, detailed config when needed
✅ **Self-documenting** - Conventions are clear, config is explicit
✅ **Helpful errors** - Suggests solutions when things go wrong

### Cons

❌ **More complex implementation** - Need to handle both conventions and config
❌ **Potential confusion** - Users might not know when to use config vs. conventions
❌ **Testing complexity** - Need to test both convention and config paths
❌ **Documentation challenge** - Need to document both approaches clearly

### Best For

- **This is the recommended approach for MongoDB's use case**
- Teams with varying technical skill levels
- Projects that start simple but may grow complex
- Organizations that value both ease-of-use and flexibility
- Documentation teams with diverse needs across different projects

---

## Option 4: Plugin-Based Architecture (Maximum Extensibility)

### Philosophy

"Core framework provides infrastructure; plugins provide functionality"

### Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                     CLI Entry Point                          │
│                  (procedure-test <file>)                     │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                    Core Framework                            │
│  • Plugin discovery and loading                              │
│  • Event bus for plugin communication                        │
│  • Minimal built-in functionality                            │
│  • Plugin lifecycle management                               │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                    Plugin Registry                           │
│  ┌────────────────────────────────────────────────┐          │
│  │  Parser Plugins                                │          │
│  │  • @procedure-test/rst-parser                  │          │
│  │  • @procedure-test/markdown-parser (future)    │          │
│  └────────────────────────────────────────────────┘          │
│  ┌────────────────────────────────────────────────┐          │
│  │  Resolver Plugins                              │          │
│  │  • @procedure-test/env-resolver                │          │
│  │  • @procedure-test/snooty-resolver             │          │
│  │  • @procedure-test/fuzzy-resolver              │          │
│  └────────────────────────────────────────────────┘          │
│  ┌────────────────────────────────────────────────┐          │
│  │  Executor Plugins                              │          │
│  │  • @procedure-test/executor-javascript         │          │
│  │  • @procedure-test/executor-python             │          │
│  │  • @procedure-test/executor-php                │          │
│  │  • @procedure-test/executor-shell              │          │
│  └────────────────────────────────────────────────┘          │
│  ┌────────────────────────────────────────────────┐          │
│  │  Reporter Plugins                              │          │
│  │  • @procedure-test/reporter-human              │          │
│  │  • @procedure-test/reporter-json               │          │
│  │  • @procedure-test/reporter-junit              │          │
│  └────────────────────────────────────────────────┘          │
│  ┌────────────────────────────────────────────────┐          │
│  │  Cleanup Plugins                               │          │
│  │  • @procedure-test/cleanup-mongodb             │          │
│  │  • @procedure-test/cleanup-files               │          │
│  └────────────────────────────────────────────────┘          │
└─────────────────────────────────────────────────────────────┘
```

### Key Features

**Minimal Core**:
The core framework only provides:
- Plugin loading and lifecycle
- Event bus for plugin communication
- Configuration management
- Error handling infrastructure

**Everything Else is a Plugin**:
```javascript
// .procedure-test.config.js
module.exports = {
  plugins: [
    // Parser
    '@procedure-test/rst-parser',

    // Resolvers (run in order)
    '@procedure-test/snooty-resolver',
    '@procedure-test/env-resolver',
    '@procedure-test/fuzzy-resolver',

    // Executors (auto-selected by language)
    '@procedure-test/executor-javascript',
    '@procedure-test/executor-python',
    '@procedure-test/executor-php',
    '@procedure-test/executor-shell',

    // Cleanup
    '@procedure-test/cleanup-mongodb',

    // Reporters
    '@procedure-test/reporter-human',

    // Custom plugins
    './my-custom-plugin.js'
  ],

  // Plugin-specific configuration
  pluginConfig: {
    '@procedure-test/fuzzy-resolver': {
      threshold: 0.8
    },
    '@procedure-test/executor-javascript': {
      runtime: 'node',
      version: '>=18'
    }
  }
};
```

**Plugin API**:
```javascript
// Example custom plugin
module.exports = {
  name: 'my-custom-resolver',
  type: 'resolver',

  // Plugin lifecycle
  async init(context) {
    // Setup
  },

  // Plugin functionality
  async resolve(placeholder, context) {
    if (placeholder.startsWith('<custom-')) {
      return process.env[placeholder.toUpperCase()];
    }
    return null; // Pass to next resolver
  },

  async cleanup() {
    // Teardown
  }
};
```

**Preset Configurations**:
```javascript
// @procedure-test/preset-mongodb
module.exports = {
  plugins: [
    '@procedure-test/rst-parser',
    '@procedure-test/snooty-resolver',
    '@procedure-test/env-resolver',
    '@procedure-test/fuzzy-resolver',
    '@procedure-test/executor-javascript',
    '@procedure-test/executor-python',
    '@procedure-test/executor-php',
    '@procedure-test/executor-shell',
    '@procedure-test/cleanup-mongodb',
    '@procedure-test/reporter-human'
  ]
};

// User config
module.exports = {
  extends: '@procedure-test/preset-mongodb',
  plugins: [
    './my-custom-plugin.js' // Add custom plugin
  ]
};
```

### Pros

✅ **Maximum extensibility** - Easy to add new functionality
✅ **Community contributions** - Others can create plugins
✅ **Separation of concerns** - Each plugin has a single responsibility
✅ **Testability** - Plugins can be tested independently
✅ **Flexibility** - Users can mix and match plugins
✅ **Future-proof** - New features don't require core changes

### Cons

❌ **Complexity** - Plugin architecture adds significant complexity
❌ **Overhead** - Plugin loading and communication has performance cost
❌ **Learning curve** - Writers need to understand plugin system
❌ **Maintenance burden** - More packages to maintain and version
❌ **Overkill for initial PoC** - Too much infrastructure for current needs
❌ **Dependency management** - More packages means more dependencies

### Best For

- Large-scale projects with diverse needs
- Organizations building a platform for multiple teams
- Projects expecting significant community contributions
- Long-term projects with evolving requirements
- **Not recommended for initial PoC** - too much overhead

---

## Comparison Matrix

| Criteria | Option 1: Convention | Option 2: Configuration | Option 3: Hybrid | Option 4: Plugin | Option 5: Hybrid + Plugin Ready ⭐ |
|----------|---------------------|------------------------|------------------|------------------|------------------------------------|
| **Writer Learning Curve** | ⭐⭐⭐⭐⭐ Minimal | ⭐⭐ Moderate | ⭐⭐⭐⭐ Low | ⭐⭐ Moderate | ⭐⭐⭐⭐⭐ Minimal |
| **Setup Time** | ⭐⭐⭐⭐⭐ Instant | ⭐⭐ Requires config | ⭐⭐⭐⭐⭐ Instant | ⭐⭐ Requires config | ⭐⭐⭐⭐⭐ Instant |
| **Flexibility** | ⭐⭐ Limited | ⭐⭐⭐⭐⭐ Maximum | ⭐⭐⭐⭐ High | ⭐⭐⭐⭐⭐ Maximum | ⭐⭐⭐⭐⭐ Maximum |
| **Customization** | ⭐⭐ Escape hatches | ⭐⭐⭐⭐⭐ Everything | ⭐⭐⭐⭐ Progressive | ⭐⭐⭐⭐⭐ Everything | ⭐⭐⭐⭐⭐ Progressive to Everything |
| **Debuggability** | ⭐⭐⭐ Good | ⭐⭐⭐⭐ Explicit | ⭐⭐⭐⭐ Very good | ⭐⭐⭐ Complex | ⭐⭐⭐⭐⭐ Excellent |
| **Maintenance** | ⭐⭐⭐⭐⭐ Minimal | ⭐⭐ High | ⭐⭐⭐⭐ Low | ⭐⭐ High | ⭐⭐⭐⭐⭐ Minimal to Low |
| **Implementation Complexity** | ⭐⭐⭐ Moderate | ⭐⭐⭐ Moderate | ⭐⭐⭐⭐ Higher | ⭐ Very high | ⭐⭐⭐ Moderate+ |
| **Extensibility** | ⭐⭐ Limited | ⭐⭐⭐ Good | ⭐⭐⭐⭐ Very good | ⭐⭐⭐⭐⭐ Maximum | ⭐⭐⭐⭐⭐ Maximum |
| **Error Messages** | ⭐⭐⭐⭐ Helpful | ⭐⭐⭐ Good | ⭐⭐⭐⭐⭐ Very helpful | ⭐⭐⭐ Good | ⭐⭐⭐⭐⭐ Very helpful |
| **Suitable for PoC** | ⭐⭐⭐⭐ Yes | ⭐⭐⭐⭐ Yes | ⭐⭐⭐⭐⭐ Ideal | ⭐⭐ Overkill | ⭐⭐⭐⭐⭐ Ideal |
| **Long-term Viability** | ⭐⭐⭐ Good | ⭐⭐⭐⭐ Very good | ⭐⭐⭐⭐ Excellent | ⭐⭐⭐⭐⭐ Excellent | ⭐⭐⭐⭐⭐ Excellent |
| **Refactoring Risk** | ⭐⭐ High | ⭐⭐⭐ Moderate | ⭐⭐⭐ Moderate | ⭐⭐⭐⭐⭐ None | ⭐⭐⭐⭐⭐ Minimal |
| **Community Extensibility** | ⭐ Very limited | ⭐⭐ Limited | ⭐⭐⭐ Moderate | ⭐⭐⭐⭐⭐ Maximum | ⭐⭐⭐⭐⭐ Maximum (future) |
| **PoC to Production Path** | ⭐⭐ Requires refactor | ⭐⭐⭐⭐ Smooth | ⭐⭐⭐⭐ Smooth | ⭐⭐ Overbuilt | ⭐⭐⭐⭐⭐ Seamless |

## Detailed Comparison

### For Technical Writers (Primary Users)

**Easiest to Use**:
1. Option 1 (Convention) - Zero config, just run
2. Option 3 (Hybrid) - Zero config initially, add as needed
3. Option 2 (Configuration) - Requires understanding config file
4. Option 4 (Plugin) - Requires understanding plugins

**Most Helpful When Things Go Wrong**:
1. Option 3 (Hybrid) - Smart suggestions based on conventions + config
2. Option 1 (Convention) - Good suggestions based on conventions
3. Option 2 (Configuration) - Clear but requires config knowledge
4. Option 4 (Plugin) - May be unclear which plugin is failing

**Lowest Maintenance Burden**:
1. Option 1 (Convention) - No config to maintain
2. Option 3 (Hybrid) - Minimal config for most cases
3. Option 2 (Configuration) - Config must be kept in sync
4. Option 4 (Plugin) - Config + plugin versions to manage

### For Complex Scenarios

**Most Flexible**:
1. Option 4 (Plugin) - Can replace any component
2. Option 2 (Configuration) - Everything is configurable
3. Option 3 (Hybrid) - Most things configurable
4. Option 1 (Convention) - Limited escape hatches

**Best for Custom Workflows**:
1. Option 2 (Configuration) - Hooks and custom functions
2. Option 4 (Plugin) - Custom plugins
3. Option 3 (Hybrid) - Optional hooks
4. Option 1 (Convention) - Limited customization

### For Implementation

**Fastest to Build (PoC)**:
1. Option 1 (Convention) - Straightforward implementation
2. Option 2 (Configuration) - Straightforward with config layer
3. Option 3 (Hybrid) - Requires both convention + config logic
4. Option 4 (Plugin) - Requires plugin infrastructure

**Easiest to Test**:
1. Option 2 (Configuration) - Explicit behavior
2. Option 4 (Plugin) - Isolated plugins
3. Option 1 (Convention) - Predictable conventions
4. Option 3 (Hybrid) - Must test both paths

**Most Maintainable Long-term**:
1. Option 3 (Hybrid) - Balance of simplicity and flexibility
2. Option 4 (Plugin) - Separation of concerns
3. Option 2 (Configuration) - Explicit and documented
4. Option 1 (Convention) - May become limiting

---

## Option 5: Hybrid + Plugin Ready (Recommended)

### Philosophy

"Start simple, grow naturally, architect for the future"

This option combines the best aspects of Option 3 (Hybrid) with the forward-looking architecture of Option 4 (Plugin), but implements the plugin system only when needed rather than upfront.

### Architecture Overview

```
┌─────────────────────────────────────────────────────────────┐
│                     CLI Entry Point                          │
│         (procedure-test <file> [--config <file>])            │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              Smart Configuration Manager                     │
│  • Convention-based defaults (zero config)                   │
│  • Optional config file (progressive disclosure)             │
│  • Plugin registry (future-ready, not required)              │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                  Component Interfaces                        │
│  • Well-defined interfaces for all components                │
│  • Built-in implementations (convention + config)            │
│  • Pluggable architecture (can swap implementations)         │
│  • No plugin infrastructure required initially               │
└────────────────────────┬────────────────────────────────────┘
                         │
        ┌────────────────┼────────────────┬────────────────┐
        ▼                ▼                ▼                ▼
┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐
│   Parser    │  │  Resolver   │  │  Executor   │  │  Reporter   │
│ Interface   │  │ Interface   │  │ Interface   │  │ Interface   │
└──────┬──────┘  └──────┬──────┘  └──────┬──────┘  └──────┬──────┘
       │                │                │                │
       ▼                ▼                ▼                ▼
┌─────────────┐  ┌─────────────┐  ┌─────────────┐  ┌─────────────┐
│   Built-in  │  │   Built-in  │  │   Built-in  │  │   Built-in  │
│     RST     │  │   Layered   │  │  Language   │  │   Human     │
│   Parser    │  │  Resolver   │  │  Executors  │  │  Reporter   │
└─────────────┘  └─────────────┘  └─────────────┘  └─────────────┘
```

### Key Architectural Principles

**1. Interface-Driven Design** (Plugin-Ready Without Plugin Overhead)

Define clear interfaces for all major components:

```typescript
// Core interfaces that enable future plugin system
interface Parser {
  parse(content: string, context: Context): Promise<ProcedureAST>;
}

interface PlaceholderResolver {
  resolve(placeholder: string, context: Context): Promise<string | null>;
  getSuggestions(placeholder: string, context: Context): string[];
}

interface Executor {
  canExecute(language: string): boolean;
  execute(code: string, context: ExecutionContext): Promise<ExecutionResult>;
}

interface Reporter {
  report(results: TestResults): Promise<void>;
}
```

**2. Built-in Implementations** (No Plugins Required)

Provide production-quality built-in implementations:

```typescript
// Built-in implementations work without any plugin system
class RSTParser implements Parser { /* ... */ }
class LayeredPlaceholderResolver implements PlaceholderResolver { /* ... */ }
class JavaScriptExecutor implements Executor { /* ... */ }
class HumanReporter implements Reporter { /* ... */ }

// Simple registration (no plugin infrastructure needed)
const components = {
  parser: new RSTParser(),
  resolver: new LayeredPlaceholderResolver([
    new ExactEnvResolver(),
    new FuzzyEnvResolver(),
    new SnootyConstantResolver()
  ]),
  executors: [
    new JavaScriptExecutor(),
    new PythonExecutor(),
    new PHPExecutor(),
    new ShellExecutor()
  ],
  reporters: [new HumanReporter()]
};
```

**3. Configuration-Based Customization** (Before Plugins)

Allow customization through configuration before requiring plugins:

```javascript
// .procedure-test.config.js
module.exports = {
  // Convention-based defaults work with zero config

  // Progressive configuration
  placeholders: {
    // Custom resolver function (no plugin needed)
    resolver: (placeholder, context) => {
      if (placeholder.match(/<cluster-\d+>/)) {
        return process.env.ATLAS_CLUSTER_NAME;
      }
      return null; // Fall back to built-in resolvers
    }
  },

  // Hooks for custom behavior (no plugin needed)
  hooks: {
    beforeEach: async (procedure) => {
      // Custom setup logic
    },
    afterEach: async (procedure, result) => {
      // Custom cleanup logic
    }
  },

  // Future: Plugin support (when needed)
  // plugins: ['./my-custom-plugin.js']
};
```

**4. Future Plugin System** (When Complexity Justifies It)

The architecture supports plugins when needed, but doesn't require them:

```javascript
// Future plugin support (Phase 3+)
module.exports = {
  // Existing config still works
  testFiles: ['procedures/**/*.txt'],

  // Add plugins when needed
  plugins: [
    // Custom parser for different format
    { type: 'parser', implementation: './custom-parser.js' },

    // Custom executor for new language
    { type: 'executor', implementation: '@company/kotlin-executor' },

    // Custom reporter for CI integration
    { type: 'reporter', implementation: '@company/ci-reporter' }
  ]
};
```

### Progressive Complexity Model

**Level 0: Zero Config** (90% of users)
```bash
# Just works
procedure-test myfile.txt
```

**Level 1: Environment Setup** (80% of users)
```bash
# Create .env file
echo "MONGODB_URI=mongodb://localhost" > .env
procedure-test myfile.txt
```

**Level 2: Basic Config** (30% of users)
```javascript
// .procedure-test.json
{
  "testFiles": ["procedures/**/*.txt"],
  "placeholders": {
    "<my-special-value>": "MY_ENV_VAR"
  }
}
```

**Level 3: Advanced Config** (10% of users)
```javascript
// .procedure-test.config.js
module.exports = {
  testFiles: ['procedures/**/*.txt'],
  placeholders: {
    resolver: (placeholder, context) => { /* custom logic */ }
  },
  hooks: {
    beforeEach: async (procedure) => { /* setup */ }
  }
};
```

**Level 4: Custom Plugins** (5% of users, future)
```javascript
// .procedure-test.config.js
module.exports = {
  extends: '@mongodb/procedure-test-preset',
  plugins: ['./my-custom-plugin.js']
};
```

### Implementation Phases

**Phase 1: PoC**
- ✅ Convention-based discovery
- ✅ Smart placeholder resolution
- ✅ Basic executors (JS, Python, PHP, Shell)
- ✅ Human-friendly reporter
- ✅ CLI with basic flags
- ✅ **Interface definitions** (plugin-ready architecture)
- ❌ No configuration file support yet
- ❌ No plugin system yet

**Phase 2: Production-Ready**
- ✅ Optional configuration file support
- ✅ Progressive configuration levels
- ✅ `--init` command for guided setup
- ✅ Configuration validation
- ✅ Optional hooks (beforeEach, afterEach, etc.)
- ✅ Multiple reporter formats (JSON, JUnit)
- ✅ **Component registry** (preparation for plugins)
- ❌ No plugin loading system yet

**Phase 3: Advanced Features**
- ✅ UI testing support (headless browser)
- ✅ Advanced cleanup strategies
- ✅ Parallel execution (optional)
- ✅ Preset configurations
- ✅ **Plugin system** (if needed based on user feedback)

**Phase 4: Ecosystem**
- ✅ Language-specific executors as plugins
- ✅ Custom reporter plugins
- ✅ Integration plugins (CI/CD, etc.)

### Why This Approach Wins

**Near-Term Benefits** (PoC & Initial Adoption):
- ✅ Writers can start immediately with zero config
- ✅ Fast implementation (no plugin infrastructure overhead)
- ✅ Simple mental model for users
- ✅ Easy to test and debug
- ✅ Quick wins demonstrate value

**Mid-Term Benefits** (Production Use):
- ✅ Configuration handles 95% of customization needs
- ✅ No plugin complexity for most users
- ✅ Hooks provide escape hatches for complex scenarios
- ✅ Interface-driven design keeps codebase clean
- ✅ Easy to add new built-in features

**Long-Term Benefits** (Scale & Extensibility):
- ✅ Plugin system available when complexity justifies it
- ✅ No architectural refactoring needed to add plugins
- ✅ Different teams can customize without forking
- ✅ Future-proof architecture

### Comparison: Hybrid vs. Hybrid + Plugin Ready

| Aspect | Option 3: Hybrid | Option 5: Hybrid + Plugin Ready |
|--------|------------------|----------------------------------|
| **Initial Complexity** | Low | Low (same) |
| **PoC Timeline** | 3-4 weeks | 3-4 weeks (same) |
| **Code Architecture** | Component-based | Interface-driven |
| **Extensibility** | Add features to core | Add features OR plugins |
| **Refactoring Risk** | Medium (if plugins needed) | Low (already architected) |
| **Community Contributions** | Requires core changes | Can contribute plugins |
| **Team Customization** | Fork or request features | Create custom plugins |
| **Long-term Maintenance** | May need refactoring | Stable architecture |

### Risk Mitigation

**Risk: Over-engineering for PoC**
- **Mitigation**: Interface definitions add minimal overhead (~1-2 days)
- **Benefit**: Prevents costly refactoring later
- **Evidence**: Interface-driven design is a best practice regardless of plugins

**Risk: Plugin system never needed**
- **Mitigation**: Don't build plugin loading until Phase 3+
- **Benefit**: Clean interfaces improve code quality anyway
- **Evidence**: YAGNI principle - only build when needed

**Risk: Writers confused by architecture**
- **Mitigation**: Writers never see interfaces/plugins unless they need them
- **Benefit**: User experience identical to Option 3
- **Evidence**: Implementation detail, not user-facing

### Decision Framework

Choose **Option 5 (Hybrid + Plugin Ready)** if:
- ✅ You want to avoid refactoring later
- ✅ You anticipate diverse customization needs
- ✅ You value long-term architectural stability
- ✅ You can invest a little extra time in PoC for interface design

Choose **Option 3 (Hybrid)** if:
- ✅ You want the absolute fastest PoC
- ✅ You're certain plugins will never be needed
- ✅ You're willing to refactor if needs change
- ✅ You want to minimize initial complexity

Choose **Option 1 (Convention)** if:
- ✅ Speed is the only priority
- ✅ This is a throwaway prototype
- ✅ You're willing to rebuild from scratch later

---

## Recommendation Summary

### Primary Recommendation: **Option 5 (Hybrid + Plugin Ready)**

**For MongoDB's procedural testing framework, this option provides the optimal balance of:**

1. **Immediate Usability** - Zero config for writers, just like Option 1
2. **Progressive Complexity** - Add configuration only when needed, like Option 3
3. **Future-Proof Architecture** - Plugin-ready without plugin overhead, inspired by Option 4
4. **Minimal Risk** - Small upfront investment prevents costly refactoring

### Implementation Recommendation

**Phase 1 (PoC) - 3-4 weeks**:
Focus on proving the concept with convention-based approach:
- Convention-based discovery (`.env`, `snooty.toml`)
- Smart placeholder resolution with fuzzy matching
- Basic executors for JS, Python, PHP, Shell
- Human-friendly error reporting
- **Define clean interfaces** (adds 1-2 days, prevents months of refactoring)

**Success Criteria for PoC**:
- Writers can test a simple procedure with zero config
- Placeholder resolution works for 80%+ of common patterns
- Error messages are actionable and helpful
- Tests execute reliably and clean up properly

**Phase 2 (Production) - 4-6 weeks**:
Add configuration support for complex scenarios:
- Optional configuration file (`.procedure-test.json` or `.procedure-test.js`)
- Progressive configuration levels
- Hooks for custom behavior
- Multiple output formats
- `--init` command for guided setup

**Success Criteria for Production**:
- 90% of writers use zero config
- 10% of writers use configuration for complex scenarios
- No requests for features that require architectural changes
- Writers report high satisfaction with usability

**Phase 3 (Advanced) - As Needed**:
Add advanced features based on user feedback:
- UI testing support
- Parallel execution
- Plugin system (only if justified by user needs)
- Preset configurations
- Advanced cleanup strategies

### Alternative Recommendations

**If speed is absolutely critical**: Start with **Option 1 (Convention)** for PoC, but plan to refactor to Option 5 for production.

**If you have limited resources**: **Option 3 (Hybrid)** provides most benefits with slightly less upfront architectural planning.

**If this is a long-term platform play**: **Option 4 (Plugin)** provides maximum extensibility, but delays time-to-value.

---

## Decision-Making Guide for Team Discussion

### Key Questions to Consider

When discussing these options with your team, consider these questions to guide your decision:

#### 1. Timeline & Resources

**Q: How quickly do you need a working PoC?**
- **1 week**: Option 1 (Convention)
- **2 weeks**: Option 5 (Hybrid + Plugin Ready) or Option 3 (Hybrid)
- **3 weeks**: Option 2 (Configuration) or Option 4 (Plugin)

**Q: How much development time can you invest upfront?**
- **Minimal**: Option 1
- **Moderate**: Option 3 or Option 5
- **Substantial**: Option 2 or Option 4

**Q: Do you have dedicated engineering resources or is this a side project?**
- **Side project**: Option 1 or Option 3 (simpler to maintain)
- **Dedicated resources**: Option 5 or Option 2 (can handle complexity)

#### 2. User Base & Adoption

**Q: What percentage of your technical writers are comfortable with configuration files?**
- **< 30%**: Prioritize Option 1 or Option 5 (zero config default)
- **30-70%**: Option 3 or Option 5 (progressive disclosure)
- **> 70%**: Option 2 (explicit configuration)

**Q: How diverse are the testing needs across different documentation projects?**
- **Mostly similar**: Option 1 (conventions work for everyone)
- **Somewhat varied**: Option 3 or Option 5 (config for edge cases)
- **Very diverse**: Option 2 or Option 4 (maximum customization)

**Q: Will different teams need different workflows?**
- **No, standardized**: Option 1 or Option 3
- **Yes, customized**: Option 5 or Option 4

#### 3. Long-Term Vision

**Q: Is this a one-time validation or a long-term tool?**
- **One-time validation**: Option 1 (fastest to prove concept)
- **Long-term tool**: Option 5 or Option 3 (sustainable architecture)

**Q: How likely is it that you'll need to support additional languages or formats?**
- **Unlikely**: Option 1 or 3
- **Likely**: Option 5 (extensible architecture)
- **Certain**: Option 4 (plugin system)

#### 4. Risk Tolerance

**Q: How comfortable are you with potential refactoring in 6-12 months?**
- **Very uncomfortable**: Option 5 (future-proof)
- **Somewhat uncomfortable**: Option 3 (balanced)
- **Comfortable**: Option 1 (iterate fast)

**Q: What's the cost of downtime or breaking changes for writers?**
- **Very high**: Option 5 (stable architecture)
- **Moderate**: Option 3 (careful evolution)
- **Low**: Option 1 (can iterate freely)

### Effort vs. Value Analysis

```
High Value
    │
    │   ┌─────────────┐
    │   │  Option 5   │  ← Recommended
    │   │  Hybrid +   │
    │   │   Plugin    │
    │   └─────────────┘
    │        ┌─────────────┐
    │        │  Option 3   │
    │        │   Hybrid    │
    │        └─────────────┘
    │   ┌─────────────┐
    │   │  Option 1   │
    │   │ Convention  │
    │   └─────────────┘
    │                    ┌─────────────┐
    │                    │  Option 2   │
    │                    │   Config    │
    │                    └─────────────┘
    │                                  ┌─────────────┐
    │                                  │  Option 4   │
    │                                  │   Plugin    │
    │                                  └─────────────┘
Low Value
    └────────────────────────────────────────────────────
         Low Effort                          High Effort
```

### Risk vs. Flexibility Analysis

```
High Flexibility
    │
    │                                  ┌─────────────┐
    │                                  │  Option 4   │
    │                                  │   Plugin    │
    │                                  └─────────────┘
    │                    ┌─────────────┐
    │                    │  Option 5   │  ← Recommended
    │                    │  Hybrid +   │
    │                    │   Plugin    │
    │                    └─────────────┘
    │        ┌─────────────┐
    │        │  Option 3   │
    │        │   Hybrid    │
    │        └─────────────┘
    │                    ┌─────────────┐
    │                    │  Option 2   │
    │                    │   Config    │
    │                    └─────────────┘
    │   ┌─────────────┐
    │   │  Option 1   │
    │   │ Convention  │
    │   └─────────────┘
Low Flexibility
    └────────────────────────────────────────────────────
         Low Risk                            High Risk
         (Refactoring)                       (Refactoring)
```

### Team Discussion Framework

Use this framework to structure your team discussion:

#### Round 1: Constraints
- What's our timeline for PoC?
- What's our timeline for production?
- What resources do we have?
- What's our risk tolerance?

#### Round 2: User Needs
- Who are our primary users?
- What's their technical skill level?
- How diverse are their needs?
- What's the adoption strategy?

#### Round 3: Long-Term Vision
- Is this a strategic tool or tactical solution?
- Will we need community contributions?
- How will requirements evolve?
- What's our maintenance capacity?

#### Round 4: Decision
- Which option best fits our constraints?
- Which option best serves our users?
- Which option aligns with our vision?
- What's our confidence level?

### Recommended Decision Process

1. **Review all options** (30 minutes)
   - Each team member reads the architecture document
   - Note questions and concerns

2. **Discuss constraints** (15 minutes)
   - Timeline, resources, risk tolerance
   - Identify hard constraints vs. preferences

3. **Discuss user needs** (15 minutes)
   - Writer skill levels and diversity
   - Adoption strategy and rollout plan

4. **Discuss long-term vision** (15 minutes)
   - Strategic importance
   - Evolution and extensibility needs

5. **Narrow to 2-3 options** (10 minutes)
   - Eliminate options that don't fit constraints
   - Focus on viable candidates

6. **Deep dive on finalists** (20 minutes)
   - Compare pros/cons in detail
   - Discuss implementation specifics
   - Consider migration paths

7. **Make decision** (10 minutes)
   - Vote or consensus
   - Document rationale
   - Identify success criteria

### Success Criteria by Option

**If you choose Option 1 (Convention)**:
- ✅ PoC completed in < 1 weeks
- ✅ 90%+ of test cases work with zero config
- ✅ Writers can run tests without training
- ⚠️ Plan for refactoring in 6-12 months

**If you choose Option 3 (Hybrid)**:
- ✅ PoC completed in ~2 weeks
- ✅ 90%+ of writers use zero config
- ✅ 10% of writers successfully use config for edge cases
- ✅ No major refactoring needed for 12+ months

**If you choose Option 5 (Hybrid + Plugin Ready)** ⭐:
- ✅ PoC completed in ~2 weeks
- ✅ 90%+ of writers use zero config
- ✅ 10% of writers successfully use config for edge cases
- ✅ Architecture supports plugins without refactoring
- ✅ No major refactoring needed for 24+ months
- ✅ Community contributions possible (if desired)

### Red Flags to Watch For

**During PoC**:
- 🚩 Writers need config for basic use cases → Architecture may be wrong
- 🚩 Common patterns require workarounds → Need better conventions
- 🚩 Error messages aren't actionable → Need better reporting
- 🚩 Tests are flaky or unreliable → Need better isolation

**During Production**:
- 🚩 Frequent requests for features requiring refactoring → Chose too simple
- 🚩 Writers struggling with configuration → Chose too complex
- 🚩 Different teams forking the tool → Need more flexibility
- 🚩 Maintenance burden increasing → Architecture issues

### Migration Paths

If you need to change approaches later:

**Option 1 → Option 3**: Moderate effort
- Add configuration layer
- Maintain backward compatibility with conventions
- Estimated: 2-3 weeks

**Option 1 → Option 5**: Significant effort
- Refactor to interface-driven design
- Add configuration layer
- Estimated: 4-6 weeks

**Option 3 → Option 5**: Low effort
- Define interfaces for existing components
- Minimal code changes
- Estimated: 1-2 weeks

**Option 5 → Option 4**: Low effort
- Add plugin loading system
- Already architected for plugins
- Estimated: 1-2 weeks

---

## Concrete Implementation Examples

To help visualize the differences, here are concrete examples of how each option would handle a common scenario.

### Scenario: Writer Needs to Test a Procedure with Custom Placeholder

**The Procedure** (RST):
```rst
.. procedure::

   .. step:: Connect to your cluster

      Replace ``<cluster-id>`` with your Atlas cluster identifier:

      .. code-block:: shell

         mongosh "mongodb+srv://<username>:<password>@<cluster-id>.mongodb.net"
```

**The Challenge**: The placeholder `<cluster-id>` doesn't match the environment variable `ATLAS_CLUSTER_NAME`.

---

### Option 1: Convention-Over-Configuration

**Writer Experience**:
```bash
$ procedure-test connect-procedure.txt

❌ Test Failed: connect-procedure.txt
Error: Could not resolve placeholder: <cluster-id>

Suggestions:
  • Add CLUSTER_ID to your .env file
  • Rename ATLAS_CLUSTER_NAME to CLUSTER_ID in .env
  • Similar environment variables found: ATLAS_CLUSTER_NAME
```

**Solution**: Writer renames their environment variable
```bash
# .env
CLUSTER_ID=my-cluster-123
```

**Pros**: Simple, no configuration needed
**Cons**: Writer must adapt to tool's conventions

---

### Option 2: Configuration-First

**Writer Experience**:
```bash
$ procedure-test --config procedure-test.config.js connect-procedure.txt

✅ Test Passed: connect-procedure.txt
```

**Solution**: Writer creates configuration file
```javascript
// procedure-test.config.js
module.exports = {
  testFiles: ['connect-procedure.txt'],
  envFiles: ['.env'],
  placeholders: {
    '<cluster-id>': 'ATLAS_CLUSTER_NAME'
  }
};
```

**Pros**: Explicit mapping, no environment changes needed
**Cons**: Requires configuration file for simple case

---

### Option 3: Hybrid

**Writer Experience (First Attempt)**:
```bash
$ procedure-test connect-procedure.txt

❌ Test Failed: connect-procedure.txt
Error: Could not resolve placeholder: <cluster-id>

Suggestions:
  • Add CLUSTER_ID to your .env file
  • Add to config: { "placeholders": { "<cluster-id>": "ATLAS_CLUSTER_NAME" } }
  • Similar environment variables found: ATLAS_CLUSTER_NAME

Run 'procedure-test --init' to create a configuration file.
```

**Solution Option A**: Rename environment variable (like Option 1)
```bash
# .env
CLUSTER_ID=my-cluster-123
```

**Solution Option B**: Create minimal config
```json
{
  "placeholders": {
    "<cluster-id>": "ATLAS_CLUSTER_NAME"
  }
}
```

**Pros**: Works with zero config for most cases, config available when needed
**Cons**: Two ways to solve the problem might be confusing

---

### Option 5: Hybrid + Plugin Ready

**Writer Experience (Identical to Option 3)**:
```bash
$ procedure-test connect-procedure.txt

❌ Test Failed: connect-procedure.txt
Error: Could not resolve placeholder: <cluster-id>

Suggestions:
  • Add CLUSTER_ID to your .env file
  • Add to config: { "placeholders": { "<cluster-id>": "ATLAS_CLUSTER_NAME" } }
  • Similar environment variables found: ATLAS_CLUSTER_NAME

Run 'procedure-test --init' to create a configuration file.
```

**Solution (Same as Option 3)**: Create minimal config
```json
{
  "placeholders": {
    "<cluster-id>": "ATLAS_CLUSTER_NAME"
  }
}
```

**Future Enhancement (No Refactoring Needed)**:
If the team later needs a custom resolver for complex patterns:

```javascript
// custom-resolver-plugin.js
module.exports = {
  name: 'atlas-cluster-resolver',
  type: 'resolver',

  resolve(placeholder, context) {
    // Custom logic for Atlas cluster patterns
    if (placeholder.match(/<cluster-\w+>/)) {
      return process.env.ATLAS_CLUSTER_NAME;
    }
    return null;
  }
};

// procedure-test.config.js
module.exports = {
  plugins: ['./custom-resolver-plugin.js']
};
```

**Pros**: Same user experience as Option 3, but architecture supports plugins without refactoring
**Cons**: Slightly more complex implementation (1-2 extra days)

---

### Side-by-Side Comparison: Adding a New Language

**Scenario**: Team needs to add Kotlin support

#### Option 1: Convention
```typescript
// Must modify core code
// src/executors/kotlin-executor.ts
export class KotlinExecutor extends BaseExecutor {
  // Implementation
}

// src/executors/index.ts
import { KotlinExecutor } from './kotlin-executor';
executors.push(new KotlinExecutor());
```
**Effort**: Modify core, rebuild, redeploy
**Timeline**: 1-2 days + deployment

#### Option 3: Hybrid
```typescript
// Must modify core code (same as Option 1)
// src/executors/kotlin-executor.ts
export class KotlinExecutor extends BaseExecutor {
  // Implementation
}

// src/executors/index.ts
import { KotlinExecutor } from './kotlin-executor';
executors.push(new KotlinExecutor());
```
**Effort**: Modify core, rebuild, redeploy
**Timeline**: 1-2 days + deployment

#### Option 5: Hybrid + Plugin Ready
```typescript
// Option A: Add to core (same as above)
// Option B: Create plugin (no core changes)

// kotlin-executor-plugin.js
module.exports = {
  name: 'kotlin-executor',
  type: 'executor',

  canExecute(language) {
    return language === 'kotlin';
  },

  async execute(code, context) {
    // Implementation
  }
};

// procedure-test.config.js
module.exports = {
  plugins: ['./kotlin-executor-plugin.js']
};
```
**Effort**: Create plugin file, update config
**Timeline**: 1-2 days, no deployment needed
**Benefit**: Different teams can use different plugins

---

### Real-World Scenario: Multiple Documentation Teams

**Situation**: MongoDB has 3 documentation teams with different needs:
- **Atlas Team**: Needs UI testing for Atlas console
- **Drivers Team**: Needs support for 10+ programming languages
- **Server Team**: Needs custom mongosh validation

#### Option 1: Convention
- ❌ All teams must use same conventions
- ❌ Edge cases require workarounds
- ❌ Teams request features, creating backlog
- ❌ Core becomes bloated with team-specific code

#### Option 3: Hybrid
- ✅ Teams can use config for customization
- ⚠️ Complex needs require core changes
- ⚠️ Teams must wait for features to be added
- ⚠️ Risk of core becoming bloated

#### Option 5: Hybrid + Plugin Ready
- ✅ Teams can use config for simple customization
- ✅ Teams can create plugins for complex needs
- ✅ Plugins don't affect other teams
- ✅ Core stays focused and maintainable
- ✅ Teams can share plugins if desired

**Example**:
```javascript
// Atlas team config
module.exports = {
  plugins: [
    '@mongodb/procedure-test-ui-testing'  // Atlas-specific plugin
  ]
};

// Drivers team config
module.exports = {
  plugins: [
    '@mongodb/procedure-test-executor-kotlin',
    '@mongodb/procedure-test-executor-rust',
    '@mongodb/procedure-test-executor-swift'
  ]
};

// Tools team config
module.exports = {
  plugins: [
    './custom-cli-validator.js'  // Team-specific plugin
  ]
};
```

---

## Cost-Benefit Analysis

### Development Cost (PoC Phase)

| Option | Initial Dev Time | Complexity | Risk |
|--------|-----------------|------------|------|
| Option 1 | 1-2 weeks | Low | High refactoring risk |
| Option 2 | 3-4 weeks | Medium | Medium refactoring risk |
| Option 3 | 3-4 weeks | Medium | Medium refactoring risk |
| Option 4 | 6-8 weeks | High | Low refactoring risk (overbuilt) |
| Option 5 | 3-4 weeks (+1-2 days) | Medium | Low refactoring risk |

### Maintenance Cost (Annual)

| Option | Config Maintenance | Code Maintenance | Support Burden | Total |
|--------|-------------------|------------------|----------------|-------|
| Option 1 | Low | High (many core changes) | High (workarounds) | **High** |
| Option 2 | High (complex configs) | Medium | Medium | **Medium-High** |
| Option 3 | Low-Medium | Medium | Low-Medium | **Medium** |
| Option 4 | Medium | Low (plugins isolated) | Medium (plugin system) | **Medium** |
| Option 5 | Low-Medium | Low (plugins optional) | Low | **Low-Medium** |

### Total Cost of Ownership (3 Years)

| Option | Year 1 | Year 2 | Year 3 | Total | Notes |
|--------|--------|--------|--------|-------|-------|
| Option 1 | Low | High (refactor) | Medium | **High** | Refactoring in Year 2 |
| Option 2 | Medium | Medium | Medium | **Medium** | Stable but config-heavy |
| Option 3 | Medium | Medium | Medium-High | **Medium-High** | May need refactor Year 3 |
| Option 4 | High | Low | Low | **Medium** | High upfront, low ongoing |
| Option 5 | Medium | Low | Low | **Low-Medium** | Best long-term value |

### Value Delivered

| Option | Writer Satisfaction | Flexibility | Extensibility | Overall Value |
|--------|-------------------|-------------|---------------|---------------|
| Option 1 | High (simple) | Low | Low | **Medium** |
| Option 2 | Medium (complex) | High | Medium | **Medium-High** |
| Option 3 | High (progressive) | High | Medium | **High** |
| Option 4 | Medium (complex) | Very High | Very High | **High** (if needed) |
| Option 5 | High (progressive) | Very High | Very High | **Very High** |

---

## Additional Architectural Considerations

Regardless of which option you choose, these architectural decisions apply to all approaches:

### 1. State Management Architecture

**Challenge**: Code blocks within a step need to share state, but procedures should be isolated.

**Proposed Solution**:
```
┌─────────────────────────────────────────────────────────────┐
│                    Procedure Context                         │
│  • Unique ID                                                 │
│  • Isolated environment                                      │
│  • Cleanup registry                                          │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│                     Step Context                             │
│  • Step number and description                               │
│  • Language-specific state accumulators                      │
│  • Variable registry                                         │
└────────────────────────┬────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────┐
│              Language State Accumulator                      │
│  JavaScript: Accumulate code, execute in same context        │
│  Python: Accumulate code, execute in same interpreter        │
│  PHP: Accumulate code, execute in same process               │
│  Shell: Maintain shell session                               │
└─────────────────────────────────────────────────────────────┘
```

**Implementation Approaches**:

**JavaScript**:
```javascript
// Accumulate code within a step
const stepCode = [];
stepCode.push('const client = new MongoClient(uri);');
stepCode.push('await client.connect();');

// Execute as single script
const fullCode = stepCode.join('\n');
await executeJavaScript(fullCode);
```

**Python**:
```python
# Use exec() with shared globals
step_globals = {}
exec('client = MongoClient(uri)', step_globals)
exec('db = client.get_database("test")', step_globals)
```

**Shell**:
```bash
# Maintain persistent shell session
# Write commands to temp file, execute with source
```

### 2. Error Handling and Reporting Architecture

**Goal**: Provide actionable error messages that help writers fix issues quickly.

**Error Context Structure**:
```typescript
interface TestError {
  // Location
  file: string;
  line: number;
  procedure: string;
  step: number;
  stepDescription: string;

  // Error details
  type: 'parse' | 'execution' | 'placeholder' | 'cleanup';
  message: string;
  originalError: Error;

  // Context
  codeSnippet: string;
  language: string;

  // Suggestions
  suggestions: string[];
  relatedDocs: string[];
}
```

**Error Reporter Output**:
```
❌ Test Failed: testdata/atlas/source/connect-to-database-deployment.txt

Procedure: "Connect to Database Deployment"
Step 2: "Connect using mongosh"
Line 78-82

Code:
  mongosh "mongodb+srv://<username>:<password>@<cluster>.mongodb.net"

Error: Could not resolve placeholder: <cluster>

Suggestions:
  • Add CLUSTER to your .env file
  • Add "cluster": "ATLAS_CLUSTER_NAME" to .procedure-test.json
  • Similar environment variables found: ATLAS_CLUSTER, CLUSTER_NAME

Documentation: https://docs.example.com/procedure-test/placeholders
```

### 3. Execution Isolation Architecture

**Challenge**: Prevent tests from interfering with each other and the writer's environment.

**Isolation Strategies**:

**Database Isolation**:
```javascript
// Generate unique database names per test run
const testRunId = generateUniqueId();
const dbName = `test_${testRunId}_${procedureId}`;

// Track for cleanup
context.cleanup.register('database', dbName);
```

**File System Isolation**:
```javascript
// Create temp directory per procedure
const tempDir = `/tmp/procedure-test/${testRunId}/${procedureId}`;
process.chdir(tempDir);

// Track for cleanup
context.cleanup.register('directory', tempDir);
```

**Environment Isolation**:
```javascript
// Clone environment, add test-specific vars
const testEnv = {
  ...process.env,
  NODE_ENV: 'test',
  MONGODB_DATABASE: dbName
};

// Execute with isolated environment
await execute(code, { env: testEnv });
```

### 4. Cleanup Architecture

**Challenge**: Reliably clean up resources even when tests fail.

**Cleanup Registry Pattern**:
```typescript
class CleanupRegistry {
  private cleanupTasks: CleanupTask[] = [];

  register(type: string, resource: any, cleanup: () => Promise<void>) {
    this.cleanupTasks.push({ type, resource, cleanup });
  }

  async executeAll() {
    // Execute in reverse order (LIFO)
    for (const task of this.cleanupTasks.reverse()) {
      try {
        await task.cleanup();
      } catch (error) {
        // Log but don't throw - try to clean up everything
        console.warn(`Cleanup failed for ${task.type}:`, error);
      }
    }
  }
}

// Usage
context.cleanup.register('database', dbName, async () => {
  await client.db(dbName).dropDatabase();
});

context.cleanup.register('collection', collName, async () => {
  await db.collection(collName).drop();
});

// Always execute cleanup
try {
  await executeProcedure(procedure);
} finally {
  await context.cleanup.executeAll();
}
```

**Auto-Detection Strategy**:
```javascript
// Track MongoDB operations
const mongoProxy = new Proxy(mongoClient, {
  get(target, prop) {
    if (prop === 'db') {
      return (name) => {
        // Register database for cleanup
        context.cleanup.register('database', name, ...);
        return target.db(name);
      };
    }
    return target[prop];
  }
});
```

### 5. Placeholder Resolution Architecture

**Multi-Strategy Resolver**:
```typescript
interface PlaceholderResolver {
  resolve(placeholder: string, context: Context): string | null;
  priority: number;
}

class PlaceholderResolutionEngine {
  private resolvers: PlaceholderResolver[] = [];

  addResolver(resolver: PlaceholderResolver) {
    this.resolvers.push(resolver);
    this.resolvers.sort((a, b) => b.priority - a.priority);
  }

  resolve(placeholder: string, context: Context): string {
    for (const resolver of this.resolvers) {
      const result = resolver.resolve(placeholder, context);
      if (result !== null) {
        return result;
      }
    }

    throw new UnresolvedPlaceholderError(
      placeholder,
      this.getSuggestions(placeholder, context)
    );
  }

  private getSuggestions(placeholder: string, context: Context): string[] {
    // Fuzzy match against environment variables
    // Suggest similar placeholders from placeholder-consistency.md
    // Suggest configuration options
  }
}

// Built-in resolvers
class ExactEnvResolver implements PlaceholderResolver {
  priority = 100;
  resolve(placeholder: string, context: Context): string | null {
    const envKey = placeholder.replace(/[<>]/g, '').toUpperCase();
    return process.env[envKey] || null;
  }
}

class FuzzyEnvResolver implements PlaceholderResolver {
  priority = 50;
  resolve(placeholder: string, context: Context): string | null {
    // Normalize: <connection-string> -> CONNECTION_STRING
    // Try variations: connectionString, connection_string, etc.
  }
}

class SnootyConstantResolver implements PlaceholderResolver {
  priority = 75;
  resolve(placeholder: string, context: Context): string | null {
    // Resolve {+constant+} from snooty.toml
  }
}
```

### 6. RST Parser Architecture

**Streaming Parser vs. AST Builder**:

**Option A: AST Builder** (Recommended)
```typescript
interface ASTNode {
  type: string;
  children: ASTNode[];
  metadata: Record<string, any>;
}

interface ProcedureNode extends ASTNode {
  type: 'procedure';
  title: string;
  style: 'normal' | 'connected';
  steps: StepNode[];
}

interface StepNode extends ASTNode {
  type: 'step';
  headline: string;
  content: ContentNode[];
}

// Build full AST, then traverse
const ast = parseRST(content);
const procedures = extractProcedures(ast);
```

**Benefits**:
- Easier to implement complex logic (tabs, composable tutorials)
- Can traverse multiple times for different purposes
- Easier to debug and test
- Better error messages with full context

**Option B: Streaming Parser**
```typescript
// Parse and execute on-the-fly
for await (const node of parseRSTStream(content)) {
  if (node.type === 'code-block') {
    await executeCodeBlock(node);
  }
}
```

**Benefits**:
- Lower memory usage
- Faster for simple cases
- Can start execution sooner

**Recommendation**: Use AST Builder (Option A) for better maintainability and debugging.

### 7. Dependency Detection Architecture

**Challenge**: Determine what tools/languages are needed before execution.

**Static Analysis Approach**:
```typescript
interface ProcedureDependencies {
  languages: Set<string>;        // ['javascript', 'shell']
  cliTools: Set<string>;         // ['mongosh', 'atlas-cli']
  runtimes: Set<string>;         // ['node', 'python3']
  packages: Map<string, string[]>; // { javascript: ['mongodb'] }
}

function analyzeDependencies(procedure: ProcedureNode): ProcedureDependencies {
  const deps = new ProcedureDependencies();

  // Walk AST and collect dependencies
  for (const step of procedure.steps) {
    for (const codeBlock of step.codeBlocks) {
      deps.languages.add(codeBlock.language);

      // Detect CLI tools from code content
      if (codeBlock.language === 'shell') {
        if (codeBlock.code.includes('mongosh')) {
          deps.cliTools.add('mongosh');
        }
        if (codeBlock.code.includes('atlas')) {
          deps.cliTools.add('atlas-cli');
        }
      }

      // Detect package imports
      if (codeBlock.language === 'javascript') {
        const imports = extractImports(codeBlock.code);
        deps.packages.get('javascript').push(...imports);
      }
    }
  }

  return deps;
}

// Validate before execution
async function validateEnvironment(deps: ProcedureDependencies) {
  const missing = [];

  for (const tool of deps.cliTools) {
    if (!await isInstalled(tool)) {
      missing.push(tool);
    }
  }

  if (missing.length > 0) {
    throw new MissingDependenciesError(missing);
  }
}
```

---

## Next Steps

Once you've selected an architecture option, the next steps would be:

1. **Create detailed technical specification** for the chosen architecture
2. **Design component interfaces** and APIs
3. **Create data flow diagrams** showing how data moves through the system
4. **Define error handling strategies** for each component
5. **Plan testing strategy** for the framework itself
6. **Build proof-of-concept** implementing core functionality

---

## Final Recommendation & Next Steps

### The Clear Winner: Option 5 (Hybrid + Plugin Ready) ⭐

After comprehensive analysis across multiple dimensions, **Option 5** emerges as the optimal choice for MongoDB's procedural testing framework.

#### Why Option 5 Wins

**1. Best User Experience**
- Zero config for 90% of users (identical to Option 1)
- Progressive disclosure for complex cases (identical to Option 3)
- Helpful error messages with actionable suggestions
- Minimal learning curve for technical writers

**2. Best Long-Term Value**
- Small upfront investment (1-2 extra days) prevents months of refactoring
- Plugin-ready architecture without plugin overhead
- Stable foundation for 24+ months
- Lowest total cost of ownership over 3 years

**3. Best Flexibility**
- Conventions work for common cases
- Configuration handles edge cases
- Plugins enable team-specific customization (future)
- Community contributions possible without core changes

**4. Best Risk Profile**
- Minimal refactoring risk
- No architectural dead-ends
- Can evolve naturally as needs grow
- Proven pattern (used by ESLint, Babel, Webpack, etc.)

#### The Investment

**Upfront Cost**: 1-2 extra days during PoC to define clean interfaces
**Benefit**: Avoid 4-6 weeks of refactoring in 6-12 months
**ROI**: 20-30x return on investment

#### Success Metrics

**PoC Success** (Week 4):
- ✅ 90%+ of test procedures work with zero config
- ✅ Writers can run tests without training
- ✅ Error messages are actionable
- ✅ Tests execute reliably and clean up properly

**Production Success** (Month 3):
- ✅ 90%+ of writers never create a config file
- ✅ 10% of writers successfully use config for edge cases
- ✅ No requests for features requiring architectural changes
- ✅ High writer satisfaction scores

**Long-Term Success** (Year 1+):
- ✅ Architecture supports new features without refactoring
- ✅ Different teams can customize without forking
- ✅ Community contributions (if desired)
- ✅ Maintenance burden remains low

### Immediate Next Steps

1. **Team Decision** (This Week)
   - Review this document with your team
   - Discuss using the Decision-Making Guide (page 1251)
   - Confirm Option 5 or select alternative
   - Document decision rationale

2. **Technical Specification** (Week 1-2)
   - Create detailed tech spec for chosen option
   - Define component interfaces
   - Design data flow
   - Plan error handling strategy
   - Identify technical risks

3. **Proof of Concept** (Week 3-6)
   - Implement core parsing (RST → AST)
   - Build smart placeholder resolver
   - Create basic executors (JS, Python, PHP, Shell)
   - Implement human-friendly reporter
   - Add CLI with basic flags
   - **Define clean interfaces** (Option 5 only)

4. **Validation** (Week 7-8)
   - Test with real documentation files
   - Gather feedback from technical writers
   - Iterate on error messages
   - Refine conventions based on usage
   - Validate success criteria

5. **Production Readiness** (Week 9-12)
   - Add optional configuration support
   - Implement `--init` command
   - Create comprehensive documentation
   - Add multiple output formats
   - Prepare for rollout

### Questions for Your Team

Before proceeding, ensure alignment on:

1. **Timeline**: Is 3-4 weeks acceptable for PoC? (vs. 1-2 weeks for Option 1)
2. **Investment**: Is 1-2 extra days for interface design worth avoiding future refactoring?
3. **Vision**: Is this a long-term strategic tool or short-term tactical solution?
4. **Flexibility**: Will different teams need different customizations?
5. **Community**: Do you want to enable community contributions in the future?

If the answer to questions 3, 4, or 5 is "yes," Option 5 is the clear choice.
If all answers are "no" and timeline is critical, Option 1 may be sufficient.

### Ready to Proceed?

I'm ready to create a detailed technical specification for **Option 5 (Hybrid + Plugin Ready)** that includes:

- **Component Architecture**: Detailed design of all components with interfaces
- **Data Flow Diagrams**: How data moves through the system
- **API Specifications**: Interface definitions for all major components
- **Implementation Plan**: Phased approach with milestones
- **Testing Strategy**: How to test the framework itself
- **Risk Mitigation**: Technical risks and mitigation strategies
- **Code Examples**: Sample implementations of key components

Would you like me to proceed with the technical specification for Option 5, or would you prefer to discuss any aspects of the architecture options in more detail first?



# Procedural Testing Framework - Usage Guide

**For MongoDB Documentation Team**

This guide shows you how to use the procedural testing framework to validate that your documentation procedures work as written.

---

## Table of Contents

1. [Quick Start](#quick-start)
2. [Local Development](#local-development)
3. [Understanding Test Results](#understanding-test-results)
4. [Configuration](#configuration)
5. [CI Integration](#ci-integration)
6. [Monorepo Conventions](#monorepo-conventions)
7. [Common Patterns](#common-patterns)
8. [Troubleshooting](#troubleshooting)

---

## Quick Start

### Installation

```bash
# Install globally (recommended for local development)
npm install -g @grove-platform/proctest

# Or install as dev dependency in your docs repo
npm install --save-dev @grove-platform/proctest
```

### Your First Test

The framework works with **zero configuration** for most cases. Just point it at your documentation file:

```bash
# Test a single file
proctest test source/tutorial/getting-started.txt

# Test all files in a directory
proctest test source/tutorial/

# Test specific files matching a pattern
proctest test "source/tutorial/install-*.txt"
```

That's it! The framework will:
- ✅ Parse your RST file
- ✅ Extract procedures and code blocks
- ✅ Resolve placeholders from `.env` and `snooty.toml`
- ✅ Execute the code
- ✅ Clean up test resources automatically

---

## Local Development

### Setting Up Your Environment

Create a `.env` file in your repository root with your test credentials:

```bash
# .env
MONGODB_URI=mongodb://localhost:27017
# More credentials as needed...
```

**Important**: Add `.env` to your `.gitignore` to avoid committing secrets!

### Running Tests

```bash
# Test a single procedure
proctest test source/tutorial/create-database.txt

# Test with verbose output
proctest test source/tutorial/create-database.txt --verbose

# Test without cleanup (useful for debugging)
proctest test source/tutorial/create-database.txt --no-cleanup

# Test and output JSON for programmatic use
proctest test source/tutorial/ --reporter json > results.json
```

### Debugging with the Parse Command

Before running tests, you can preview what the framework will extract:

```bash
# See the procedure structure as a tree
proctest parse source/tutorial/getting-started.txt

# Output as JSON for detailed inspection
proctest parse source/tutorial/getting-started.txt --format json

# Output as YAML
proctest parse source/tutorial/getting-started.txt --format yaml

# Save to a file
proctest parse source/tutorial/getting-started.txt --format json --output debug.json
```

**Example output (tree format)**:

```
Procedure: Install MongoDB Driver
├─ Step 1: Create a new project directory
│  └─ Shell: mkdir myproject && cd myproject
├─ Step 2: Initialize npm
│  └─ Shell: npm init -y
└─ Step 3: Install the MongoDB driver
   └─ Shell: npm install mongodb

Placeholders found:
  - {+api-key+} (2 occurrences)
  - {+project-id+} (1 occurrence)

Summary:
  Procedures: 1
  Steps: 3
  Testable Actions: 3 (3 Shell)
```

This helps you verify the framework is parsing your documentation correctly before running tests,
and gives you the information you need to populate any placeholders in your `.env`.

---

## Understanding Test Results

### Human-Readable Output (Default)

```
Testing: source/tutorial/install-driver.txt

✓ Install MongoDB Driver (Python)
  Duration: 2.3s
  Steps: 3/3 passed

✗ Install MongoDB Driver (Node.js)
  Duration: 1.8s
  Steps: 2/3 passed

  Error in Step 3: Install the MongoDB driver
    Command failed: npm install mongodb
    Exit code: 1

    Suggestion: Check that npm is installed and accessible

✓ Install MongoDB Driver (Java)
  Duration: 3.1s
  Steps: 3/3 passed

Summary:
  3 procedures tested (from 1 base procedure with variants)
  2 passed, 1 failed
  Total duration: 7.2s
```

### Exit Codes

- `0` - All tests passed
- `1` - One or more tests failed
- `2` - Configuration or parsing error

Use exit codes in scripts:

```bash
if proctest test source/tutorial/; then
  echo "All tests passed!"
else
  echo "Tests failed"
  exit 1
fi
```

---

## Configuration

### When Do You Need Configuration?

**90% of procedures work with zero configuration.** You only need a config file for:

- Custom UI navigation mappings
- Team-specific placeholder values
- Custom cleanup strategies
- Non-standard file patterns

### Creating a Configuration File

Create `.proctest.js` in your repository root:

```javascript
// .proctest.js
module.exports = {
  // Custom test file patterns (optional)
  testMatch: [
    'source/tutorial/**/*.txt',
    'source/guides/**/*.txt'
  ],

  // Exclude certain files (optional)
  testIgnore: [
    'source/includes/**',
    'source/tutorial/draft-*.txt'
  ],

  // UI navigation mappings (for UI testable actions)
  ui: {
    navigationMappings: {
      'click the {button} button': async (page, button) => {
        await page.click(`button:has-text("${button}")`);
      },
      'navigate to {page}': async (page, pageName) => {
        const routes = {
          'Database': '/databases',
          'Clusters': '/clusters',
          'Users': '/security/users'
        };
        await page.goto(`https://cloud.mongodb.com${routes[pageName]}`);
      }
    },

    // User-specific values for testing
    userValues: {
      'your project': 'test-project-123',
      'your cluster': 'test-cluster-0',
      'your database': 'test_db'
    }
  },

  // Cleanup configuration (optional)
  cleanup: {
    enabled: true, // Default
    databases: { enabled: true },
    collections: { enabled: true },
    files: { enabled: true }
  }
};
```

### JSON Configuration Alternative

If you prefer JSON, create `.proctest.json`:

```json
{
  "testMatch": [
    "source/tutorial/**/*.txt",
    "source/guides/**/*.txt"
  ],
  "testIgnore": [
    "source/includes/**",
    "source/tutorial/draft-*.txt"
  ],
  "cleanup": {
    "enabled": true
  }
}
```

**Note**: JavaScript config files are more powerful (you can define functions for UI mappings), but JSON works for simple configuration.

---

## CI Integration

### GitHub Actions

Create `.github/workflows/test-procedures.yml`:

```yaml
name: Test Documentation Procedures

on:
  pull_request:
    paths:
      - 'source/**/*.txt'
      - 'source/**/*.rst'
  push:
    branches:
      - main

jobs:
  test-procedures:
    runs-on: ubuntu-latest

    services:
      mongodb:
        image: mongodb/mongodb-community-server:latest
        ports:
          - 27017:27017

    steps:
      - name: Checkout code
        uses: actions/checkout@v4

      - name: Setup Node.js
        uses: actions/setup-node@v4
        with:
          node-version: '24'

      - name: Install proctest
        run: npm install -g @grove-platform/proctest

      - name: Create .env file
        run: |
          cat > .env << EOF
          MONGODB_URI=mongodb://localhost:27017
          EOF

      - name: Run procedure tests
        run: proctest test source/tutorial/ --reporter junit --output test-results.xml

      - name: Publish test results
        uses: EnricoMi/publish-unit-test-result-action@v2
        if: always()
        with:
          files: test-results.xml
```

### CI Best Practices

1. **Use JUnit XML output** for CI integration:
   ```bash
   proctest test source/ --reporter junit --output results.xml
   ```

2. **Set up MongoDB service** in CI for database tests

3. **Store secrets in CI environment variables**, not in code:
   - `ATLAS_PROJECT_ID`

4. **Run tests on documentation changes** to catch broken procedures early

5. **Use test results to block PRs** with failing procedures

---

## Monorepo Conventions

For teams working in a monorepo with multiple documentation projects, establish these conventions:

### Directory Structure

```
mongodb-docs/
├── .github/
│   └── workflows/
│       └── test-procedures.yml          # CI workflow
├── code-example-tests/
│   └── procedures/                      # Shared procedural testing content
│       ├── navigation-mappings.js       # Shared UI navigation mappings
│       ├── base-config.js               # Base configuration
│       ├── env.template                 # Template for .env file
│       └── test-registry.json           # Registry of verified procedures
├── content/
│   ├── atlas/
│   │   ├── source/                      # Atlas documentation
│   │   ├── snooty.toml
│   │   ├── .env                         # Local only (gitignored)
│   │   └── .proctest.js                 # Extends shared config (optional)
│   └── drivers/
│       ├── source/                      # Drivers documentation
│       ├── snooty.toml
│       ├── .env                         # Local only (gitignored)
│       └── .proctest.js                 # Extends shared config (optional)
└── .gitignore                           # Ignore all .env files
```

### Shared Configuration

**`code-example-tests/procedures/base-config.js`**:

```javascript
// Base configuration shared across all projects
module.exports = {
  cleanup: {
    enabled: true,
    databases: { enabled: true },
    collections: { enabled: true },
    files: { enabled: true }
  },

  // Common test patterns
  testIgnore: [
    '**/includes/**',
    '**/draft-*',
    '**/*.template.txt'
  ]
};
```

**`code-example-tests/procedures/navigation-mappings.js`**:

```javascript
// Shared UI navigation mappings for Atlas UI
module.exports = {
  'click the {button} button': async (page, button) => {
    await page.click(`button:has-text("${button}")`);
  },

  'navigate to the {page} page': async (page, pageName) => {
    const routes = {
      'Database': '/databases',
      'Clusters': '/clusters',
      'Users': '/security/users',
      'Network Access': '/security/network',
      'Database Access': '/security/database'
    };
    await page.goto(`https://cloud.mongodb.com${routes[pageName]}`);
  },

  'select {option} from the {dropdown} dropdown': async (page, option, dropdown) => {
    await page.selectOption(`select[aria-label="${dropdown}"]`, option);
  }
};
```

**`code-example-tests/procedures/env.template`**:

```bash
# Template for .env file
# Copy this to .env in your project root and fill in your values

# Local MongoDB
MONGODB_URI=mongodb://localhost:27017

# Atlas API credentials
ATLAS_PUBLIC_KEY=your-public-key-here
ATLAS_PRIVATE_KEY=your-private-key-here
ATLAS_PROJECT_ID=your-project-id-here

# Optional: Atlas cluster for testing
ATLAS_CLUSTER_NAME=test-cluster-0
```

### Project-Specific Configuration

You can optionally share configuration across a specific project by committing
a shared config file for the project:

**`content/atlas/.proctest.js`**:

```javascript
const baseConfig = require('../../code-example-tests/procedures/base-config');
const sharedNavigationMappings = require('../../code-example-tests/procedures/navigation-mappings');

module.exports = {
  ...baseConfig,

  testMatch: [
    'source/tutorial/**/*.txt',
    'source/how-to/**/*.txt'
  ],

  ui: {
    navigationMappings: {
      ...sharedNavigationMappings,

      // Atlas-specific mappings
      'create a new cluster': async (page) => {
        await page.click('button:has-text("Create")');
        await page.click('text=Deploy a cluster');
      }
    },

    userValues: {
      'your project': process.env.ATLAS_PROJECT_ID || 'test-project',
      'your cluster': process.env.ATLAS_CLUSTER_NAME || 'test-cluster-0'
    }
  }
};
```

### Sharing Navigation Mappings

Teams can contribute to shared navigation mappings:

1. **Add new mappings** to `code-example-tests/procedures/navigation-mappings.js`
2. **Test locally** in your project
3. **Submit PR** for review
4. **All projects benefit** from the new mapping

This creates a **library of reusable UI automation** that grows over time.

### Test Registry for CI

The test registry (`code-example-tests/procedures/test-registry.json`) maintains a list of procedures that have been verified and should run in CI:

**`code-example-tests/procedures/test-registry.json`**:

```json
{
  "version": "1.0",
  "description": "Registry of verified procedural tests for automated CI runs",
  "tests": [
    {
      "id": "atlas-create-cluster",
      "path": "content/atlas/source/tutorial/create-cluster.txt",
      "owner": "atlas-docs-team",
      "addedDate": "2024-01-15",
      "variants": ["atlas-ui", "mongosh"],
      "tags": ["atlas", "clusters", "tutorial"],
      "notes": "Tests cluster creation via UI and mongosh"
    },
    {
      "id": "drivers-connect-nodejs",
      "path": "content/drivers/source/quick-start/nodejs.txt",
      "owner": "drivers-team",
      "addedDate": "2024-01-20",
      "variants": ["nodejs"],
      "tags": ["drivers", "nodejs", "quick-start"],
      "notes": "Node.js driver connection quick start"
    }
  ]
}
```

**Using the Test Registry in CI**:

```yaml
# .github/workflows/test-procedures.yml
name: Test Verified Procedures

on:
  schedule:
    - cron: '0 0 * * *'  # Daily at midnight
  push:
    branches:
      - main
    paths:
      - 'code-example-tests/procedures/test-registry.json'
      - 'content/**/source/**/*.txt'

jobs:
  test-registry:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3

      - name: Setup Node.js
        uses: actions/setup-node@v3
        with:
          node-version: '24'

      - name: Install proctest
        run: npm install -g @grove-platform/proctest

      - name: Create .env file
        run: |
          echo "MONGODB_URI=${{ secrets.MONGODB_URI }}" >> .env

      - name: Run registry tests
        run: proctest test --registry code-example-tests/procedures/test-registry.json --reporter junit --output results.xml

      - name: Publish test results
        uses: EnricoMi/publish-unit-test-result-action@v2
        if: always()
        with:
          files: results.xml
```

**Adding Tests to the Registry**:

Writers can add procedures to the registry after verifying they work:

```bash
# 1. Test the procedure locally
proctest test content/atlas/source/tutorial/my-new-procedure.txt

# 2. If it passes, add it to the registry
# Edit code-example-tests/procedures/test-registry.json and add:
{
  "id": "atlas-my-new-feature",
  "path": "content/atlas/source/tutorial/my-new-procedure.txt",
  "owner": "your-team",
  "addedDate": "2024-11-24",
  "variants": ["atlas-ui"],
  "tags": ["atlas", "tutorial"],
  "notes": "Brief description of what this tests"
}

# 3. Submit PR with both the procedure and registry update
```

**Benefits of the Test Registry**:

- ✅ **Curated test suite** - Only verified, working procedures run in CI
- ✅ **Ownership tracking** - Know who to contact when tests fail
- ✅ **Metadata** - Tags and notes help organize and understand tests
- ✅ **Gradual adoption** - Add procedures incrementally as they're verified
- ✅ **Variant tracking** - Document which variants are tested
- ✅ **Change detection** - CI runs when registry or procedures change

---

## Common Patterns

### Prerequisites and Requirements

The framework automatically detects and validates prerequisites:

```rst
Prerequisites
-------------

- MongoDB Server 6.0 or later
- Node.js 18 or later
- An Atlas account with API access

Procedure
---------

1. Install the MongoDB driver...
```

If prerequisites aren't met, the test is **skipped** (not failed):

```
⊘ Install MongoDB Driver
  Skipped: Missing prerequisite
  - MongoDB Server 6.0 or later: Not found (detected: 5.0.14)
```

### Working with Placeholders

The framework resolves placeholders from multiple sources:

**Priority order**:
1. Environment variables (`.env` file)
2. Source constants (`snooty.toml`)
3. Configuration file values

**Example**:

```rst
.. code-block:: javascript

   const client = new MongoClient('{+mongodb-uri+}');
   await client.connect();
```

The framework looks for:
1. `MONGODB_URI` in `.env`
2. `mongodb-uri` in `snooty.toml` constants
3. Fuzzy matches if exact match not found

**Helpful errors** when placeholders can't be resolved:

```
✗ Connect to MongoDB
  Error: Unresolved placeholder: {+mongodb-uri+}

  Suggestions:
    - Add MONGODB_URI to your .env file
    - Add mongodb-uri to snooty.toml constants
    - Did you mean: MONGODB_URL, MONGO_URI?
```

### Tabs and Variants

The framework automatically tests **all tabs** and **composable tutorials**
in your documentation:

```rst
.. tabs::

   .. tab:: Python
      :tabid: python

      .. code-block:: python

         from pymongo import MongoClient
         client = MongoClient()

   .. tab:: Node.js
      :tabid: nodejs

      .. code-block:: javascript

         const { MongoClient } = require('mongodb');
         const client = new MongoClient();
```

This generates **2 test cases**:
- "Install MongoDB Driver (Python)"
- "Install MongoDB Driver (Node.js)"

Each variant is tested independently and reported separately.

### Cleanup Behavior

By default, the framework cleans up:
- ✅ Test databases (names starting with `test_`)
- ✅ Test collections
- ✅ Temporary files created during tests
- ✅ Background processes

**Disable cleanup for debugging**:

```bash
proctest test source/tutorial/create-database.txt --no-cleanup
```

After the test, you can inspect the database:

```bash
mongosh
> show dbs
test_mydb  0.000GB  # Still exists for inspection
```

**Custom cleanup patterns** in configuration:

```javascript
module.exports = {
  cleanup: {
    databases: {
      enabled: true,
      pattern: /^(test_|temp_|demo_)/  // Clean up databases matching this pattern
    }
  }
};
```

---

## Troubleshooting

### Common Issues

#### "Unresolved placeholder" errors

**Problem**: `Error: Unresolved placeholder: {+api-key+}`

**Solution**:
1. Check your `.env` file has the variable: `API_KEY=your-key`
2. Check `snooty.toml` for the constant
3. Look at the suggestions in the error message

#### "Command not found" errors

**Problem**: `Error: Command failed: mongosh`

**Solution**:
1. Install the missing tool (`mongosh`, `npm`, `python`, etc.)
2. Ensure it's in your PATH
3. Check prerequisites in the documentation

#### Tests pass locally but fail in CI

**Problem**: Tests work on your machine but fail in GitHub Actions

**Solution**:
1. Check CI has all required secrets configured
2. Verify MongoDB service is running in CI
3. Check Node.js version matches (use Node.js 24+)
4. Review CI logs for missing environment variables

#### Parser doesn't detect my procedure

**Problem**: `No procedures found in file`

**Solution**:
1. Use `proctest parse` to see what's detected
2. Check your RST syntax (ordered lists or `.. procedure::` directive)
3. Ensure proper indentation (RST is whitespace-sensitive)

#### Cleanup doesn't work

**Problem**: Test databases remain after tests

**Solution**:
1. Check database names start with `test_` (default pattern)
2. Verify cleanup is enabled (it is by default)
3. Check for errors during cleanup in verbose output:
   ```bash
   proctest test source/tutorial/ --verbose
   ```

### Getting Help

1. **Check the parse output** first:
   ```bash
   proctest parse your-file.txt --format json
   ```

2. **Run with verbose logging**:
   ```bash
   proctest test your-file.txt --verbose
   ```

3. **Check the technical specification** for detailed behavior

4. **Ask in #ask-devdocs** Slack channel

---

## Summary

### Quick Reference

```bash
# Test a file
proctest test source/tutorial/getting-started.txt

# Test a directory
proctest test source/tutorial/

# Debug parsing
proctest parse source/tutorial/getting-started.txt

# Output JSON for CI
proctest test source/ --reporter json

# Output JUnit XML for CI
proctest test source/ --reporter junit --output results.xml

# Test without cleanup
proctest test source/tutorial/ --no-cleanup

# Verbose output
proctest test source/tutorial/ --verbose
```

### Key Principles

1. ✅ **Zero configuration** for most cases
2. ✅ **Automatic cleanup** by default
3. ✅ **Helpful error messages** with suggestions
4. ✅ **All variants tested** automatically
5. ✅ **Prerequisites validated** before running tests
6. ✅ **Share configurations** across teams in monorepos

### Next Steps

1. Install the framework: `npm install -g @grove-platform/proctest`
2. Create a `.env` file with your credentials
3. Run your first test: `proctest test source/tutorial/`
4. Set up CI integration for your repository
5. Contribute shared navigation mappings for your team

Happy testing! 🎉

- `0` - All tests passed
- `1` - One or more tests failed
- `2` - Configuration or parsing error

Use exit codes in scripts:

```bash
if proctest test source/tutorial/; then
  echo "All tests passed!"
else
  echo "Tests failed"
  exit 1
fi
```

---

## Configuration

### When Do You Need Configuration?

**90% of procedures work with zero configuration.** You only need a config file for:

- Custom UI navigation mappings
- Team-specific placeholder values
- Custom cleanup strategies
- Non-standard file patterns

### Creating a Configuration File

Create `.proctest.js` in your repository root:

```javascript
// .proctest.js
module.exports = {
  // Custom test file patterns (optional)
  testMatch: [
    'source/tutorial/**/*.txt',
    'source/guides/**/*.txt'
  ],

  // Exclude certain files (optional)
  testIgnore: [
    'source/includes/**',
    'source/tutorial/draft-*.txt'
  ],

  // UI navigation mappings (for UI testable actions)
  ui: {
    navigationMappings: {
      'click the {button} button': async (page, button) => {
        await page.click(`button:has-text("${button}")`);
      },
      'navigate to {page}': async (page, pageName) => {
        const routes = {
          'Database': '/databases',
          'Clusters': '/clusters',
          'Users': '/security/users'
        };
        await page.goto(`https://cloud.mongodb.com${routes[pageName]}`);
      }
    },

    // User-specific values for testing
    userValues: {
      'your project': 'test-project-123',
      'your cluster': 'test-cluster-0',
      'your database': 'test_db'
    }
  },

  // Cleanup configuration (optional)
  cleanup: {
    enabled: true, // Default
    databases: { enabled: true },
    collections: { enabled: true },
    files: { enabled: true }
  }
};
```



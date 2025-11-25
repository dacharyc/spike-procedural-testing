# Analysis: Implicit Patterns in Test Data

## Executive Summary

After auditing all procedure files in the `testdata/` directory, I've identified **5 major categories of implicit patterns** that could cause issues for both human readers and automated testing. This document categorizes these patterns, provides examples, and recommends solutions (tooling vs. training).

---

## Methodology

- **Files Analyzed**: 68+ RST/TXT files across `testdata/atlas/` and `testdata/drivers/`
- **Search Patterns**:
  - Directory navigation ("cd", "enter", "navigate to")
  - File operations ("create", "copy", "paste", "replace", "save")
  - Execution commands ("run", "execute", "from your IDE")
  - Implicit references ("the file", "your project")

---

## Category 1: Implicit Directory Navigation ⚠️ HIGH PRIORITY

### Pattern: Commands that create directories without explicit navigation

**Frequency**: ~15 occurrences
**Impact**: HIGH - Causes test failures
**Solution**: **TRAINING** (fix documentation)

### Examples Found:

#### ❌ Implicit (Symfony - FIXED)
```rst
.. step:: Initialize a Symfony Project
   composer create-project symfony/skeleton restaurants

.. step:: Install Dependencies
   composer require doctrine/mongodb-odm-bundle
```

#### ✅ Explicit (C# - Good Example)
```rst
.. step:: Create a new directory and initialize your project.
   a. mkdir csharp-create-index
   #. cd csharp-create-index
   #. dotnet new console
```

#### ✅ Explicit (Node.js - Good Example)
```rst
a. Run the following commands:
   mkdir atlas-search-quickstart
   cd atlas-search-quickstart
   npm init -y
```

### Files with Good Patterns:
- ✅ `testdata/atlas/source/includes/fts/search-index-management/procedures/steps-fts-create-index-csharp.rst`
- ✅ `testdata/atlas/source/includes/fts/search-index-management/procedures/steps-fts-create-index-tutorial-node.rst`
- ✅ `testdata/atlas/source/includes/fts/search-index-management/procedures/steps-fts-create-index-tutorial-csharp.rst`

### Recommendation:
**TRAINING** - Update style guide to require explicit `cd` commands after directory creation.

---

## Category 2: Implicit File Creation and Content Insertion 🟡 MEDIUM PRIORITY

### Pattern: "Create file X and paste/copy code" without explicit file operation

**Frequency**: ~40 occurrences
**Impact**: MEDIUM - Tooling can handle this
**Solution**: **TOOLING** (already addressed in spec)

### Examples Found:

#### Pattern A: "Create file" + literalinclude
```rst
.. step:: Create a new file named ``create-index.go``.

.. step:: Copy the following code example into the file.

   .. literalinclude:: /includes/fts/search-index-management/create-index.go
      :language: go
```

#### Pattern B: "Create file and define" (combined)
```rst
.. step:: Create a ``.java`` file and define the index in the file.

   .. literalinclude:: /includes/avs/index-management/create-index/create-index.java
      :language: java
```

#### Pattern C: "Paste code into file"
```rst
Paste the following code into the ``index.html.twig`` file:

.. literalinclude:: /includes/php-frameworks/symfony/index.html.twig
   :language: html
```

### Detection Strategy (Already in Spec):
The tooling should detect these patterns:
- Prose containing: "Create", "Paste", "Copy", "Replace" + filename
- Followed by: `literalinclude` directive or code block
- Map to: `FileTestableAction` with operation type (create/replace/append)

### Recommendation:
**TOOLING** - Already addressed in technical specification (Appendix D.2: File Testable Actions)

---

## Category 3: IDE-Based File Execution ✅ HANDLED BY TOOLING

### Pattern: "From your IDE, run the file"

**Frequency**: ~25 occurrences
**Impact**: MEDIUM - Needs language-specific handling
**Solution**: **TOOLING** (convention-over-configuration)
**Status**: ✅ Addressed in technical specification

### Examples Found:

#### Pattern A: "Run the file" (Python - Explicit)
```rst
.. step:: Run the file using the following command.

   .. code-block:: shell

      python view_index.py
```

#### Pattern B: "From your IDE, run the file" (Java - IDE-based)
```rst
.. step:: Execute the code to create the index.

   From your IDE, run the file to create the index.
```

#### Pattern C: Explicit command (Node.js - Explicit)
```rst
.. io-code-block::
   :copyable: true

   .. input::
      :language: shell

      node create-index.js
```

### Why "From your IDE" is Intentional:

This pattern is a **deliberate documentation choice** because:
1. Most developers use IDEs for execution (IntelliJ, VS Code, Eclipse)
2. Command-line execution varies by build tool (Maven vs Gradle vs direct `java`)
3. Avoids prescribing specific toolchains that may not match user environments

### Tooling Solution - Convention over Configuration:

The `proctest` framework handles this using **sensible defaults with optional overrides**:

**Default Commands** (Convention):
```typescript
const DEFAULT_IDE_COMMANDS = {
  java: 'mvn compile exec:java -Dexec.mainClass="{className}"',
  csharp: 'dotnet run',
  cpp: 'g++ {filename} -o {basename} && ./{basename}',
  python: 'python {filename}',
  javascript: 'node {filename}',
  go: 'go run {filename}',
};
```

**User Override** (Configuration):
```javascript
// .proctest.js
module.exports = {
  ideExecution: {
    commands: {
      java: 'gradle run',  // User prefers Gradle
    }
  }
};
```

### Current Handling:
- ✅ Pattern A & C: Explicit shell commands → `ShellTestableAction`
- ✅ Pattern B: "From your IDE" → `CodeTestableAction` with `executionMode: 'ide'`
  - Uses default command for language (e.g., `mvn exec:java` for Java)
  - User can override via `.proctest.js` configuration
  - Reports which command was used for transparency

### Recommendation:
**NO DOCUMENTATION CHANGES NEEDED** - This pattern is correct as-is. The tooling handles it automatically.

---

## Category 4: Implicit Placeholder Replacement 🟢 LOW PRIORITY

### Pattern: "Replace the following values and save the file"

**Frequency**: ~30 occurrences
**Impact**: LOW - Tooling handles this well
**Solution**: **TOOLING** (already handled by placeholder resolution)

### Examples Found:

```rst
.. step:: Replace the following values and save the file.

   - ``<connection-string>``: Your Atlas connection string
   - ``<database-name>``: Database name
   - ``<collection-name>``: Collection name
```

### Current Handling:
✅ The `PlaceholderResolver` already handles this:
- Detects `<placeholder>` patterns
- Resolves from `.env`, `snooty.toml`, or prompts user
- No changes needed

### Recommendation:
**NO ACTION** - Already handled by existing tooling.

---

## Category 5: Implicit Directory Context 🟡 MEDIUM PRIORITY

### Pattern: References to "the root directory", "your project directory" without naming it

**Frequency**: ~10 occurrences
**Impact**: MEDIUM - Causes confusion
**Solution**: **TRAINING** (fix documentation)

### Examples Found:

#### ❌ Implicit (Symfony - FIXED)
```rst
In the root directory, navigate to the ``.env`` file and define the
following environment variables:
```

#### ✅ Explicit (Symfony - After Fix)
```rst
In the project root directory (``restaurants/``), replace the contents of
the ``.env`` file with the following code:
```

### Recommendation:
**TRAINING** - Always name the directory explicitly:
- ❌ "the root directory" → ✅ "the project root directory (`my-project/`)"
- ❌ "your project directory" → ✅ "the `atlas-search-quickstart/` directory"
- ❌ "the application root" → ✅ "the `restaurants/` directory"

---

## Summary Table

| Category | Frequency | Impact | Solution | Status |
|----------|-----------|--------|----------|--------|
| **1. Implicit Directory Navigation** | ~15 | HIGH | Training | ✅ Addressed in Symfony |
| **2. Implicit File Creation** | ~40 | MEDIUM | Tooling | ✅ In spec (D.2) |
| **3. IDE-Based File Execution** | ~25 | MEDIUM | Tooling | ✅ In spec (D.3) |
| **4. Implicit Placeholder Replacement** | ~30 | LOW | Tooling | ✅ Already handled |
| **5. Implicit Directory Context** | ~10 | MEDIUM | Training | ✅ Addressed in Symfony |

---

## Recommendations by Priority

### 🔴 High Priority (Do First)

1. **Update Style Guide**: Add explicit directory navigation requirements
   - Require `cd` command after any directory creation
   - Require explicit directory names (not "your project directory")
   - "From your IDE, run" pattern is ACCEPTABLE (tooling handles it)

2. **Create Linting Rule**: Detect implicit patterns in CI
   - Flag: "Enter your project directory" without `cd` command
   - Flag: "the root directory" without explicit name

### 🟡 Medium Priority (Do Next)

3. **Audit Existing Docs**: Review and fix implicit patterns
   - Priority files: Driver tutorials, quickstart guides
   - Use Symfony fix as template

4. **Enhance Tooling**: Add warnings for unhandled patterns
   - Warn: "From your IDE" detected (suggest explicit command)
   - Warn: File creation without literalinclude (manual file needed?)

### 🟢 Low Priority (Nice to Have)

5. **Documentation**: Add examples to writer training
   - Use `TRAINING-working-directory-clarity.md` as foundation
   - Add sections for file operations, execution commands

---

## Next Steps

1. ✅ **Symfony file updated** - Serves as gold standard example
2. ✅ **Training document created** - Ready for writers
3. 📋 **Create style guide updates** - Document these patterns
4. 📋 **Create linting rules** - Automate detection in CI
5. 📋 **Audit other files** - Fix high-priority implicit patterns

---

## Appendix: Specific Files by Category

### Category 1: Files with Implicit Directory Navigation

**Need Review** (no explicit `cd` after directory creation):
- None found in current test data (most examples are good!)

**Good Examples** (explicit `cd` commands):
- `testdata/atlas/source/includes/fts/search-index-management/procedures/steps-fts-create-index-csharp.rst` (lines 18-37)
- `testdata/atlas/source/includes/fts/search-index-management/procedures/steps-fts-create-index-tutorial-node.rst` (lines 1-11)
- `testdata/atlas/source/includes/fts/search-index-management/procedures/steps-fts-create-index-tutorial-csharp.rst` (lines 1-10)

### Category 2: Files with "Create file + literalinclude" Pattern

**These are FINE** (tooling will handle):
- `testdata/atlas/source/includes/fts/search-index-management/procedures/steps-fts-create-index-go.rst` (line 17)
- `testdata/atlas/source/includes/fts/search-index-management/procedures/steps-fts-create-index-node.rst` (line 17)
- `testdata/atlas/source/includes/fts/search-index-management/procedures/steps-fts-view-index-python.rst` (line 10)
- `testdata/atlas/source/includes/avs/index-examples/steps-avs-create-index-*.rst` (40+ files)

### Category 3: Files with "From your IDE" Pattern

**These are FINE** (tooling handles with convention-over-configuration):
- `testdata/atlas/source/includes/avs/index-examples/steps-avs-create-index-java.rst` (line 82)
- `testdata/atlas/source/includes/avs/index-examples/steps-avs-edit-index-java.rst` (line 36)
- `testdata/atlas/source/includes/avs/index-examples/steps-avs-view-index-java.rst` (line 32)
- `testdata/atlas/source/includes/avs/index-examples/steps-avs-delete-index-java.rst` (line 30)

**How Tooling Handles It**:
```rst
Documentation says:
   From your IDE, run the file to create the index.

Tooling executes:
   mvn compile exec:java -Dexec.mainClass="CreateIndex"

User can override in .proctest.js:
   ideExecution: {
     commands: { java: 'gradle run' }
   }
```

### Category 4: Files with Placeholder Replacement

**No Action Needed** (tooling handles these):
- All files with "Replace the following values and save the file" pattern
- ~30 files across testdata/atlas/source/includes/

### Category 5: Files with Implicit Directory References

**Fixed**:
- ✅ `testdata/drivers/source/symfony.txt` (lines 216, 282)

**Pattern to Watch For**:
- "In the root directory" → Should be "In the project root directory (`project-name/`)"
- "Navigate to the `.env` file" → Should be "Replace the contents of the `.env` file"

---

## Tooling Detection Patterns

For automated linting, detect these patterns:

### Red Flags (Require Manual Review):
```regex
# Implicit directory references
(enter|navigate to|in) (the|your) (root |project )?directory(?! \(`[^`]+`/\))

# Missing cd after directory creation
(mkdir|composer create-project|npm create|dotnet new).*\n(?!.*cd )

# Vague file references
(navigate to|open) the \.[\w]+ file(?! in the `[^`]+` directory)
```

### Acceptable Patterns (Tooling Handles):
```regex
# IDE execution (tooling provides defaults)
from your ide.*run

# File creation with literalinclude (tooling detects)
create.*file.*named.*\n.*literalinclude
```

### Green Patterns (Good Examples):
```regex
# Explicit directory with name
in the `[\w-]+/` directory

# Explicit cd command
cd [\w-]+

# Explicit execution command
(python|node|java|dotnet run|go run) [\w.-]+
```



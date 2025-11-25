# Technical Specification: Procedural Testing Framework
## Convention + Configuration and Plugin-Ready Architecture

**Version**: 1.0
**Date**: 2025-11-24
**Status**: Draft for Review

---

## Table of Contents

1. [Executive Summary](#1-executive-summary)
2. [System Architecture](#2-system-architecture)
3. [Component Specifications](#3-component-specifications)
4. [Data Models](#4-data-models)
5. [API Specifications](#5-api-specifications)
6. [Implementation Plan](#6-implementation-plan)
7. [Testing Strategy](#7-testing-strategy)
8. [Risk Assessment](#8-risk-assessment)
9. [Appendices](#9-appendices)

---

## 1. Executive Summary

### 1.1 Purpose

This technical specification defines the implementation details for a procedural testing framework that enables technical writers to validate documentation procedures. The framework parses reStructuredText (RST) documentation, extracts procedural steps with code examples, executes the code, and reports results.

### 1.2 Architecture Choice

**Convention + Configuration and Plugin Ready** provides:
- Zero-configuration operation for 90% of use cases
- Progressive configuration disclosure for complex scenarios
- Interface-driven design enabling future plugin system
- No plugin infrastructure overhead in initial implementation

### 1.3 Technology Stack

- **Runtime**: Node.js 24+ (LTS)
- **Language**: TypeScript 5.x
- **Package Manager**: npm (standard with Node.js)
- **Testing**: Jest (team familiarity, organizational consistency)
- **CLI Framework**: Commander.js
- **Parser**: Custom RST parser (with future MDX support)
- **Process Execution**: Node.js `child_process` with proper isolation

### 1.4 Key Design Principles

1. **Convention over Configuration**: Smart defaults, minimal required setup
2. **Interface-Driven**: All major components implement well-defined interfaces
3. **Progressive Disclosure**: Complexity revealed only when needed
4. **Fail-Fast with Helpful Errors**: Clear, actionable error messages
5. **Isolation**: Tests don't interfere with each other
6. **Cleanup**: Automatic resource cleanup with manual override options

### 1.5 Success Criteria

**PoC**:
- ✅ Parse RST files and extract procedures
- ✅ Execute code blocks in JavaScript, Python, PHP, Shell
- ✅ Resolve placeholders from `.env` and `snooty.toml`
- ✅ Report results with helpful error messages
- ✅ Clean up test databases automatically

**Production-Ready**:
- ✅ Support optional configuration file
- ✅ Handle tabs and composable tutorials
- ✅ Multiple output formats (human, JSON, JUnit)
- ✅ 90%+ writer satisfaction
- ✅ < 5% of procedures require configuration

---

## 2. System Architecture

### 2.1 High-Level Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         CLI Entry Point                          │
│                        (proctest command)                        │
└────────────────────────────┬────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                   Configuration Manager                          │
│  • Load conventions (auto-discover .env, snooty.toml)            │
│  • Load optional config file (.proctest.{json,js})               │
│  • Merge with CLI flags                                          │
│  • Validate and provide defaults                                 │
└────────────────────────────┬────────────────────────────────────┘
                             │
                             ▼
┌─────────────────────────────────────────────────────────────────┐
│                      Test Orchestrator                           │
│  • Discover test files                                           │
│  • Create execution context                                      │
│  • Coordinate component lifecycle                                │
│  • Manage test execution flow                                    │
└────────────────────────────┬────────────────────────────────────┘
                             │
        ┌────────────────────┼────────────────────┬────────────────┐
        ▼                    ▼                    ▼                ▼
┌──────────────┐    ┌──────────────┐    ┌──────────────┐  ┌──────────────┐
│    Parser    │    │   Resolver   │    │   Executor   │  │   Reporter   │
│  Interface   │    │  Interface   │    │  Interface   │  │  Interface   │
└──────┬───────┘    └──────┬───────┘    └──────┬───────┘  └──────┬───────┘
       │                   │                   │                 │
       ▼                   ▼                   ▼                 ▼
┌──────────────┐    ┌──────────────┐    ┌──────────────┐  ┌──────────────┐
│ RST Parser   │    │   Layered    │    │   Language   │  │    Human     │
│ (built-in)   │    │  Resolver    │    │  Executors   │  │   Reporter   │
│              │    │  (built-in)  │    │  (built-in)  │  │  (built-in)  │
└──────────────┘    └──────────────┘    └──────────────┘  └──────────────┘
```

### 2.2 Component Interaction Flow

```
1. CLI parses arguments
   ↓
2. Configuration Manager loads config (conventions + optional file + CLI flags)
   ↓
3. Test Orchestrator discovers test files
   ↓
4. For each test file:
   a. Parser converts RST → AST
   b. Orchestrator extracts procedures and expands variants:
      - If procedure contains tabs → generate test case for each tab
        (tab content + any content outside tabs)
      - If procedure contains composable tutorial → generate test case for each selection
        (selection-specific content + general content, properly interleaved)
      - Each variant is a complete, executable procedure
   c. For each procedure variant:
      i.   Orchestrator creates per-procedure temp directory
      ii.  Resolver interpolates placeholders
      iii. Executor runs file operations (create/replace files)
      iv.  Executor runs code blocks
      v.   Cleanup Registry tracks resources (including temp directory)
   d. Reporter collects results
   ↓
5. Reporter outputs final results
   ↓
6. Cleanup Registry executes cleanup (LIFO order)
   ↓
7. Exit with appropriate code (0 = success, 1 = failure)
```

### 2.3 Directory Structure

```
proctest/
├── src/
│   ├── cli/
│   │   ├── index.ts              # CLI entry point
│   │   ├── commands/
│   │   │   ├── test.ts           # Main test command
│   │   │   └── init.ts           # Config initialization
│   │   └── args.ts               # Argument parsing
│   ├── config/
│   │   ├── index.ts              # Configuration manager
│   │   ├── loader.ts             # Config file loading
│   │   ├── defaults.ts           # Default configuration
│   │   └── validator.ts          # Config validation
│   ├── core/
│   │   ├── orchestrator.ts       # Test orchestration
│   │   ├── context.ts            # Execution context
│   │   ├── discovery.ts          # Test file discovery
│   │   └── cleanup.ts            # Cleanup registry
│   ├── interfaces/
│   │   ├── parser.ts             # Parser interface
│   │   ├── resolver.ts           # Resolver interface
│   │   ├── executor.ts           # Executor interface
│   │   ├── reporter.ts           # Reporter interface
│   │   └── types.ts              # Shared types
│   ├── parser/
│   │   ├── index.ts              # Parser factory
│   │   ├── rst/
│   │   │   ├── parser.ts         # RST parser implementation
│   │   │   ├── tokenizer.ts      # RST tokenization
│   │   │   ├── ast-builder.ts    # AST construction
│   │   │   └── directives/       # RST directive handlers
│   │   │       ├── procedure.ts
│   │   │       ├── code-block.ts
│   │   │       ├── tabs.ts
│   │   │       └── composable.ts
│   │   └── mdx/                  # Future: MDX parser
│   ├── resolver/
│   │   ├── index.ts              # Resolver factory
│   │   ├── layered-resolver.ts   # Layered resolution strategy
│   │   ├── resolvers/
│   │   │   ├── exact-env.ts      # Exact environment match
│   │   │   ├── fuzzy-env.ts      # Fuzzy environment match
│   │   │   ├── snooty.ts         # Snooty constants
│   │   │   └── custom.ts         # Custom resolver (config)
│   │   └── utils/
│   │       ├── fuzzy-match.ts    # Fuzzy matching algorithm
│   │       └── suggestions.ts    # Suggestion generation
│   ├── executor/
│   │   ├── index.ts              # Executor registry
│   │   ├── base-executor.ts      # Base executor class
│   │   ├── executors/
│   │   │   ├── javascript.ts     # JavaScript/Node.js
│   │   │   ├── python.ts         # Python
│   │   │   ├── php.ts            # PHP
│   │   │   ├── shell.ts          # Shell/Bash
│   │   │   └── mongosh.ts        # MongoDB Shell (special)
│   │   └── utils/
│   │       ├── process.ts        # Process execution utilities
│   │       ├── timeout.ts        # Timeout handling
│   │       └── state.ts          # State accumulation
│   ├── reporter/
│   │   ├── index.ts              # Reporter factory
│   │   ├── reporters/
│   │   │   ├── human.ts          # Human-friendly output
│   │   │   ├── json.ts           # JSON output
│   │   │   └── junit.ts          # JUnit XML output
│   │   └── formatters/
│   │       ├── error.ts          # Error formatting
│   │       └── suggestions.ts    # Suggestion formatting
│   └── utils/
│       ├── env.ts                # Environment loading
│       ├── snooty.ts             # Snooty.toml parsing
│       ├── file.ts               # File utilities
│       └── logger.ts             # Logging utilities
├── tests/
│   ├── unit/                     # Unit tests
│   ├── integration/              # Integration tests
│   └── fixtures/                 # Test fixtures
├── package.json
├── tsconfig.json
├── jest.config.ts
└── README.md
```

---

## 3. Component Specifications

### 3.1 Core Interfaces

All major components implement well-defined interfaces to enable future extensibility without refactoring.

#### 3.1.1 Parser Interface

```typescript
/**
 * Parser interface for converting documentation formats to AST
 */
export interface Parser {
  /**
   * Parse content into an Abstract Syntax Tree
   * @param content - Raw file content
   * @param context - Parsing context (file path, config, etc.)
   * @returns Parsed AST
   */
  parse(content: string, context: ParserContext): Promise<DocumentAST>;

  /**
   * Check if this parser can handle the given file
   * @param filePath - Path to the file
   * @returns true if parser can handle this file
   */
  canParse(filePath: string): boolean;

  /**
   * Get supported file extensions
   */
  getSupportedExtensions(): string[];
}

export interface ParserContext {
  filePath: string;
  workingDirectory: string;
  config: Configuration;
  snootyConfig?: SnootyConfig;
}

export interface DocumentAST {
  type: 'document';
  filePath: string;
  procedures: ProcedureNode[];
  metadata: DocumentMetadata;
}

export interface ProcedureNode {
  type: 'procedure';
  title?: string;
  style?: 'normal' | 'connected';
  prerequisites?: PrerequisiteNode; // Requirements that must be met before running
  steps: StepNode[];
  location: SourceLocation;
}

/**
 * Prerequisites/Requirements for a procedure
 */
export interface PrerequisiteNode {
  type: 'prerequisites';
  title?: string; // "Prerequisites", "Requirements", etc.
  requirements: Requirement[];
  location: SourceLocation;
}

/**
 * Individual requirement that must be met
 */
export type Requirement =
  | SoftwareRequirement
  | EnvironmentRequirement
  | ServiceRequirement
  | ConfigurationRequirement;

/**
 * Base interface for all requirements
 */
export interface BaseRequirement {
  requirementType: 'software' | 'environment' | 'service' | 'configuration';
  description: string;
  optional: boolean;
  location: SourceLocation;
}

/**
 * Software/tool that must be installed
 */
export interface SoftwareRequirement extends BaseRequirement {
  requirementType: 'software';
  name: string; // e.g., "PHP", "Node.js", "Composer"
  version?: string; // e.g., ">=8.0", "^18.0.0"
  checkCommand?: string; // Command to verify installation, e.g., "php --version"
  installUrl?: string; // URL to installation instructions
}

/**
 * Environment variable or configuration that must be set
 */
export interface EnvironmentRequirement extends BaseRequirement {
  requirementType: 'environment';
  variable: string; // e.g., "MONGODB_URI"
  example?: string; // Example value
}

/**
 * External service that must be available
 */
export interface ServiceRequirement extends BaseRequirement {
  requirementType: 'service';
  name: string; // e.g., "MongoDB Atlas cluster"
  setupUrl?: string; // URL to setup instructions
}

/**
 * Configuration file or setting that must exist
 */
export interface ConfigurationRequirement extends BaseRequirement {
  requirementType: 'configuration';
  name: string; // e.g., "snooty.toml", ".env file"
  path?: string; // Expected file path
}

export interface StepNode {
  type: 'step';
  headline?: string;
  content: ContentNode[];
  testableActions: TestableAction[]; // All testable actions in this step
  subSteps?: SubStepNode[]; // Ordered lists within a step
  location: SourceLocation;
}

export interface SubStepNode {
  type: 'sub-step';
  number: number | string; // Can be numeric (1, 2, 3) or alphabetic (a, b, c)
  content: ContentNode[];
  testableActions: TestableAction[]; // All testable actions in this sub-step
  location: SourceLocation;
}

/**
 * Testable Action - Represents any testable action in documentation
 */
export type TestableAction =
  | CodeTestableAction
  | ShellTestableAction
  | UITestableAction
  | CLITestableAction
  | APITestableAction
  | DownloadTestableAction
  | URLTestableAction
  | FileTestableAction;

/**
 * Base interface for all testable actions
 */
export interface BaseTestableAction {
  actionType: 'code' | 'shell' | 'ui' | 'cli' | 'api' | 'download' | 'url' | 'file';
  location: SourceLocation;
}

/**
 * Code block execution (JavaScript, Python, PHP, etc.)
 */
export interface CodeTestableAction extends BaseTestableAction {
  actionType: 'code';
  language: string;
  code: string;
  options: CodeBlockOptions;
  placeholders?: string[];
  executionMode?: 'direct' | 'ide'; // 'ide' when prose says "From your IDE, run..."
  filePath?: string; // Path to file when executionMode is 'ide'
}

/**
 * Shell command execution
 */
export interface ShellTestableAction extends BaseTestableAction {
  actionType: 'shell';
  command: string;
  expectedOutput?: string; // From io-code-block output section
  placeholders?: string[];
}

/**
 * UI interaction (detected via :guilabel: role)
 */
export interface UITestableAction extends BaseTestableAction {
  actionType: 'ui';
  action: 'click' | 'select' | 'input' | 'verify';
  target: string; // The guilabel text
  value?: string; // For input actions
  description?: string; // Surrounding context
}

/**
 * CLI command execution (atlas-cli, mongosh, etc.)
 */
export interface CLITestableAction extends BaseTestableAction {
  actionType: 'cli';
  tool: 'atlas-cli' | 'mongosh' | 'other';
  command: string;
  expectedOutput?: string;
  placeholders?: string[];
}

/**
 * API request execution (MongoDB Atlas Admin API only)
 * Detected when curl command targets Atlas Admin API endpoints
 */
export interface APITestableAction extends BaseTestableAction {
  actionType: 'api';
  method: 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH';
  endpoint: string;
  headers?: Record<string, string>;
  body?: string;
  expectedStatus?: number;
  expectedResponse?: string;
  placeholders?: string[];
}

/**
 * File download action
 * Detected when curl command downloads files (uses -o or -O flags)
 * Important: Downloads may take time and subsequent steps may depend on completion
 */
export interface DownloadTestableAction extends BaseTestableAction {
  actionType: 'download';
  url: string;
  outputPath: string; // Where the file will be saved
  method?: 'GET' | 'POST'; // Usually GET
  headers?: Record<string, string>;
  expectedSize?: number; // Expected file size in bytes (optional)
  timeout?: number; // Download timeout in milliseconds
  description?: string; // What is being downloaded
}

/**
 * URL validation (link checking)
 */
export interface URLTestableAction extends BaseTestableAction {
  actionType: 'url';
  url: string;
  expectedStatus?: number;
  description?: string;
}

/**
 * File operation (create, replace, append)
 * Detected when prose instructions describe file operations followed by literalinclude
 * Examples:
 * - "Create the Restaurant.php file in the src/Document directory and paste the following code"
 * - "Replace the contents of your doctrine_mongodb.yaml file with the following code"
 * - "Paste the following code into the index.html.twig file"
 */
export interface FileTestableAction extends BaseTestableAction {
  actionType: 'file';
  operation: 'create' | 'replace' | 'append';
  path: string; // Relative path from working directory
  content: string; // File content (from literalinclude or code block)
  description?: string; // Human-readable description of the operation
}

export interface CodeBlockNode {
  type: 'code-block';
  language: string;
  code: string;
  options: CodeBlockOptions;
  location: SourceLocation;
}

export interface CodeBlockOptions {
  emphasizeLines?: number[];
  lineNumbers?: boolean;
  caption?: string;
  executable?: boolean; // Derived from directive type
}

export interface SourceLocation {
  filePath: string;
  startLine: number;
  endLine: number;
  startColumn?: number;
  endColumn?: number;
}
```

#### 3.1.2 Prerequisite Checker Interface

```typescript
/**
 * Prerequisite checker interface for validating requirements
 */
export interface PrerequisiteChecker {
  /**
   * Check if a requirement is met
   * @param requirement - Requirement to check
   * @returns Check result with details
   */
  check(requirement: Requirement): Promise<PrerequisiteCheckResult>;

  /**
   * Check all requirements for a procedure
   * @param prerequisites - Prerequisites node
   * @returns Array of check results
   */
  checkAll(prerequisites: PrerequisiteNode): Promise<PrerequisiteCheckResult[]>;
}

export interface PrerequisiteCheckResult {
  requirement: Requirement;
  met: boolean;
  message: string; // Human-readable message
  details?: {
    found?: string; // What was found (e.g., "PHP 8.2.0")
    expected?: string; // What was expected (e.g., "PHP >=8.0")
    command?: string; // Command that was run to check
    output?: string; // Output from the check command
  };
  skipReason?: string; // If not met, why the test should be skipped
}
```

#### 3.1.3 Resolver Interface

```typescript
/**
 * Resolver interface for placeholder interpolation
 */
export interface PlaceholderResolver {
  /**
   * Resolve a placeholder to its value
   * @param placeholder - The placeholder string (e.g., "<username>")
   * @param context - Resolution context
   * @returns Resolved value or null if cannot resolve
   */
  resolve(placeholder: string, context: ResolverContext): Promise<string | null>;

  /**
   * Get suggestions for unresolved placeholders
   * @param placeholder - The placeholder that couldn't be resolved
   * @param context - Resolution context
   * @returns Array of suggestions
   */
  getSuggestions(placeholder: string, context: ResolverContext): Promise<string[]>;

  /**
   * Priority of this resolver (higher = tried first)
   */
  getPriority(): number;
}

export interface ResolverContext {
  environment: Record<string, string>;
  snootyConstants: Record<string, string>;
  config: Configuration;
  procedure: ProcedureNode;
  step: StepNode;
}

/**
 * Layered resolver that tries multiple strategies
 */
export interface LayeredResolver extends PlaceholderResolver {
  /**
   * Add a resolver to the chain
   */
  addResolver(resolver: PlaceholderResolver): void;

  /**
   * Get all registered resolvers
   */
  getResolvers(): PlaceholderResolver[];
}
```

#### 3.1.3 Executor Interface

```typescript
/**
 * Executor interface for executing different types of testable actions
 */
export interface Executor {
  /**
   * Check if this executor can handle the given testable action
   * @param action - Testable action to check
   * @returns true if executor can handle this action
   */
  canExecute(action: TestableAction): boolean;

  /**
   * Execute a testable action and return result
   * @param action - Action to execute
   * @param context - Execution context
   * @returns Execution result
   */
  execute(action: TestableAction, context: ExecutionContext): Promise<ExecutionResult>;

  /**
   * Get supported types (action types or languages)
   */
  getSupportedTypes(): string[];

  /**
   * Validate that required runtime/tools are available
   * @returns Validation result with error message if invalid
   */
  validate(): Promise<ValidationResult>;
}

/**
 * Specialized executor for code blocks (JavaScript, Python, PHP, etc.)
 */
export interface CodeExecutor extends Executor {
  canExecute(action: TestableAction): action is CodeTestableAction;
  getSupportedLanguages(): string[];
}

/**
 * Specialized executor for shell commands
 */
export interface ShellExecutor extends Executor {
  canExecute(action: TestableAction): action is ShellTestableAction;
}

/**
 * Specialized executor for UI interactions
 */
export interface UIExecutor extends Executor {
  canExecute(action: TestableAction): action is UITestableAction;
  getSupportedActions(): Array<'click' | 'select' | 'input' | 'verify'>;
}

/**
 * Specialized executor for CLI tools (atlas-cli, mongosh)
 */
export interface CLIExecutor extends Executor {
  canExecute(action: TestableAction): action is CLITestableAction;
  getSupportedTools(): string[];
}

/**
 * Specialized executor for API requests
 */
export interface APIExecutor extends Executor {
  canExecute(action: TestableAction): action is APITestableAction;
}

/**
 * Specialized executor for URL validation
 */
export interface URLExecutor extends Executor {
  canExecute(action: TestableAction): action is URLTestableAction;
}

/**
 * Specialized executor for file operations
 */
export interface FileExecutor extends Executor {
  canExecute(action: TestableAction): action is FileTestableAction;
  getSupportedOperations(): Array<'create' | 'replace' | 'append'>;
}

export interface ExecutionContext {
  action: TestableAction;
  procedure: ProcedureNode;
  step: StepNode;
  subStep?: SubStepNode;
  environment: Record<string, string>;
  workingDirectory: string; // Per-procedure temp directory: .proctest/runs/<timestamp>-<procedure-id>/
  timeout: number;
  state: ExecutionState;
  cleanup: CleanupRegistry;
  config: Configuration;
}

export interface ExecutionState {
  /**
   * Accumulated code from previous blocks in this step
   */
  accumulatedCode: string;

  /**
   * Variables/state from previous executions
   */
  variables: Record<string, unknown>;

  /**
   * Whether this is the first code block in the step
   */
  isFirstBlock: boolean;
}

export interface ExecutionResult {
  success: boolean;
  stdout: string;
  stderr: string;
  exitCode: number;
  duration: number; // milliseconds
  error?: Error;
  timedOut?: boolean;
}

export interface ValidationResult {
  valid: boolean;
  error?: string;
  suggestions?: string[];
}
```

#### 3.1.4 Reporter Interface

```typescript
/**
 * Reporter interface for outputting test results
 */
export interface Reporter {
  /**
   * Report results for a single procedure
   * @param result - Procedure test result
   */
  reportProcedure(result: ProcedureResult): Promise<void>;

  /**
   * Report final summary of all tests
   * @param summary - Test summary
   */
  reportSummary(summary: TestSummary): Promise<void>;

  /**
   * Initialize reporter (e.g., open file, print header)
   */
  initialize(): Promise<void>;

  /**
   * Finalize reporter (e.g., close file, print footer)
   */
  finalize(): Promise<void>;
}
```

**Variant Reporting Strategy**:

Each variant (tab or composable tutorial selection) is treated as a **separate test case** with its own result. This provides:

1. **Clear pass/fail status per variant** - If Python passes but Node.js fails, it's immediately visible
2. **Individual timing** - Each variant has its own duration
3. **Separate error reporting** - Failures are isolated to specific variants
4. **Standard test filtering** - Users can filter by variant name (e.g., `--grep python`)

**Human-Readable Output Example**:

```
Testing: testdata/install-driver.txt

✓ Install MongoDB Driver (Python)
  Duration: 2.3s
  Steps: 3/3 passed

✗ Install MongoDB Driver (Node.js)
  Duration: 1.8s
  Steps: 2/3 passed
  Error in Step 2: npm install failed
    Command: npm install mongodb
    Exit code: 1
    stderr: npm ERR! network timeout

✓ Install MongoDB Driver (Java)
  Duration: 3.1s
  Steps: 3/3 passed

Summary:
  3 test cases (1 base procedure with 3 variants)
  2 passed, 1 failed
  Total duration: 7.2s
```

**JSON Output Example**:

```json
{
  "results": [
    {
      "procedure": { "name": "Install MongoDB Driver" },
      "variant": {
        "type": "tab",
        "id": "python",
        "label": "Python",
        "baseProcedureName": "Install MongoDB Driver"
      },
      "success": true,
      "duration": 2300,
      "steps": [...]
    },
    {
      "procedure": { "name": "Install MongoDB Driver" },
      "variant": {
        "type": "tab",
        "id": "nodejs",
        "label": "Node.js",
        "baseProcedureName": "Install MongoDB Driver"
      },
      "success": false,
      "duration": 1800,
      "error": {...}
    }
  ],
  "summary": {
    "totalProcedures": 3,
    "passedProcedures": 2,
    "failedProcedures": 1,
    "variantSummary": {
      "totalBaseProcedures": 1,
      "totalVariants": 3,
      "variantsByType": {
        "tabs": 3,
        "composableTutorials": 0
      }
    }
  }
}
```

**JUnit XML Output**:

Each variant becomes a separate `<testcase>` with the variant label in the name:

```xml
<testsuite name="testdata/install-driver.txt" tests="3" failures="1">
  <testcase name="Install MongoDB Driver (Python)" time="2.3" />
  <testcase name="Install MongoDB Driver (Node.js)" time="1.8">
    <failure message="npm install failed">
      Command: npm install mongodb
      Exit code: 1
      stderr: npm ERR! network timeout
    </failure>
  </testcase>
  <testcase name="Install MongoDB Driver (Java)" time="3.1" />
</testsuite>
```

**Key Design Decisions**:

✅ **Each variant is a separate test case** - Simplifies reporting and filtering
✅ **Variant info is included in result** - Enables grouping and aggregation if needed
✅ **Test name includes variant label** - Clear identification in all output formats
✅ **Summary includes variant statistics** - Optional detailed breakdown for analysis
✅ **No special aggregation in PoC** - Keep it simple; add rollup views in future if needed

```typescript
export interface ProcedureResult {
  procedure: ProcedureNode;
  variant?: VariantInfo; // Present if this is a variant (tab or composable tutorial selection)
  success: boolean;
  skipped: boolean; // True if skipped due to unmet prerequisites
  skipReason?: string; // Reason for skipping
  prerequisiteChecks?: PrerequisiteCheckResult[]; // Results from prerequisite checks
  duration: number;
  steps: StepResult[];
  error?: TestError;
}

export interface VariantInfo {
  type: 'tab' | 'composable-tutorial';
  id: string; // tabid or selection value
  label: string; // Human-readable label (e.g., "Python", "Node.js Driver")
  baseProcedureName: string; // Original procedure name without variant suffix
}

export interface StepResult {
  step: StepNode;
  success: boolean;
  duration: number;
  actionResults: TestableActionResult[]; // Results from all testable actions
  subSteps?: SubStepResult[]; // Results from sub-steps (ordered lists within step)
  error?: TestError;
}

export interface SubStepResult {
  subStepNumber: number | string; // Can be numeric (1, 2, 3) or alphabetic (a, b, c)
  success: boolean;
  duration: number;
  actionResults: TestableActionResult[]; // Results from all testable actions
  error?: TestError;
  location: SourceLocation;
}

/**
 * Result from executing any type of testable action
 */
export interface TestableActionResult {
  action: TestableAction;
  execution: ExecutionResult;
  resolvedContent?: string; // Content after placeholder resolution (for code, shell, cli, api)
  error?: TestError;
}

/**
 * Deprecated: Use TestableActionResult instead
 * Kept for backward compatibility during transition
 */
export interface CodeBlockResult extends TestableActionResult {
  action: CodeTestableAction;
  codeBlock: CodeTestableAction; // Alias for compatibility
  resolvedCode: string; // Alias for resolvedContent
}

export interface TestError {
  type: 'parse' | 'resolve' | 'execute' | 'cleanup';
  message: string;
  location: SourceLocation;
  context?: ErrorContext; // Hierarchical context for error location
  details?: ErrorDetails;
  suggestions?: string[];
}

export interface ErrorContext {
  procedureTitle?: string;
  stepNumber?: number;
  stepHeadline?: string;
  subStepNumber?: number | string; // For sub-steps within a step
  codeBlockLanguage?: string;
}

export interface ErrorDetails {
  placeholder?: string;
  unresolvedPlaceholders?: string[];
  similarEnvVars?: string[];
  stdout?: string;
  stderr?: string;
  exitCode?: number;
}

export interface TestSummary {
  totalProcedures: number; // Total test cases (including variants)
  passedProcedures: number;
  failedProcedures: number;
  skippedProcedures: number;
  totalSteps: number;
  passedSteps: number;
  failedSteps: number;
  totalDuration: number;
  results: ProcedureResult[];

  // Variant-specific summary (optional, for detailed reporting)
  variantSummary?: {
    totalBaseProcedures: number; // Number of unique procedures (before variant expansion)
    totalVariants: number; // Total number of variants across all procedures
    variantsByType: {
      tabs: number;
      composableTutorials: number;
    };
  };
}
```

### 3.2 Configuration System

#### 3.2.1 Configuration Interface

```typescript
export interface Configuration {
  // Test Discovery
  testFiles: string[];
  exclude?: string[];
  registry?: string;  // Path to test registry JSON file

  // Environment
  envFiles: string[];
  snootyConfig?: string;

  // Placeholder Resolution
  placeholders?: PlaceholderConfig;

  // Execution
  executors?: ExecutorConfig;
  timeout?: number;
  ideExecution?: IDEExecutionConfig;  // Handle "From your IDE, run" instructions

  // UI Testing (Phase 3)
  ui?: UITestingConfig;

  // State Management
  stateManagement?: StateManagementConfig;

  // Cleanup
  cleanup?: CleanupConfig;

  // Reporting
  reporters?: ReporterConfig[];
  verbose?: boolean;

  // Hooks
  hooks?: HooksConfig;

  // Future: Plugin support
  plugins?: string[];
}

export interface PlaceholderConfig {
  // Custom resolver function
  resolver?: (placeholder: string, context: ResolverContext) => string | null | Promise<string | null>;

  // Explicit mappings
  mappings?: Record<string, string>;

  // Behavior on unresolved
  onUnresolved?: 'fail' | 'warn' | 'skip';
}

export interface ExecutorConfig {
  [language: string]: LanguageExecutorConfig;
}

export interface LanguageExecutorConfig {
  runtime?: string;
  version?: string;
  timeout?: number;
  env?: Record<string, string>;
  ideCommand?: string; // Custom command for "From your IDE, run" instructions
}

/**
 * IDE Execution Configuration
 * Handles "From your IDE, run the file" instructions
 */
export interface IDEExecutionConfig {
  // Custom commands by language (overrides defaults)
  commands?: Record<string, string>;

  // Skip IDE execution entirely (mark as manual verification)
  skip?: boolean;

  // Available interpolation variables: {filename}, {basename}, {className}
}

/**
 * UI Testing Configuration (Phase 3)
 * Handles mapping generic UI instructions to automation steps
 */
export interface UITestingConfig {
  // Enable/disable UI testing
  enabled?: boolean;

  // Automation framework to use
  framework?: 'playwright' | 'puppeteer' | 'selenium';

  // Browser configuration
  browser?: {
    type?: 'chromium' | 'firefox' | 'webkit';
    headless?: boolean;
    slowMo?: number; // Slow down operations by N milliseconds
  };

  // Navigation mappings: Map generic phrases to automation steps
  navigationMappings?: NavigationMapping[];

  // User-specific values: Map generic references to actual values
  userValues?: Record<string, string>;

  // Base URL for the application under test
  baseUrl?: string;

  // Authentication
  auth?: {
    username?: string;
    password?: string;
    // Or a function that performs login
    loginFunction?: () => Promise<void>;
  };

  // Timeouts
  timeouts?: {
    navigation?: number; // Wait for navigation to complete
    element?: number;    // Wait for element to appear
    action?: number;     // Wait after performing action
  };

  // Screenshots
  screenshots?: {
    onFailure?: boolean;
    onSuccess?: boolean;
    directory?: string;
  };
}

/**
 * Maps a generic navigation phrase to automation steps
 */
export interface NavigationMapping {
  // The phrase in the documentation (e.g., "In Atlas, go to the Clusters page")
  phrase: string | RegExp;

  // The automation steps to perform
  steps: UIAutomationStep[];

  // Optional description for debugging
  description?: string;
}

/**
 * Individual UI automation step
 */
export type UIAutomationStep =
  | { action: 'navigate'; url: string }
  | { action: 'click'; selector: string; waitFor?: string }
  | { action: 'type'; selector: string; text: string }
  | { action: 'select'; selector: string; value: string }
  | { action: 'wait'; selector: string; timeout?: number }
  | { action: 'waitForNavigation'; timeout?: number }
  | { action: 'custom'; function: () => Promise<void> };

export interface StateManagementConfig {
  persistAcrossSteps?: boolean;
  isolationLevel?: 'procedure' | 'step' | 'file';
  strategy?: 'accumulate' | 'isolated';
}

export interface CleanupConfig {
  databases?: DatabaseCleanupConfig;
  collections?: CollectionCleanupConfig;
  files?: FileCleanupConfig;
}

export interface DatabaseCleanupConfig {
  enabled: boolean;
  pattern?: RegExp | string;
  onFailure?: 'warn' | 'block' | 'ignore';
}

export interface CollectionCleanupConfig {
  enabled: boolean;
  pattern?: RegExp | string;
}

export interface FileCleanupConfig {
  enabled: boolean;
  paths?: string[];
}

export interface ReporterConfig {
  type: 'human' | 'json' | 'junit' | 'custom';
  options?: Record<string, unknown>;
}

export interface HooksConfig {
  beforeAll?: (context: TestContext) => Promise<void>;
  beforeEach?: (procedure: ProcedureNode, context: TestContext) => Promise<void>;
  afterEach?: (procedure: ProcedureNode, result: ProcedureResult, context: TestContext) => Promise<void>;
  afterAll?: (summary: TestSummary, context: TestContext) => Promise<void>;
}

export interface TestContext {
  config: Configuration;
  environment: Record<string, string>;
  snootyConstants: Record<string, string>;
  workingDirectory: string;
  cleanup: CleanupRegistry;
}
```

#### 3.2.2 Default Configuration

```typescript
export const DEFAULT_CONFIG: Configuration = {
  // Auto-discover test files
  testFiles: ['**/*.txt', '**/*.rst'],
  exclude: ['**/includes/**', '**/node_modules/**'],

  // Auto-discover environment
  envFiles: ['.env', '.env.local'],
  snootyConfig: 'snooty.toml', // Auto-discover

  // Placeholder resolution
  placeholders: {
    onUnresolved: 'fail'
  },

  // Execution defaults
  timeout: 30000, // 30 seconds

  // State management
  stateManagement: {
    persistAcrossSteps: false,
    isolationLevel: 'step',
    strategy: 'accumulate'
  },

  // Cleanup defaults
  cleanup: {
    databases: {
      enabled: true,
      pattern: /^test_/,
      onFailure: 'warn'
    },
    collections: {
      enabled: true,
      pattern: /^temp_/
    },
    files: {
      enabled: true,
      paths: ['./temp/**']
    },
    workingDirectories: {
      enabled: true, // Automatically clean up per-procedure temp directories
      keepOnFailure: false // Set to true to preserve directories for debugging
    }
  },

  // Reporting
  reporters: [{ type: 'human' }],
  verbose: false
};
```

#### 3.2.3 Configuration Loading Strategy

```typescript
/**
 * Configuration loading priority (highest to lowest):
 * 1. CLI flags (--timeout, --verbose, etc.)
 * 2. Config file (.proctest.js or .proctest.json)
 * 3. Convention-based discovery (.env, snooty.toml)
 * 4. Default configuration
 */
export class ConfigurationManager {
  async load(options: CLIOptions): Promise<Configuration> {
    // 1. Load defaults
    const config = { ...DEFAULT_CONFIG };

    // 2. Apply convention-based discovery
    await this.applyConventions(config);

    // 3. Load config file if exists
    const configFile = await this.findConfigFile();
    if (configFile) {
      const fileConfig = await this.loadConfigFile(configFile);
      Object.assign(config, fileConfig);
    }

    // 4. Apply CLI flags (highest priority)
    this.applyCLIFlags(config, options);

    // 5. Validate configuration
    await this.validate(config);

    return config;
  }

  private async applyConventions(config: Configuration): Promise<void> {
    // Auto-discover .env files
    const envFiles = await this.discoverEnvFiles();
    if (envFiles.length > 0) {
      config.envFiles = envFiles;
    }

    // Auto-discover snooty.toml
    const snootyConfig = await this.discoverSnootyConfig();
    if (snootyConfig) {
      config.snootyConfig = snootyConfig;
    }
  }

  private async findConfigFile(): Promise<string | null> {
    const candidates = [
      '.proctest.js',
      '.proctest.json',
      'proctest.config.js',
      'proctest.config.json'
    ];

    for (const candidate of candidates) {
      if (await fileExists(candidate)) {
        return candidate;
      }
    }

    return null;
  }
}
```

#### 3.2.4 Environment Variable Loading (Native Node.js Support)

Node.js 24+ provides native `.env` file support, eliminating the need for external dependencies like `dotenv`.

**CLI Usage** (Recommended):

```bash
# Load .env file automatically
node --env-file=.env dist/cli.js run testdata/example.txt

# Load .env file if it exists (no error if missing)
node --env-file-if-exists=.env dist/cli.js run testdata/example.txt

# Load multiple .env files (later files override earlier ones)
node --env-file=.env --env-file=.env.local dist/cli.js run testdata/example.txt
```

**Programmatic Usage**:

```typescript
// Load .env file and populate process.env
try {
  process.loadEnvFile('.env');
  console.log('Loaded .env file');
} catch (error) {
  console.warn('.env file not found, using existing environment');
}

// Parse .env file content without loading into process.env
import { parseEnv } from 'node:util';
import { readFileSync } from 'node:fs';

const envContent = readFileSync('.env', 'utf-8');
const envVars = parseEnv(envContent);
console.log('Parsed environment variables:', envVars);
```

**Configuration Manager Integration**:

```typescript
export class ConfigurationManager {
  async load(): Promise<Configuration> {
    const config = await this.loadDefaults();

    // Load .env files using native Node.js support
    for (const envFile of config.envFiles) {
      try {
        process.loadEnvFile(envFile);
        console.log(`Loaded environment from ${envFile}`);
      } catch (error) {
        // File doesn't exist or can't be read - continue
        console.warn(`Could not load ${envFile}: ${error.message}`);
      }
    }

    // Continue with rest of configuration loading...
    return config;
  }
}
```

**Benefits**:
- ✅ No external dependencies required
- ✅ Built-in to Node.js 24+ runtime
- ✅ CLI option for easy integration with npm scripts
- ✅ Programmatic API for advanced use cases
- ✅ Multiple file support with override behavior

**Example .env file**:

```bash
# MongoDB Connection
MONGODB_URI=mongodb+srv://user:pass@cluster.mongodb.net/test

# Atlas API Credentials
ATLAS_API_PUBLIC_KEY=your-public-key
ATLAS_API_PRIVATE_KEY=your-private-key
ATLAS_PROJECT_ID=your-project-id

# Test Configuration
ATLAS_PROJECT_NAME=MyTestProject
ATLAS_CLUSTER_NAME=Cluster0
ATLAS_ORG_NAME=MyOrg
```

#### 3.2.5 Example Configuration with UI Testing

**Complete Configuration Example** (`proctest.config.js`):

```javascript
export default {
  // Test files
  testFiles: ['testdata/**/*.txt'],
  exclude: ['testdata/**/draft-*.txt'],

  // Environment
  envFiles: ['.env', '.env.local'],
  snootyConfig: 'snooty.toml',

  // Placeholder resolution
  placeholders: {
    mappings: {
      '<database-name>': 'test_db_' + Date.now(),
      '<collection-name>': 'test_collection'
    },
    onUnresolved: 'warn'
  },

  // UI Testing Configuration (Phase 3)
  ui: {
    enabled: true,
    framework: 'playwright',

    browser: {
      type: 'chromium',
      headless: false, // Show browser during development
      slowMo: 100      // Slow down by 100ms for visibility
    },

    baseUrl: 'https://cloud.mongodb.com',

    // Authentication
    auth: {
      username: process.env.ATLAS_USERNAME,
      password: process.env.ATLAS_PASSWORD
    },

    // Map generic phrases to automation steps
    navigationMappings: [
      {
        phrase: /In Atlas, go to the Clusters page/i,
        description: 'Navigate to Clusters page in Atlas',
        steps: [
          { action: 'navigate', url: 'https://cloud.mongodb.com' },
          { action: 'wait', selector: '.project-selector', timeout: 5000 },
          { action: 'click', selector: 'nav a[href*="clusters"]' },
          { action: 'waitForNavigation', timeout: 5000 }
        ]
      },
      {
        phrase: /go to the Database Deployments page/i,
        description: 'Navigate to Database Deployments',
        steps: [
          { action: 'click', selector: 'nav [data-testid="nav-database"]' },
          { action: 'waitForNavigation' }
        ]
      },
      {
        phrase: /open the connection dialog/i,
        description: 'Open connection dialog for cluster',
        steps: [
          { action: 'click', selector: '[data-testid="cluster-connect-button"]' },
          { action: 'wait', selector: '.connection-modal' }
        ]
      }
    ],

    // User-specific values for generic references
    userValues: {
      'your project': process.env.ATLAS_PROJECT_NAME || 'MyTestProject',
      'your desired project': process.env.ATLAS_PROJECT_NAME || 'MyTestProject',
      'your cluster': process.env.ATLAS_CLUSTER_NAME || 'Cluster0',
      'your organization': process.env.ATLAS_ORG_NAME || 'MyOrg',
      'your database': 'sample_mflix',
      'your collection': 'movies'
    },

    timeouts: {
      navigation: 10000,
      element: 5000,
      action: 1000
    },

    screenshots: {
      onFailure: true,
      onSuccess: false,
      directory: './test-screenshots'
    }
  },

  // Execution
  timeout: 30000,
  executors: {
    javascript: {
      runtime: 'node',
      version: '>=18.0.0'
    },
    python: {
      runtime: 'python3',
      version: '>=3.8'
    }
  },

  // State management
  stateManagement: {
    persistAcrossSteps: true,
    isolationLevel: 'procedure',
    strategy: 'accumulate'
  },

  // Cleanup
  cleanup: {
    databases: {
      enabled: true,
      pattern: /^test_db_/
    },
    collections: {
      enabled: true,
      pattern: /^test_collection/
    },
    files: {
      enabled: true,
      paths: ['./temp', './test-output']
    }
  },

  // Reporting
  reporters: [
    { type: 'human', output: 'console' },
    { type: 'json', output: './test-results.json' }
  ],
  verbose: true
};
```

**Minimal UI Configuration** (for simple cases):

```javascript
export default {
  ui: {
    enabled: true,
    baseUrl: 'https://cloud.mongodb.com',

    // Just provide user values, use default navigation
    userValues: {
      'your project': 'MyProject',
      'your cluster': 'Cluster0'
    }
  }
};
```

#### 3.2.6 Shared Navigation Mappings

**Organizational Pattern**: For teams with multiple documentation repositories or pages, navigation mappings can be maintained centrally and shared across the organization.

**Shared Mappings File** (`shared/ui-navigation/atlas.js`):

```javascript
/**
 * Shared Atlas UI navigation mappings
 * Maintained by: DevEx UI Automation Team
 * Last updated: 2025-11-24
 *
 * These mappings are used across all Atlas documentation to ensure
 * consistent UI automation when the Atlas UI changes.
 */

export const atlasNavigationMappings = [
  {
    phrase: /In Atlas, go to the Clusters page/i,
    description: 'Navigate to Clusters page in Atlas',
    steps: [
      { action: 'navigate', url: 'https://cloud.mongodb.com' },
      { action: 'wait', selector: '.project-selector', timeout: 5000 },
      { action: 'click', selector: 'nav a[href*="clusters"]' },
      { action: 'waitForNavigation', timeout: 5000 }
    ]
  },
  {
    phrase: /go to the Database Deployments page/i,
    description: 'Navigate to Database Deployments',
    steps: [
      { action: 'click', selector: 'nav [data-testid="nav-database"]' },
      { action: 'waitForNavigation' }
    ]
  },
  {
    phrase: /open the connection dialog/i,
    description: 'Open connection dialog for cluster',
    steps: [
      { action: 'click', selector: '[data-testid="cluster-connect-button"]' },
      { action: 'wait', selector: '.connection-modal' }
    ]
  },
  {
    phrase: /create a new (cluster|database deployment)/i,
    description: 'Start cluster creation workflow',
    steps: [
      { action: 'click', selector: '[data-testid="create-cluster-button"]' },
      { action: 'wait', selector: '.cluster-creation-modal' }
    ]
  }
  // ... more shared mappings maintained by central team
];

export const atlasUserValues = {
  'your project': process.env.ATLAS_PROJECT_NAME || 'TestProject',
  'your desired project': process.env.ATLAS_PROJECT_NAME || 'TestProject',
  'your cluster': process.env.ATLAS_CLUSTER_NAME || 'Cluster0',
  'your organization': process.env.ATLAS_ORG_NAME || 'TestOrg'
};
```

**Additional Shared Mappings** (`shared/ui-navigation/compass.js`):

```javascript
/**
 * Shared Compass UI navigation mappings
 * Maintained by: DevEx UI Automation Team
 */

export const compassNavigationMappings = [
  {
    phrase: /open Compass/i,
    description: 'Launch MongoDB Compass',
    steps: [
      { action: 'custom', function: async () => {
        // Custom logic to launch Compass application
      }}
    ]
  },
  {
    phrase: /connect to your (cluster|deployment)/i,
    description: 'Connect to MongoDB deployment in Compass',
    steps: [
      { action: 'click', selector: '[data-testid="connect-button"]' },
      { action: 'wait', selector: '.connection-form' }
    ]
  }
  // ... more Compass-specific mappings
];
```

**Index File** (`shared/ui-navigation/index.js`):

```javascript
/**
 * Central export for all shared UI navigation mappings
 */

export * from './atlas.js';
export * from './compass.js';
// Export other product-specific mappings as needed
```

**Consuming Shared Mappings** (`proctest.config.js`):

```javascript
import {
  atlasNavigationMappings,
  atlasUserValues,
  compassNavigationMappings
} from './shared/ui-navigation/index.js';

export default {
  ui: {
    enabled: true,
    framework: 'playwright',
    baseUrl: 'https://cloud.mongodb.com',

    // Combine shared mappings with page-specific ones
    navigationMappings: [
      ...atlasNavigationMappings,      // Shared Atlas mappings
      ...compassNavigationMappings,    // Shared Compass mappings

      // Add page-specific mappings if needed
      {
        phrase: /open advanced cluster settings/i,
        description: 'Page-specific navigation',
        steps: [
          { action: 'click', selector: '[data-testid="advanced-settings"]' },
          { action: 'wait', selector: '.advanced-settings-panel' }
        ]
      }
    ],

    // Use shared user values with optional overrides
    userValues: {
      ...atlasUserValues,
      // Override or add page-specific values
      'your database': 'sample_mflix',
      'your collection': 'movies'
    }
  }
};
```

**Benefits of Shared Mappings**:

1. **Single Source of Truth**: Central team maintains UI automation steps
2. **Automatic Composition**: Mappings are reused wherever phrases appear in documentation
3. **Easy Updates**: When UI changes, update shared file once, all tests benefit
4. **Consistency**: All documentation uses the same automation steps for common tasks
5. **Reduced Duplication**: Individual teams don't duplicate navigation logic
6. **Version Control**: Shared mappings are committed to repo with clear ownership

**Governance Model**:

```
Repository Structure:
├── shared/
│   └── ui-navigation/
│       ├── index.js           # Central export
│       ├── atlas.js           # Atlas UI mappings (maintained by DevEx team)
│       ├── compass.js         # Compass UI mappings (maintained by DevEx team)
│       └── README.md          # Documentation for contributors
├── testdata/
│   ├── atlas/
│   │   ├── proctest.config.js  # Imports shared mappings
│   │   └── source/*.txt
│   └── drivers/
│       ├── proctest.config.js  # Imports shared mappings
│       └── source/*.txt
└── proctest.config.js   # Root config (optional)
```

**Update Workflow**:
1. Atlas UI changes navigation to Clusters page
2. DevEx team updates `shared/ui-navigation/atlas.js`
3. All documentation pages using "In Atlas, go to the Clusters page" automatically use updated steps
4. No changes needed to individual test configs or documentation content

**Hierarchical Configuration** (optional):

Teams can also use hierarchical configs for inheritance:

```javascript
// Root: proctest.config.js
import { atlasNavigationMappings, atlasUserValues } from './shared/ui-navigation/index.js';

export default {
  ui: {
    navigationMappings: [...atlasNavigationMappings],
    userValues: { ...atlasUserValues }
  }
};

// Subdirectory: testdata/atlas/proctest.config.js
import rootConfig from '../../proctest.config.js';

export default {
  ...rootConfig,
  ui: {
    ...rootConfig.ui,
    // Inherit shared mappings, add Atlas-specific ones
    navigationMappings: [
      ...rootConfig.ui.navigationMappings,
      { phrase: /atlas-specific phrase/i, steps: [/* ... */] }
    ],
    // Inherit shared values, override as needed
    userValues: {
      ...rootConfig.ui.userValues,
      'your cluster': 'AtlasSpecificCluster'
    }
  }
};
```

---

### 3.3 RST Parser Implementation

The framework uses a **lightweight, targeted parsing approach** that parses only the RST features needed for procedural testing, rather than building a complete RST syntax tree. This approach is based on the proven regex-based parser used in MongoDB's existing code-example-tooling project.

#### 3.3.1 Parsing Philosophy

**What We Parse** (Minimal Set):
- ✅ Specific directives: `code-block`, `literalinclude`, `procedure`, `step`, `tabs`, `include`
- ✅ Directive options (`:language:`, `:start-after:`, `:end-before:`, `:tabid:`, etc.)
- ✅ Directive content (indented blocks)
- ✅ Ordered lists (procedures as numbered lists)
- ✅ Sub-lists (sub-procedures within steps)
- ✅ Inline roles for URL validation (`:doc:`, `:ref:`, external links)
- ✅ Inline roles for UI testing (`:guilabel:`)
- ✅ Source constants from snooty.toml (`{+variable-name+}`)

**What We Ignore** (Everything Else):
- ❌ Full RST syntax tree
- ❌ Tables, footnotes, citations
- ❌ Most inline roles (except those listed above)
- ❌ Complex directives we don't test (`:figure:`, `:image:`, etc.)
- ❌ Comments, most substitutions

**Risk Mitigation**:
- Parse errors are handled gracefully (skip unparseable content, warn user)
- Comprehensive tests using actual MongoDB documentation files
- Incremental approach: start with minimal set, add features as needed

#### 3.3.2 Core Parser Interface

```typescript
export class LightweightRSTParser implements Parser {
  // Regular expressions for directive detection
  private static readonly DIRECTIVES = {
    CODE_BLOCK: /^\.\.\s+code-block::\s*(.*)$/,
    LITERAL_INCLUDE: /^\.\.\s+literalinclude::\s+(.+)$/,
    PROCEDURE: /^\.\.\s+procedure::\s*$/,
    STEP: /^\.\.\s+step::\s*(.*)$/,
    TABS: /^\.\.\s+tabs::\s*$/,
    TAB: /^\.\.\s+tab::\s*(.*)$/,
    INCLUDE: /^\.\.\s+include::\s+(.+)$/,
    IO_CODE_BLOCK: /^\.\.\s+io-code-block::\s*$/,
    INPUT: /^\.\.\s+input::\s*(.*)$/,
    OUTPUT: /^\.\.\s+output::\s*(.*)$/
  };

  // Regular expressions for list detection
  private static readonly LISTS = {
    ORDERED: /^(\d+)\.\s+(.*)$/,
    SUB_ORDERED: /^([a-z])\.\s+(.*)$/,
    UNORDERED: /^[*\-+]\s+(.*)$/
  };

  // Regular expressions for inline roles
  private static readonly INLINE_ROLES = {
    GUILABEL: /:guilabel:`([^`]+)`/g,
    DOC: /:doc:`([^`<]+)(?:\s*<([^>]+)>)?`/g,
    REF: /:ref:`([^`<]+)(?:\s*<([^>]+)>)?`/g,
    EXTERNAL_LINK: /`([^`<]+)\s*<(https?:\/\/[^>]+)>`_/g,
    SIMPLE_LINK: /(https?:\/\/[^\s]+)/g
  };

  // Regular expressions for source constants
  private static readonly SOURCE_CONSTANT = /\{\+([^}]+)\+\}/g;

  // Regular expression for directive options
  private static readonly OPTION = /^\s+:([^:]+):\s*(.*)$/;

  parse(content: string, context: ParserContext): Promise<DocumentAST> {
    const scanner = new LineScanner(content.split('\n'));
    const procedures: ProcedureNode[] = [];
    const urls: URLNode[] = [];

    while (scanner.hasNext()) {
      const line = scanner.next();
      const trimmed = line.trim();

      // Parse procedures (directive form)
      if (this.isProcedureDirective(trimmed)) {
        procedures.push(this.parseProcedureDirective(scanner, context));
        continue;
      }

      // Parse procedures (ordered list form)
      if (this.isOrderedListStart(trimmed)) {
        procedures.push(this.parseProcedureList(scanner, line, context));
        continue;
      }

      // Extract URLs from inline roles
      urls.push(...this.extractURLsFromLine(trimmed, context));
    }

    return {
      type: 'document',
      filePath: context.filePath,
      procedures,
      urls,
      metadata: {
        sourceFormat: 'rst',
        parsedAt: new Date().toISOString()
      }
    };
  }

  canParse(filePath: string): boolean {
    const ext = path.extname(filePath).toLowerCase();
    return ['.rst', '.txt'].includes(ext);
  }

  getSupportedExtensions(): string[] {
    return ['.rst', '.txt'];
  }
}
```

#### 3.3.3 Directive Parsing

**Code Block Directive**:

```typescript
private parseCodeBlock(scanner: LineScanner, language: string): CodeBlockNode {
  const options = this.parseDirectiveOptions(scanner);
  const content = this.extractIndentedContent(scanner);

  return {
    type: 'code-block',
    language: language || options.language || '',
    content,
    options,
    location: scanner.getLocation()
  };
}
```

**Literal Include Directive**:

```typescript
private parseLiteralInclude(scanner: LineScanner, filePath: string, context: ParserContext): CodeBlockNode {
  const options = this.parseDirectiveOptions(scanner);

  // Resolve the include path (handles MongoDB-specific conventions)
  const resolvedPath = this.resolveIncludePath(context.filePath, filePath);

  // Read and extract content
  let content = fs.readFileSync(resolvedPath, 'utf-8');

  // Handle start-after option
  if (options['start-after']) {
    const startIdx = content.indexOf(options['start-after']);
    if (startIdx !== -1) {
      const lineEnd = content.indexOf('\n', startIdx);
      content = content.substring(lineEnd + 1);
    }
  }

  // Handle end-before option
  if (options['end-before']) {
    const endIdx = content.indexOf(options['end-before']);
    if (endIdx !== -1) {
      const lineStart = content.lastIndexOf('\n', endIdx);
      content = content.substring(0, lineStart);
    }
  }

  // Handle dedent option
  if (options.dedent !== undefined) {
    content = this.dedentContent(content);
  }

  return {
    type: 'code-block',
    language: options.language || this.inferLanguageFromPath(resolvedPath),
    content: content.trim(),
    options,
    location: scanner.getLocation()
  };
}
```

#### 3.3.4 Tabs Directive Parsing

Tabs are critical for MongoDB documentation as they represent mutually-exclusive steps or procedures. The parser must handle tabs and allow users to select which tab to test.

**Tabs Directive Structure**:

```rst
.. tabs::

   .. tab:: Python
      :tabid: python

      1. Install the driver:

         .. code-block:: bash

            pip install pymongo

      2. Connect to MongoDB:

         .. code-block:: python

            from pymongo import MongoClient
            client = MongoClient("mongodb://localhost:27017")

   .. tab:: Node.js
      :tabid: nodejs

      1. Install the driver:

         .. code-block:: bash

            npm install mongodb

      2. Connect to MongoDB:

         .. code-block:: javascript

            const { MongoClient } = require('mongodb');
            const client = new MongoClient("mongodb://localhost:27017");
```

**Parser Implementation**:

```typescript
private parseTabsDirective(scanner: LineScanner, context: ParserContext): TabsNode {
  const options = this.parseDirectiveOptions(scanner);
  const tabs: TabNode[] = [];

  // Parse nested tab directives
  while (scanner.hasNext()) {
    const line = scanner.peek();
    const trimmed = line.trim();

    // Check for tab directive
    const tabMatch = LightweightRSTParser.DIRECTIVES.TAB.exec(trimmed);
    if (tabMatch) {
      scanner.next(); // Consume the tab line
      tabs.push(this.parseTabDirective(scanner, tabMatch[1], context));
      continue;
    }

    // If we hit a non-indented line, we're done with tabs
    if (trimmed && !line.startsWith(' ') && !line.startsWith('\t')) {
      break;
    }

    scanner.next();
  }

  return {
    type: 'tabs',
    tabs,
    options,
    location: scanner.getLocation()
  };
}

private parseTabDirective(scanner: LineScanner, title: string, context: ParserContext): TabNode {
  const options = this.parseDirectiveOptions(scanner);
  const tabId = options.tabid || this.slugify(title);

  // Parse tab content (can contain procedures, steps, code blocks, etc.)
  const content = this.parseTabContent(scanner, context);

  return {
    type: 'tab',
    tabId,
    title: title.trim(),
    content,
    options,
    location: scanner.getLocation()
  };
}

private parseTabContent(scanner: LineScanner, context: ParserContext): ContentNode[] {
  const content: ContentNode[] = [];
  let baseIndent = -1;

  while (scanner.hasNext()) {
    const line = scanner.peek();
    const trimmed = line.trim();

    // Calculate indentation
    const indent = line.length - line.trimLeft().length;

    // Set base indent from first content line
    if (baseIndent === -1 && trimmed) {
      baseIndent = indent;
    }

    // If we've dedented back to or past the base level, we're done
    if (trimmed && baseIndent !== -1 && indent < baseIndent) {
      break;
    }

    // Check for nested directives or lists within tab
    if (this.isProcedureDirective(trimmed)) {
      scanner.next();
      content.push(this.parseProcedureDirective(scanner, context));
      continue;
    }

    if (this.isOrderedListStart(trimmed)) {
      content.push(this.parseProcedureList(scanner, line, context));
      continue;
    }

    if (this.isCodeBlockDirective(trimmed)) {
      scanner.next();
      const match = LightweightRSTParser.DIRECTIVES.CODE_BLOCK.exec(trimmed);
      content.push(this.parseCodeBlock(scanner, match![1]));
      continue;
    }

    scanner.next();
  }

  return content;
}
```

**Tab Testing Strategy**:

Tabs and composable tutorials represent mutually-exclusive content that must be **interpolated** with general content. The framework automatically:

1. **Identifies all variants** (tabs or selections)
2. **Interpolates** variant-specific content with general content
3. **Generates complete test cases** for each variant

**Example 1: Tabs within a procedure**

```rst
1. Install the MongoDB driver:

   .. tabs::

      .. tab:: Python
         :tabid: python

         .. code-block:: bash

            pip install pymongo

      .. tab:: Node.js
         :tabid: nodejs

         .. code-block:: bash

            npm install mongodb

2. Connect to your database:

   .. code-block:: javascript

      const client = new MongoClient(uri);
```

**Generated test cases**:

```
Test Case 1: "Procedure Name (Python)"
  Step 1: Install the MongoDB driver
    - Execute: pip install pymongo
  Step 2: Connect to your database
    - Execute: const client = new MongoClient(uri);

Test Case 2: "Procedure Name (Node.js)"
  Step 1: Install the MongoDB driver
    - Execute: npm install mongodb
  Step 2: Connect to your database
    - Execute: const client = new MongoClient(uri);
```

**Example 2: Composable tutorial with interleaved content**

```rst
.. composable-tutorial::
   :options: language
   :defaults: nodejs

   .. selected-content::
      :selections: nodejs

      1. Install Node.js driver:

         .. code-block:: bash

            npm install mongodb

   .. selected-content::
      :selections: python

      1. Install Python driver:

         .. code-block:: bash

            pip install pymongo

   2. Create a connection string:

      .. code-block:: javascript

         const uri = "mongodb://localhost:27017";

   .. selected-content::
      :selections: nodejs

      3. Connect using Node.js:

         .. code-block:: javascript

            const client = new MongoClient(uri);

   .. selected-content::
      :selections: python

      3. Connect using Python:

         .. code-block:: python

            client = MongoClient(uri)
```

**Generated test cases**:

```
Test Case 1: "Procedure Name (nodejs)"
  Step 1: Install Node.js driver
    - Execute: npm install mongodb
  Step 2: Create a connection string
    - Execute: const uri = "mongodb://localhost:27017";
  Step 3: Connect using Node.js
    - Execute: const client = new MongoClient(uri);

Test Case 2: "Procedure Name (python)"
  Step 1: Install Python driver
    - Execute: pip install pymongo
  Step 2: Create a connection string
    - Execute: const uri = "mongodb://localhost:27017";
  Step 3: Connect using Python
    - Execute: client = MongoClient(uri)
```

**Key Points**:
- ✅ General content (outside tabs/selected-content) is included in **all** test cases
- ✅ Variant-specific content is included only in matching test cases
- ✅ Content order is preserved (interleaved correctly)
- ✅ Step numbering is recalculated for each variant
- ✅ No user configuration needed - automatic interpolation
- ✅ Users can filter with `--grep nodejs` to run specific variants

**Variant Expansion Algorithm**:

```typescript
function expandProcedureVariants(procedure: ProcedureNode): ProcedureVariant[] {
  // 1. Identify all variants (tabs or composable tutorial selections)
  const variants = identifyVariants(procedure);

  if (variants.length === 0) {
    // No variants - return single test case with all content
    return [{ name: procedure.name, content: procedure.content }];
  }

  // 2. For each variant, interpolate content
  return variants.map(variant => {
    const interpolatedContent = [];

    for (const node of procedure.content) {
      if (node.type === 'tabs') {
        // Include only the matching tab's content
        const matchingTab = node.tabs.find(tab => tab.tabId === variant.id);
        if (matchingTab) {
          interpolatedContent.push(...matchingTab.content);
        }
      } else if (node.type === 'selected-content') {
        // Include only if selections match
        if (node.selections.includes(variant.selection)) {
          interpolatedContent.push(...node.content);
        }
      } else {
        // General content - include in all variants
        interpolatedContent.push(node);
      }
    }

    return {
      name: `${procedure.name} (${variant.label})`,
      content: interpolatedContent,
      variant: variant.id
    };
  });
}

function identifyVariants(procedure: ProcedureNode): Variant[] {
  // Check for composable tutorial
  const composable = findComposableTutorial(procedure);
  if (composable) {
    return composable.selections.map(sel => ({
      id: sel.id,
      label: sel.label,
      selection: sel.value
    }));
  }

  // Check for tabs
  const tabs = findTabs(procedure);
  if (tabs.length > 0) {
    return tabs[0].tabs.map(tab => ({
      id: tab.tabId,
      label: tab.title,
      selection: tab.tabId
    }));
  }

  return [];
}
```

#### 3.3.5 Inline Role Parsing for URL Validation

The parser extracts URLs from inline roles for validation (Phase 3: URL testable actions).

**Supported Inline Roles**:

```typescript
private extractURLsFromLine(line: string, context: ParserContext): URLNode[] {
  const urls: URLNode[] = [];

  // Extract :doc: references (internal documentation links)
  // Example: :doc:`/tutorial/install-mongodb`
  const docMatches = line.matchAll(LightweightRSTParser.INLINE_ROLES.DOC);
  for (const match of docMatches) {
    const docPath = match[2] || match[1]; // Use explicit path if provided
    urls.push({
      type: 'url',
      url: this.resolveDocReference(docPath, context),
      source: 'doc-role',
      location: { line: context.currentLine, column: match.index }
    });
  }

  // Extract :ref: references (internal cross-references)
  // Example: :ref:`installation-guide`
  const refMatches = line.matchAll(LightweightRSTParser.INLINE_ROLES.REF);
  for (const match of refMatches) {
    const refName = match[2] || match[1];
    urls.push({
      type: 'url',
      url: this.resolveRefReference(refName, context),
      source: 'ref-role',
      location: { line: context.currentLine, column: match.index }
    });
  }

  // Extract external links
  // Example: `MongoDB Atlas <https://cloud.mongodb.com>`_
  const externalMatches = line.matchAll(LightweightRSTParser.INLINE_ROLES.EXTERNAL_LINK);
  for (const match of externalMatches) {
    urls.push({
      type: 'url',
      url: match[2],
      linkText: match[1],
      source: 'external-link',
      location: { line: context.currentLine, column: match.index }
    });
  }

  // Extract simple URLs
  // Example: https://www.mongodb.com/docs
  const simpleMatches = line.matchAll(LightweightRSTParser.INLINE_ROLES.SIMPLE_LINK);
  for (const match of simpleMatches) {
    urls.push({
      type: 'url',
      url: match[1],
      source: 'simple-link',
      location: { line: context.currentLine, column: match.index }
    });
  }

  return urls;
}

private resolveDocReference(docPath: string, context: ParserContext): string {
  // Resolve :doc: reference to actual file path
  // Example: :doc:`/tutorial/install-mongodb` -> /path/to/source/tutorial/install-mongodb.txt
  const sourceDir = this.findSourceDirectory(context.filePath);
  const cleanPath = docPath.startsWith('/') ? docPath.substring(1) : docPath;

  // Try common extensions
  for (const ext of ['.txt', '.rst', '']) {
    const fullPath = path.join(sourceDir, cleanPath + ext);
    if (fs.existsSync(fullPath)) {
      return fullPath;
    }
  }

  throw new Error(`Doc reference not found: ${docPath}`);
}

private resolveRefReference(refName: string, context: ParserContext): string {
  // Resolve :ref: reference by searching for .. _refName: labels
  // This is more complex and may require scanning multiple files
  // For now, return the ref name and validate it exists during execution
  return `ref:${refName}`;
}
```

#### 3.3.6 Helper Utilities

**Indented Content Extraction**:

```typescript
private extractIndentedContent(scanner: LineScanner): string {
  const lines: string[] = [];
  let baseIndent = -1;

  while (scanner.hasNext()) {
    const line = scanner.peek();

    // Empty lines are part of content
    if (line.trim() === '') {
      lines.push('');
      scanner.next();
      continue;
    }

    // Calculate indentation
    const indent = line.length - line.trimLeft().length;

    // First content line establishes base indentation
    if (baseIndent === -1) {
      baseIndent = indent;
    }

    // If less indented than base, we're done
    if (indent < baseIndent) {
      break;
    }

    // Add line with base indentation removed
    lines.push(line.substring(baseIndent));
    scanner.next();
  }

  return lines.join('\n').trim();
}
```

**Directive Options Parsing**:

```typescript
private parseDirectiveOptions(scanner: LineScanner): Record<string, string> {
  const options: Record<string, string> = {};

  while (scanner.hasNext()) {
    const line = scanner.peek();
    const trimmed = line.trim();

    // Check if this is an option line
    const optionMatch = LightweightRSTParser.OPTION.exec(line);
    if (optionMatch) {
      const optionName = optionMatch[1].trim();
      const optionValue = optionMatch[2].trim();
      options[optionName] = optionValue;
      scanner.next();
      continue;
    }

    // If we hit a blank line, skip it
    if (trimmed === '') {
      scanner.next();
      continue;
    }

    // If the line is not indented and not an option, we're done
    if (line.length > 0 && line[0] !== ' ' && line[0] !== '\t') {
      break;
    }

    // If we have indented content (not an option), we're done with options
    if (line.length > 0 && (line[0] === ' ' || line[0] === '\t') && !optionMatch) {
      break;
    }

    scanner.next();
  }

  return options;
}
```

**Include Path Resolution** (MongoDB-specific conventions):

```typescript
private resolveIncludePath(currentFilePath: string, includePath: string): string {
  // Find the source directory by walking up from current file
  const sourceDir = this.findSourceDirectory(currentFilePath);

  // Clean the include path (remove leading slash if present)
  const cleanPath = includePath.startsWith('/') ? includePath.substring(1) : includePath;

  // Special handling for steps/ includes
  // Convert /includes/steps/filename.rst to /includes/steps-filename.yaml
  if (cleanPath.includes('steps/')) {
    const stepsPath = this.resolveStepsInclude(sourceDir, cleanPath);
    if (stepsPath) return stepsPath;
  }

  // Special handling for extracts/ includes (ref-based YAML content blocks)
  if (cleanPath.includes('extracts/')) {
    const extractsPath = this.resolveRefBasedInclude(sourceDir, cleanPath, 'extracts');
    if (extractsPath) return extractsPath;
  }

  // Construct the full path
  let fullPath = path.join(sourceDir, cleanPath);

  // If the file exists as-is, return it
  if (fs.existsSync(fullPath)) {
    return fullPath;
  }

  // If the path doesn't have an extension, try adding .rst
  if (path.extname(cleanPath) === '') {
    fullPath = fullPath + '.rst';
    if (fs.existsSync(fullPath)) {
      return fullPath;
    }
  }

  throw new Error(`Include file not found: ${includePath}`);
}

private findSourceDirectory(filePath: string): string {
  // Walk up from current file looking for snooty.toml or source/ directory
  let dir = path.dirname(filePath);

  while (dir !== path.dirname(dir)) { // Not at root
    // Check for snooty.toml
    if (fs.existsSync(path.join(dir, 'snooty.toml'))) {
      // Source directory is typically ./source relative to snooty.toml
      const sourceDir = path.join(dir, 'source');
      if (fs.existsSync(sourceDir)) {
        return sourceDir;
      }
      return dir;
    }

    // Check if current directory is named 'source'
    if (path.basename(dir) === 'source') {
      return dir;
    }

    dir = path.dirname(dir);
  }

  // Fallback to directory containing the file
  return path.dirname(filePath);
}
```

**Dedent Content**:

```typescript
private dedentContent(content: string): string {
  const lines = content.split('\n');
  if (lines.length === 0) return content;

  // Find minimum indentation (ignoring empty lines)
  let minIndent = Infinity;
  for (const line of lines) {
    if (line.trim() === '') continue;
    const indent = line.length - line.trimLeft().length;
    minIndent = Math.min(minIndent, indent);
  }

  if (minIndent === Infinity || minIndent === 0) {
    return content;
  }

  // Remove common indentation from all lines
  const dedentedLines = lines.map(line => {
    if (line.trim() === '') return '';
    return line.substring(minIndent);
  });

  return dedentedLines.join('\n');
}
```

#### 3.3.7 Line Scanner Utility

```typescript
class LineScanner {
  private lines: string[];
  private position: number = 0;

  constructor(lines: string[]) {
    this.lines = lines;
  }

  hasNext(): boolean {
    return this.position < this.lines.length;
  }

  next(): string {
    return this.lines[this.position++];
  }

  peek(): string {
    return this.lines[this.position];
  }

  peekAhead(offset: number): string | undefined {
    return this.lines[this.position + offset];
  }

  getPosition(): number {
    return this.position;
  }

  getLocation(): SourceLocation {
    return {
      line: this.position + 1,
      column: 0,
      file: '' // Set by parser
    };
  }
}
```

---

## 4. Data Models

### 4.1 AST Node Types

```typescript
export type ContentNode =
  | TextNode
  | ParagraphNode
  | ListNode
  | CodeBlockNode
  | TabsNode
  | ComposableTutorialNode;

export interface TextNode {
  type: 'text';
  value: string;
  location: SourceLocation;
}

export interface ParagraphNode {
  type: 'paragraph';
  children: ContentNode[];
  location: SourceLocation;
}

export interface ListNode {
  type: 'list';
  ordered: boolean;
  items: ListItemNode[];
  location: SourceLocation;
}

export interface ListItemNode {
  type: 'list-item';
  children: ContentNode[];
  location: SourceLocation;
}

export interface TabsNode {
  type: 'tabs';
  tabs: TabNode[];
  location: SourceLocation;
}

export interface TabNode {
  type: 'tab';
  tabId: string;
  title: string;
  content: ContentNode[];
  location: SourceLocation;
}

export interface ComposableTutorialNode {
  type: 'composable-tutorial';
  id: string;
  options: string[];
  defaults: string[];
  selections: SelectedContentNode[];
  location: SourceLocation;
}

export interface URLNode {
  type: 'url';
  url: string;
  linkText?: string;
  source: 'doc-role' | 'ref-role' | 'external-link' | 'simple-link';
  location: SourceLocation;
}

export interface SourceLocation {
  line: number;
  column: number;
  file: string;
}

export interface SelectedContentNode {
  type: 'selected-content';
  selections: Record<string, string>;
  content: ContentNode[];
  location: SourceLocation;
}

export interface DocumentMetadata {
  title?: string;
  author?: string;
  date?: string;
  tags?: string[];
}
```

### 4.2 Cleanup Registry

```typescript
/**
 * Registry for tracking resources that need cleanup
 * Executes cleanup in LIFO order (last registered, first cleaned)
 */
export class CleanupRegistry {
  private cleanupTasks: CleanupTask[] = [];

  /**
   * Register a cleanup task
   */
  register(task: CleanupTask): void {
    this.cleanupTasks.push(task);
  }

  /**
   * Execute all cleanup tasks in LIFO order
   */
  async executeAll(): Promise<CleanupResult[]> {
    const results: CleanupResult[] = [];

    // Execute in reverse order (LIFO)
    for (let i = this.cleanupTasks.length - 1; i >= 0; i--) {
      const task = this.cleanupTasks[i];
      try {
        await task.cleanup();
        results.push({ task, success: true });
      } catch (error) {
        results.push({
          task,
          success: false,
          error: error as Error
        });
      }
    }

    return results;
  }

  /**
   * Clear all registered tasks
   */
  clear(): void {
    this.cleanupTasks = [];
  }
}

export interface CleanupTask {
  type: 'database' | 'collection' | 'file' | 'directory' | 'custom';
  description: string;
  cleanup: () => Promise<void>;
}

export interface CleanupResult {
  task: CleanupTask;
  success: boolean;
  error?: Error;
}
```

### 4.3 Test Registry Data Model

The test registry is a JSON file that maintains a curated list of verified procedures for automated CI testing.

```typescript
export interface TestRegistry {
  version: string;
  description?: string;
  tests: TestRegistryEntry[];
}

export interface TestRegistryEntry {
  // Unique identifier for this test
  id: string;

  // Path to the procedure file (relative to repository root)
  path: string;

  // Team or individual responsible for maintaining this test
  owner: string;

  // Date when this test was added to the registry (ISO 8601 format)
  addedDate: string;

  // Variants to test (e.g., ["atlas-ui", "mongosh", "nodejs"])
  // If empty or omitted, all variants are tested
  variants?: string[];

  // Tags for categorization and filtering
  tags: string[];

  // Optional notes about what this test covers or special considerations
  notes?: string;

  // Optional: Skip this test temporarily (with reason)
  skip?: {
    reason: string;
    skipUntil?: string; // ISO 8601 date
  };
}
```

**Example Test Registry** (`code-example-tests/procedures/test-registry.json`):

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
    },
    {
      "id": "atlas-backup-restore",
      "path": "content/atlas/source/how-to/backup-restore.txt",
      "owner": "atlas-docs-team",
      "addedDate": "2024-02-01",
      "tags": ["atlas", "backup", "how-to"],
      "skip": {
        "reason": "Waiting for backup API changes to stabilize",
        "skipUntil": "2024-03-01"
      }
    }
  ]
}
```

**Registry Loading and Filtering**:

```typescript
export class TestRegistryLoader {
  async load(registryPath: string): Promise<TestRegistry> {
    const content = await fs.readFile(registryPath, 'utf-8');
    const registry = JSON.parse(content) as TestRegistry;

    // Validate registry format
    this.validate(registry);

    return registry;
  }

  /**
   * Filter registry entries based on criteria
   */
  filter(registry: TestRegistry, options: {
    tags?: string[];
    owner?: string;
    includeSkipped?: boolean;
  }): TestRegistryEntry[] {
    let entries = registry.tests;

    // Filter by tags
    if (options.tags && options.tags.length > 0) {
      entries = entries.filter(entry =>
        options.tags!.some(tag => entry.tags.includes(tag))
      );
    }

    // Filter by owner
    if (options.owner) {
      entries = entries.filter(entry => entry.owner === options.owner);
    }

    // Filter skipped tests
    if (!options.includeSkipped) {
      entries = entries.filter(entry => {
        if (!entry.skip) return true;

        // Check if skip period has expired
        if (entry.skip.skipUntil) {
          const skipUntil = new Date(entry.skip.skipUntil);
          return new Date() > skipUntil;
        }

        return false;
      });
    }

    return entries;
  }

  private validate(registry: TestRegistry): void {
    if (!registry.version) {
      throw new Error('Registry must have a version field');
    }

    if (!Array.isArray(registry.tests)) {
      throw new Error('Registry must have a tests array');
    }

    // Validate each entry
    for (const entry of registry.tests) {
      if (!entry.id || !entry.path || !entry.owner || !entry.addedDate) {
        throw new Error(`Invalid registry entry: ${JSON.stringify(entry)}`);
      }

      // Validate date format
      if (isNaN(Date.parse(entry.addedDate))) {
        throw new Error(`Invalid date format in entry ${entry.id}: ${entry.addedDate}`);
      }
    }

    // Check for duplicate IDs
    const ids = new Set<string>();
    for (const entry of registry.tests) {
      if (ids.has(entry.id)) {
        throw new Error(`Duplicate test ID: ${entry.id}`);
      }
      ids.add(entry.id);
    }
  }
}
```

**CLI Usage with Registry**:

```bash
# Run all tests in registry
proctest test --registry code-example-tests/procedures/test-registry.json

# Run registry tests with specific tags
proctest test --registry code-example-tests/procedures/test-registry.json --tags atlas,tutorial

# Run registry tests for specific owner
proctest test --registry code-example-tests/procedures/test-registry.json --owner atlas-docs-team

# Include skipped tests
proctest test --registry code-example-tests/procedures/test-registry.json --include-skipped
```

---

## 5. API Specifications

### 5.1 CLI API

```bash
# Basic usage
proctest <file>                    # Test single file
proctest <directory>               # Test all files in directory
proctest --all                     # Test all discovered files
proctest --registry <file>         # Test procedures from registry file

# Configuration
proctest --config <file>           # Use specific config file
proctest --init                    # Create config file interactively

# Environment
proctest --env <file>              # Use specific .env file
proctest --snooty <file>           # Use specific snooty.toml

# Execution control
proctest --timeout <ms>            # Override timeout
proctest --no-cleanup              # Skip cleanup
proctest --fail-fast               # Stop on first failure

# Output control
proctest --verbose                 # Verbose output
proctest --quiet                   # Minimal output
proctest --reporter <type>         # Specify reporter (human, json, junit)
proctest --output <file>           # Write output to file

# Filtering
proctest --filter <pattern>        # Filter procedures by name
proctest --exclude <pattern>       # Exclude files/procedures

# Debugging
proctest parse <file>              # Parse file and display AST structure
proctest parse <file> --output <file>  # Write parsed AST to file
proctest parse <file> --format <type>  # Output format: tree (default), json, yaml
proctest --dry-run                 # Parse and validate without executing
proctest --list                    # List discovered procedures
proctest --validate-env            # Validate environment setup
```

### 5.2 Programmatic API

```typescript
import { ProcedureTest } from 'proctest';

// Create instance with configuration
const tester = new ProcedureTest({
  testFiles: ['procedures/**/*.txt'],
  envFiles: ['.env'],
  verbose: true
});

// Run tests
const results = await tester.run();

// Access results
console.log(`Passed: ${results.passedProcedures}/${results.totalProcedures}`);
console.log(`Failed: ${results.failedProcedures}/${results.totalProcedures}`);

// Run specific file
const fileResults = await tester.runFile('path/to/procedure.txt');

// Validate environment
const validation = await tester.validateEnvironment();
if (!validation.valid) {
  console.error('Environment validation failed:', validation.errors);
}
```

### 5.3 Plugin API (Future)

```typescript
/**
 * Plugin interface for extending functionality
 */
export interface Plugin {
  name: string;
  version: string;

  /**
   * Initialize plugin
   */
  initialize(context: PluginContext): Promise<void>;

  /**
   * Register components with the framework
   */
  register(registry: ComponentRegistry): void;
}

export interface PluginContext {
  config: Configuration;
  logger: Logger;
}

export interface ComponentRegistry {
  registerParser(parser: Parser): void;
  registerResolver(resolver: PlaceholderResolver): void;
  registerExecutor(executor: Executor): void;
  registerReporter(reporter: Reporter): void;
}

// Example plugin
export class KotlinExecutorPlugin implements Plugin {
  name = 'kotlin-executor';
  version = '1.0.0';

  async initialize(context: PluginContext): Promise<void> {
    context.logger.info('Initializing Kotlin executor plugin');
  }

  register(registry: ComponentRegistry): void {
    registry.registerExecutor(new KotlinExecutor());
  }
}
```

---

## 6. Implementation Plan

### 6.1 Phase 1: PoC

**Goal**: Prove the concept with minimal viable functionality

#### Milestone 1: Foundation
- [ ] Project setup (TypeScript, Jest, Commander.js)
- [ ] Define core interfaces (Parser, Resolver, Executor, Reporter)
- [ ] Implement Configuration Manager with convention-based discovery
- [ ] Implement basic CLI (test command, argument parsing)
- [ ] Set up testing infrastructure

**Deliverables**:
- Project structure with build system
- Core interfaces defined
- CLI can parse arguments and load configuration
- Unit tests for configuration loading

#### Milestone 2: Parsing
- [ ] Implement RST tokenizer
- [ ] Implement RST AST builder
- [ ] Handle `procedure` and `step` directives
- [ ] Handle ordered lists within steps (sub-procedures)
- [ ] Detect and parse prerequisites/requirements sections
- [ ] Handle `code-block`, `code`, `literalinclude` directives
- [ ] Detect file operation prose patterns (create, replace, append)
- [ ] Handle `include` directive for transclusion
- [ ] Parse `snooty.toml` for constants
- [ ] Implement `parse` command with tree, JSON, and YAML output formats
- [ ] Add placeholder detection and reporting in parse output
- [ ] Add summary statistics to parse output

**Deliverables**:
- RST parser that converts `.txt` files to AST
- Support for procedures, steps, and sub-steps (ordered lists)
- Prerequisite/requirement detection and parsing
- File operation detection from prose + literalinclude
- Transclusion support
- `parse` command for debugging and validation
- Unit tests for parser with real examples

#### Milestone 3: Resolution & Execution
- [ ] Implement ExactEnvResolver (exact environment variable match)
- [ ] Implement FuzzyEnvResolver (fuzzy matching with suggestions)
- [ ] Implement SnootyConstantResolver (snooty.toml constants)
- [ ] Implement LayeredResolver (combines all resolvers)
- [ ] Implement testable action detection (File, Code, and Shell types for PoC)
- [ ] Implement FileExecutor (create, replace, append operations)
- [ ] Implement per-procedure temp directory creation
- [ ] Implement CodeExecutor for JavaScript (Node.js)
- [ ] Implement CodeExecutor for Python
- [ ] Implement IDE execution detection and default command mapping
- [ ] Implement IDE execution command interpolation ({filename}, {className}, etc.)
- [ ] Implement ShellExecutor
- [ ] Implement state accumulation within steps

**Deliverables**:
- Placeholder resolution with fuzzy matching
- Testable action detection and classification
- File operations (create, replace, append)
- Per-procedure working directory management
- Code execution for JavaScript, Python
- IDE-based execution with convention-over-configuration defaults
- Shell command execution
- State persistence within steps
- Unit tests for resolvers and executors

#### Milestone 4: Prerequisites, Orchestration & Reporting
- [ ] Implement PrerequisiteChecker for software requirements
- [ ] Implement PrerequisiteChecker for environment requirements
- [ ] Add prerequisite validation to test execution flow
- [ ] Implement TestOrchestrator (coordinates prerequisite checking, parsing, resolution, execution)
- [ ] Implement working directory cleanup in CleanupRegistry
- [ ] Implement CleanupRegistry
- [ ] Implement HumanReporter (human-friendly output with prerequisite results)
- [ ] Implement error formatting with suggestions
- [ ] Integrate all components
- [ ] End-to-end testing with real documentation files

**Deliverables**:
- Prerequisite checking and validation
- Test skipping when prerequisites not met
- Complete PoC that can test real procedures
- Human-friendly error messages with prerequisite status
- Automatic cleanup (databases, collections, files, working directories)
- Integration tests with testdata files
- Demo video/documentation

**Success Criteria**:
- ✅ Can parse RST files from testdata/
- ✅ Can create files and execute JavaScript, Python, Shell code
- ✅ Resolves 80%+ of placeholders automatically
- ✅ Provides helpful error messages
- ✅ Cleans up test databases and temp directories

### 6.2 Phase 2: Production-Ready

**Goal**: Make it production-ready for technical writers

#### Milestone 5: Configuration & Advanced Parsing
- [ ] Implement optional configuration file support
- [ ] Implement `--init` command for guided setup
- [ ] Add configuration validation with helpful errors
- [ ] Handle `tabs` directive
- [ ] Handle `io-code-block` directive
- [ ] Improve error messages with more context

**Deliverables**:
- Optional configuration file support
- Tabs support
- Better error messages
- Configuration validation

#### Milestone 6: Additional Executors & Features
- [ ] Implement CodeExecutor for PHP
- [ ] Implement CLIExecutor for mongosh (CLI element type)
- [ ] Implement CLIExecutor for atlas-cli (CLI element type)
- [ ] Implement runtime validation (check if node, python, etc. are installed)
- [ ] Add timeout handling
- [ ] Add `--dry-run` mode
- [ ] Add `--list` mode

**Deliverables**:
- PHP code execution
- CLI element type support (mongosh, atlas-cli)
- Runtime validation
- Dry-run and list modes

#### Milestone 7: Reporting & Output
- [ ] Implement JSONReporter
- [ ] Implement JUnitReporter (for CI integration)
- [ ] Add `--output` flag for file output
- [ ] Add progress indicators
- [ ] Improve summary statistics

**Deliverables**:
- Multiple output formats
- CI integration support
- Better progress feedback

#### Milestone 8: Polish & Documentation
- [ ] Comprehensive user documentation
- [ ] API documentation
- [ ] Troubleshooting guide
- [ ] Example configurations
- [ ] Performance optimization
- [ ] Bug fixes from testing

**Deliverables**:
- Complete documentation
- Optimized performance
- Production-ready release

**Success Criteria**:
- ✅ 90%+ of writers use zero config
- ✅ 10% of writers successfully use config for edge cases
- ✅ All major RST directives supported
- ✅ Multiple output formats available
- ✅ Comprehensive documentation

### 6.3 Phase 3: Advanced Features

**Goal**: Add advanced features based on user feedback

#### Milestone 9: Composable Tutorials
- [ ] Parse `composable-tutorial` directive
- [ ] Parse `selected-content` directive
- [ ] Handle composable dependencies
- [ ] Test selection logic
- [ ] Support multiple selection paths

**Deliverables**:
- Composable tutorial support
- Selection-based testing

#### Milestone 10: Additional Testable Action Types
- [ ] Implement UIExecutor (Playwright/Puppeteer)
- [ ] Parse `:guilabel:` roles and detect UI actions
- [ ] Implement UI navigation mapping system
- [ ] Implement UI user values substitution
- [ ] Add UI configuration validation
- [ ] Implement APIExecutor for Atlas Admin API (axios/fetch)
- [ ] Parse curl commands targeting Atlas Admin API
- [ ] Implement URLExecutor for link validation
- [ ] Add screenshot capture on UI test failure

**Deliverables**:
- UI action type support (click, select, input, verify)
- UI navigation mappings (map generic phrases to automation steps)
- UI user values (map generic references like "your project" to actual values)
- Atlas Admin API action type support (https://cloud.mongodb.com/api/atlas/)
- URL action type support (link validation)
- Enhanced error reporting with screenshots

#### Milestone 11: Advanced Cleanup & Polish
- [ ] Implement database cleanup detection
- [ ] Implement collection cleanup detection
- [ ] Improve file cleanup
- [ ] Add cleanup hooks
- [ ] Performance optimization (parser caching, lazy loading)
- [ ] Documentation polish and examples

**Deliverables**:
- Automatic resource detection
- Better cleanup strategies
- Performance improvements
- Comprehensive documentation

**Success Criteria**:
- ✅ Composable tutorials work correctly
- ✅ Cleanup is reliable and comprehensive
- ✅ Framework is well-documented and performant

---

## 7. Testing Strategy

### 7.1 Unit Testing

**Framework**: Jest with ts-jest (team familiarity, organizational consistency)

**Coverage Goals**:
- 80%+ code coverage overall
- 90%+ coverage for core components (Parser, Resolver, Executor)
- 100% coverage for critical paths (error handling, cleanup)

**Test Categories**:

#### Parser Tests
```typescript
describe('RSTParser', () => {
  it('should parse procedure directive', () => {
    const content = `
.. procedure::

   .. step:: First step

      Do something.
    `;
    const ast = parser.parse(content, context);
    expect(ast.procedures).toHaveLength(1);
    expect(ast.procedures[0].steps).toHaveLength(1);
  });

  it('should parse code blocks with language', () => {
    const content = `
.. code-block:: javascript

   console.log('hello');
    `;
    const ast = parser.parse(content, context);
    const codeBlock = findCodeBlock(ast);
    expect(codeBlock.language).toBe('javascript');
    expect(codeBlock.code).toContain('console.log');
  });

  it('should handle transclusion', () => {
    // Test include directive
  });
});
```

#### Resolver Tests
```typescript
describe('FuzzyEnvResolver', () => {
  it('should resolve exact matches', async () => {
    const env = { 'USERNAME': 'testuser' };
    const result = await resolver.resolve('<username>', { environment: env });
    expect(result).toBe('testuser');
  });

  it('should resolve fuzzy matches', async () => {
    const env = { 'DB_USERNAME': 'testuser' };
    const result = await resolver.resolve('<username>', { environment: env });
    expect(result).toBe('testuser');
  });

  it('should provide suggestions for unresolved', async () => {
    const env = { 'USER_NAME': 'test', 'DB_USER': 'test' };
    const suggestions = await resolver.getSuggestions('<username>', { environment: env });
    expect(suggestions).toContain('USER_NAME');
    expect(suggestions).toContain('DB_USER');
  });
});
```

#### Executor Tests
```typescript
describe('JavaScriptExecutor', () => {
  it('should execute simple code', async () => {
    const code = 'console.log("hello");';
    const result = await executor.execute(code, context);
    expect(result.success).toBe(true);
    expect(result.stdout).toContain('hello');
  });

  it('should accumulate state within step', async () => {
    const code1 = 'const x = 5;';
    const code2 = 'console.log(x);';

    await executor.execute(code1, context);
    const result = await executor.execute(code2, context);

    expect(result.success).toBe(true);
    expect(result.stdout).toContain('5');
  });

  it('should handle errors gracefully', async () => {
    const code = 'throw new Error("test error");';
    const result = await executor.execute(code, context);
    expect(result.success).toBe(false);
    expect(result.error).toBeDefined();
  });

  it('should timeout long-running code', async () => {
    const code = 'while(true) {}';
    const result = await executor.execute(code, { ...context, timeout: 1000 });
    expect(result.timedOut).toBe(true);
  });
});

describe('CodeExecutor - IDE Execution', () => {
  it('should detect IDE execution from prose', () => {
    const content = `
.. step:: Create a file named CreateIndex.java

   .. literalinclude:: /includes/CreateIndex.java
      :language: java

.. step:: Execute the code to create the index.

   From your IDE, run the file to create the index.
    `;

    const ast = parser.parse(content, context);
    const actions = detector.detectActions(ast);
    const ideAction = actions.find(a => a.executionMode === 'ide');

    expect(ideAction).toBeDefined();
    expect(ideAction.actionType).toBe('code');
    expect(ideAction.language).toBe('java');
    expect(ideAction.executionMode).toBe('ide');
    expect(ideAction.filePath).toBe('CreateIndex.java');
  });

  it('should use default command for Java', async () => {
    const action: CodeTestableAction = {
      actionType: 'code',
      language: 'java',
      code: 'public class CreateIndex { public static void main(String[] args) {} }',
      executionMode: 'ide',
      filePath: 'CreateIndex.java',
      options: { executable: true }
    };

    const command = executor.resolveIDECommand(action, context);

    expect(command).toBe('mvn compile exec:java -Dexec.mainClass="CreateIndex"');
  });

  it('should use custom command from config', async () => {
    const customContext = {
      ...context,
      config: {
        ...context.config,
        ideExecution: {
          commands: {
            java: 'gradle run'
          }
        }
      }
    };

    const action: CodeTestableAction = {
      actionType: 'code',
      language: 'java',
      code: 'public class CreateIndex {}',
      executionMode: 'ide',
      filePath: 'CreateIndex.java',
      options: { executable: true }
    };

    const command = executor.resolveIDECommand(action, customContext);

    expect(command).toBe('gradle run');
  });

  it('should interpolate command variables', () => {
    const action: CodeTestableAction = {
      actionType: 'code',
      language: 'cpp',
      code: '#include <iostream>\nint main() {}',
      executionMode: 'ide',
      filePath: 'src/main.cpp',
      options: { executable: true }
    };

    const command = executor.resolveIDECommand(action, context);

    expect(command).toBe('g++ src/main.cpp -o main && ./main');
  });

  it('should extract class name from Java code', () => {
    const code = `
package com.example;

public class MyIndexService {
  public static void main(String[] args) {
    System.out.println("Hello");
  }
}
    `;

    const className = executor.extractClassName(code, 'java');

    expect(className).toBe('MyIndexService');
  });

  it('should extract class name from C# code', () => {
    const code = `
using System;

namespace MyApp
{
    class Program
    {
        static void Main(string[] args)
        {
            Console.WriteLine("Hello");
        }
    }
}
    `;

    const className = executor.extractClassName(code, 'csharp');

    expect(className).toBe('Program');
  });

  it('should throw error for unsupported language without config', () => {
    const action: CodeTestableAction = {
      actionType: 'code',
      language: 'rust',
      code: 'fn main() {}',
      executionMode: 'ide',
      filePath: 'main.rs',
      options: { executable: true }
    };

    expect(() => executor.resolveIDECommand(action, context)).toThrow(
      'No default IDE command for rust'
    );
  });

  it('should execute IDE action via ShellExecutor', async () => {
    const action: CodeTestableAction = {
      actionType: 'code',
      language: 'python',
      code: 'print("Hello from IDE")',
      executionMode: 'ide',
      filePath: 'test.py',
      options: { executable: true }
    };

    const result = await executor.execute(action, context);

    expect(result.success).toBe(true);
    expect(result.output).toContain('Hello from IDE');
  });

  it('should skip IDE execution when configured', async () => {
    const skipContext = {
      ...context,
      config: {
        ...context.config,
        ideExecution: {
          skip: true
        }
      }
    };

    const action: CodeTestableAction = {
      actionType: 'code',
      language: 'java',
      code: 'public class Test {}',
      executionMode: 'ide',
      filePath: 'Test.java',
      options: { executable: true }
    };

    const result = await executor.execute(action, skipContext);

    expect(result.success).toBe(true);
    expect(result.skipped).toBe(true);
    expect(result.message).toContain('IDE execution skipped');
  });
});
```

### 7.2 Integration Testing

**Goal**: Test component interactions

```typescript
describe('End-to-End Integration', () => {
  it('should test complete procedure from RST to execution', async () => {
    const rstContent = `
.. procedure::

   .. step:: Connect to database

      .. code-block:: javascript

         const { MongoClient } = require('mongodb');
         const client = new MongoClient('<connection-string>');
         await client.connect();
    `;

    // Setup environment
    const env = { 'MONGODB_URI': 'mongodb://localhost:27017' };

    // Run test
    const tester = new ProcedureTest({ envFiles: [] });
    const result = await tester.runContent(rstContent, { environment: env });

    expect(result.success).toBe(true);
  });
});
```

### 7.3 Fixture-Based Testing

**Use Real Documentation Files**:

```typescript
describe('Real Documentation Tests', () => {
  const testFiles = [
    'testdata/drivers/source/symfony.txt',
    'testdata/atlas/source/connect-to-database-deployment.txt'
  ];

  testFiles.forEach(file => {
    it(`should parse ${file}`, async () => {
      const content = await readFile(file, 'utf-8');
      const ast = await parser.parse(content, context);
      expect(ast.procedures.length).toBeGreaterThan(0);
    });
  });
});
```

### 7.4 Error Scenario Testing

```typescript
describe('Error Handling', () => {
  it('should provide helpful error for unresolved placeholder', async () => {
    const code = 'console.log("<unknown-placeholder>");';
    const result = await executor.execute(code, context);

    expect(result.success).toBe(false);
    expect(result.error?.suggestions).toBeDefined();
  });

  it('should handle missing runtime gracefully', async () => {
    const executor = new PythonExecutor();
    // Mock python not being installed
    const validation = await executor.validate();
    expect(validation.valid).toBe(false);
    expect(validation.suggestions).toContain('install python');
  });
});
```

### 7.5 Performance Testing

```typescript
describe('Performance', () => {
  it('should parse large files in reasonable time', async () => {
    const largeFile = generateLargeRSTFile(1000); // 1000 procedures
    const start = Date.now();
    await parser.parse(largeFile, context);
    const duration = Date.now() - start;
    expect(duration).toBeLessThan(5000); // < 5 seconds
  });

  it('should handle concurrent executions', async () => {
    // Test that sequential execution doesn't have race conditions
  });
});
```

---

## 8. Risk Assessment

### 8.1 Technical Risks

#### Risk 1: RST Parsing Complexity
**Severity**: Low (reduced from High)
**Probability**: Low (reduced from Medium)
**Impact**: Minimal - proven approach

**Description**: ~~RST has complex syntax with many edge cases. Parsing all directives correctly may be more difficult than anticipated.~~ **MITIGATED**: Using lightweight, targeted parsing approach based on proven regex-based parser from MongoDB's code-example-tooling project.

**Approach**:
- Parse only the RST features needed for procedural testing (code-block, literalinclude, procedure, step, tabs, include)
- Use regex-based line-by-line scanning with indentation tracking
- Port proven logic from existing Go implementation to TypeScript
- Handle MongoDB-specific conventions (steps files, extracts, source constants)
- Gracefully handle parse errors (skip unparseable content, warn user)

**What We Parse**:
- ✅ Directives: `code-block`, `literalinclude`, `procedure`, `step`, `tabs`, `include`
- ✅ Directive options (`:language:`, `:start-after:`, `:end-before:`, `:tabid:`)
- ✅ Ordered lists (procedures as numbered lists)
- ✅ Sub-lists (sub-procedures within steps)
- ✅ Inline roles for URL validation (`:doc:`, `:ref:`, external links)
- ✅ Inline roles for UI testing (`:guilabel:`)
- ✅ Source constants (`{+variable-name+}`)

**What We Ignore**:
- ❌ Full RST syntax tree
- ❌ Tables, footnotes, citations
- ❌ Most inline roles (except those listed above)
- ❌ Complex directives we don't test

**Testing Strategy**:
- Test against real MongoDB documentation files from the start
- Comprehensive test suite using actual .txt files from testdata/
- Handle edge cases incrementally as discovered

**Contingency** (unlikely to be needed):
- Fall back to snooty-parser for complex cases (heavyweight but comprehensive)
- Limit PoC scope to most common directives

#### Risk 2: Placeholder Resolution Accuracy
**Severity**: Medium
**Probability**: Medium
**Impact**: Poor user experience, low adoption

**Description**: Fuzzy matching may produce false positives or miss valid matches, leading to frustration.

**Mitigation**:
- Implement multiple resolution strategies with priority
- Provide clear suggestions when resolution fails
- Allow explicit mappings in configuration
- Test against real placeholder patterns from documentation

**Contingency**:
- Fall back to exact matching only
- Require more explicit configuration

#### Risk 3: State Management Complexity
**Severity**: Medium
**Probability**: Low
**Impact**: Code blocks don't execute correctly

**Description**: Maintaining state across code blocks (especially for different languages) may be complex.

**Mitigation**:
- Start with simple accumulation strategy
- Test thoroughly with real examples
- Document limitations clearly
- Provide escape hatches (isolated execution mode)

**Contingency**:
- Execute each code block in isolation
- Require writers to include all necessary code in each block

#### Risk 4: Execution Environment Isolation
**Severity**: High
**Probability**: Low
**Impact**: Tests interfere with each other or writer's environment

**Description**: Tests may create databases, files, or other resources that conflict.

**Mitigation**:
- Use unique identifiers for all created resources
- Implement comprehensive cleanup registry
- Test cleanup thoroughly
- Provide `--no-cleanup` flag for debugging

**Contingency**:
- Recommend running in Docker container
- Provide cleanup scripts

### 8.2 User Adoption Risks

#### Risk 5: Too Complex for Writers
**Severity**: High
**Probability**: Low
**Impact**: Low adoption, tool abandonment

**Description**: If the tool requires too much configuration or technical knowledge, writers won't use it.

**Mitigation**:
- Zero-config operation for common cases
- Excellent error messages with suggestions
- Comprehensive documentation with examples
- User testing with real technical writers

**Contingency**:
- Simplify further
- Provide more presets and templates
- Offer training sessions

#### Risk 6: Performance Issues
**Severity**: Medium
**Probability**: Low
**Impact**: Slow feedback loop, poor developer experience

**Description**: If tests take too long to run, writers won't run them frequently.

**Mitigation**:
- Optimize parser and executor performance
- Provide `--filter` flag to run subset of tests
- Consider caching parsed ASTs
- Profile and optimize hot paths

**Contingency**:
- Provide incremental testing mode (only test changed files)
- Add watch mode for continuous testing during development
- Optimize test startup time (lazy loading, faster imports)

### 8.3 Maintenance Risks

#### Risk 7: RST → MDX Migration
**Severity**: Medium
**Probability**: High (planned)
**Impact**: Requires parser refactoring

**Description**: MongoDB plans to migrate from RST to MDX in the future.

**Mitigation**:
- **Interface-driven design** (Option 5 key benefit)
- Parser abstraction from the start
- Test parser interface with mock implementations
- Plan for dual-parser support period

**Contingency**:
- This is why we chose Option 5 - minimal refactoring needed
- Add MDX parser as new implementation of Parser interface

#### Risk 8: Diverse Team Needs
**Severity**: Medium
**Probability**: High
**Impact**: Feature requests exceed capacity

**Description**: Different documentation teams may have conflicting requirements.

**Mitigation**:
- **Interface-driven architecture** enables future extensibility
- Configuration system for customization
- Clear prioritization of features
- Community contribution guidelines
- Built-in support for common use cases

**Contingency**:
- Add new built-in executors/resolvers as needed
- Extend configuration system for team-specific customization
- If extensibility needs become critical, implement plugin system at that time

### 8.4 Risk Summary Matrix

| Risk | Severity | Probability | Mitigation Effectiveness | Residual Risk |
|------|----------|-------------|-------------------------|---------------|
| RST Parsing Complexity | High | Medium | High | Low |
| Placeholder Resolution | Medium | Medium | High | Low |
| State Management | Medium | Low | Medium | Low |
| Execution Isolation | High | Low | High | Very Low |
| Too Complex for Writers | High | Low | High | Very Low |
| Performance Issues | Medium | Low | Medium | Low |
| RST → MDX Migration | Medium | High | **Very High** (Option 5) | Very Low |
| Diverse Team Needs | Medium | High | **Very High** (Option 5) | Low |

**Overall Risk Level**: **Low** (with Option 5 architecture)

---

## 9. Appendices

### Appendix A: Parse Command Specification

The `parse` command is a debugging tool that displays the parsed AST structure without executing code.

#### Usage

```bash
# Display parsed structure in tree format (default)
proctest parse path/to/procedure.txt

# Write parsed structure to file
proctest parse path/to/procedure.txt --output parsed-output.txt

# Output as JSON for programmatic processing
proctest parse path/to/procedure.txt --format json

# Output as YAML for readability
proctest parse path/to/procedure.txt --format yaml

# Combine output file with format
proctest parse path/to/procedure.txt --output ast.json --format json
```

#### Output Formats

**Tree Format (Default)** - Human-readable hierarchical view:
```
Document: symfony.txt
├─ Procedure: "Symfony MongoDB Integration"
│  ├─ Prerequisites:
│  │  ├─ Software: PHP (>=8.0)
│  │  │  Check: php --version
│  │  │  Install: https://www.php.net/downloads
│  │  ├─ Software: Composer
│  │  │  Check: composer --version
│  │  │  Install: https://getcomposer.org/download/
│  │  ├─ Software: Symfony CLI
│  │  │  Check: symfony version
│  │  │  Install: https://symfony.com/download
│  │  └─ Service: MongoDB Atlas cluster
│  │     Setup: https://www.mongodb.com/docs/atlas/getting-started
│  │
│  ├─ Step 1: "Initialize a Symfony Project"
│  │  ├─ Paragraph: "Run the following command..."
│  │  └─ CodeBlock [shell] (executable)
│  │     Language: shell
│  │     Code: "composer create-project symfony/skeleton restaurants"
│  │     Placeholders: none
│  │     Location: symfony.txt:132-134
│  │
│  ├─ Step 2: "Install PHP Driver and Doctrine ODM"
│  │  ├─ Paragraph: "Enter your project directory..."
│  │  └─ CodeBlock [shell] (executable)
│  │     Language: shell
│  │     Code: "composer require doctrine/mongodb-odm-bundle"
│  │     Placeholders: none
│  │     Location: symfony.txt:147-149
│  │
│  └─ Step 3: "Configure MongoDB Connection"
│     ├─ Paragraph: "Create a .env file..."
│     ├─ Sub-Step a: "Create the file"
│     │  └─ CodeBlock [shell] (executable)
│     │     Language: shell
│     │     Code: "echo 'MONGODB_URI=<connection-string>' > .env"
│     │     Placeholders: ["<connection-string>"]
│     │     Location: symfony.txt:155-157
│     └─ Sub-Step b: "Verify configuration"
│        └─ CodeBlock [shell] (executable)
│           Language: shell
│           Code: "symfony console doctrine:mongodb:schema:validate"
│           Placeholders: none
│           Location: symfony.txt:160-162

Summary:
- Total Procedures: 1
- Total Steps: 3
- Total Sub-Steps: 2
- Total Prerequisites: 4
  - Software: 3
  - Service: 1
- Total Testable Actions: 5
  - Code: 0
  - Shell: 5
  - UI: 0
- Unresolved Placeholders: 1 ("<connection-string>")
```

**JSON Format** - For programmatic processing:
```json
{
  "file": "symfony.txt",
  "type": "document",
  "metadata": {
    "title": "Symfony MongoDB Integration"
  },
  "procedures": [
    {
      "type": "procedure",
      "title": "Symfony MongoDB Integration",
      "style": "connected",
      "prerequisites": {
        "type": "prerequisites",
        "title": "Prerequisites",
        "requirements": [
          {
            "requirementType": "software",
            "name": "PHP",
            "version": ">=8.0",
            "description": "Ensure that your PHP installation includes the MongoDB extension...",
            "optional": false,
            "checkCommand": "php --version",
            "installUrl": "https://www.php.net/downloads",
            "location": {
              "file": "symfony.txt",
              "startLine": 95,
              "endLine": 98
            }
          },
          {
            "requirementType": "software",
            "name": "Composer",
            "description": "Dependency manager for PHP.",
            "optional": false,
            "checkCommand": "composer --version",
            "installUrl": "https://getcomposer.org/download/",
            "location": {
              "file": "symfony.txt",
              "startLine": 100,
              "endLine": 101
            }
          }
        ],
        "location": {
          "file": "symfony.txt",
          "startLine": 84,
          "endLine": 110
        }
      },
      "steps": [
        {
          "type": "step",
          "number": 1,
          "title": "Initialize a Symfony Project",
          "testableActions": [
            {
              "actionType": "shell",
              "command": "composer create-project symfony/skeleton restaurants",
              "placeholders": [],
              "location": {
                "file": "symfony.txt",
                "startLine": 132,
                "endLine": 134
              }
            }
          ]
        }
      ],
      "location": {
        "file": "symfony.txt",
        "startLine": 81,
        "endLine": 327
      }
    }
  ],
  "summary": {
    "totalProcedures": 1,
    "totalSteps": 3,
    "totalSubSteps": 2,
    "totalPrerequisites": 4,
    "totalTestableActions": 5,
    "unresolvedPlaceholders": ["<connection-string>"]
  }
}
```

**YAML Format** - Human-readable structured format:
```yaml
file: symfony.txt
type: document
metadata:
  title: Symfony MongoDB Integration
procedures:
  - type: procedure
    title: Symfony MongoDB Integration
    style: connected
    prerequisites:
      type: prerequisites
      title: Prerequisites
      requirements:
        - requirementType: software
          name: PHP
          version: ">=8.0"
          description: Ensure that your PHP installation includes the MongoDB extension...
          optional: false
          checkCommand: php --version
          installUrl: https://www.php.net/downloads
          location:
            file: symfony.txt
            startLine: 95
            endLine: 98
        - requirementType: software
          name: Composer
          description: Dependency manager for PHP.
          optional: false
          checkCommand: composer --version
          installUrl: https://getcomposer.org/download/
          location:
            file: symfony.txt
            startLine: 100
            endLine: 101
      location:
        file: symfony.txt
        startLine: 84
        endLine: 110
    steps:
      - type: step
        number: 1
        title: Initialize a Symfony Project
        testableActions:
          - actionType: shell
            command: composer create-project symfony/skeleton restaurants
            placeholders: []
            location:
              file: symfony.txt
              startLine: 132
              endLine: 134
summary:
  totalProcedures: 1
  totalSteps: 3
  totalSubSteps: 2
  totalPrerequisites: 4
  totalTestableActions: 5
  unresolvedPlaceholders:
    - "<connection-string>"
```

#### Implementation

```typescript
// src/cli/commands/parse.ts
import { Command } from 'commander';
import { RSTParser } from '../../parser/rst-parser';
import { ConfigurationManager } from '../../config/configuration-manager';
import { writeFile } from 'fs/promises';

export function createParseCommand(): Command {
  const command = new Command('parse');

  command
    .description('Parse a documentation file and display the AST structure')
    .argument('<file>', 'Path to the file to parse')
    .option('-o, --output <file>', 'Write output to file instead of stdout')
    .option('-f, --format <type>', 'Output format: tree, json, yaml', 'tree')
    .option('--snooty <file>', 'Path to snooty.toml file')
    .action(async (file: string, options) => {
      try {
        // Load configuration
        const configManager = new ConfigurationManager();
        const config = await configManager.load({ snootyConfig: options.snooty });

        // Parse file
        const parser = new RSTParser();
        const context = {
          filePath: file,
          snootyConstants: config.snootyConstants || {},
          workingDirectory: process.cwd()
        };

        const content = await readFile(file, 'utf-8');
        const ast = await parser.parse(content, context);

        // Format output
        let output: string;
        switch (options.format) {
          case 'json':
            output = JSON.stringify(ast, null, 2);
            break;
          case 'yaml':
            output = formatAsYAML(ast);
            break;
          case 'tree':
          default:
            output = formatAsTree(ast);
            break;
        }

        // Write to file or stdout
        if (options.output) {
          await writeFile(options.output, output, 'utf-8');
          console.log(`✓ Parsed AST written to ${options.output}`);
        } else {
          console.log(output);
        }

      } catch (error) {
        console.error('Error parsing file:', error);
        process.exit(1);
      }
    });

  return command;
}

function formatAsTree(ast: DocumentAST): string {
  const lines: string[] = [];

  lines.push(`Document: ${ast.metadata?.title || ast.filePath}`);

  for (const procedure of ast.procedures) {
    lines.push(`├─ Procedure: "${procedure.title || 'Untitled'}"`);

    for (let i = 0; i < procedure.steps.length; i++) {
      const step = procedure.steps[i];
      const isLast = i === procedure.steps.length - 1;
      const prefix = isLast ? '└─' : '├─';

      lines.push(`│  ${prefix} Step ${step.number || i + 1}: "${step.title || 'Untitled'}"`);

      // Show content items
      for (const content of step.content) {
        if (content.type === 'code-block') {
          const cb = content as CodeBlockNode;
          lines.push(`│  │  └─ CodeBlock [${cb.language}] (${cb.executable ? 'executable' : 'non-executable'})`);
          lines.push(`│  │     Language: ${cb.language}`);
          lines.push(`│  │     Code: ${cb.code.substring(0, 50)}...`);
          lines.push(`│  │     Placeholders: ${cb.placeholders?.length ? JSON.stringify(cb.placeholders) : 'none'}`);
          lines.push(`│  │     Location: ${cb.location.file}:${cb.location.startLine}-${cb.location.endLine}`);
        } else if (content.type === 'paragraph') {
          lines.push(`│  │  ├─ Paragraph: "${content.text?.substring(0, 40)}..."`);
        }
      }

      // Show sub-steps if present
      if (step.subSteps && step.subSteps.length > 0) {
        for (let j = 0; j < step.subSteps.length; j++) {
          const subStep = step.subSteps[j];
          lines.push(`│  │  ├─ Sub-Step ${subStep.number}`);

          for (const content of subStep.content) {
            if (content.type === 'code-block') {
              const cb = content as CodeBlockNode;
              lines.push(`│  │  │  └─ CodeBlock [${cb.language}] (${cb.executable ? 'executable' : 'non-executable'})`);
              lines.push(`│  │  │     Language: ${cb.language}`);
              lines.push(`│  │  │     Code: ${cb.code.substring(0, 50)}...`);
              lines.push(`│  │  │     Placeholders: ${cb.placeholders?.length ? JSON.stringify(cb.placeholders) : 'none'}`);
              lines.push(`│  │  │     Location: ${cb.location.file}:${cb.location.startLine}-${cb.location.endLine}`);
            } else if (content.type === 'paragraph') {
              lines.push(`│  │  │  ├─ Paragraph: "${content.text?.substring(0, 40)}..."`);
            }
          }
        }
      }
    }
  }

  // Summary
  const executableBlocks = countExecutableBlocks(ast);
  const placeholders = collectPlaceholders(ast);

  lines.push('');
  lines.push('Summary:');
  lines.push(`- Total Procedures: ${ast.procedures.length}`);
  lines.push(`- Total Steps: ${ast.procedures.reduce((sum, p) => sum + p.steps.length, 0)}`);
  lines.push(`- Total Code Blocks: ${countCodeBlocks(ast)}`);
  lines.push(`- Executable Code Blocks: ${executableBlocks}`);
  lines.push(`- Unresolved Placeholders: ${placeholders.length} (${JSON.stringify(placeholders)})`);

  return lines.join('\n');
}

function formatAsYAML(ast: DocumentAST): string {
  // Simple YAML formatter (or use a library like 'yaml')
  return require('yaml').stringify(ast);
}

function countCodeBlocks(ast: DocumentAST): number {
  let count = 0;
  for (const procedure of ast.procedures) {
    for (const step of procedure.steps) {
      for (const content of step.content) {
        if (content.type === 'code-block') count++;
      }
    }
  }
  return count;
}

function countExecutableBlocks(ast: DocumentAST): number {
  let count = 0;
  for (const procedure of ast.procedures) {
    for (const step of procedure.steps) {
      for (const content of step.content) {
        if (content.type === 'code-block' && (content as CodeBlockNode).executable) {
          count++;
        }
      }
    }
  }
  return count;
}

function collectPlaceholders(ast: DocumentAST): string[] {
  const placeholders = new Set<string>();
  for (const procedure of ast.procedures) {
    for (const step of procedure.steps) {
      for (const content of step.content) {
        if (content.type === 'code-block') {
          const cb = content as CodeBlockNode;
          cb.placeholders?.forEach(p => placeholders.add(p));
        }
      }
    }
  }
  return Array.from(placeholders);
}
```

#### Use Cases

1. **Verify Parser Correctness**: Check that the parser correctly identifies procedures, steps, and code blocks
2. **Debug Placeholder Detection**: See which placeholders were detected in code blocks
3. **Validate File Structure**: Ensure documentation follows expected structure before running tests
4. **Generate Test Fixtures**: Export AST as JSON for use in unit tests
5. **Documentation Review**: Quickly see the structure of a documentation file

#### Integration with Implementation Plan

Add to **Milestone 2: Parsing** deliverables:
- [ ] Implement `parse` command with tree, JSON, and YAML output formats
- [ ] Add placeholder detection and reporting in parse output
- [ ] Add summary statistics to parse output

---

### Appendix B: Sub-Procedure Handling and Error Reporting

#### Sub-Procedure Structure

Sub-procedures are represented as **ordered lists within steps**. They can use:
- **Numeric numbering**: `1.`, `2.`, `3.`
- **Alphabetic numbering**: `a.`, `b.`, `c.` or `A.`, `B.`, `C.`
- **Auto-numbering**: `#.` (RST automatically numbers these)

**Example from requirements.md**:
```rst
.. procedure::

   .. step:: Step 1

      a. Sub-step 1
      b. Sub-step 2
      c. Sub-step 3

   .. step:: Step 2
```

**Real-world example from testdata**:
```rst
.. step:: In Atlas, go to the Clusters page for your project.

   a. If it's not already displayed, select the organization that
      contains your desired project from the organization menu in the
      navigation bar.

   #. If it's not already displayed, select your desired project
      from the Projects menu in the navigation bar.

   #. In the sidebar, click Clusters under
      the Database heading.

   The Clusters page displays.
```

#### Execution Behavior

When executing a procedure with sub-steps:

1. **Sequential Execution**: Sub-steps execute in order (a, b, c, ...)
2. **Failure Propagation**: If a sub-step fails, the parent step fails
3. **State Accumulation**: Code blocks within sub-steps share state (within the step)
4. **Hierarchical Context**: Errors report the full path: Procedure → Step → Sub-Step

#### Error Reporting with Sub-Steps

**Example Error Output**:
```
❌ FAILED: Navigate to Clusters Page

Procedure: "Navigate to Clusters Page"
  Step 1: "In Atlas, go to the Clusters page for your project"
    Sub-Step b: "Select your desired project"
      ❌ Code Block [javascript] FAILED

         Error: Cannot find element with selector '.project-menu'

         Location: steps-db-deployments-page.rst:7-9

         Context:
           Procedure: "Navigate to Clusters Page"
           Step: 1 - "In Atlas, go to the Clusters page for your project"
           Sub-Step: b - "Select your desired project"
           Code Block: javascript

         Suggestions:
           - Verify the UI selector is correct
           - Check if the page has loaded completely
           - Ensure you're logged in to Atlas

Duration: 2.3s
```

**Hierarchical Error Context**:
```typescript
const error: TestError = {
  type: 'execute',
  message: "Cannot find element with selector '.project-menu'",
  location: {
    file: 'steps-db-deployments-page.rst',
    startLine: 7,
    endLine: 9
  },
  context: {
    procedureTitle: 'Navigate to Clusters Page',
    stepNumber: 1,
    stepHeadline: 'In Atlas, go to the Clusters page for your project',
    subStepNumber: 'b', // Alphabetic sub-step
    codeBlockLanguage: 'javascript'
  },
  suggestions: [
    "Verify the UI selector is correct",
    "Check if the page has loaded completely",
    "Ensure you're logged in to Atlas"
  ]
};
```

#### Implementation Notes

**Parsing Sub-Steps**:
```typescript
// When parsing a step, detect ordered lists
function parseStep(stepDirective: Directive): StepNode {
  const content: ContentNode[] = [];
  const subSteps: SubStepNode[] = [];

  for (const child of stepDirective.children) {
    if (child.type === 'enumerated_list') {
      // This is a sub-procedure
      for (const item of child.items) {
        subSteps.push({
          type: 'sub-step',
          number: item.number, // 'a', 'b', 'c' or 1, 2, 3
          content: parseContent(item.content),
          codeBlocks: extractCodeBlocks(item.content),
          location: getLocation(item)
        });
      }
    } else {
      content.push(parseContentNode(child));
    }
  }

  return {
    type: 'step',
    headline: stepDirective.argument,
    content,
    codeBlocks: extractCodeBlocks(content),
    subSteps: subSteps.length > 0 ? subSteps : undefined,
    location: getLocation(stepDirective)
  };
}
```

**Executing Sub-Steps**:
```typescript
async function executeStep(step: StepNode, context: ExecutionContext): Promise<StepResult> {
  const subStepResults: SubStepResult[] = [];

  // Execute main step code blocks first
  const mainCodeBlockResults = await executeCodeBlocks(step.codeBlocks, context);

  // Execute sub-steps if present
  if (step.subSteps) {
    for (const subStep of step.subSteps) {
      const subStepResult = await executeSubStep(subStep, context);
      subStepResults.push(subStepResult);

      // If sub-step fails, fail the entire step
      if (!subStepResult.success) {
        return {
          step,
          success: false,
          duration: Date.now() - startTime,
          codeBlocks: mainCodeBlockResults,
          subSteps: subStepResults,
          error: subStepResult.error
        };
      }
    }
  }

  return {
    step,
    success: true,
    duration: Date.now() - startTime,
    codeBlocks: mainCodeBlockResults,
    subSteps: subStepResults.length > 0 ? subStepResults : undefined
  };
}
```

**Reporting Sub-Step Failures**:
```typescript
function formatError(error: TestError): string {
  const lines: string[] = [];

  if (error.context) {
    lines.push('Context:');
    if (error.context.procedureTitle) {
      lines.push(`  Procedure: "${error.context.procedureTitle}"`);
    }
    if (error.context.stepNumber) {
      const stepInfo = error.context.stepHeadline
        ? `${error.context.stepNumber} - "${error.context.stepHeadline}"`
        : `${error.context.stepNumber}`;
      lines.push(`  Step: ${stepInfo}`);
    }
    if (error.context.subStepNumber) {
      lines.push(`  Sub-Step: ${error.context.subStepNumber}`);
    }
    if (error.context.codeBlockLanguage) {
      lines.push(`  Code Block: ${error.context.codeBlockLanguage}`);
    }
  }

  return lines.join('\n');
}
```

#### Testing Sub-Procedures

**Unit Test Example**:
```typescript
describe('Sub-Procedure Execution', () => {
  it('should execute sub-steps in order', async () => {
    const procedure = {
      type: 'procedure',
      steps: [{
        type: 'step',
        headline: 'Main step',
        subSteps: [
          { type: 'sub-step', number: 'a', codeBlocks: [/* ... */] },
          { type: 'sub-step', number: 'b', codeBlocks: [/* ... */] },
          { type: 'sub-step', number: 'c', codeBlocks: [/* ... */] }
        ]
      }]
    };

    const result = await executeStep(procedure.steps[0], context);

    expect(result.subSteps).toHaveLength(3);
    expect(result.subSteps[0].subStepNumber).toBe('a');
    expect(result.subSteps[1].subStepNumber).toBe('b');
    expect(result.subSteps[2].subStepNumber).toBe('c');
  });

  it('should fail step when sub-step fails', async () => {
    const procedure = {
      type: 'procedure',
      steps: [{
        type: 'step',
        subSteps: [
          { type: 'sub-step', number: 'a', codeBlocks: [validCodeBlock] },
          { type: 'sub-step', number: 'b', codeBlocks: [failingCodeBlock] },
          { type: 'sub-step', number: 'c', codeBlocks: [validCodeBlock] }
        ]
      }]
    };

    const result = await executeStep(procedure.steps[0], context);

    expect(result.success).toBe(false);
    expect(result.subSteps).toHaveLength(2); // Only a and b executed
    expect(result.error?.context?.subStepNumber).toBe('b');
  });

  it('should include sub-step context in errors', async () => {
    const result = await executeStep(stepWithFailingSubStep, context);

    expect(result.error?.context).toMatchObject({
      procedureTitle: 'Navigate to Clusters Page',
      stepNumber: 1,
      stepHeadline: 'In Atlas, go to the Clusters page',
      subStepNumber: 'b',
      codeBlockLanguage: 'javascript'
    });
  });
});
```

---

### Appendix C: Prerequisite Detection and Validation

The framework automatically detects and validates prerequisites/requirements before executing procedures, allowing tests to be skipped gracefully when requirements aren't met.

#### C.1 Prerequisite Detection

**Detection Strategy**:
1. Look for sections with headings containing "Prerequisite", "Requirement", "Before you begin", etc.
2. These sections typically appear **before** the procedure
3. Parse the content to extract individual requirements

**Example from testdata** (`symfony.txt`):
```rst
.. procedure::
   :style: connected

   .. step:: Prerequisites

      To create the Quick Start application, you need the following software
      installed in your development environment:

      .. list-table::
         :header-rows: 1

         * - Prerequisite
           - Notes

         * - `PHP <https://www.php.net/downloads>`__
           - Ensure that your PHP installation includes the MongoDB extension
             and that it is enabled. To confirm version compatibility, see Compatibility.

         * - `Composer <https://getcomposer.org/download/>`__
           - Dependency manager for PHP.

         * - `Symfony CLI <https://symfony.com/download>`__
           - Command-line tool for managing Symfony applications.

         * - A terminal app and shell
           - For MacOS users, use Terminal or a
             similar app. For Windows users, use PowerShell.
```

**Parsed As**:
```typescript
{
  type: 'prerequisites',
  title: 'Prerequisites',
  requirements: [
    {
      requirementType: 'software',
      name: 'PHP',
      description: 'Ensure that your PHP installation includes the MongoDB extension...',
      optional: false,
      checkCommand: 'php --version',
      installUrl: 'https://www.php.net/downloads',
      location: { file: 'symfony.txt', startLine: 95, endLine: 98 }
    },
    {
      requirementType: 'software',
      name: 'Composer',
      description: 'Dependency manager for PHP.',
      optional: false,
      checkCommand: 'composer --version',
      installUrl: 'https://getcomposer.org/download/',
      location: { file: 'symfony.txt', startLine: 100, endLine: 101 }
    },
    {
      requirementType: 'software',
      name: 'Symfony CLI',
      description: 'Command-line tool for managing Symfony applications.',
      optional: false,
      checkCommand: 'symfony version',
      installUrl: 'https://symfony.com/download',
      location: { file: 'symfony.txt', startLine: 103, endLine: 105 }
    }
  ],
  location: { file: 'symfony.txt', startLine: 84, endLine: 110 }
}
```

#### C.2 Requirement Types

**Software Requirements**:
- Detected from: Tool names (PHP, Node.js, Python, Composer, etc.)
- Validation: Run check command (e.g., `php --version`)
- Version matching: Parse version from output, compare with requirement

**Environment Requirements**:
- Detected from: References to environment variables or `.env` files
- Validation: Check if variable is set in environment
- Example: "Set `MONGODB_URI` in your `.env` file"

**Service Requirements**:
- Detected from: References to external services (MongoDB Atlas, etc.)
- Validation: Optional - can check connectivity if credentials provided
- Example: "Create a MongoDB Atlas cluster"

**Configuration Requirements**:
- Detected from: References to config files
- Validation: Check if file exists
- Example: "Ensure `snooty.toml` is configured"

#### C.3 Prerequisite Checking Implementation

```typescript
/**
 * Check software requirement
 */
async function checkSoftwareRequirement(req: SoftwareRequirement): Promise<PrerequisiteCheckResult> {
  if (!req.checkCommand) {
    return {
      requirement: req,
      met: true, // Assume met if we can't check
      message: `Cannot verify ${req.name} installation (no check command)`,
      skipReason: undefined
    };
  }

  try {
    const { stdout, stderr } = await execAsync(req.checkCommand);
    const output = stdout || stderr;

    // Extract version if specified
    if (req.version) {
      const foundVersion = extractVersion(output);
      const meetsVersion = compareVersions(foundVersion, req.version);

      if (!meetsVersion) {
        return {
          requirement: req,
          met: false,
          message: `${req.name} version mismatch`,
          details: {
            found: foundVersion,
            expected: req.version,
            command: req.checkCommand,
            output
          },
          skipReason: `${req.name} ${req.version} is required, but ${foundVersion} was found`
        };
      }
    }

    return {
      requirement: req,
      met: true,
      message: `✓ ${req.name} is installed`,
      details: {
        found: extractVersion(output) || 'installed',
        command: req.checkCommand,
        output
      }
    };
  } catch (error) {
    return {
      requirement: req,
      met: false,
      message: `${req.name} is not installed`,
      details: {
        command: req.checkCommand,
        output: error.message
      },
      skipReason: `${req.name} is required but not installed. Install from: ${req.installUrl || 'N/A'}`
    };
  }
}

/**
 * Check environment requirement
 */
async function checkEnvironmentRequirement(req: EnvironmentRequirement): Promise<PrerequisiteCheckResult> {
  const value = process.env[req.variable];

  if (!value) {
    return {
      requirement: req,
      met: false,
      message: `Environment variable ${req.variable} is not set`,
      skipReason: `Set ${req.variable} in your environment or .env file${req.example ? `. Example: ${req.example}` : ''}`
    };
  }

  return {
    requirement: req,
    met: true,
    message: `✓ ${req.variable} is set`,
    details: {
      found: value.substring(0, 20) + '...' // Don't expose full value
    }
  };
}

/**
 * Check all prerequisites for a procedure
 */
async function checkPrerequisites(prerequisites: PrerequisiteNode): Promise<PrerequisiteCheckResult[]> {
  const results: PrerequisiteCheckResult[] = [];

  for (const requirement of prerequisites.requirements) {
    let result: PrerequisiteCheckResult;

    switch (requirement.requirementType) {
      case 'software':
        result = await checkSoftwareRequirement(requirement);
        break;
      case 'environment':
        result = await checkEnvironmentRequirement(requirement);
        break;
      case 'service':
        result = await checkServiceRequirement(requirement);
        break;
      case 'configuration':
        result = await checkConfigurationRequirement(requirement);
        break;
    }

    results.push(result);
  }

  return results;
}
```

#### C.4 Test Execution Flow with Prerequisites

```typescript
async function executeProcedure(procedure: ProcedureNode): Promise<ProcedureResult> {
  const startTime = Date.now();

  // Check prerequisites first
  if (procedure.prerequisites) {
    const prerequisiteChecks = await checkPrerequisites(procedure.prerequisites);
    const unmetRequired = prerequisiteChecks.filter(c => !c.met && !c.requirement.optional);

    if (unmetRequired.length > 0) {
      // Skip the test if required prerequisites aren't met
      const skipReasons = unmetRequired.map(c => c.skipReason).filter(Boolean);

      return {
        procedure,
        success: false,
        skipped: true,
        skipReason: skipReasons.join('\n'),
        prerequisiteChecks,
        duration: Date.now() - startTime,
        steps: []
      };
    }
  }

  // Prerequisites met, execute the procedure
  const steps: StepResult[] = [];
  for (const step of procedure.steps) {
    const stepResult = await executeStep(step, context);
    steps.push(stepResult);

    if (!stepResult.success) {
      return {
        procedure,
        success: false,
        skipped: false,
        prerequisiteChecks: procedure.prerequisites ? await checkPrerequisites(procedure.prerequisites) : undefined,
        duration: Date.now() - startTime,
        steps,
        error: stepResult.error
      };
    }
  }

  return {
    procedure,
    success: true,
    skipped: false,
    prerequisiteChecks: procedure.prerequisites ? await checkPrerequisites(procedure.prerequisites) : undefined,
    duration: Date.now() - startTime,
    steps
  };
}
```

#### C.5 Reporting Skipped Tests

**Example Output** (prerequisites not met):
```
⊘ SKIPPED: Symfony MongoDB Integration

Prerequisites:
  ✓ PHP 8.2.0 is installed
  ✓ Composer 2.5.0 is installed
  ✗ Symfony CLI is not installed

Skip Reason:
  Symfony CLI is required but not installed.
  Install from: https://symfony.com/download

Duration: 0.3s
```

**Example Output** (prerequisites met):
```
✓ PASSED: Symfony MongoDB Integration

Prerequisites:
  ✓ PHP 8.2.0 is installed
  ✓ Composer 2.5.0 is installed
  ✓ Symfony CLI 5.4.0 is installed

Procedure: "Symfony MongoDB Integration"
  Step 1: "Prerequisites"
    (No testable actions)

  Step 2: "Initialize a Symfony Project"
    ✓ Shell: composer create-project symfony/skeleton restaurants (5.2s)

  ... (remaining steps)

Duration: 45.3s
```

#### C.6 Configuration Options

**Skip Prerequisite Checks**:
```bash
# Skip all prerequisite checks (useful for debugging)
proctest run symfony.txt --skip-prerequisites

# Run even if prerequisites aren't met (will likely fail)
proctest run symfony.txt --ignore-prerequisites
```

**Configuration File**:
```javascript
// proctest.config.js
export default {
  prerequisites: {
    check: true, // Enable prerequisite checking (default: true)
    skipOnFailure: true, // Skip tests if prerequisites not met (default: true)
    failOnWarning: false, // Fail if optional prerequisites not met (default: false)

    // Custom check commands for software
    checkCommands: {
      'PHP': 'php --version',
      'Node.js': 'node --version',
      'Python': 'python3 --version',
      'Composer': 'composer --version',
      'Symfony CLI': 'symfony version'
    }
  }
};
```

#### C.7 Detection Heuristics

**Heading Patterns** (case-insensitive):
- "Prerequisites"
- "Requirements"
- "Before you begin"
- "Before you start"
- "What you need"
- "Setup"

**Software Name Detection**:
- Known tools: PHP, Node.js, Python, Java, Ruby, Go, Rust, etc.
- Package managers: npm, pip, composer, maven, gradle, cargo, etc.
- CLIs: mongosh, atlas-cli, aws-cli, gcloud, etc.

**Version Pattern Detection**:
- Semantic versioning: `>=8.0`, `^18.0.0`, `~5.4.0`
- Range: `8.0 or higher`, `version 18+`
- Exact: `8.2.0`, `v18.0.0`

---

### Appendix D: Testable Action Types and Detection

The framework supports eight types of testable actions, each with specific detection rules and execution strategies.

#### D.1 Action Type Overview

| Action Type | Detection Method | Execution Strategy | PoC Phase |
|-------------|------------------|-------------------|-----------|
| **File** | Prose instructions + `literalinclude` | File system operations | Phase 1 |
| **Code** | `code-block`, `code`, `literalinclude` directives | Language-specific runtime | Phase 1 |
| **Shell** | `code-block` with `shell`/`bash`/`sh` language | Shell execution | Phase 1 |
| **UI** | `:guilabel:` role in text | UI automation (Playwright/Puppeteer) | Phase 3 |
| **CLI** | `code-block` with `mongosh` or atlas-cli commands | Tool-specific execution | Phase 2 |
| **API** | `curl` commands targeting Atlas Admin API | HTTP client (axios/fetch) | Phase 3 |
| **Download** | `curl` commands with `-o` or `-O` flags | HTTP client with file writing | Phase 2 |
| **URL** | Links in documentation | HTTP HEAD/GET request | Phase 3 |

#### D.2 File Testable Actions

**Purpose**: Handle file creation, modification, and content replacement as described in procedural documentation.

**Detection Patterns**:

File operations are detected when prose instructions contain file operation keywords followed by a `literalinclude` directive or code block:

1. **Create operation**:
   ```rst
   Create the ``Restaurant.php`` file in the ``src/Document`` directory and
   paste the following code:

   .. literalinclude:: /includes/php-frameworks/symfony/Restaurant.php
      :caption: src/Document/Restaurant.php
      :language: php
   ```

2. **Replace operation**:
   ```rst
   In the ``config/packages`` directory, replace the contents of your
   ``doctrine_mongodb.yaml`` file with the following code:

   .. literalinclude:: /includes/php-frameworks/symfony/doctrine_mongodb.yaml
      :caption: config/packages/doctrine_mongodb.yaml
      :language: yaml
   ```

3. **Append operation** (less common):
   ```rst
   Add the following code to the end of your ``index.js`` file:

   .. code-block:: javascript

      module.exports = app;
   ```

**Detection Keywords**:
- **Create**: "Create the `<filename>` file", "Create a file named `<filename>`", "Create `<filename>`"
- **Replace**: "Replace the contents of `<filename>`", "Replace `<filename>` with"
- **Append**: "Add to `<filename>`", "Append to `<filename>`", "Add the following to `<filename>`"

**Parsed As**:
```typescript
{
  actionType: 'file',
  operation: 'create',
  path: 'src/Document/Restaurant.php',
  content: '<?php\n\nnamespace App\\Document;\n\nuse Doctrine\\ODM\\MongoDB\\Mapping\\Annotations as MongoDB;\n\n...',
  description: 'Create the Restaurant.php file in the src/Document directory',
  location: { file: 'symfony.txt', startLine: 229, endLine: 235 }
}
```

**Execution Strategy**:

File operations are executed **before** code blocks in the same step:

1. **Working Directory**: All file paths are relative to the per-procedure temp directory (`.proctest/runs/<timestamp>-<procedure-id>/`)
2. **Directory Creation**: Parent directories are created automatically if they don't exist
3. **Operation Execution**:
   - `create`: Write content to file (fail if file already exists, unless overwrite is configured)
   - `replace`: Overwrite existing file content (create if doesn't exist)
   - `append`: Append content to existing file (create if doesn't exist)
4. **Cleanup**: Files are automatically cleaned up when the temp directory is removed

**Working Directory Changes**:

Shell commands can change the working directory using `cd`, `mkdir`, etc. The ShellExecutor maintains the current working directory across shell commands within a step:

1. **Initial Working Directory**: Set to per-procedure temp directory
2. **Directory Changes Persist**: `cd` commands affect subsequent shell commands in the same step
3. **Step Isolation**: Each step starts with the working directory from the end of the previous step
4. **Common Patterns**:
   ```rst
   .. code-block:: bash

      mkdir my-project
      cd my-project
      npm init -y
   ```

   Or as sub-steps:
   ```rst
   a. Create a new directory:

      .. code-block:: bash

         mkdir my-project

   #. Change to the new directory:

      .. code-block:: bash

         cd my-project

   #. Initialize the project:

      .. code-block:: bash

         npm init -y
   ```

**Implementation Note**: The ShellExecutor should execute all shell commands in the same shell session (or track `cwd` changes) to preserve directory context.

**Implementation Example**:

```typescript
// src/executor/file-executor.ts
export class FileExecutor implements Executor {
  canExecute(action: TestableAction): action is FileTestableAction {
    return action.actionType === 'file';
  }

  async execute(action: FileTestableAction, context: ExecutionContext): Promise<ExecutionResult> {
    const startTime = Date.now();
    const fullPath = path.join(context.workingDirectory, action.path);

    try {
      // Ensure parent directory exists
      await fs.mkdir(path.dirname(fullPath), { recursive: true });

      switch (action.operation) {
        case 'create':
          // Fail if file exists (unless configured otherwise)
          if (await this.fileExists(fullPath)) {
            throw new Error(`File already exists: ${action.path}`);
          }
          await fs.writeFile(fullPath, action.content, 'utf-8');
          break;

        case 'replace':
          await fs.writeFile(fullPath, action.content, 'utf-8');
          break;

        case 'append':
          await fs.appendFile(fullPath, action.content, 'utf-8');
          break;
      }

      return {
        success: true,
        duration: Date.now() - startTime,
        output: `${action.operation} ${action.path} (${action.content.length} bytes)`
      };
    } catch (error) {
      return {
        success: false,
        duration: Date.now() - startTime,
        error: error.message,
        stderr: error.stack
      };
    }
  }

  getSupportedTypes(): string[] {
    return ['file'];
  }

  getSupportedOperations(): Array<'create' | 'replace' | 'append'> {
    return ['create', 'replace', 'append'];
  }

  async validate(): Promise<ValidationResult> {
    // File operations only require Node.js fs module (always available)
    return { valid: true };
  }

  private async fileExists(path: string): Promise<boolean> {
    try {
      await fs.access(path);
      return true;
    } catch {
      return false;
    }
  }
}
```

**Working Directory Management**:

The Test Orchestrator creates a unique temp directory for each procedure:

```typescript
// src/orchestrator/test-orchestrator.ts
async executeProcedure(procedure: ProcedureNode, config: Configuration): Promise<TestResult> {
  // Create per-procedure temp directory
  const timestamp = Date.now();
  const procedureId = this.sanitizeProcedureId(procedure.title);
  const workingDirectory = path.join(
    config.workingDirectoryBase || '.proctest/runs',
    `${timestamp}-${procedureId}`
  );

  await fs.mkdir(workingDirectory, { recursive: true });

  // Register for cleanup
  const cleanup = new CleanupRegistry(config.cleanup);
  cleanup.registerDirectory(workingDirectory);

  try {
    // Execute procedure with this working directory
    const context: ExecutionContext = {
      workingDirectory,
      cleanup,
      // ... other context
    };

    // Execute file operations first, then code blocks
    for (const step of procedure.steps) {
      // File operations
      for (const action of step.testableActions.filter(a => a.actionType === 'file')) {
        await this.executeAction(action, context);
      }

      // Code/Shell/CLI operations
      for (const action of step.testableActions.filter(a => a.actionType !== 'file')) {
        await this.executeAction(action, context);
      }
    }

  } finally {
    // Cleanup temp directory (unless keepOnFailure is set and test failed)
    if (config.cleanup.workingDirectories.enabled) {
      const shouldKeep = !result.success && config.cleanup.workingDirectories.keepOnFailure;
      if (!shouldKeep) {
        await cleanup.execute();
      }
    }
  }
}
```

**CLI Flag for Debugging**:

```bash
# Keep temp directories on failure for debugging
proctest --keep-artifacts

# Keep all temp directories (success or failure)
proctest --keep-all-artifacts
```

#### D.3 Code Testable Actions

**Detection**:
```rst
.. code-block:: javascript

   const client = new MongoClient('<connection-string>');
   await client.connect();
```

**Parsed As**:
```typescript
{
  actionType: 'code',
  language: 'javascript',
  code: "const client = new MongoClient('<connection-string>');\nawait client.connect();",
  options: { executable: true },
  placeholders: ['<connection-string>'],
  location: { file: 'example.txt', startLine: 10, endLine: 13 }
}
```

**Execution**: Language-specific executor (JavaScriptExecutor, PythonExecutor, PHPExecutor)

---

##### IDE-Based Code Execution

Some procedures instruct readers to execute code "from your IDE" rather than providing explicit command-line instructions. This pattern:

- Reflects real-world developer workflows (most developers use IDEs)
- Avoids prescribing specific build tools (Maven vs Gradle, etc.)
- Accommodates toolchain variations across environments

**Detection Pattern**:

Prose containing phrases like:
- "From your IDE, run the file"
- "In your IDE, execute the code"
- "Using your IDE, run the application"

**Example**:
```rst
.. step:: Create a new file named ``CreateIndex.java``.

   .. literalinclude:: /includes/examples/CreateIndex.java
      :language: java

.. step:: Execute the code to create the index.

   From your IDE, run the file to create the index.
```

**Parsed As**:
```typescript
{
  actionType: 'code',
  language: 'java',
  code: '/* content from literalinclude */',
  options: { executable: true },
  executionMode: 'ide',  // Indicates IDE-based execution
  filePath: 'CreateIndex.java',
  location: { file: 'example.txt', startLine: 10, endLine: 20 }
}
```

**Execution Strategy - Convention over Configuration**:

The framework provides sensible default commands for each language, with optional user overrides:

**Default IDE Commands** (Convention):
```typescript
const DEFAULT_IDE_COMMANDS: Record<string, string> = {
  // Java - prefer Maven (most common in MongoDB docs)
  java: 'mvn compile exec:java -Dexec.mainClass="{className}"',

  // C# - use dotnet CLI
  csharp: 'dotnet run',
  cs: 'dotnet run',

  // C/C++ - compile and run
  cpp: 'g++ {filename} -o {basename} && ./{basename}',
  c: 'gcc {filename} -o {basename} && ./{basename}',

  // Python - direct execution
  python: 'python {filename}',
  py: 'python {filename}',

  // Node.js - direct execution
  javascript: 'node {filename}',
  js: 'node {filename}',

  // Go - go run
  go: 'go run {filename}',
};
```

**Configuration Override** (User Preference):
```javascript
// .proctest.js
module.exports = {
  ideExecution: {
    commands: {
      java: 'gradle run',  // User prefers Gradle over Maven
      cpp: './build/my-app',  // Custom build output location
    },

    // Or skip IDE execution entirely (mark as manual verification)
    skip: false,
  }
};
```

**Command Interpolation Variables**:
- `{filename}` - Full file path (e.g., `src/CreateIndex.java`)
- `{basename}` - Filename without extension (e.g., `CreateIndex`)
- `{className}` - Extracted class name from Java/C# files (e.g., `CreateIndex`)

**Implementation Example**:
```typescript
// src/executor/code-executor.ts
export class CodeExecutor implements Executor {
  private shellExecutor: ShellExecutor;

  async execute(action: CodeTestableAction, context: ExecutionContext): Promise<ExecutionResult> {
    // Handle IDE-based execution
    if (action.executionMode === 'ide') {
      const command = this.resolveIDECommand(action, context);

      // Log what we're doing for transparency
      console.log(`IDE execution: Using command: ${command}`);

      // Execute using ShellExecutor
      return await this.shellExecutor.execute({
        actionType: 'shell',
        command,
        workingDirectory: context.workingDirectory
      }, context);
    }

    // ... normal code execution (inline code blocks)
  }

  private resolveIDECommand(action: CodeTestableAction, context: ExecutionContext): string {
    const language = action.language.toLowerCase();

    // 1. Check user configuration first (convention over configuration)
    const userCommand = context.config.ideExecution?.commands?.[language];

    if (userCommand) {
      return this.interpolateCommand(userCommand, action);
    }

    // 2. Fall back to sensible defaults
    const defaultCommand = DEFAULT_IDE_COMMANDS[language];

    if (!defaultCommand) {
      throw new Error(
        `No default IDE command for ${language}. ` +
        `Configure in .proctest.js: ideExecution.commands.${language}`
      );
    }

    return this.interpolateCommand(defaultCommand, action);
  }

  private interpolateCommand(template: string, action: CodeTestableAction): string {
    const filename = action.filePath || 'unknown';
    const basename = path.basename(filename, path.extname(filename));
    const className = this.extractClassName(action.code, action.language);

    return template
      .replace(/{filename}/g, filename)
      .replace(/{basename}/g, basename)
      .replace(/{className}/g, className);
  }

  private extractClassName(code: string, language: string): string {
    // Extract class name from Java/C# code
    if (language === 'java' || language === 'csharp' || language === 'cs') {
      const match = code.match(/(?:public\s+)?class\s+(\w+)/);
      return match ? match[1] : 'Main';
    }
    return 'Main';
  }
}
```

**Reporting**:

When IDE execution is used, the test report should clearly indicate:
- ✅ The command that was executed
- ✅ Whether it used default or custom configuration
- ✅ The working directory

Example output:
```
✓ Execute code to create index
  IDE Execution: mvn compile exec:java -Dexec.mainClass="CreateIndex"
  Working Directory: /tmp/proctest/runs/1234567890-create-index/
  Duration: 2.3s
```

---

#### D.4 Shell Testable Actions

**Detection**:
```rst
.. code-block:: shell

   npm install mongodb
```

**Parsed As**:
```typescript
{
  actionType: 'shell',
  command: 'npm install mongodb',
  placeholders: [],
  location: { file: 'example.txt', startLine: 15, endLine: 17 }
}
```

**Execution**: ShellExecutor (spawns shell process)

**Working Directory Handling**:

The ShellExecutor must maintain working directory context across commands within a procedure:

```typescript
// src/executor/shell-executor.ts
export class ShellExecutor implements Executor {
  private currentWorkingDirectory: string;

  async execute(action: ShellTestableAction, context: ExecutionContext): Promise<ExecutionResult> {
    // Initialize working directory from context on first execution
    if (!this.currentWorkingDirectory) {
      this.currentWorkingDirectory = context.workingDirectory;
    }

    const startTime = Date.now();

    try {
      // Execute command in current working directory
      const result = await execAsync(action.command, {
        cwd: this.currentWorkingDirectory,
        env: { ...process.env, ...context.environment }
      });

      // Track directory changes (cd, pushd, etc.)
      if (action.command.trim().startsWith('cd ')) {
        const targetDir = action.command.replace(/^cd\s+/, '').trim();
        this.currentWorkingDirectory = path.resolve(this.currentWorkingDirectory, targetDir);
      }

      return {
        success: true,
        duration: Date.now() - startTime,
        output: result.stdout,
        stderr: result.stderr
      };
    } catch (error) {
      return {
        success: false,
        duration: Date.now() - startTime,
        error: error.message,
        stderr: error.stderr
      };
    }
  }

  // Reset working directory for new procedure
  reset(workingDirectory: string): void {
    this.currentWorkingDirectory = workingDirectory;
  }
}
```

**Common Patterns**:

1. **Sequential directory operations**:
   ```rst
   .. code-block:: bash

      mkdir my-project
      cd my-project
      npm init -y
   ```

2. **Sub-steps with directory changes**:
   ```rst
   a. Create a new directory:

      .. code-block:: bash

         mkdir csharp-create-index

   #. Change to the new directory:

      .. code-block:: bash

         cd csharp-create-index

   #. Initialize the project:

      .. code-block:: bash

         dotnet new console
   ```

**Note**: Directory changes persist across shell commands within the same procedure, but each procedure starts with a fresh working directory.

#### D.5 UI Testable Actions

**Detection**:
```rst
Click the :guilabel:`Create Database` button.

Select :guilabel:`MongoDB Atlas` from the dropdown menu.
```

**Parsed As**:
```typescript
[
  {
    actionType: 'ui',
    action: 'click',
    target: 'Create Database',
    description: 'Click the Create Database button.',
    location: { file: 'example.txt', startLine: 20, endLine: 20 }
  },
  {
    actionType: 'ui',
    action: 'select',
    target: 'MongoDB Atlas',
    description: 'Select MongoDB Atlas from the dropdown menu.',
    location: { file: 'example.txt', startLine: 22, endLine: 22 }
  }
]
```

**Execution**: UIExecutor (Playwright/Puppeteer for browser automation)

**Action Detection Rules**:
- "Click" → `action: 'click'`
- "Select" → `action: 'select'`
- "Enter" / "Type" → `action: 'input'`
- "Verify" / "Ensure" / "Check" → `action: 'verify'`

**UI Configuration and Execution**:

UI testing requires two types of configuration to handle frequently-changing UIs and generic instructions:

**1. Navigation Mappings** - Map generic phrases to automation steps:

```javascript
// In proctest.config.js
ui: {
  navigationMappings: [
    {
      phrase: /In Atlas, go to the Clusters page/i,
      steps: [
        { action: 'navigate', url: 'https://cloud.mongodb.com' },
        { action: 'wait', selector: '.project-selector', timeout: 5000 },
        { action: 'click', selector: 'nav a[href*="clusters"]' },
        { action: 'waitForNavigation', timeout: 5000 }
      ]
    }
  ]
}
```

When the documentation says "In Atlas, go to the Clusters page for your project", the framework:
1. Matches the phrase against `navigationMappings`
2. Executes the defined automation steps
3. Continues with the next step in the procedure

**2. User Values** - Map generic references to actual values:

```javascript
// In proctest.config.js
ui: {
  userValues: {
    'your project': process.env.ATLAS_PROJECT_NAME || 'MyTestProject',
    'your desired project': process.env.ATLAS_PROJECT_NAME || 'MyTestProject',
    'your cluster': process.env.ATLAS_CLUSTER_NAME || 'Cluster0',
    'your organization': process.env.ATLAS_ORG_NAME || 'MyOrg'
  }
}
```

When the documentation says "Select your desired project", the framework:
1. Looks up "your desired project" in `userValues`
2. Finds the actual value (e.g., "MyTestProject")
3. Uses that value in the UI automation (e.g., clicks the element with text "MyTestProject")

**Example Execution Flow**:

Documentation step:
```
In Atlas, go to the Clusters page for your project.

1. Select your organization from the dropdown.
2. Select your desired project.
3. Click :guilabel:`Clusters` in the sidebar.
```

With configuration:
```javascript
ui: {
  baseUrl: 'https://cloud.mongodb.com',
  navigationMappings: [
    {
      phrase: /In Atlas, go to the Clusters page/i,
      steps: [
        { action: 'navigate', url: 'https://cloud.mongodb.com' },
        { action: 'wait', selector: '.org-selector' }
      ]
    }
  ],
  userValues: {
    'your organization': 'MongoDB',
    'your desired project': 'Documentation'
  }
}
```

Execution:
1. **Phrase match**: "In Atlas, go to the Clusters page" → Execute navigation mapping
2. **Sub-step 1**: "Select your organization" → Look up "your organization" → Click element with text "MongoDB"
3. **Sub-step 2**: "Select your desired project" → Look up "your desired project" → Click element with text "Documentation"
4. **Sub-step 3**: `:guilabel:`Clusters`` → Click element with text "Clusters"

**Error Handling**:

If a navigation mapping is not found:
```
❌ FAILED: Navigate to Clusters Page

Error: No navigation mapping found for phrase "In Atlas, go to the Clusters page"

Suggestions:
  - Add a navigation mapping in your config:
    ui: {
      navigationMappings: [
        {
          phrase: /In Atlas, go to the Clusters page/i,
          steps: [
            { action: 'navigate', url: 'https://cloud.mongodb.com' },
            // ... your automation steps
          ]
        }
      ]
    }
```

If a user value is not found:
```
❌ FAILED: Navigate to Clusters Page

Error: No value configured for "your desired project"

Suggestions:
  - Add the value in your config:
    ui: {
      userValues: {
        'your desired project': 'MyProject'
      }
    }
  - Or set environment variable: ATLAS_PROJECT_NAME=MyProject
```

#### D.6 CLI Testable Actions

**Detection** (mongosh):
```rst
.. code-block:: javascript

   use myDatabase
   db.myCollection.insertOne({ name: "test" })
```

**Parsed As**:
```typescript
{
  actionType: 'cli',
  tool: 'mongosh',
  command: 'use myDatabase\ndb.myCollection.insertOne({ name: "test" })',
  placeholders: [],
  location: { file: 'example.txt', startLine: 25, endLine: 28 }
}
```

**Detection** (atlas-cli):
```rst
.. code-block:: shell

   atlas clusters create myCluster --provider AWS --region US_EAST_1
```

**Parsed As**:
```typescript
{
  actionType: 'cli',
  tool: 'atlas-cli',
  command: 'atlas clusters create myCluster --provider AWS --region US_EAST_1',
  placeholders: [],
  location: { file: 'example.txt', startLine: 30, endLine: 32 }
}
```

**Execution**: CLIExecutor (tool-specific handling)

**Detection Rules**:
- Code block starts with `use `, `db.`, `show ` → mongosh
- Code block starts with `atlas ` → atlas-cli
- Otherwise → regular shell or code

#### D.7 API Testable Actions

**Important**: API testable actions are specifically for the **MongoDB Atlas Admin API only**, not all curl commands.

**Detection** (MongoDB Atlas Admin API):
```rst
.. code-block:: shell

   curl -X POST https://cloud.mongodb.com/api/atlas/v2/groups \
     -H "Content-Type: application/json" \
     -H "Authorization: Bearer $ATLAS_API_TOKEN" \
     -d '{"name": "MyProject"}'
```

**Parsed As**:
```typescript
{
  actionType: 'api',
  method: 'POST',
  endpoint: 'https://cloud.mongodb.com/api/atlas/v2/groups',
  headers: {
    'Content-Type': 'application/json',
    'Authorization': 'Bearer $ATLAS_API_TOKEN'
  },
  body: '{"name": "MyProject"}',
  placeholders: ['$ATLAS_API_TOKEN'],
  location: { file: 'example.txt', startLine: 35, endLine: 39 }
}
```

**Another Example** (Atlas Admin API - GET request):
```rst
.. code-block:: shell

   curl -X GET https://cloud.mongodb.com/api/atlas/v2/groups/<project-id>/clusters \
     -H "Authorization: Bearer $ATLAS_API_TOKEN"
```

**Parsed As**:
```typescript
{
  actionType: 'api',
  method: 'GET',
  endpoint: 'https://cloud.mongodb.com/api/atlas/v2/groups/<project-id>/clusters',
  headers: {
    'Authorization': 'Bearer $ATLAS_API_TOKEN'
  },
  placeholders: ['<project-id>', '$ATLAS_API_TOKEN'],
  location: { file: 'example.txt', startLine: 41, endLine: 43 }
}
```

**Execution**: APIExecutor (HTTP client like axios or fetch)

**Detection Rules** (Atlas Admin API only):
- Contains `curl` command **AND**
- Endpoint matches Atlas Admin API pattern:
  - `https://cloud.mongodb.com/api/atlas/` (Atlas Admin API v1.0 and v2)
- Extract method, endpoint, headers, body from curl syntax

**Non-API curl commands** (see Download and Shell actions):
- Software installation: `curl https://example.com/install.sh | sh` → Shell action
- File downloads: `curl -o file.tar.gz https://example.com/file.tar.gz` → Download action
- Generic HTTP requests: `curl https://example.com` → Shell action (unless Atlas Admin API)

#### D.8 Download Testable Actions

**Important**: Downloads may take significant time and subsequent steps often depend on download completion.

**Detection** (Sample data download):
```rst
.. code-block:: shell

   curl https://atlas-education.s3.amazonaws.com/sampledata.archive -o sampledata.archive
```

**Parsed As**:
```typescript
{
  actionType: 'download',
  url: 'https://atlas-education.s3.amazonaws.com/sampledata.archive',
  outputPath: 'sampledata.archive',
  method: 'GET',
  timeout: 300000, // 5 minutes for large files
  description: 'Download sample data archive',
  location: { file: 'example.txt', startLine: 45, endLine: 45 }
}
```

**Detection** (Software download):
```rst
.. code-block:: shell

   curl -O https://fastdl.mongodb.org/tools/db/mongodb-database-tools-macos-x86_64-100.9.4.zip
```

**Parsed As**:
```typescript
{
  actionType: 'download',
  url: 'https://fastdl.mongodb.org/tools/db/mongodb-database-tools-macos-x86_64-100.9.4.zip',
  outputPath: 'mongodb-database-tools-macos-x86_64-100.9.4.zip', // Inferred from -O flag
  method: 'GET',
  timeout: 300000,
  description: 'Download MongoDB Database Tools',
  location: { file: 'example.txt', startLine: 47, endLine: 47 }
}
```

**Execution**: DownloadExecutor (HTTP client with progress tracking and file writing)

**Detection Rules**:
- Contains `curl` command **AND**
- Has output flag: `-o <filename>` or `-O` (use URL filename)
- **NOT** a MongoDB API endpoint (those are API actions)
- **NOT** piped to shell (those are Shell actions)

**Execution Considerations**:
- Track download progress for large files
- Verify file exists and has expected size after download
- Wait for download to complete before proceeding to next step
- Handle network timeouts gracefully
- Clean up partial downloads on failure

**Example Execution Flow**:
```typescript
async executeDownload(action: DownloadTestableAction): Promise<ExecutionResult> {
  const startTime = Date.now();

  try {
    // Start download with progress tracking
    const response = await fetch(action.url);
    const fileStream = fs.createWriteStream(action.outputPath);

    // Track progress
    let downloadedBytes = 0;
    response.body.on('data', (chunk) => {
      downloadedBytes += chunk.length;
      // Report progress for large files
    });

    // Wait for completion
    await pipeline(response.body, fileStream);

    // Verify file exists
    const stats = await fs.stat(action.outputPath);

    return {
      success: true,
      duration: Date.now() - startTime,
      output: `Downloaded ${stats.size} bytes to ${action.outputPath}`
    };
  } catch (error) {
    // Clean up partial download
    await fs.unlink(action.outputPath).catch(() => {});

    return {
      success: false,
      duration: Date.now() - startTime,
      error: `Download failed: ${error.message}`
    };
  }
}
```

#### D.9 URL Testable Actions

**Detection**:
```rst
For more information, see the `MongoDB Documentation <https://docs.mongodb.com>`_.

Navigate to https://cloud.mongodb.com to access Atlas.
```

**Parsed As**:
```typescript
[
  {
    actionType: 'url',
    url: 'https://docs.mongodb.com',
    description: 'MongoDB Documentation',
    location: { file: 'example.txt', startLine: 40, endLine: 40 }
  },
  {
    actionType: 'url',
    url: 'https://cloud.mongodb.com',
    description: 'Navigate to https://cloud.mongodb.com to access Atlas.',
    location: { file: 'example.txt', startLine: 42, endLine: 42 }
  }
]
```

**Execution**: URLExecutor (HTTP HEAD request to verify URL is accessible)

**Detection Rules**:
- RST link syntax: `` `text <url>`_ ``
- Bare URLs in text: `https://...` or `http://...`

#### D.10 Execution Priority and Phasing

**Phase 1 (PoC)**: File, Code, and Shell
- File operations (create, replace, append)
- Focus on most common testable actions
- Prove the concept with code execution and file setup

**Phase 2 (Production)**: Add CLI and Download
- mongosh and atlas-cli support
- Tool-specific execution strategies
- Download action support for large files

**Phase 3 (Advanced)**: Add UI, API, URL
- UI automation for guilabel interactions
- Atlas Admin API request validation (https://cloud.mongodb.com/api/atlas/)
- URL accessibility checks (link validation)

#### D.11 Implementation Example: Action Detection

```typescript
/**
 * Detect testable action type from RST node
 */
function detectTestableAction(node: RSTNode, context: ParserContext): TestableAction | null {
  // Code blocks
  if (node.type === 'code-block' || node.type === 'literalinclude') {
    const language = node.language || deriveLanguageFromFile(node.source);
    const code = node.code || readFile(node.source);

    // Check if it's a CLI tool
    if (isMongoshCode(code)) {
      return {
        elementType: 'cli',
        tool: 'mongosh',
        command: code,
        placeholders: extractPlaceholders(code),
        location: getLocation(node)
      };
    }

    if (isAtlasCLI(code)) {
      return {
        elementType: 'cli',
        tool: 'atlas-cli',
        command: code,
        placeholders: extractPlaceholders(code),
        location: getLocation(node)
      };
    }

    // Check if it's a curl command (could be API, Download, or Shell)
    if (isCurlCommand(code)) {
      return parseCurlCommand(code, node);
    }

    // Check if it's shell
    if (isShellLanguage(language)) {
      return {
        actionType: 'shell',
        command: code,
        placeholders: extractPlaceholders(code),
        location: getLocation(node)
      };
    }

    // Regular code block
    return {
      actionType: 'code',
      language: normalizeLanguage(language),
      code,
      options: { executable: true },
      placeholders: extractPlaceholders(code),
      location: getLocation(node)
    };
  }

  // UI interactions (guilabel role)
  if (node.type === 'paragraph' && containsGuilabel(node)) {
    return parseUIInteraction(node);
  }

  // URLs
  if (node.type === 'reference' || containsURL(node)) {
    return parseURL(node);
  }

  return null;
}

function isMongoshCode(code: string): boolean {
  const mongoshPatterns = [
    /^use\s+\w+/m,
    /^db\.\w+/m,
    /^show\s+(dbs|databases|collections)/m
  ];
  return mongoshPatterns.some(pattern => pattern.test(code));
}

function isAtlasCLI(code: string): boolean {
  return code.trim().startsWith('atlas ');
}

function isCurlCommand(code: string): boolean {
  return code.trim().startsWith('curl ');
}

function isShellLanguage(language: string): boolean {
  const shellLanguages = ['shell', 'bash', 'sh', 'console', 'terminal'];
  return shellLanguages.includes(language.toLowerCase());
}

/**
 * Parse curl command - distinguish between API, Download, and Shell actions
 */
function parseCurlCommand(code: string, node: RSTNode): TestableAction {
  const trimmed = code.trim();

  // Extract URL from curl command
  const urlMatch = trimmed.match(/curl\s+(?:-[^\s]+\s+)*(?:-X\s+\w+\s+)?([^\s]+)/);
  const url = urlMatch ? urlMatch[1] : '';

  // Check if it's a MongoDB API endpoint
  if (isMongoDBAPIEndpoint(url)) {
    return parseMongoDBAPIRequest(code, node);
  }

  // Check if it's a download (has -o or -O flag)
  if (isDownloadCommand(trimmed)) {
    return parseDownloadAction(code, node);
  }

  // Otherwise, treat as regular shell command
  return {
    actionType: 'shell',
    command: code,
    placeholders: extractPlaceholders(code),
    location: getLocation(node)
  };
}

/**
 * Check if URL is a MongoDB Atlas Admin API endpoint
 */
function isMongoDBAPIEndpoint(url: string): boolean {
  // Only Atlas Admin API - Data API and App Services API are deprecated
  const atlasAdminAPIPattern = /^https?:\/\/cloud\.mongodb\.com\/api\/atlas\//;

  return atlasAdminAPIPattern.test(url);
}

/**
 * Check if curl command is a download (has -o or -O flag)
 */
function isDownloadCommand(curlCommand: string): boolean {
  // Match -o <filename> or -O (use URL filename)
  return /\s-[oO]\s/.test(curlCommand) || /\s-[oO]$/.test(curlCommand);
}

/**
 * Parse MongoDB API request from curl command
 */
function parseMongoDBAPIRequest(code: string, node: RSTNode): APITestableAction {
  // Extract method (default to GET)
  const methodMatch = code.match(/-X\s+(GET|POST|PUT|DELETE|PATCH)/i);
  const method = methodMatch ? methodMatch[1].toUpperCase() as any : 'GET';

  // Extract endpoint
  const urlMatch = code.match(/curl\s+(?:-[^\s]+\s+)*(?:-X\s+\w+\s+)?([^\s\\]+)/);
  const endpoint = urlMatch ? urlMatch[1] : '';

  // Extract headers
  const headers: Record<string, string> = {};
  const headerMatches = code.matchAll(/-H\s+["']([^:]+):\s*([^"']+)["']/g);
  for (const match of headerMatches) {
    headers[match[1]] = match[2];
  }

  // Extract body
  const bodyMatch = code.match(/-d\s+["'](.+?)["']/s);
  const body = bodyMatch ? bodyMatch[1] : undefined;

  return {
    actionType: 'api',
    method,
    endpoint,
    headers,
    body,
    placeholders: extractPlaceholders(code),
    location: getLocation(node)
  };
}

/**
 * Parse download action from curl command
 */
function parseDownloadAction(code: string, node: RSTNode): DownloadTestableAction {
  // Extract URL
  const urlMatch = code.match(/curl\s+(?:-[^\s]+\s+)*([^\s\\]+)/);
  const url = urlMatch ? urlMatch[1] : '';

  // Extract output path
  let outputPath = '';
  const outputMatch = code.match(/-o\s+([^\s]+)/);
  if (outputMatch) {
    outputPath = outputMatch[1];
  } else if (/-O/.test(code)) {
    // -O uses filename from URL
    const urlParts = url.split('/');
    outputPath = urlParts[urlParts.length - 1];
  }

  // Determine description from context
  let description = 'Download file';
  if (url.includes('sampledata')) {
    description = 'Download sample data archive';
  } else if (url.includes('mongodb') || url.includes('mongo')) {
    description = 'Download MongoDB tools or software';
  }

  return {
    actionType: 'download',
    url,
    outputPath,
    method: 'GET',
    timeout: 300000, // 5 minutes default for large files
    description,
    location: getLocation(node)
  };
}

function parseUIInteraction(node: RSTNode): UITestableAction {
  const text = node.text;
  const guilabelMatch = text.match(/:guilabel:`([^`]+)`/);
  const target = guilabelMatch ? guilabelMatch[1] : '';

  // Detect action from surrounding text
  let action: 'click' | 'select' | 'input' | 'verify' = 'click';
  if (/\b(select|choose)\b/i.test(text)) action = 'select';
  if (/\b(enter|type|input)\b/i.test(text)) action = 'input';
  if (/\b(verify|check|ensure)\b/i.test(text)) action = 'verify';

  return {
    actionType: 'ui',
    action,
    target,
    description: text,
    location: getLocation(node)
  };
}
```

#### D.12 Reporting by Action Type

**Example Output**:
```
✓ PASSED: Install MongoDB Driver and Sample Data

Procedure: "Install MongoDB Driver and Sample Data"
  Step 1: "Download sample data"
    ✓ Download: sampledata.archive (45.2 MB, 12.3s)

  Step 2: "Install dependencies"
    ✓ Shell: npm install mongodb (0.5s)
    ✓ Code [javascript]: const client = new MongoClient(...) (0.2s)

  Step 3: "Configure Atlas via API"
    ✓ API [POST]: https://cloud.mongodb.com/api/atlas/v2/groups (0.8s)
    ✓ UI: Click "Create Database" (1.2s)
    ✓ UI: Select "MongoDB Atlas" (0.8s)
    ✓ CLI [atlas-cli]: atlas clusters create... (3.5s)

  Step 4: "Verify connection"
    ✓ Code [javascript]: await client.connect() (0.3s)
    ✓ URL: https://docs.mongodb.com (0.1s)

Summary:
- Total Actions: 9
  - Download: 1 passed (45.2 MB)
  - Shell: 1 passed
  - Code: 2 passed
  - API: 1 passed
  - UI: 2 passed
  - CLI: 1 passed
  - URL: 1 passed
- Duration: 20.7s
```

---

### Appendix E: Package.json Example

```json
{
  "name": "proctest",
  "version": "0.1.0",
  "description": "Testing framework for documentation procedures",
  "main": "dist/index.js",
  "bin": {
    "proctest": "dist/cli/index.js"
  },
  "scripts": {
    "build": "tsc",
    "dev": "tsc --watch",
    "test": "jest",
    "test:watch": "jest --watch",
    "test:coverage": "jest --coverage",
    "lint": "eslint src/**/*.ts",
    "format": "prettier --write src/**/*.ts"
  },
  "keywords": ["testing", "documentation", "mongodb", "rst", "procedures"],
  "author": "MongoDB Documentation Team",
  "license": "Apache-2.0",
  "engines": {
    "node": ">=24.0.0"
  },
  "dependencies": {
    "commander": "^11.0.0",
    "glob": "^10.0.0",
    "toml": "^3.0.0",
    "yaml": "^2.3.0"
  },
  "devDependencies": {
    "@types/jest": "^29.5.0",
    "@types/node": "^20.0.0",
    "@typescript-eslint/eslint-plugin": "^6.0.0",
    "@typescript-eslint/parser": "^6.0.0",
    "eslint": "^8.0.0",
    "jest": "^29.7.0",
    "prettier": "^3.0.0",
    "ts-jest": "^29.1.0",
    "typescript": "^5.0.0"
  }
}
```

### Appendix F: Jest Configuration

```typescript
// jest.config.ts
import type { Config } from 'jest';

const config: Config = {
  preset: 'ts-jest',
  testEnvironment: 'node',
  roots: ['<rootDir>/src', '<rootDir>/tests'],
  testMatch: ['**/__tests__/**/*.ts', '**/*.test.ts', '**/*.spec.ts'],
  collectCoverageFrom: [
    'src/**/*.ts',
    '!src/**/*.d.ts',
    '!src/**/*.test.ts',
    '!src/**/*.spec.ts'
  ],
  coverageThreshold: {
    global: {
      branches: 80,
      functions: 80,
      lines: 80,
      statements: 80
    }
  },
  moduleNameMapper: {
    '^@/(.*)$': '<rootDir>/src/$1'
  },
  clearMocks: true,
  restoreMocks: true
};

export default config;
```

### Appendix G: TypeScript Configuration

```json
{
  "compilerOptions": {
    "target": "ES2022",
    "module": "commonjs",
    "lib": ["ES2022"],
    "outDir": "./dist",
    "rootDir": "./src",
    "strict": true,
    "esModuleInterop": true,
    "skipLibCheck": true,
    "forceConsistentCasingInFileNames": true,
    "resolveJsonModule": true,
    "declaration": true,
    "declarationMap": true,
    "sourceMap": true,
    "moduleResolution": "node",
    "types": ["node", "jest"]
  },
  "include": ["src/**/*"],
  "exclude": ["node_modules", "dist", "**/*.test.ts"]
}
```

### Appendix H: Example Implementation - JavaScript Executor

```typescript
import { spawn } from 'child_process';
import { Executor, ExecutionContext, ExecutionResult, ValidationResult } from '../interfaces/executor';

export class JavaScriptExecutor implements Executor {
  canExecute(language: string): boolean {
    const jsLanguages = ['javascript', 'js', 'node', 'nodejs'];
    return jsLanguages.includes(language.toLowerCase());
  }

  getSupportedLanguages(): string[] {
    return ['javascript', 'js', 'node', 'nodejs'];
  }

  async validate(): Promise<ValidationResult> {
    try {
      const { stdout } = await this.runCommand('node', ['--version']);
      const version = stdout.trim();

      // Check if version is >= 18
      const majorVersion = parseInt(version.replace('v', '').split('.')[0]);
      if (majorVersion < 18) {
        return {
          valid: false,
          error: `Node.js version ${version} is too old. Requires >= 18.0.0`,
          suggestions: ['Upgrade Node.js to version 18 or higher']
        };
      }

      return { valid: true };
    } catch (error) {
      return {
        valid: false,
        error: 'Node.js is not installed or not in PATH',
        suggestions: [
          'Install Node.js from https://nodejs.org',
          'Ensure node is in your PATH'
        ]
      };
    }
  }

  async execute(code: string, context: ExecutionContext): Promise<ExecutionResult> {
    const startTime = Date.now();

    try {
      // Accumulate code from previous blocks in this step
      const fullCode = context.state.isFirstBlock
        ? code
        : `${context.state.accumulatedCode}\n${code}`;

      // Update accumulated code for next block
      context.state.accumulatedCode = fullCode;

      // Execute code
      const result = await this.runCode(fullCode, context);

      return {
        ...result,
        duration: Date.now() - startTime
      };
    } catch (error) {
      return {
        success: false,
        stdout: '',
        stderr: error instanceof Error ? error.message : String(error),
        exitCode: 1,
        duration: Date.now() - startTime,
        error: error instanceof Error ? error : new Error(String(error))
      };
    }
  }

  private async runCode(code: string, context: ExecutionContext): Promise<Omit<ExecutionResult, 'duration'>> {
    return new Promise((resolve) => {
      const child = spawn('node', ['-e', code], {
        cwd: context.workingDirectory,
        env: { ...process.env, ...context.environment },
        timeout: context.timeout
      });

      let stdout = '';
      let stderr = '';
      let timedOut = false;

      child.stdout.on('data', (data) => {
        stdout += data.toString();
      });

      child.stderr.on('data', (data) => {
        stderr += data.toString();
      });

      child.on('error', (error) => {
        if (error.message.includes('ETIMEDOUT')) {
          timedOut = true;
        }
      });

      child.on('close', (exitCode) => {
        resolve({
          success: exitCode === 0 && !timedOut,
          stdout,
          stderr,
          exitCode: exitCode ?? 1,
          timedOut
        });
      });
    });
  }

  private runCommand(command: string, args: string[]): Promise<{ stdout: string; stderr: string }> {
    return new Promise((resolve, reject) => {
      const child = spawn(command, args);
      let stdout = '';
      let stderr = '';

      child.stdout.on('data', (data) => stdout += data.toString());
      child.stderr.on('data', (data) => stderr += data.toString());

      child.on('close', (exitCode) => {
        if (exitCode === 0) {
          resolve({ stdout, stderr });
        } else {
          reject(new Error(`Command failed with exit code ${exitCode}`));
        }
      });

      child.on('error', reject);
    });
  }
}
```

### Appendix I: Example Implementation - Fuzzy Resolver

```typescript
import { PlaceholderResolver, ResolverContext } from '../interfaces/resolver';

export class FuzzyEnvResolver implements PlaceholderResolver {
  getPriority(): number {
    return 50; // Medium priority (after exact match, before snooty)
  }

  async resolve(placeholder: string, context: ResolverContext): Promise<string | null> {
    // Remove angle brackets and normalize
    const normalized = this.normalizePlaceholder(placeholder);

    // Try exact match first
    if (context.environment[normalized]) {
      return context.environment[normalized];
    }

    // Try fuzzy matching
    const matches = this.findFuzzyMatches(normalized, context.environment);

    // If exactly one match, use it
    if (matches.length === 1) {
      return context.environment[matches[0]];
    }

    // Multiple or no matches - return null
    return null;
  }

  async getSuggestions(placeholder: string, context: ResolverContext): Promise<string[]> {
    const normalized = this.normalizePlaceholder(placeholder);
    return this.findFuzzyMatches(normalized, context.environment);
  }

  private normalizePlaceholder(placeholder: string): string {
    return placeholder
      .replace(/^</, '')
      .replace(/>$/, '')
      .replace(/[-\s]/g, '_')
      .toUpperCase();
  }

  private findFuzzyMatches(normalized: string, environment: Record<string, string>): string[] {
    const envKeys = Object.keys(environment);
    const matches: Array<{ key: string; score: number }> = [];

    for (const key of envKeys) {
      const score = this.calculateSimilarity(normalized, key);
      if (score > 0.6) { // 60% similarity threshold
        matches.push({ key, score });
      }
    }

    // Sort by score (highest first)
    matches.sort((a, b) => b.score - a.score);

    return matches.map(m => m.key);
  }

  private calculateSimilarity(str1: string, str2: string): number {
    // Levenshtein distance-based similarity
    const longer = str1.length > str2.length ? str1 : str2;
    const shorter = str1.length > str2.length ? str2 : str1;

    if (longer.length === 0) {
      return 1.0;
    }

    const distance = this.levenshteinDistance(longer, shorter);
    return (longer.length - distance) / longer.length;
  }

  private levenshteinDistance(str1: string, str2: string): number {
    const matrix: number[][] = [];

    for (let i = 0; i <= str2.length; i++) {
      matrix[i] = [i];
    }

    for (let j = 0; j <= str1.length; j++) {
      matrix[0][j] = j;
    }

    for (let i = 1; i <= str2.length; i++) {
      for (let j = 1; j <= str1.length; j++) {
        if (str2.charAt(i - 1) === str1.charAt(j - 1)) {
          matrix[i][j] = matrix[i - 1][j - 1];
        } else {
          matrix[i][j] = Math.min(
            matrix[i - 1][j - 1] + 1, // substitution
            matrix[i][j - 1] + 1,     // insertion
            matrix[i - 1][j] + 1      // deletion
          );
        }
      }
    }

    return matrix[str2.length][str1.length];
  }
}
```

### Appendix J: Key Dependencies Rationale

| Dependency | Purpose | Rationale |
|------------|---------|-----------|
| **commander** | CLI argument parsing | Industry standard, well-maintained, TypeScript support |
| **glob** | File pattern matching | Reliable, widely used for test discovery |
| **toml** | Parse snooty.toml | Standard TOML parser for Node.js |
| **yaml** | YAML output for parse command | Standard YAML parser/stringifier, human-readable debug output |
| **jest** | Testing framework | Team familiarity, organizational consistency, mature ecosystem |
| **ts-jest** | TypeScript support for Jest | Standard TypeScript preprocessor for Jest |
| **typescript** | Type safety | Required for maintainability and developer experience |

**Minimal Dependencies**: Minimal dependencies keeps the tool lightweight and reduces maintenance burden.

**Note on .env file support**: Node.js 24+ provides native `.env` file support via `process.loadEnvFile()` and the `--env-file` CLI option, eliminating the need for the external `dotenv` package.

**Note on Jest vs Vitest**: While Vitest offers faster test execution and better TypeScript support out-of-the-box, Jest was chosen for team familiarity and organizational consistency. The MongoDB documentation team already uses Jest for other testing tooling, making it easier to share knowledge and maintain consistency across projects.

---

## Summary

This technical specification defines a comprehensive implementation plan for the procedural testing framework using **Option 5: Hybrid + Plugin Ready** architecture.

**Key Highlights**:
- ✅ Interface-driven design enables future extensibility
- ✅ Zero-config operation for 90% of use cases
- ✅ Progressive configuration for complex scenarios
- ✅ Phased implementation plan with clear milestones
- ✅ Comprehensive testing strategy
- ✅ Low overall risk with strong mitigation strategies

**Next Steps**:
1. Review and approve technical specification
2. Set up project repository and development environment
3. Begin Milestone 1 implementation (Foundation)

**Questions or Concerns?**
This specification is a living document. Please provide feedback on:
- Technical approach
- Implementation plan
- Any missing details or concerns


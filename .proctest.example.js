/**
 * Example proctest configuration file
 * 
 * Copy this file to `.proctest.js` in your project root and customize as needed.
 * 
 * This file demonstrates advanced configuration options. Most users won't need
 * to configure these - proctest uses sensible defaults.
 */

module.exports = {
  // ============================================================================
  // Test Discovery
  // ============================================================================
  
  // Files to test (glob patterns)
  testFiles: [
    'content/**/*.txt',
    'content/**/*.rst',
  ],
  
  // Files to exclude
  exclude: [
    '**/node_modules/**',
    '**/.git/**',
  ],

  // ============================================================================
  // Environment & Placeholders
  // ============================================================================
  
  // Environment files to load (in order)
  envFiles: [
    '.env',
    '.env.local',
  ],
  
  // Path to snooty.toml for source constants
  snootyConfig: 'snooty.toml',

  // ============================================================================
  // IDE Execution Configuration
  // ============================================================================
  // 
  // When documentation says "From your IDE, run the file", proctest needs to
  // know what command to execute. These settings override the defaults.
  // 
  // Default commands (you only need to override if different):
  //   java:       'mvn compile exec:java -Dexec.mainClass="{className}"'
  //   csharp/cs:  'dotnet run'
  //   cpp:        'g++ {filename} -o {basename} && ./{basename}'
  //   c:          'gcc {filename} -o {basename} && ./{basename}'
  //   python/py:  'python {filename}'
  //   javascript/js: 'node {filename}'
  //   go:         'go run {filename}'
  // 
  // Available interpolation variables:
  //   {filename}  - Full file path (e.g., 'src/CreateIndex.java')
  //   {basename}  - Filename without extension (e.g., 'CreateIndex')
  //   {className} - Extracted class name from Java/C# code (e.g., 'CreateIndex')
  // ============================================================================
  
  ideExecution: {
    // Override default commands for specific languages
    commands: {
      // Example: Use Gradle instead of Maven for Java
      java: 'gradle run',
      
      // Example: Use custom build output location for C++
      cpp: './build/bin/{basename}',
      
      // Example: Use python3 explicitly
      python: 'python3 {filename}',
      
      // Example: Custom command for a language without defaults
      rust: 'cargo run --bin {basename}',
    },
    
    // Skip IDE execution entirely (mark as manual verification)
    // Useful for CI environments where IDE execution isn't possible
    skip: false,
  },

  // ============================================================================
  // Execution Configuration
  // ============================================================================
  
  // Global timeout for all actions (milliseconds)
  timeout: 30000, // 30 seconds
  
  // Language-specific executor configuration
  executors: {
    javascript: {
      runtime: 'node',
      version: '>=18.0.0',
      timeout: 10000,
      env: {
        NODE_ENV: 'test',
      },
    },
    python: {
      runtime: 'python3',
      timeout: 15000,
    },
  },

  // ============================================================================
  // Cleanup Configuration
  // ============================================================================
  
  cleanup: {
    // Clean up working directories after tests
    workingDirectories: {
      enabled: true,
      keepOnFailure: true, // Keep for debugging
    },
    
    // Clean up test databases
    databases: {
      enabled: true,
      pattern: /^proctest_/,
      onFailure: 'warn',
    },
    
    // Clean up test collections
    collections: {
      enabled: true,
      pattern: /^test_/,
    },
  },

  // ============================================================================
  // Reporting
  // ============================================================================
  
  reporters: [
    { type: 'human' }, // Human-friendly console output
    { 
      type: 'json',
      options: {
        outputFile: 'test-results.json',
      },
    },
  ],
  
  verbose: false,

  // ============================================================================
  // Hooks (Advanced)
  // ============================================================================
  
  hooks: {
    // Run before all tests
    beforeAll: async (context) => {
      console.log('Setting up test environment...');
      // Example: Start local MongoDB instance
      // Example: Seed test data
    },
    
    // Run before each procedure
    beforeEach: async (procedure, context) => {
      // Example: Reset database state
    },
    
    // Run after each procedure
    afterEach: async (procedure, result, context) => {
      if (!result.success) {
        console.log(`Procedure failed: ${procedure.title}`);
      }
    },
    
    // Run after all tests
    afterAll: async (summary, context) => {
      console.log(`Tests complete: ${summary.passedProcedures}/${summary.totalProcedures} passed`);
      // Example: Tear down test environment
    },
  },
};


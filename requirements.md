# Procedure Testing Requirements

## Purpose

The purpose of the procedure testing framework is to validate that the
procedures in the documentation are accurate and up-to-date. The tooling
should execute each step in the procedure, report errors with specific
details for debugging purposes, and report success if all steps execute
successfully.

The procedures are written in plain text, with no programmatic structure. The
tooling must be able to parse the text, interpret the instructions, and
execute the steps.

## User Profile

The intended user of the procedure testing framework is a documentation writer
on the docs team. Tooling must be easy to use and understand for writers with
varying levels of technical expertise. Tooling cannot require a development
background to use.

## Required Functionality

### ReStructuredText Parsing

The tooling should be able to test procedures written in reStructuredText
files. It should be able to parse the reStructuredText files and extract the
procedures, including any code snippets, shell commands, UI interactions, and
API requests. It should not require writers to maintain separate parallel
procedures in a programmatic format in addition to the documentation written in
reStructuredText.

Any given reStructuredText file may contain multiple procedures. The tooling
should be able to identify and test each procedure independently.

Refer to Appendix B: ReStructuredText Syntax for details about relevant syntax
used in the MongoDB documentation.

### Supported Test Interactions

The tooling should be able to test procedures that include the following
types of interactions:

- Code snippets: should be able to test code snippets written in various languages
- Shell commands: should be able to test shell commands written for Unix.
  Windows is out of scope of the initial Poc.
- UI interactions: should be able to test UI interactions in a headless browser
  or other UI testing tool
- CLI requests: should be able to test CLI requests written for various CLI tools
- API requests: should be able to test API requests written for various APIs
- URL resolution: should be able to test that URLs are valid

### Interpolation

It should be able to test procedures as written, interpolating any placeholder
variables such as connection strings from .env or configuration files.

Many procedures include steps to connect to MongoDB Atlas or a local MongoDB
server. The tooling should be able to interpolate the connection string from
a .env file or other configuration file. For example:

```rst
.. step:: Connect to MongoDB Atlas

   Connect to your Atlas cluster using the following connection string:

   .. code-block:: bash

      mongosh "mongodb+srv://<username>:<password>@<cluster>.mongodb.net/admin"
```

The tooling should be able to interpolate the `<username>`, `<password>`, and
`<cluster>` variables from a .env file or other configuration file, and/or
replace the entire connection string with a variable such as `<connection-string>`.

The user's environment file might resemble:

```
CONNECTION_STRING=mongodb+srv://myAwesomeUsername:myAwesomePassword@Cluster0.mongodb.net
```

The tooling should replace the connection string in the procedure with the
value from the environment file, and correctly append the `/admin` database
name to the end of the connection string.

Additional environment variables may be required for other types of procedures,
such as API keys, credentials, or other configuration values. The tooling
should be able to interpolate these variables into the procedure as well. We
do not currently have a comprehensive list of all possible variables, but we
can start with a small set and expand as needed.

### Execution

#### Environments

The docs team should be able to automate procedure testing as part of the CI/CD
pipeline for the documentation. Writers should be able to run the tests
locally as well.

The tooling can require some setup to configure credentials and other
environment variables. For example, writers may need to provide a MongoDB
connection string or other credentials to test database interactions. Writers
may maintain a `.env` file to store credentials and other environment variables.

The tooling should attempt to verify any execution environments are available
before executing the procedure. For example, if a procedure contains a step to
execute Go code, the tooling should verify that Go is installed and available
for execution before attempting to test the procedure.

For the initial scope, assume the tooling must support execution environments
required to execute code in the following languages:

- JavaScript
- PHP
- Python
- Shell

If any required environments are not available, the tooling should report the
missing requirements and skip the test. We will not require writers to install
*all* environments to test all procedures; writers should only be required to
install environments relevant to the procedures they are testing.

CI/CD should be configured to run tests with all required environment
installed to ensure that no tests are skipped in CI.

#### Code Snippet Execution

The tooling should be able to execute code snippets written in various languages.
For each code snippet, it should be able to determine the language of the code
and execute it accordingly.

Not all code snippets are intended to be executed. For example, a code snippet
may represent output rather than input. The tooling should be able to
determine whether a code snippet is intended to be executed based on the
context of the step and the content of the code snippet.

For more details, refer to Appendix A: Code Block Types.

##### Executable examples

###### Shell commands

The tooling should be able to execute shell commands written for Unix. Windows
is out of scope of the initial Poc.

For example:

```sh
mkdir my-project
```

###### Code

The tooling should be able to execute code written in various languages. For
example:

```javascript
const assert = require('assert');
assert.equal(1, 1);
```

```python
assert 1 == 1
```

###### Combining snippets

A procedure may break down a usage example into a series of snippets. For
example, a procedure may show how to perform an aggregation query by breaking
down the query into its component stages. In this case, the tooling should
attempt to piece together the snippets into a complete usage example and test
the complete usage example.

For example, the tooling should be able to combine the following snippets:

```javascript
const MongoClient = require('mongodb').MongoClient;
```

```javascript
const uri = "<connection-string>";
```

```javascript
const client = new MongoClient(uri);
await client.connect();
```

```javascript
const db = await client.db('mydb');
```

```javascript
db.ping();
```

Into the following code to test:

```javascript
const assert = require('assert');
const MongoClient = require('mongodb').MongoClient;
const uri = process.ENV.CONNECTION_STRING;
const client = new MongoClient(uri);
await client.connect();
const db = client.db('mydb');
db.ping();
```

##### Non-executable examples

###### Abstract placeholder examples

Examples that demonstrate object shapes using field/value pairs that contain
field names and types rather than concrete values are not executable. For
example:

```json
{
  "name": "string",
  "age": "number",
  "isStudent": "boolean"
}
```

###### Output examples

Code snippets that represent output rather than input are not executable. For
example, the following snippet represents output and is not executable:

```javascript
{'plot': 'At the age of 21, Tim discovers he can travel in time and change what happens and has happened in his own life. His decision to make his world a better place by getting a girlfriend turns out not to be as easy as you might think.', 'title': 'About Time', 'score': 0.7710106372833252}
{'plot': 'A psychiatrist makes multiple trips through time to save a woman that was murdered by her brutal husband.', 'title': 'Retroactive', 'score': 0.760047972202301}
{'plot': 'A time-travel experiment in which a robot probe is sent from the year 2073 to the year 1973 goes terribly wrong thrusting one of the project scientists, a man named Nicholas Sinclair into a...', 'title': 'A.P.E.X.', 'score': 0.7576861381530762}
{'plot': 'An officer for a security agency that regulates time travel, must fend for his life against a shady politician who has a tie to his past.', 'title': 'Timecop', 'score': 0.7576561570167542}
{'plot': 'After visiting 2015, Marty McFly must repeat his visit to 1955 to prevent disastrous changes to 1985... without interfering with his first trip.', 'title': 'Back to the Future Part II', 'score': 0.7521393895149231}
{'plot': 'A reporter, learning of time travelers visiting 20th century disasters, tries to change the history they know by averting upcoming disasters.', 'title': 'Thrill Seekers', 'score': 0.7509932518005371}
{'plot': 'Lyle, a motorcycle champion is traveling the Mexican desert, when he find himself in the action radius of a time machine. So he find himself one century back in the past between rapists, ...', 'title': 'Timerider: The Adventure of Lyle Swann', 'score': 0.7502642869949341}
{'plot': 'Hoping to alter the events of the past, a 19th century inventor instead travels 800,000 years into the future, where he finds humankind divided into two warring races.', 'title': 'The Time Machine', 'score': 0.7502503395080566}
{'plot': 'A romantic drama about a Chicago librarian with a gene that causes him to involuntarily time travel, and the complications it creates for his marriage.', 'title': "The Time Traveler's Wife", 'score': 0.749496340751648}
{'plot': 'A modern aircraft carrier is thrown back in time to 1941 near Hawaii, just hours before the Japanese attack on Pearl Harbor.', 'title': 'The Final Countdown', 'score': 0.7469133734703064}
```

#### Test execution

This spec is for an initial spike on what this tooling may look like. Depending
on the outcome of the spike, we may choose to implement the tooling differently.

One option following existing patterns might be to execute the tooling within
a Node.js project that uses Jest to execute tests. Consider the following
pseudo-code for a set of procedures:

```js
const filepath = '/path/to/procedure.rst';
const procedures = parseProcedures(filepath);
procedures.forEach(procedure => {
    executeProcedure(procedure);
});
```

Each procedure should be tested as an independent test case, and report
success or failure independently. It should provide enough detail for
a documentation writer to determine which procedure failed and how to
debug the failure.

#### Environment Cleanup

The tooling should be able to clean up any environment changes made during
procedure execution. For example, if a procedure creates a new database, the
tooling should be able to drop the database after the procedure is complete.

Each procecure should run in isolation to avoid potential interference with
the execution.

If the tooling creates temporary files for testing, the tooling should clean up
after itself.

#### Dependency Management

If a procedure requires installing dependencies, we should check whether the
environment already has the required dependencies installed. If so, use the
existing dependencies.

If not, the tooling should install the dependencies in a temporary environment
for the duration of the test.

We do not want to damage the writer's local environment, and we do not want to
require writers to install dependencies that are not relevant to their work.
If a writer needs to install dependencies to test a procedure, the tooling
should provide clear instructions for the writer to install the dependencies
manually.

We may want to investigate the viability of running tests that require installing
dependencies in a virtual machine or a Docker container to avoid polluting the
writer's local environment. This is out of scope for the initial POC, as we
would prefer to keep the tooling as lightweight as possible.

### Ecosystem

The tooling should be written in TypeScript and use Node.js for execution.
Node.js is a required tool for writer setup and should be available on all
writer workstations.

The tooling *can* use libraries available on npm, but should keep external
dependencies to a minimum for maintainability.

### Conditional Execution

Some documentation pages contain branches that represent alternative ways to
execute a procedure. For example, a page may offer instructions for
connecting to MongoDB Atlas using either the MongoDB CLI or the Atlas
UI. The tooling should be able to test both branches of the procedure
independently.

Writers may use two reStructuredText directives to represent mutually-exclusive
branches:

- `.. tabs::`
- `.. composable-tutorial::`

For more details about the syntax of each of these directives, refer to
Appendix B: ReStructuredText Syntax.

#### Tabs

Writers may use tabs to represent mutually-exclusive paths within a procedure.
When possible, the tooling should attempt to test each tab within a tab set
as an alternate version of the procedure test. For example, if a procedure
contains a tab set with two tabs, the tooling should attempt to test the
procedure twice: once for each tab.

The tooling should not consider a tab set as a conditional branch within a
procedure. For example, if a procedure contains a tab set with two tabs, and
each tab contains two steps, the tooling should not consider that four steps.
Rather, it should consider that two steps, with alternate implementations
provided for each step via the tab set.

A tab may be presented within a larger procedure which may contain multiple
instances of tabs across different steps. In this case, the tooling should
consider each tab's `:tabid:` to be unique within the procedure. For example,
if a procedure contains two tab sets, each with two tabs, the tooling should
execute the procedure twice, selecting the same tab from each set:

```rst
.. procedure::

   .. step:: Step 1

      .. tabs::

         .. tab:: Tab 1
            :tabid: tab1

            Tab 1 content

         .. tab:: Tab 2
            :tabid: tab2

            Tab 2 content

   .. step:: Step 2

      .. tabs::

         .. tab:: Tab 1
            :tabid: tab1

            Tab 1 content

         .. tab:: Tab 2
            :tabid: tab2

            Tab 2 content
```

Execution should consist of two distinct tests which each execute both steps,
selecting the same tab during each step:

First test (:tabid: tab1):
- Step 1: tab1
- Step 2: tab1

Second test (:tabid: tab2):
- Step 1: tab2
- Step 2: tab2

#### Composable Tutorials

Writers may use the `.. composable-tutorial::` directive to represent different
ways to achieve the same outcome. For example, a page may offer instructions for
connecting to MongoDB Atlas using either the Atlas CLI or MongoDB Drivers.
The tooling should be able to test each combination of composable tutorial
selections as a unique test. For example, if a page offers a composable tutorial
with three distinct selections, the tooling should execute three different tests:

```rst
.. composable-tutorial::
   :options: interface, language

   .. selected-content::
      :selections: driver, nodejs

      Content for Node.js Driver.

   .. selected-content::
      :selections: driver, python

      Content for Python Driver.

   .. selected-content::
      :selections: atlas-cli, none

      Content for Atlas CLI.
```

The tooling should execute the following three tests:

First test (:selections: driver, nodejs):
- Content for Node.js Driver.

Second test (:selections: driver, python):
- Content for Python Driver.

Third test (:selections: atlas-cli, none):
- Content for Atlas CLI.

The tooling should not consider a composable tutorial as a conditional branch
within a procedure. For example, if a composable tutorial contains three
selections, and each selection contains two steps, the tooling should not
consider that six steps. Rather, it should consider that two steps, with
alternate implementations provided for each step via the composable tutorial.

If a composable tutorial is presented within a larger procedure which may
contain multiple instances of composable tutorials across different steps, the
tooling should consider each composable tutorial's `:selections:` to be unique
within the procedure. For example, if a procedure contains two composable
tutorials, each with three selections, the tooling should execute three different
tests:

```rst
.. composable-tutorial::
   :options: interface, language

   .. procedure::

      .. step:: Step 1

         .. selected-content::
             :selections: driver, nodejs

             Content for Node.js Driver.

         .. selected-content::
             :selections: driver, python

             Content for Python Driver.

         .. selected-content::
             :selections: atlas-cli, none

             Content for Atlas CLI.

      .. step:: Step 2

         .. selected-content::
             :selections: driver, nodejs

             Content for Node.js Driver.

         .. selected-content::
             :selections: driver, python

             Content for Python Driver.

         .. selected-content::
             :selections: atlas-cli, none

             Content for Atlas CLI.
```

This should result in the following three tests:

First test (:selections: driver, nodejs):
- Step 1: Content for Node.js Driver.
- Step 2: Content for Node.js Driver.

Second test (:selections: driver, python):
- Step 1: Content for Python Driver.
- Step 2: Content for Python Driver.

Third test (:selections: atlas-cli, none):
- Step 1: Content for Atlas CLI.
- Step 2: Content for Atlas CLI.

### Error Handling

The tooling should handle errors gracefully and provide as much information as
possible to help the writer debug the issue. For example, if a code snippet
fails to execute, the tooling should provide the following information:

- The procedure and step where the error occurred
- The code snippet that failed to execute
- The error message returned by the execution environment
- The line number of the code snippet where the error occurred (if applicable)

#### Parser errors

If the tooling encounters a parsing error while attempting to extract a
procedure from the reStructuredText, it should provide a clear error message
indicating the file and line number of the error.

#### Execution errors

### Test output

Test reports should contain minimal output, except when needed to debug a failure.
Consider a `--verbose` flag to enable verbose output with more details about
what was executed, and how.

If we implement something that uses Jest under the hood, we can take advantage
of its built-in reporting and assertion capabilities.

Output should be communicated in plain language for consumption by writers.

#### Parsing details

With verbose output enabled, the tooling should output details about the
parsed procedures, including:

- The number of procedures found
- The number of steps within each procedure
- A list containing each executable element within each step

For example:

```
Found 2 procedures
Procedure 1: 3 steps
   Step 1: 1 executable elements
   Step 2: 2 executable elements
   Step 3: 1 executable elements
Procedure 2: 2 steps
   Step 1: 1 executable elements
   Step 2: 1 executable elements
```

If the procedure is derived using composable tutorial selections or tab selections,
the tooling should output the selections or tab IDs used to derive the procedure.

For example:

```
Found 1 procedure
Procedure 1 (selections: driver, nodejs): 3 steps
   Step 1: 1 executable elements
   Step 2: 2 executable elements
   Step 3: 1 executable elements
```

Executable elements should contain the following information:

- The type of executable element (code, shell, ui, cli, api, url)
- The language of the executable element (if applicable)

For example:

```
Found 1 procedure
Procedure 1 (selections: driver, nodejs): 3 steps
   Step 1: 1 executable elements (code, javascript)
   Step 2: 2 executable elements (code, javascript) (code, javascript)
   Step 3: 1 executable elements (code, javascript)
```

#### Execution details

With verbose output enabled, the tooling should output details about the
result of each executable element within each step.

For example:

```
Procedure 1:
   Step 1:
      Successfully executed 1 executable elements (code, javascript)
   Step 2:
      Successfully executed 2 executable elements (code, javascript) (code, javascript)
   Step 3:
      Successfully executed 1 executable elements (code, javascript)
```

If an executable element fails, the tooling should output the error message
and stack trace. For example:

```
Procedure 1:
   Step 1:
      Successfully executed 1 executable elements (code, javascript)
   Step 2:
      Failed to execute 1 of 2 executable elements
      (code, javascript)
      Error: ReferenceError: x is not defined
         at file:///path/to/procedure.rst:10:5
   Step 3:
      Successfully executed 1 executable elements (code, javascript)
```

## Success Criteria

This procedure testing tooling should evaluate success on two axes:

- Each executable element within a procedure should execute successfully
- Each procedure should execute successfully

A procedure is considered successful if all of its executable elements execute
successfully. A procedure is considered unsuccessful if any of its executable
elements fail to execute.

The tooling should not consider a procedure test to have failed if the
procedure contains executable elements that are not tested. For example, if a
procedure contains a code block that is not tested because it is not executable,
the tooling should not consider the procedure to have failed.

For the purposes of this PoC, an executable element is considered to have
executed successfully if it does not throw an error. What this looks like in
practice depends on the type of executable element. For example, a code block
is considered to have executed successfully if the code does not throw an
error. A shell command is considered to have executed successfully if the
command returns a zero exit code. A URL is considered to have executed
successfully if it returns a <400 HTTP status code.

## Project Structure

The tooling will be used in the documentation monorepo. This will likely
reside within a `code-example-tests/procedures` directory at the root of the
repo. This will be outside the scope of the specific documentation projects,
and will be used to test procedures across the documentation.

This is one possible project structure, but the specifics will vary
based on what we decide to implement, and how we want writers to use it:

```
code-example-tests/
   procedures/
      tests/
         drivers.test.js
         atlas.test.js
      tooling/
         index.js
         procedure.js
         code-snippet.js
         shell-command.js
         ui-interaction.js
         cli-request.js
         api-request.js
         url.js
      package.json
      package-lock.json
content/
   drivers/
      source/
         index.txt
         includes/
            driver/authenticate.txt
   atlas/
      source/
         index.txt
         includes/
            atlas/create-cluster.txt
```

## Implementation Clarifications

This section captures additional clarifications and decisions made during the
requirements gathering process.

### Executable vs. Non-Executable Code Heuristics

**Context**: MongoDB documentation has been assembled over many years with many
inconsistencies. We are messy humans who have created a large corpus of
documentation with a wide range of content and many inconsistencies.

**Decision**: The tooling should **not** rely on specific heuristics for
detecting output or non-executable code, as these patterns are not consistently
applied across the documentation corpus.

**For `io-code-block` directives**: The tooling should **only** execute the
`input` directive and should **not** validate against the `output` directive
for the initial PoC. We have separate tooling for validating specific code
blocks that can handle output validation. The scope here is to evaluate whether
the procedure is executable as written.

### State Management Between Code Blocks

**Decision**: Variables from one code block **should** be available in
subsequent blocks within the same step.

**Open Question**: Whether state should persist across steps within a procedure
is unclear. We don't have a clear picture of how we use code within procedures,
so we'll have to figure this out through exploration of this space.

**Implementation Note**: Each language may need its own state accumulation
strategy. The tooling may need to create a wrapper/harness to execute
accumulated code.

### UI Interactions

**PoC Scope**: Skip UI interactions in the initial PoC, but plan/account for it
in the specification.

**Framework**: No strong preference for UI testing framework. Optimize for
well-maintained and lightweight options. Ideally, the tooling should abstract
away the details so technical writers don't have to care what we use.

**RST Representation**: A UI interaction in RST is either:
- Plain text instruction to select something from a menu or interact with an
  element
- A `:guilabel:` directive with a UI text label indicating the element that we
  should be interacting with or which should be present

### API Requests

**Scope**: The tooling should make actual API requests instead of mocking them,
to validate that procedures work as written.

**Authentication**: Most API calls will probably require authentication. If
there is an API procedure, we can ask writers to provide authentication info in
an `.env` file. The tooling should have the context to authenticate to any of
our APIs as a prerequisite for running an API procedure test, as many of our
API procedures will refer developers to somewhere else for details about
authenticating with the API, and will assume the person following the procedure
has authenticated.

**Validation**: The tooling should make HTTP requests and validate that they
succeed (appropriate status codes).

### CLI Tool Execution

**Shell Commands**: Atlas CLI and similar tools should be covered by shell
execution.

**mongosh Special Handling**: mongosh opens its own shell where we may need to
execute commands, so it may need special handling. Alternatively, we can use an
approach where we write mongosh commands to temp files and use the mongosh file
execution functionality to execute them, if we don't want to deal with
operating directly in a mongosh shell.

**Prerequisite Checks**: The tooling should verify CLI tools are installed
before attempting execution.

### URL Validation

**Scope**: The tooling should validate URLs in comments/strings **and** explicit
URL directives. If we show a URL, it should work, and if it doesn't, we should
fix it or remove it.

**Implementation**: Make HEAD requests to validate URLs. No strong opinion
regarding following redirects or timeouts for URL checks. Basically, we want to
make sure a user attempting to click a URL in our docs procedure can get to the
URL, or we should consider it a bad URL and remove it.

### Placeholder Detection and Replacement

**Reference**: See `reference-code/placeholder-consistency.md` for a
comprehensive analysis of placeholder inconsistencies across the documentation.

**Strategy**: The tooling should attempt to resolve common placeholders and
their variations. Common placeholder patterns include:
- Angle brackets: `<connection-string>`, `<username>`, `<password>`
- Multiple naming conventions: kebab-case, camelCase, space separated, snake_case
- Variations for the same concept (e.g., 25 different ways to represent
  "connection string")

**Failure Handling**: If a placeholder cannot be resolved, the test should
**fail** until writers correct the inconsistent placeholders. This is
acceptable - if writers are using inconsistent placeholders, we're okay with
failing the test until writers can correct them.

**Configuration**: For the initial implementation, writers should **not** be
required to provide a configuration file mapping placeholders to environment
variables. The tooling should use intelligent matching based on the placeholder
consistency analysis.

### Test Discovery and Execution

**Discovery Options**:
1. Maintain a separate list of files that should be tested (easiest at the
   beginning)
2. Add a new option to the reStructuredText `meta` directive to specify in page
   metadata that a page contains procedures we should test (would require
   scanning thousands of pages)

**Decision**: Use option 1 (separate list) for the initial implementation. If
we use the project structure proposed in this document with a separate
`procedures` testing directory, we can ask writers to manually specify the
procedure files we want to test.

**Philosophy**: Optimize for the technical writer audience and make it as easy
as possible for them to use the tooling. Reduce their learning curve as much as
possible to make it easier for them to get testing on procedures.

### Parallel vs. Sequential Execution

**Decision**: Tests should run **sequentially**, not in parallel.

**Rationale**: If we're executing tests against a single local environment,
like a local database, parallel execution could cause flaky or inconsistent
test behavior.

### Database/Resource Cleanup

**Automatic Cleanup**: The tooling should automatically detect and clean up
created databases/collections.

**Failure Handling**: The tooling should flag cleanup failures in the test
output so writers can manually perform cleanup steps.

**Blocking Consideration**: We may want to consider a cleanup failure blocking
further test execution, because it may create issues with tests that come
after. We may want to handle this situationally, because some types of cleanup
failures may not warrant completely failing the tests.

### Extract File Resolution

**Complexity**: The extract file resolution is complex, but it's
a reality of our docs.

**Requirements**:
- Recursively resolve nested `inherit` references
- Perform string replacement for `{{variable}}` placeholders
- There may be other extract file patterns beyond the YAML format shown in the
  examples, but we don't know what they are, so we'll have to handle them as we
  encounter them

### Composable Tutorial Dependencies

**Incomplete Information**: Composable tutorial selections do not tell the full
story on dependencies.

**Additional Dependencies**: Dependencies may also be communicated in:
- Prerequisites or requirements documented within the procedure
- Languages of code blocks independent of the composable tutorial selections
- For example, most of our Driver procedures probably also have shell commands
  in addition to the Driver programming language
- Something like the Java Driver may have additional dependencies, such as
  Maven, which aren't necessarily communicated in the composable tutorial
  selection

**Implication**: The tooling cannot rely solely on composable tutorial
dependencies to determine what needs to be installed or configured.

## Recommended PoC Scope

Based on the requirements and clarifications above, the following scope is
recommended for the proof-of-concept implementation:

### Phase 1: Core Parsing & Execution (Priority for PoC)

**RST Parsing**:
- Parse `procedure` and `step` directives
- Parse `code-block` and `literalinclude` directives
- Extract code content and metadata (language, options)
- Basic `include` directive resolution for transclusion

**Code Execution**:
- Execute code snippets in JavaScript, Python, PHP, and Shell
- Combine code snippets within a step (state persistence within step)
- Basic error capture and reporting

**Placeholder Interpolation**:
- Load environment variables from `.env` files
- Resolve common placeholders using the consistency analysis from
  `reference-code/placeholder-consistency.md`
- Support common patterns: `<connection-string>`, `<username>`, `<password>`,
  `<database>`, `<collection>`, etc.
- Handle variations (kebab-case, camelCase, space separated)

**Environment Detection**:
- Verify required tools are installed (Node.js, Python, PHP, mongosh, atlas-cli)
- Report missing dependencies before attempting execution
- Provide clear error messages for missing tools

**Error Reporting**:
- Report errors with file path and line numbers
- Include step number and description
- Show the failing code snippet
- Display error messages from execution
- Provide verbose mode for detailed output

### Phase 2: Conditional Content (Post-PoC)

**Tab Variants**:
- Parse `tabs` and `tab` directives
- Generate test variants based on `:tabid:` values
- Execute each variant independently

**Composable Tutorial Variants**:
- Parse `composable-tutorial` and `selected-content` directives
- Generate test combinations based on `:options:` and `:selections:`
- Execute each combination independently

**Advanced Transclusion**:
- Resolve nested `include` directives
- Handle `start-after` and `end-before` options
- Support `sharedinclude` directive

### Phase 3: Advanced Features (Defer to Post-PoC)

**Extract File Resolution**:
- Parse YAML extract files
- Resolve nested `inherit` references
- Perform `{{variable}}` placeholder replacement
- Handle extract file edge cases as encountered

**UI Testing**:
- Integrate headless browser (Playwright, Puppeteer, or similar)
- Detect `:guilabel:` directives
- Validate UI elements are present
- Abstract framework details from technical writers

**Advanced Output Validation**:
- Validate `io-code-block` output against expected results
- Support flexible matching (exact, regex, contains)
- Handle dynamic values (timestamps, IDs, etc.)

**Parallel Test Execution**:
- Run tests in parallel where safe
- Detect and avoid resource conflicts
- Provide configuration for parallelization level

**Sophisticated Cleanup**:
- Provide cleanup hooks/callbacks for writers to define
- Automatic detection of created resources
- Configurable cleanup failure handling (block vs. warn)

### Out of Scope for Initial Implementation

The following items are explicitly out of scope for the initial PoC:

- Windows support (Unix/macOS only)
- Advanced output validation for `io-code-block`
- UI testing with headless browsers
- Parallel test execution
- Complex extract file resolution
- Custom cleanup hooks
- Integration with existing test frameworks beyond basic reporting
- Performance optimization
- Caching of parsed RST files
- Incremental test execution

### Test Invocation

**Writer Experience**: The tooling should be as easy as possible for technical
writers to use. The learning curve should be minimal.

**Invocation Options** (to be determined during implementation):
- `npm test` or `npm run test-procedures`
- Custom CLI tool (e.g., `procedure-test <file>`)
- Integration with existing test runners

**Test File Specification**: Writers should manually specify which files to
test, either through:
- A configuration file listing test files
- Command-line arguments
- A dedicated test directory structure

**Avoid**: Requiring writers to scan thousands of pages or add metadata to
every page. Keep the barrier to entry low.

## Appendix A: Code Block Types

For the MongoDB documentation, we group code blocks into the following types:

- [Usage Examples](#usage-examples): Standalone code blocks that show how to
  perform a task, including the relevant setup and context.

- [Snippets](#snippets): Code that illustrates a specific concept or detail in
  the context of a larger example, tutorial, or reference page.

- [Sample Applications](#sample-applications): Runnable applications
  demonstrating broader use cases.

Not all code blocks are testable. For example, a code block that represents
output rather than input is not testable. A snippet that demonstrates an
object shape using field/value pairs that contain field names and types
rather than concrete values is not testable as part of this scope.

The tooling should be able to determine whether a code block is testable based
on the type of code block and the content of the code block.

If a given code block is not testable, the tooling should not *fail* the test
necessarily if the rest of the procedure can successfully execute. But it
should provide feedback that the code block was not tested.

### Usage Examples

Usage examples are self-contained, actionable code blocks that show how to
accomplish a specific task using MongoDB tools, drivers, or APIs. Usage
examples include enough information to understand, modify, and run the code
contained in the code block (for example, a single code block that contains
all declared variables and includes comments to indicate which placeholders to
update).

```csharp
using MongoDB.Driver;

// Replace the following with your MongoDB connection string
const string connectionUri = "mongodb://<db_username>:<db_password>@<hostname>:<port>/?connectTimeoutMS=10000";

var client = new MongoClient(connectionUri);
```

### Snippets

Snippets are narrowly scoped code blocks that help explain a specific concept
or detail. They are typically used as part of a broader explanation or tutorial,
and are often meaningful only within that context.

Snippets are intended to provide information. They aren't required to be valid
or runnable code. In some cases, snippets may contain intentionally incomplete
or incorrect code for demonstration purposes (for example, a snippet showing
all possible arguments for a command).

Snippets fall into one of the following subtypes:

- **Non-MongoDB command**: a command-line (CLI) command for any non-MongoDB
  tooling (for example, `mkdir`, `cd`, or `npm`), often used in the context of
  a tutorial.

  ```shell
  dotnet run MyCompany.RAG.csproj
  ```

- **Syntax example**: an example of the syntax or structure for an API method,
  an Atlas CLI command, a `mongosh` command, or other MongoDB tooling.

  ```text
  mongodb+srv://<db_username>:<db_password>@<clusterName>.<hostname>.mongodb.net
  ```

- **Example return object**: an example of an object, such as a JSON blob or
  sample document, returned after executing a corresponding piece of code.
  Commonly included as the output of an `io-code-block`.

  ```text
  A timeout occurred after 30000ms selecting a server using ...
  Client view of cluster state is
  {
      ClusterId : "1",
      State : "Disconnected",
      Servers :
      [{
        ServerId: "{ ClusterId : 1, EndPoint : "localhost:27017" }",
        EndPoint: "localhost:27017",
        State: "Disconnected"
      }]
  }
  ```

- **Example configuration object**: an example configuration object, often
  represented in YAML or JSON, enumerating parameters and their types.

  ```ini
  apiVersion: atlas.mongodb.com/v1
  kind: AtlasDeployment
  metadata:
  name: my-atlas-cluster
  spec:
  backupRef:
      name: atlas-default-backupschedule
      namespace: mongodb-atlas-system
  ```

In some cases, procedures may break down a usage example into a series of
snippets. For example, a procedure may show how to perform an aggregation
query by breaking down the query into its component stages. In this case,
the tooling should attempt to piece together the snippets into a complete
usage example and test the complete usage example.

### Sample Applications

Sample applications are complete, runnable programs that connect multiple
discrete pieces of code. Sample apps may include error handling, framework
integrations, or frontend UI elements.

Sample applications are not testable as part of this scope.

## Appendix B: ReStructuredText Syntax

The testing tooling need not support *all* reStructuredText syntax; only
the subset relevant to procedure testing in MongoDB Documentation. This includes
syntax related to:

- Headings
- Filepath parsing and transclusion
- Procedures
- Code snippets
- Tabs
- Composable tutorials
- URL links

### Headings

MongoDB documentation uses the following restructuredText heading styles:

```rst
==
H1
==

H2
--

H3
~~

H4
``

H5
++
```

H2 headings represent different sections of the documentation page, and
may correspond to distinct procedures.

For example, on a page about managing search indexes, there may be distinct
procedures for creating search indexes, viewing search indexes, updating search
indexes, and deleting search indexes. Each of these procedures would be
represented by an H2 heading.

```rst
Create a Search Index
---------------------

View Search Indexes
-------------------

Update Search Indexes
---------------------

Delete a Search Index
---------------------
```

Each heading may also contain distinct procedures nested within it. For example,
a section about creating search indexes may contain a procedure for creating a
search index using the Atlas UI, and a separate procedure for creating a search
index using the Atlas CLI.

When scaled across all of the interfaces a user may use to create a search
index, there may be many distinct procedures within a single H2 heading.

### Filepath parsing and transclusion

MongoDB documentation uses a few different reStructuredText directives to
parse and transclude code from files into documentation. The most common
directives related to transclusion are:

- `include`
- `literalinclude`
- `io-code-block`

#### Filepath resolution

Filepaths in MongoDB documentation are relative to the `source` directory of
the documentation repository. For example, if the documentation repository is
located at `/docs-mongodb-internal/content/drivers`, and the reStructuredText
file is located at `/docs-mongodb-internal/content/drivers/source/index.txt`,
then the filepath
`/docs-mongodb-internal/content/drivers/source/includes/driver/authenticate.rst`
would be written as `../includes/driver/authenticate.txt`.

From any given procedure file, trace back the directory structure until you find
the `source` directory for the documentation, and then resolve the filepath
from there.

#### `include`

The `include` directive is used to transclude text from a file into the
documentation.

```rst
.. include:: /path/to/file.rst
```

We typically use this directive to include text from a file in the `includes`
directory.

The ``include`` directive supports the following options that control what
content is included from the source file:

- `start-after`
- `end-before`

The file text between the `start-after` and `end-before` lines is rendered in
the documentation page.

```rst
.. include:: /path/to/file.rst
   :start-after: start-marker
   :end-before: end-marker
```

##### Special case: `extracts`

The MongoDB documentation contains a number of "extracts" yaml files that
contain snippets of text that are included in multiple places within the
documentation. For example, the `extracts-atlas-cli-commands.yaml` file
contains snippets for all Atlas CLI commands.

Extract files are referred to with a filepath that includes `/extracts/`, but
are typically at the root of the `includes` directory. For example, the
following line in a documentation source file references content from the
atlas-cli extracts file:

```rst
.. include:: /includes/extracts/atlas-clusters-connectionStrings-describe.rst
```

But the actual file that contains the referenced content is at:

```rst
/includes/extracts-atlas-cli-commands.yaml
```

The content within the extract file is formatted as a series of yaml documents,
each of which contains a `ref` key that is used to reference the content within
the extract file. For example, the include above refers to this content:

```yaml
ref: atlas-clusters-connectionStrings-describe
inherit:
  ref: atlas-cli-source-tabs
  file: extracts-atlas-cli-source-tabs.yaml
replacement:
  task: "return the SRV connection strings for your Atlas cluster"
  commandWithDashes: "atlas-clusters-connectionStrings-describe"
  commandWithoutDashes: "atlas clusters connectionStrings describe"
```

Which itself contains a reference to another extract file which has this content:

```yaml
ref: atlas-cli-source-tabs
content: |

  To {{task}} using the
  {+atlas-cli+}, run the following command:

  .. literalinclude:: /includes/command/{{commandWithDashes}}.rst
     :start-after: :caption: Command Syntax
     :end-before: .. Code end marker, please don't delete this comment
     :language: sh
     :dedent:

  To learn more about the command syntax and parameters, see the
  {+atlas-cli+} documentation for :atlascli:`{{commandWithoutDashes}}
  </command/{{commandWithDashes}}>`.

  {{optionalTutorialLine}}

replacement:
  task: ""
  commandWithDashes: "atlas-accessLists-create"
  commandWithoutDashes: "atlas accessLists create"
  optionalTutorialLine: ""
```

The tooling should be able to resolve the actual file path for an extract file
based on the reference in the `include` directive, and correctly resolve the
nested references within the extract file to produce the final content to be
tested.

The example above should resolve to the following content:

```rst
To return the SRV connection strings for your Atlas cluster using the
{+atlas-cli+}, run the following command:

.. literalinclude:: /includes/command/atlas-clusters-connectionStrings-describe.rst
   :start-after: :caption: Command Syntax
   :end-before: .. Code end marker, please don't delete this comment
   :language: sh
   :dedent:

To learn more about the command syntax and parameters, see the
{+atlas-cli+} documentation for :atlascli:`atlas clusters connectionStrings describe
</command/atlas-clusters-connectionStrings-describe>`.
```

#### `literalinclude`

The `literalinclude` directive is used to transclude code from a file into the
documentation.

```rst
.. literalinclude:: /path/to/file.txt
```

The `literalinclude` directive supports the following options that control what
content is included from the source file and how it's rendered:

- `language`: the language of the code block for syntax highlighting
- `start-after`: the line after which to start including content
- `end-before`: the line before which to stop including content
- `copyable`: whether the code block has a copy icon
- `caption`: the caption to display above the code block
- `dedent`: whether to remove leading whitespace from the code block
- `emphasize-lines`: lines to emphasize
- `lineno-start`: the line number to start with
- `linenos`: whether to show line numbers
- `category`: the category of the code block

```rst
.. literalinclude:: /path/to/file.txt
   :language: bash
   :start-after: start-marker
   :end-before: end-marker
   :copyable: true
   :caption: Caption
   :dedent: true
   :emphasize-lines: 1,2
   :lineno-start: 1
   :linenos: true
   :category: syntax example
```

For procedure testing, the relevant options are:

- `start-after`
- `end-before`
- `language`

The `language` option is used to determine how to execute the code block. For
example, a `language` of `bash` indicates that the code block contains shell
commands that should be executed. A `language` of `php` indicates that the code
block contains a PHP code example we should attempt to execute as PHP code.

#### `io-code-block`

The `io-code-block` directive is used to transclude code from a file into the
documentation. It is similar to the `literalinclude` directive, but is used to
pair code blocks with input and output.

It may contain inline code, or it may refer to a file using a filepath.

```rst
.. io-code-block::

   .. input:: /path/to/file.sh (optional)
      :language: bash
      :emphasize-lines: 1, 2
      :lineos:
      :category: syntax example

   .. output:: /path/to/file.txt (optional)
      :language: text
      :emphasize-lines: 3
      :lineos:
      :visible: false
```

The `io-code-block` directive supports the following options:

- `caption`: the caption to display above the code block
- `class`: the class to apply to the code block
- `source`: a URL to the source of the code block
- `input`: the input code block (required)
- `output`: the output code block

The `input` directive supports the following options:

- `language`: the language of the code block for syntax highlighting
- `emphasize-lines`: lines to emphasize
- `linenos`: whether to show line numbers
- `category`: the category of the code block

The `output` directive supports the following options:

- `language`: the language of the code block for syntax highlighting
- `emphasize-lines`: lines to emphasize
- `linenos`: whether to show line numbers
- `visible`: whether the output is visible by default in the documentation,
  rather than hidden behind a toggle

### Procedures

Procedures may be defined in one of two ways:

- Using an ordered list
- Using the `procedure` directive

#### Ordered List

An ordered list may be numbered:

```rst
1. Step 1
2. Step 2
3. Step 3
```

Use letters:

```rst
A. Step 1
B. Step 2
C. Step 3
```

Or mix numbers and letters for sub-steps:

```rst
1. Step 1
   a. Sub-step 1
   b. Sub-step 2
2. Step 2
```

#### Procedure Directive

The `procedure` directive is used to define a procedure. It is similar to an
ordered list, but allows for additional features such as steps to further clarify
the procedure.

```rst
.. procedure::

   .. step:: Step 1

   .. step:: Step 2

   .. step:: Step 3
```

Steps may contain sub-steps as ordered lists:

```rst
.. procedure::

   .. step:: Step 1

      a. Sub-step 1
      b. Sub-step 2
      c. Sub-step 3

   .. step:: Step 2

   .. step:: Step 3
```

### Code Snippets

Code snippets are defined using one of the following directives:

- `code`
- `code-block`
- `literalinclude`
- `io-code-block`

The tooling should attempt to determine the language of the code snippet
and execute it accordingly.

#### Language considerations

##### Alternative language derivation

The language value is not a required reStructuredText option - it's optional.
If the language is not specified, the tooling should attempt to derive the
language from the file extension of the `literalinclude` or `io-code-block`
directives. For example, a file with a `.js` extension should be considered
JavaScript.

If the file extension is not recognized, the tooling should consider the code
snippet `text`.

Refer to `reference-code/language-examples.go` for an example of how we map
supported languages and file extensions in a different application.

##### Sanitization

The language of the code snippet is a string, and writers may use different
values to represent the same language. The tooling should sanitize the language
to a common set of values for consistent execution. For example, `js` and
`javascript` should both be sanitized to `javascript`.

Some writers may use a value that cannot be sanitized. In this case, the
tooling should consider the code snippet `text`.

Refer to `reference-code/language-examples.go` for an example of how we map
supported languages and file extensions in a different application.

##### Bad values

Writers may use a value for syntax highlighting that does not represent the
executable language of the code snippet. Common cases may include:

- `bash` for a snippet that contains text or output, rather than shell commands
- `javascript` for a snippet that contains MongoDB BSON output
- `none` for a snippet that contains text or output

In these cases, the tooling should attempt to determine whether the code snippet
is executable. For example, if the code snippet contains shell commands, it
should be executed as shell commands regardless of the language specified for
syntax highlighting. If it is output, it should not be executed.

#### `code`

The `code` directive is used to define a code snippet.

```rst
.. code:: bash

   echo "Hello, world!"
```

The code directive takes an optional argument that represents the language of
the code block for syntax highlighting.

The `code` directive supports the following options:

- `copyable`: whether the code block has a copy icon
- `caption`: the caption to display above the code block
- `emphasize-lines`: lines to emphasize
- `class`: the class to apply to the code block
- `lineos`: whether to show line numbers
- `category`: the category of the code block
- `source`: a URL to the source of the code block

The tooling should execute the code snippet based on the language specified. For
example, a `language` of `bash` indicates that the code block contains shell
commands that should be executed. A `language` of `php` indicates that the code
block contains a PHP code example we should attempt to execute as PHP code.

#### `code-block`

The `code-block` directive is used to define a code snippet.

```rst
.. code-block:: bash

   echo "Hello, world!"
```

The code block directive takes an optional argument that represents the
language of the code block for syntax highlighting.

The `code-block` directive supports the following options:

- `language`: the language of the code block for syntax highlighting
- `caption`: the caption to display above the code block
- `copyable`: whether the code block has a copy icon
- `emphasize-lines`: lines to emphasize
- `class`: the class to apply to the code block
- `lineos`: whether to show line numbers
- `source`: a URL to the source of the code block
- `category`: the category of the code block

For procedure testing, the relevant option is:

- `language`

The tooling should execute the code snippet based on the language specified. For
example, a `language` of `bash` indicates that the code block contains shell
commands that should be executed. A `language` of `php` indicates that the code
block contains a PHP code example we should attempt to execute as PHP code.

#### `literalinclude`

See the [Filepath parsing and transclusion](#filepath-parsing-and-transclusion)
section for more information on the `literalinclude` directive.

### Tabs

The tooling should be able to parse tabs created with the `.. tabs::` directive.

The `.. tabs::` directive supports the following options:

- `hidden`: whether the tab titles and chrome are hidden by default in the
  documentation, rather than visible
- `tabset`: an optional pre-defined tab set with a specific set of tabs

Each tab within the tab set is defined using the `.. tab::` directive.

The `.. tab::` directive supports the following options:

- `tabid`: a unique identifier for the tab

It can take an optional argument that represents the string title of the tab.

```rst
.. tabs::

   .. tab:: Tab 1
      :tabid: tab1

      Tab 1 content

   .. tab:: Tab 2
      :tabid: tab2

      Tab 2 content
```

### Composable Tutorials

The tooling should be able to parse composable tutorials created with the
`.. composable-tutorial::` directive.

The `.. composable-tutorial::` directive supports the following options:

- `options`: a comma-separated list of option names that are used to define
  the composable tutorial selection fields
- `defaults`: a comma-separated list of default selections for each option
  in the `options` list. The order of the defaults must match the order of the
  options in the `options` list.

```rst
.. composable-tutorial::
   :options: interface, language
   :defaults: driver, nodejs
```

The content in the composable tutorial is defined using the
`.. selected-content::` directive.

The `.. selected-content::` directive supports the following options:

- `selections`: a comma-separated list of selections that are used to define
  the content to be displayed for the given composable tutorial selections

```rst
.. composable-tutorial::
   :options: interface, language
   :defaults: driver, nodejs

   .. selected-content::
      :selections: driver, nodejs

      Content for Node.js Driver.

   .. selected-content::
      :selections: driver, python

      Content for Python Driver.
```

Each unique set of selection options in the composable tutorial should be
tested independently. For example, the following composable tutorial:

```rst
.. composable-tutorial::
   :options: interface, language
   :defaults: driver, nodejs

   .. selected-content::
      :selections: driver, nodejs

      Content for Node.js Driver.

   .. selected-content::
      :selections: driver, python

      Content for Python Driver.
```

Should be tested twice: once for the `driver, nodejs` selection and once for
the `driver, python` selection.

A page that contains a composable tutorial may contain multiple blocks for a
given selected content combination, interleaved with content that is not
specific to a given selection. For example:

```rst
.. composable-tutorial::
   :options: interface, language
   :defaults: driver, nodejs

   .. selected-content::
      :selections: driver, nodejs

      Content for Node.js Driver.

   General content that applies to all selections.

   .. selected-content::
      :selections: driver, nodejs

      A second piece of content that is specific to the Node.js Driver.
```

The tooling should be able to interpolate non-specific content and content
that is specific to a given selection into one continuous block of content to
test. For example, the following composable tutorial:

```rst
.. composable-tutorial::
   :options: interface, language
   :defaults: driver, nodejs

   .. selected-content::
      :selections: driver, nodejs

      Content for Node.js Driver.

   .. selected-content::
      :selections: driver, python

      Content for Python Driver.

   General content that applies to all selections.

   .. selected-content::
      :selections: driver, nodejs

      A second piece of content that is specific to the Node.js Driver.

   .. selected-content::
      :selections: driver, python

      A second piece of content that is specific to the Python Driver.
```

Should be tested as if it were written as two separate procedures:

Node.js Driver:
```rst
Content for Node.js Driver.

General content that applies to all selections.

A second piece of content that is specific to the Node.js Driver.
```

And Python Driver:

```rst
Content for Python Driver.

General content that applies to all selections.

A second piece of content that is specific to the Python Driver.
```

### URL Links

The tooling should be able to parse URL links created with any of these
methods:

- External link syntax:
  ```rst
   `Link text <https://example.com>`__
  ```
- Using reStructuredText roles:
  ```rst
   :driver:`Link text </some-page>`
  ```
- As a source constant:
  ```rst
   {+link+}
  ```
- With inline source constants or substitutions:
  ```rst
  `Link text <http://mongodb.com/{+version+}/some/url>`__
  ```

#### Role resolution

The tooling should be able to resolve reStructuredText roles to their
corresponding URLs. For example, the `:driver:` role should be resolved to
the URL for the MongoDB driver documentation.

For development, the tooling should use the local `rstspec.toml` file to
resolve roles to their corresponding URLs.

In production, the tooling should refer to the production `rstspec.toml` file at
`https://raw.githubusercontent.com/mongodb/snooty-parser/refs/heads/main/snooty/rstspec.toml`
for up-to-date role resolution at execution time.

Example:

Given this role definition in `rstspec.toml`:

```toml
[role.driver]
type = {link = "https://www.mongodb.com/docs/drivers/%s", ensure_trailing_slash = true}
```

And this reStructuredText content:

```rst
:driver:`Node.js </node/current>`
```

The tooling should resolve the role to the URL
`https://www.mongodb.com/docs/drivers/node/current/`.

#### Source constants

URLs may contain source constants as part of the path, or may be defined
entirely as a source constant. For example:

```rst
{+link+}
```
```rst
http://mongodb.com/{+version+}/some/url
```

The tooling should be able to resolve source constants to their corresponding
values. For example, the `{+version+}` source constant should be resolved to
the version of the MongoDB server.

Source constants are defined in the `snooty.toml` file for a given project. For
example:

```toml
[constants]
version = "5.0"
```

The tooling should refer to the `snooty.toml` file for the project being tested
for up-to-date source constant resolution at execution time. The `snooty.toml`
file is located as a peer of the `source` directory for a given project.

For example:

```
--- content/
    --- project-name/
        --- snooty.toml
        --- source/
            --- index.txt
```

## Appendix C: Example Pages

This repository contains a `testdata` directory that contains example pages
for testing the tooling. The `testdata` directory contains the following
test pages:

- `atlas/source/atlas-search/manage-indexes.txt`: a page with a composable tutorial
- `atlas/source/connect-to-database-deployment.txt`: a page with a tab set
- `drivers/source/symfony.txt`: a page with a relatively simple procedure

These pages are paired with `snooty.toml` files that contain the source
constant definitions for their projects.

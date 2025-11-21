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

## Appendix B: ReStructuredText Syntax

The testing tooling need not support *all* reStructuredText syntax; only
the subset relevant to procedure testing in MongoDB Documentation. This includes
syntax related to:

- Filepath parsing and transclusion
- Procedures
- Code snippets
- Tabs
- Composable tutorials
- URL links

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

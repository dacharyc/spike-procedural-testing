# Training: Working Directory Clarity in Procedures

## Purpose

This document demonstrates best practices for writing procedures that involve directory navigation and file operations. Clear, explicit instructions about working directories help both human readers and automated testing tools understand exactly where commands should be executed.

## The Problem: Implicit Directory References

### ❌ Before (Implicit - Confusing)

```rst
.. step:: Initialize a Symfony Project

   Run the following command in your terminal to create a skeleton Symfony project:

   .. code-block:: bash

      composer create-project symfony/skeleton restaurants

.. step:: Install PHP Driver and Doctrine ODM

   Enter your project directory, then add the MongoDB PHP driver and the
   Doctrine ODM bundle to your application.

   Run the following commands to install the ODM:

   .. code-block:: bash

      composer require doctrine/mongodb-odm-bundle
```

**Problems**:
1. ❌ "Enter your project directory" - Which directory? The name isn't mentioned in this step
2. ❌ No explicit `cd` command - Readers must infer they need to change directories
3. ❌ Not testable - Automated tools can't execute implicit instructions
4. ❌ Error-prone - Readers might forget to change directories and get confusing errors

---

## The Solution: Explicit Directory Navigation

### ✅ After (Explicit - Clear)

```rst
.. step:: Create a Symfony Project

   a. Run the following command in your terminal to create a skeleton Symfony project called
      ``restaurants``:

      .. code-block:: bash

         composer create-project symfony/skeleton restaurants

   #. Change to the new project directory:

      .. code-block:: bash

         cd restaurants

.. step:: Install PHP Driver and Doctrine ODM

   Add the MongoDB PHP driver and the Doctrine ODM bundle to your application.
   The bundle integrates the ODM library into Symfony so that you can read from
   and write objects to MongoDB.

   Run the following command to install the ODM:

   .. code-block:: bash

      composer require doctrine/mongodb-odm-bundle
```

**Benefits**:
1. ✅ **Explicit** - The `cd restaurants` command is clearly shown
2. ✅ **Testable** - Automated tools can execute the exact commands
3. ✅ **Clear** - No ambiguity about which directory to use
4. ✅ **Consistent** - Matches patterns used in other MongoDB documentation

---

## Additional Best Practices

### 1. Name Directories Explicitly When Referencing Files

❌ **Before (Vague)**:
```rst
In the root directory, navigate to the ``.env`` file and define the
following environment variables:
```

✅ **After (Specific)**:
```rst
In the project root directory (``restaurants/``), replace the contents of
the ``.env`` file with the following code to define your connection string
and target database:
```

### 2. Remind Readers of Context When Needed

❌ **Before (Assumes Context)**:
```rst
Run the following command from the application root directory to start
your PHP built-in web server:
```

✅ **After (Reinforces Context)**:
```rst
From the project root directory (``restaurants/``), run the following
command to start your PHP built-in web server:
```

---

## Pattern Reference

### Pattern 1: Create Directory + Navigate

Use this pattern when a command creates a new directory that subsequent commands must run from:

```rst
.. step:: Create and initialize a project

   a. Run the following command to create a new directory called ``my-project``:

      .. code-block:: bash

         mkdir my-project

   #. Change to the new directory:

      .. code-block:: bash

         cd my-project

   #. Initialize the project:

      .. code-block:: bash

         npm init -y
```

### Pattern 2: Project Scaffolding Command + Navigate

Use this pattern when a scaffolding command (like `composer create-project`, `dotnet new`, `npm create`) creates a project directory:

```rst
.. step:: Create a new project

   a. Run the following command to create a project called ``my-app``:

      .. code-block:: bash

         composer create-project framework/skeleton my-app

   #. Change to the new project directory:

      .. code-block:: bash

         cd my-app
```

---

## Real-World Example

See the updated `testdata/drivers/source/symfony.txt` file for a complete example of these patterns in practice.

**Key Changes Made**:
1. Split "Initialize a Symfony Project" into sub-steps with explicit `cd` command
2. Changed "Enter your project directory" to explicit directory name
3. Updated "In the root directory, navigate to" to "In the project root directory (``restaurants/``)"
4. Updated "from the application root directory" to "From the project root directory (``restaurants/``)"

---

## Why This Matters for Testing

The `proctest` framework will:
1. Execute each shell command in sequence
2. Track working directory changes from `cd` commands
3. Execute subsequent commands in the correct directory
4. Fail with clear errors if files aren't found in expected locations

**Without explicit `cd` commands**, the testing framework would execute all commands from the initial working directory, causing failures that don't reflect real-world usage.

---

## Checklist for Writers

When writing procedures that involve directory navigation:

- [ ] Is there an explicit `cd` command when changing directories?
- [ ] Are directory names mentioned by name (not just "your project directory")?
- [ ] Do file operation instructions specify which directory the file is in?
- [ ] Would a reader who follows the steps exactly (without inferring anything) succeed?
- [ ] Could an automated tool execute these steps without additional context?

---

## Questions?

If you're unsure whether your procedure needs explicit directory navigation, ask:
1. "Does a command create a new directory?"
2. "Do subsequent commands need to run from inside that directory?"

If both answers are "yes", add an explicit `cd` command as a sub-step.


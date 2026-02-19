# Project Onboarding Guide: Go Backend

Welcome to the dis awesome small project! This guide is designed to help new developers quickly (or meself incase i forgot) get up to speed with the project's structure, tools, and development practices.

## Table of Contents

1.  [Introduction](#1-introduction)
2.  [Project Structure Overview](#2-project-structure-overview)
3.  [Development Dependencies](#3-development-dependencies)
4.  [Configuration](#4-configuration)
5.  [Frameworks and Template Engine](#5-frameworks-and-template-engine)
6.  [Mocking](#6-mocking)
7.  [Testing](#7-testing)
8.  [Pre-commit Checklist](#8-pre-commit-checklist)
9.  [Database Design](#9-database-design)
10. [Makefile](#10-makefile)
11. [VSCode Launch Configurations](#11-vscode-launch-configurations)

---

## 1. Introduction

This document serves as your primary resource for understanding the foundational aspects of this Go backend application. It covers essential setup, core technologies, development workflows, and best practices to ensure a smooth onboarding experience.

**Onboarding Tips:**

- **Start with the `main.go`:** Begin by exploring `cmd/app/main.go` to understand the application's entry point and initialization flow.
- **Follow Existing Patterns:** When adding new features or fixing bugs, always observe and adhere to the existing code patterns, naming conventions, and architectural decisions.

## 2. Project Structure Overview

This project adheres to a modular and layered architecture, heavily inspired by Clean Architecture principles. This approach emphasizes separation of concerns, making the codebase more maintainable, testable, and scalable.

The core idea is to organize code into distinct layers, where inner layers define interfaces that outer layers implement. This promotes **Dependency Inversion**, ensuring that business logic (inner layers) remains independent of external concerns like databases or web frameworks (outer layers).

Here's a breakdown of the key directories and their responsibilities:

- **`cmd/` (Commands/Applications)**
  - **Responsibility:** Contains the main entry points for different applications or services within the project. Each subdirectory here typically represents a distinct executable.
  - **What belongs here:** `main.go` files for various services (e.g., `cmd/app/main.go` for the main web application).
  - **What should NOT be here:** Business logic, data access code, or shared utilities. This layer should be thin and primarily responsible for wiring up dependencies and starting the application.
  - **Interaction:** This layer orchestrates the application by initializing configurations, setting up dependencies (handlers, services, repositories), and starting the server. It depends on `internal/config`, `internal/handler`, `internal/service`, and `internal/repository`.

- **`internal/` (Internal Application Code)**
  - **`internal/handler/` (HTTP/Presentation Layer)**
    - **Responsibility:** Handles incoming HTTP requests, parses request data, calls the appropriate service layer methods, and formats the response. This is the outermost layer of our application logic.
    - **What belongs here:** HTTP handlers, controllers, and presentation-specific logic (e.g., request/response DTO mapping).
    - **What should NOT be here:** Business logic, database queries, or direct interaction with external services.
    - **Interaction:** Depends on the `internal/service` layer (via interfaces) to execute business operations. It receives requests from the web framework (e.g., Echo) and sends responses back.

  - **`internal/service/` (Business Logic Layer)**
    - **Responsibility:** Encapsulates the application's core business rules and orchestrates operations. It defines _what_ the application does.
    - **What belongs here:** Business logic, use cases, and coordination of multiple repository operations. Service interfaces are defined here, and their implementations reside in this package or subpackages.
    - **What should NOT be here:** HTTP-specific details, database access logic, or direct interaction with external frameworks.
    - **Interaction:** Depends on the `internal/repository` layer (via interfaces) to perform data persistence operations. It is called by the `internal/handler` layer. This layer is designed to be framework-agnostic and highly testable.

  - **`internal/repository/` (Data Access Layer)**
    - **Responsibility:** Provides an abstraction over data persistence mechanisms (e.g., database, external APIs). It defines _how_ data is stored and retrieved.
    - **What belongs here:** Repository interfaces (defining data access contracts) and their concrete implementations (e.g., SQL queries, ORM calls, API client code).
    - **What should NOT be here:** Business logic or HTTP handling.
    - **Interaction:** Implements interfaces defined by the `internal/service` layer. It interacts directly with the database or external data sources.

  - **`internal/config/` (Configuration Management)**
    - **Responsibility:** Handles loading, parsing, and providing application configuration settings.
    - **What belongs here:** Configuration structures, functions for loading configuration from files (e.g., `config.yml`), environment variables, or other sources.
    - **What should NOT be here:** Business logic or data access.
    - **Interaction:** Provides configuration to other layers, typically initialized early in the application lifecycle (e.g., in `cmd/`).

  - **`internal/db/migrations/` (Database Migrations)**
    - **Responsibility:** Manages schema changes for the database. Each file represents a versioned change to the database structure.
    - **What belongs here:** SQL scripts or Go migration files (depending on the migration tool used) for creating, altering, or dropping database tables and columns.
    - **What should NOT be here:** Application logic or data.
    - **Interaction:** Executed by a migration tool (e.g., Goose) to evolve the database schema.

- **`docs/` (Documentation)**
  - **Responsibility:** Stores all project documentation, including architectural diagrams, design decisions, and onboarding guides.
  - **What belongs here:** Markdown files, diagrams (like `db_design.vuerd`), API specifications, etc.
  - **What should NOT be here:** Executable code.
  - **Interaction:** Primarily for human consumption, providing context and understanding of the project.

- **`pages/` (Templates)**
  - **Responsibility:** Contains the HTML templates used by the web application for rendering dynamic content.
  - **What belongs here:** `.jet` files that define the structure and presentation of web pages.
  - **What should NOT be here:** Go backend logic or static assets.
  - **Interaction:** Used by the `internal/handler` layer to render responses to the client.

- **`public/` (Static Assets)**
  - **Responsibility:** Serves static files directly to clients (e.g., web browsers).
  - **What belongs here:** CSS files, JavaScript files, images, fonts, and other client-side assets.
  - **What should NOT be here:** Server-side code or sensitive information.
  - **Interaction:** Directly served by the web server (e.g., Echo's static file serving capabilities).

## 3. Development Dependencies

Before you begin development, please ensure you have the following tools installed:

- **Mockery v3 (v3.2.5 or higher):**
  - **Purpose:** Used for automatically generating mock objects for interfaces, crucial for writing isolated and efficient unit tests.
  - **Installation:** Follow the instructions on the official Mockery GitHub page: [https://vektra.github.io/mockery/v3/](https://vektra.github.io/mockery/v3/)
- **YQ v4:**
  - **Purpose:** A lightweight and portable command-line YAML processor. It's used within the `Makefile` for parsing and manipulating `config.yml`.
  - **Installation:** Download from the official YQ releases page: [https://github.com/mikefarah/yq/releases](https://github.com/mikefarah/yq/releases)
- **Goose:**
  - **Purpose:** A database migration tool for Go. It manages schema changes for our SQLite database.
  - **Installation:** Refer to the Goose GitHub repository: [https://github.com/pressly/goose](https://github.com/pressly/goose)
- **GCC Compiler:**
  - **Purpose:** Required for compiling the `mattn/go-sqlite3` package, which is a CGo dependency.
  - **Installation:**
    - **Ubuntu/Debian:** `sudo apt install gcc`
    - **macOS:** Install Xcode Command Line Tools (`xcode-select --install`)
    - **Windows:** Install MinGW or MSYS2.

## 4. Configuration

The project uses a `config.yml` file for environment-specific settings.

- **Location:** The configuration file is expected to be at the root of the project directory.
- **Setup:** You **must** copy `config.yml.example` to `config.yml` before running the application:
  ```bash
  cp config.yml.example config.yml
  ```

## 5. Frameworks and Template Engine

This project leverages robust and widely-used Go libraries:

- **[Echo](https://github.com/labstack/echo) (HTTP Framework):**
- **[Jet](https://github.com/CloudyKit/jet) (Template Engine):**

## 6. Mocking

Mocking is essential for writing isolated and reliable unit tests.

- **Tool:** [Mockery](https://vektra.github.io/mockery/) (v3.2.5 or higher) for automatic mock generation.
- **Configuration:** Mockery's configuration is defined in the `.mockery.yml` file located at the project root.
- **Generating Mocks:** To generate or update mocks for interfaces in the project, run:
  ```bash
  mockery
  ```

## 7. Testing

Robust testing is a cornerstone of this project.

- **Framework:** We utilize [Testify](https://github.com/stretchr/testify) for its assertion package and mocking capabilities, which complements Go's built-in testing framework.
- **Best Practices:**
  - **Clean and Deterministic:** Tests should be independent, repeatable, and produce the same results every time. Avoid reliance on external state or timing.
  - **Fast:** Unit tests should run quickly. If a test is slow, consider if it's truly a unit test or if it should be an integration test.
  - **Coverage:** Aim for good test coverage, focusing on critical business logic and edge cases.
- **Table-Driven Tests:** Prefer table-driven tests for functions with multiple input/output scenarios. This approach improves readability, maintainability, and makes it easy to add new test cases.

## 8. Pre-commit Checklist

Before committing your changes, always perform the following checks:

- **Run All Tests:**
  ```bash
  go test ./...
  ```
  **Why it matters:** This command executes all tests in your project. It's crucial to ensure that your changes haven't introduced regressions or broken existing functionality. A green test suite provides confidence in your code.
- **Run Linter:**
  ```bash
  make lint
  ```
  **Why it matters:** Linting checks your code for stylistic errors, potential bugs, and adherence to coding standards. Running the linter pre-commit helps maintain code quality, consistency across the codebase, and catches issues early before they become harder to fix.

## 9. Database Design

Understanding the database schema is vital for backend development.
It's the single source of truth for our database structure.

**Location:** The database design is documented in `./docs/db_design.vuerd`.

## 10. Makefile

The `Makefile` at the project root provides a convenient way to automate common development tasks.

- **Purpose:** It orchestrates tasks such as:
  - Running database migrations
  - Executing linters (`make lint`)
  - And other utility functions.
- **Usage:** We encourage you to read the `Makefile` itself (`cat Makefile`) to discover all available commands and their specific usage. Each target typically includes comments explaining its function.

## 11. VSCode Launch Configurations

For VSCode users, the following `launch.json` configurations are provided for debugging:

```json
{
  "version": "0.2.0",
  "configurations": [
    {
      "name": "dev",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "program": "${workspaceFolder}/cmd/app/main.go",
      "output": "${workspaceFolder}/bin/debug_web",
      "cwd": "${workspaceFolder}/bin/",
      "args": ["--env=development", "--cfg=../config.yml"]
    },
    {
      "name": "prod",
      "type": "go",
      "request": "launch",
      "mode": "debug",
      "program": "${workspaceFolder}/cmd/app/main.go",
      "output": "${workspaceFolder}/bin/debug_web",
      "cwd": "${workspaceFolder}/bin/",
      "args": ["--env=production", "--cfg=../config.yml"]
    }
  ]
}
```

These configurations allow you to launch the application in either `development` or `production` mode directly from VSCode, with debugging capabilities. The `cwd` is set to `bin/` to ensure correct resolution of relative paths for `public/`, `pages/`, and `config.yml`.
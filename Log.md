# Update log

---

> _Note: Timezone is set to Asia/Bangkok_

#### Fix update test request path for generate handler

> 1/2/2569 17:42:21

#### Implement schemaground module for database schema comparison

> 1/2/2569 17:38:19

```
- Remove unused go.mod and go.sum files from etc/schemaground
- Add schemaground module with necessary domain, handler, service, and routes
- Implement database configuration loading and schema comparison logic
- Create API endpoint for comparing schemas between multiple databases
- Add tests for handler and service functionalities
- Include example .env file for database configuration
```

#### Restructure code generation and refundable date modules

> 1/2/2569 17:18:23

```
- Removed obsolete `compare-columns-database` module and its dependencies.
- Introduced `generate_code` module with handlers, services, and interfaces for generating various codes.
- Added `refundable_date` module to check refundable dates with appropriate handlers and services.
- Implemented unit tests for both `generate_code` and `refundable_date` modules to ensure functionality.
- Updated main routes to include new modules for code generation and refundable date checks.
```

#### Add VAT module tests and enhance README with testing instructions

> 31/1/2569 14:39:58

#### Refactor RSA and VAT modules: implement interfaces, enhance error handling, and update service methods

> 31/1/2569 14:03:48

```
- Refactor VAT handling to use Calculator interface and update service methods for error handling
- Refactor RSA module: update service interface, enhance ToHex method, and adjust handler initialization
```

#### Enhance VAT calculation logic, add RSA key handling, and update .gitignore

> 31/1/2569 12:40:08

```
- Implement VAT calculation logic and introduce Calculation struct
- Add RSA private key and update service to convert PEM to hex and write to file
```

#### Refactor response handling and update .gitignore for VSCode files

> 30/1/2569 18:00:02

#### Refactor application structure and enhance dependency injection

> 30/1/2569 16:21:22

```
- Introduced new `infra` package for dependency injection and logging.
- Updated server initialization to include environment configuration and logging.
- Modified route registration to utilize dependency injection.
- Refactored handlers and services to support dependency injection.
- Added logger middleware for HTTP request logging.
- Updated `go.mod` to include new dependencies.
```

</br>

---

_© 2026 Proundmhee. Released under the MIT License._

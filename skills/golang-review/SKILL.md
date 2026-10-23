---
name: golang-review
description: Use this skill when the user requests a Go code review, architecture review, package review, microservice review, performance analysis, concurrency analysis, security assessment, or general code quality evaluation of a Go project.
---
 
# Golang Code Review
 
## Purpose
 
Perform a comprehensive review of Go codebases and provide actionable feedback focused on:
 
- Correctness
- Maintainability
- Readability
- Performance
- Security
- Concurrency
- Testability
- Cloud-native best practices
 
The goal is to identify risks, technical debt, architectural weaknesses, and opportunities for improvement while preserving idiomatic Go principles.
 
---
 
# Review Process
 
## 1. Understand the Project
 
Identify:
 
- Project purpose
- Go version
- Application type
  - REST API
  - gRPC service
  - CLI application
  - Worker
  - Event processor
  - Microservice
  - Library
- Main dependencies
- External integrations
 
Review:
 
- README
- go.mod
- go.sum
- Project structure
- CI/CD configuration
 
Provide a brief project summary before starting the review.
 
---
 
## 2. Project Structure Review
 
Evaluate:
 
- Package organization
- Separation of concerns
- Module boundaries
- Dependency management
- Folder structure consistency
 
Identify:
 
- Circular dependencies
- Overly large packages
- Mixed responsibilities
- Tight coupling
 
Assess whether the structure aligns with Go community practices.
 
---
 
## 3. Go Idioms and Code Quality
 
Verify adherence to idiomatic Go principles.
 
Review:
 
- Package naming
- Function naming
- Interface naming
- Exported symbols
- Simplicity of implementation
 
Prefer:
 
```go
type Reader interface
```
 
over:
 
```go
type IReader interface
```
 
Identify:
 
- Unnecessary abstractions
- Over-engineering
- Excessive interfaces
- Repeated code
- Complex control flows
 
Encourage clear and maintainable code.
 
---
 
## 4. Error Handling
 
Review:
 
- Error propagation
- Error wrapping
- Error classification
- Error logging
 
Prefer:
 
```go
fmt.Errorf("failed to load configuration: %w", err)
```
 
Verify proper usage of:
 
```go
errors.Is(...)
errors.As(...)
```
 
Identify:
 
- Ignored errors
- Lost stack context
- Inconsistent handling
- Unnecessary panic calls
 
---
 
## 5. Context Usage
 
Verify proper use of:
 
```go
context.Context
```
 
Check:
 
- Context passed as first parameter
- Context propagation
- Cancellation support
- Timeout handling
 
Identify:
 
- Missing context propagation
- Long-running operations without timeout
- Context ignored by downstream operations
 
---
 
## 6. Concurrency Review
 
Analyze:
 
- Goroutines
- Channels
- Synchronization primitives
- Worker pools
 
Look for:
 
- Race conditions
- Deadlocks
- Goroutine leaks
- Channel misuse
- Resource contention
 
Review the use of:
 
- sync.Mutex
- sync.RWMutex
- sync.WaitGroup
- context cancellation
 
Highlight all concurrency risks and explain their impact.
 
---
 
## 7. Performance Analysis
 
Evaluate:
 
- Allocation patterns
- Memory usage
- CPU usage
- Algorithmic complexity
 
Look for:
 
- Unnecessary allocations
- Expensive conversions
- Reflection overuse
- Excessive copying
- Inefficient loops
 
Review:
 
- Slice management
- Map usage
- Buffer reuse
- Serialization overhead
 
Clearly distinguish between proven issues and potential optimizations.
 
---
 
## 8. Security Assessment
 
Review:
 
- Secret management
- Authentication
- Authorization
- Input validation
- Output sanitization
 
Identify risks related to:
 
- Injection attacks
- Path traversal
- Sensitive data exposure
- Weak configuration handling
- Insecure dependency usage
 
Assign severity levels to all findings.
 
---
 
## 9. Testing Review
 
Evaluate:
 
- Unit tests
- Integration tests
- End-to-end tests
 
Review:
 
- Coverage quality
- Edge cases
- Error cases
- Test maintainability
 
Prefer:
 
- Table-driven tests
- Independent tests
- Deterministic tests
 
Identify critical areas lacking adequate testing.
 
---
 
## 10. Observability
 
Review:
 
- Structured logging
- Metrics
- Distributed tracing
- Error reporting
 
Evaluate the usage of:
 
- slog
- OpenTelemetry
- Prometheus
 
Identify visibility gaps that could impact production support.
 
---
 
## 11. Cloud-Native Review
 
For services and microservices, evaluate:
 
- Graceful shutdown
- Readiness probes
- Liveness probes
- Configuration management
- Resilience patterns
- Retry strategies
- Backoff strategies
 
Review interaction with external systems and failure handling.
 
---
 
# Severity Levels
 
## Critical
 
Can cause:
 
- Security vulnerabilities
- Data loss
- Service outages
- Production incidents
 
## High
 
Likely to impact:
 
- Reliability
- Scalability
- Maintainability
 
## Medium
 
Recommended improvements that reduce technical debt.
 
## Low
 
Minor improvements or stylistic refinements.
 
---
 
# Output Format
 
## Executive Summary
 
Provide a concise summary of:
 
- Project purpose
- Overall quality
- Main concerns
 
---
 
## Positive Findings
 
List strengths observed in the codebase.
 
Example:
 
- Clear package separation
- Consistent error handling
- Good test coverage
 
---
 
## Findings
 
| Severity | Area | Issue | Recommendation |
|-----------|--------|--------|----------------|
| High | Concurrency | Potential goroutine leak | Add cancellation handling |
 
---
 
## Architecture Assessment
 
Describe:
 
- Architectural strengths
- Architectural weaknesses
- Scalability concerns
- Dependency concerns
 
---
 
## Refactoring Priorities
 
Rank the top improvements by business impact.
 
1. Highest priority
2. Medium priority
3. Nice to have
 
---
 
## Final Score
 
| Category | Score |
|-----------|---------|
| Architecture | x/10 |
| Go Idioms | x/10 |
| Error Handling | x/10 |
| Concurrency | x/10 |
| Security | x/10 |
| Testing | x/10 |
| Performance | x/10 |
| Maintainability | x/10 |
 
---
 
## Final Verdict
 
Select one:
 
- Production Ready
- Production Ready with Minor Improvements
- Requires Significant Refactoring
- High Technical Risk
 
Justify the decision with clear evidence from the review.
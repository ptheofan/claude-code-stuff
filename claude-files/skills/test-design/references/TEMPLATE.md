# [Feature Name] - Test Design

**Author:** [Name]  
**Date:** [YYYY-MM-DD]  
**Status:** Draft | In Review | Approved  
**Feature Spec:** `./docs/features/<NNN>-<feature-name>.feature.md`  
**TDD:** `./docs/features/<NNN>-<feature-name>.tdd.md`

---

## 1. Test Strategy Overview

| Layer | Count | Focus |
|-------|-------|-------|
| Unit Tests | | Smallest functions → services → handlers → controllers |
| Integration Tests | | Database ops, module interactions |
| E2E Tests | | Critical business paths only |

---

## 2. Unit Tests

### 2.1 Utilities / Helpers

| ID | Description | Given | When | Then | Priority |
|----|-------------|-------|------|------|----------|
| UT-001 | | | | | |

### 2.2 Services

| ID | Description | Given | When | Then | Mocks | Priority |
|----|-------------|-------|------|------|-------|----------|
| UT-010 | | | | | | |

### 2.3 Command / Query Handlers

| ID | Description | Given | When | Then | Mocks | Priority |
|----|-------------|-------|------|------|-------|----------|
| UT-020 | | | | | | |

### 2.4 Controllers / Resolvers (Input/Output Only)

| ID | Description | Input | Expected Output | Expected Status | Priority |
|----|-------------|-------|-----------------|-----------------|----------|
| UT-030 | Valid input | | | 200/201 | |
| UT-031 | Invalid input | | | 400 | |

---

## 3. Integration Tests

| ID | Description | Given | When | Then | Real Dependencies | Priority |
|----|-------------|-------|------|------|-------------------|----------|
| IT-001 | | | | | Database | |

---

## 4. E2E Tests

### 4.1 Critical Happy Paths

| ID | Description | Steps | Expected Outcome | Priority |
|----|-------------|-------|------------------|----------|
| E2E-001 | | 1. <br> 2. <br> 3. | | High |

### 4.2 Critical Failure Paths

| ID | Description | Steps | Expected Outcome | Priority |
|----|-------------|-------|------------------|----------|
| E2E-010 | | | | High |

---

## 5. Domain Error Tests

| Error Class | Trigger Condition | Test ID |
|-------------|-------------------|---------|
| | | |

---

## 6. Edge Cases

| ID | Scenario | Type | Expected Behavior | Priority |
|----|----------|------|-------------------|----------|
| EC-001 | | Unit | | |

---

## 7. Test Data Requirements

| Data | Purpose | Setup Method |
|------|---------|--------------|
| | | Fixture / Factory / Seed |

---

## 8. Out of Scope (Covered Elsewhere)

- 

---

## Sign-off

- [ ] Test coverage is meaningful
- [ ] Critical paths covered
- [ ] Ready for Implementation (`/coder`)

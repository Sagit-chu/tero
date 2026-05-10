# Flvx Node Monitor Testing Strategy - Design Specification

## 1. Overview
This document outlines the testing strategy for the Flvx Node Monitor application. It focuses on expanding the unit test coverage for the existing foundational components (Go backend and React frontend) and establishing an End-to-End (E2E) testing pipeline using Playwright.

## 2. End-to-End (E2E) Testing
- **Framework:** Playwright (TypeScript/Node.js).
- **Location:** A new `e2e/` directory at the project root.
- **Infrastructure:**
  - Playwright will manage the lifecycle of the application during tests via its `webServer` configuration.
  - The Go backend will be compiled and run on a designated test port (e.g., `8080`).
  - The Vite React frontend will be run in preview or dev mode on a designated port (e.g., `5173`), proxying API requests to the Go backend.
- **Core Test Cases:**
  1. **Dashboard Load:** Navigate to the frontend URL. Verify the page title is "Flvx Monitor Dashboard".
  2. **Status Fetch Flow:** Verify the initial "Loading..." state is visible, followed by the final "System Status: ok" state after the backend API responds.

## 3. Expanded Unit Testing

### 3.1 Frontend (React + Vitest)
The frontend currently has a basic component test. This will be expanded to thoroughly test the `App.tsx` component states by mocking the native `fetch` API.
- **Loading State:** Assert that the text "Loading..." is rendered immediately upon component mount.
- **Success State:** Mock a successful `fetch` response (`{"status": "ok"}`). Wait for the DOM to update and assert that "System Status: ok" is displayed.
- **Error State:** Mock a failed `fetch` response (e.g., Network Error or 500 status). Wait for the DOM to update and assert that "System Status: Error" is displayed.

### 3.2 Backend (Go)
The backend currently has standard TDD coverage for the `nodes` repository and DB setup. This will be expanded to ensure robustness.
- **Config Repository CRUD:** 
  - Create a new `backend/repository/config.go` to handle the `config` table created during DB init.
  - Write tests for `SetConfig(key, value)` and `GetConfig(key)` methods.
- **API Server Edge Cases:**
  - Update `backend/api/server_test.go` to use table-driven testing.
  - Add a test case ensuring that a `POST /api/status` request returns a `405 Method Not Allowed` or handles it appropriately, enforcing correct HTTP methods.
  - Verify the JSON payload body matches exactly `{"status": "ok"}`.

## 4. Execution Flow
Once implemented, the CI/CD pipeline or local developer workflow will execute tests in the following order:
1. `go test ./backend/...` (Backend Unit Tests)
2. `npm run test` inside `frontend/` (Frontend Unit Tests)
3. `npx playwright test` inside `e2e/` (Full System E2E Tests)

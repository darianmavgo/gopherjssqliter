# GopherJS SQLite AG Grid Viewer Plan

## Goal
Create a web interface that displays the contents of a local SQLite database using AG Grid. 
The system uses GopherJS to handle data fetching and communication between the backend and the AG Grid component.
Key constraints: Finite memory usage (streaming/paging), immediate display of initial 50 rows.

## Architecture

### 1. Backend (Go)
- **Responsibility**: access SQLite database, serve static files, provide JSON API.
- **Components**:
    - `net/http` server.
    - `database/sql` with `modernc.org/sqlite` (pure Go) or `github.com/mattn/go-sqlite3` (requires CGO). *Plan: Use modernc.org/sqlite for ease of build if possible, else standard.*
    - **Endpoint**: `GET /api/rows?table={mytable}&start={start}&end={end}`.
- **Memory Strategy**:
    - Never read the full table into memory.
    - Use `LIMIT` and `OFFSET` in SQL queries based on client request.

### 2. Frontend Logic (GopherJS)
- **Responsibility**: Bridge between the browser (React/AG Grid) and the Server.
- **Components**:
    - Compiled to `client.js`.
    - Exposes a global API for the React App: `window.Backend.getRows(params, successCallback)`.
    - Handles HTTP requests to the Go backend.

### 3. Frontend UI (React + AG Grid)
- **Responsibility**: Render the Grid.
- **Components**:
    - `index.html`: Loads React, AG Grid (Community), and `client.js`.
    - `App.js` (or inline babel script): Configures AG Grid.
    - **Datasource**: Implements AG Grid `IDatasource` (Infinite Scroll Model) which calls `window.Backend.getRows`.

## Implementation Steps

1.  **Project Setup**
    - Initialize Go module.
    - Create directory structure: `server/`, `client/`, `public/`.

2.  **Backend Implementation**
    - Create `server/main.go`.
    - Implement `main` to open SQLite DB.
    - Implement HTTP handler for `/api/rows`.
    - Serve `public/` directory.

3.  **GopherJS Client Implementation**
    - Create `client/main.go`.
    - Use `gopherjs/gopherjs/js` to export functions to `window`.
    - Implement `FetchData` function that calls `fetch('/api/rows...')`.

4.  **Frontend Integration**
    - Create `public/index.html` with React and AG Grid CDN links.
    - Implement the AG Grid setup, connecting the `serverSideDatasource` (or infinite) to the GopherJS function.

5.  **Build & Run**
    - Ensure `gopherjs` is installed (`go install github.com/gopherjs/gopherjs@latest`).
    - Build GopherJS: `gopherjs build -o public/client.js ./client`
    - Run Server: `go run ./server`

## Verification
- **Manual Test**:
    - Start server.
    - Open `http://localhost:8080`.
    - Verify Grid loads first 50 rows immediately.
    - Scroll down to verify dynamic loading of subsequent rows.
    - Check server logs to confirm "finite" queries (LIMIT 50).

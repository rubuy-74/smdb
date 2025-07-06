# smDB - Simple Memory Database

<p align="center">
  <img src="https://img.shields.io/badge/build-passing-brightgreen" alt="Build Status" />
  <img src="https://img.shields.io/badge/license-MIT-blue.svg" alt="License" />
  <img src="https://img.shields.io/badge/go-1.20%2B-blue" alt="Go Version" />
</p>


<p align="center">
  <b>smDB</b> is a minimal in-memory key-value store and database engine written in Go, designed for learning and experimentation with database internals.
</p>
---

## ℹ️ Project Status & Development Story

smDB began as an ambitious attempt to build a full-featured database from scratch in Go, including components like B+Trees, paging, and SQL-like parsing. The project served as a deep dive into database internals and low-level data structures. Over time, as priorities shifted and available free time became limited, development slowed and eventually paused. As a result, only the basic in-memory key-value store is currently functional, and many advanced features remain as work-in-progress or design sketches.

This repository is left open as a resource for others interested in database implementation, Go programming, or systems design. Contributions and forks are welcome!

---

## 🚀 Features

- 🗝️ In-memory key-value store
- 🌳 B+Tree-based storage engine (WIP)
- 📝 Simple SQL-like statement parsing (insert/select)
- 🧪 Easy to extend and experiment

> **Note:** This project is unfinished. Only a basic key-value store is implemented. B+Tree and other database features are under development. smDB is currently paused, but contributions and forks are welcome!

---

## 🛠 Tech Stack

- [Go](https://golang.org/) (1.20+)

---

## 🏗️ Getting Started

### Prerequisites

- Go 1.20 or newer

### ▶️ Running smDB

1. **Clone the repository:**
   ```bash
   git clone https://github.com/rubuy-74/smDB.git
   cd smDB
   ```

2. **Run the database server:**
   ```bash
   go run ./cmd/main.go
   ```

3. **Interact with the database:**
   - Use `.exit` to quit
   - Use `insert <id> <username> <email>` to insert a row (WIP)
   - Use `select` to list all rows (WIP)

---

## 📁 Project Structure

```text
├── cmd/
│   └── main.go              # Application entry point
├── internal/
│   ├── b3/                  # B+Tree implementation (WIP)
│   ├── handlers/            # Statement handlers
│   ├── models/              # Data models (Row, Page, Statement)
│   └── utils/               # Utility functions
├── go.mod                   # Go module definition
```

---

## 📜 License

This project is licensed under the MIT License.

---

<p align="center">
  <sub>Made with ❤️ by <a href="https://github.com/rubuy-74">rubuy-74</a></sub>
</p>

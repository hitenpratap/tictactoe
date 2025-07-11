# Go Tic-Tac-Toe (Terminal Edition)

[![Go CI](https://github.com/hitenpratap/tictactoe/actions/workflows/ci.yml/badge.svg)](https://github.com/hitenpratap/tictactoe/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/hitenpratap/tictactoe/branch/master/graph/badge.svg)](https://codecov.io/gh/hitenpratap/tictactoe)

A classic Tic-Tac-Toe game built to run in your terminal, developed using Go and the charming Bubble Tea TUI framework. This project is fully containerized with Docker for easy setup and deployment.

![Screenshot of Tic-Tac-Toe Game](/assets/img/demo.jpeg)

---

## Features

- **Interactive TUI:** Clean and responsive terminal user interface.
- **Vim Keybindings:** Move the cursor with arrow keys or `h/j/k/l`.
- **Player Turns:** Alternates between Player 'X' and Player 'O'.
- **Win/Draw Detection:** Automatically detects and announces a win or a draw.
- **Cross-Platform:** Runs anywhere Go can run.
- **Containerized:** Includes `Dockerfile` and `docker-compose.yaml` for a hassle-free setup.
- **Managed with Make:** A `Makefile` provides simple commands for building, running, testing, and cleaning the project.

---

## Getting Started

Follow these instructions to get a copy of the project up and running on your local machine.

### Prerequisites

You need one of the following setups:

1.  **Go Environment (Local Development)**
    * [Go](https://golang.org/doc/install) (version 1.24 or newer)

2.  **Docker (Containerized Development)**
    * [Docker](https://docs.docker.com/get-docker/)
    * [Docker Compose](https://docs.docker.com/compose/install/)

### Installation

1.  **Clone the repository:**
    ```sh
    git clone <your-repository-url>
    cd <repository-directory>
    ```

2.  **Download Go modules (for local development):**
    ```sh
    go mod tidy
    ```

---

## Usage

You can run the game directly on your machine or within a Docker container. The included `Makefile` simplifies all common tasks.

### Using the Makefile (Recommended)

The `Makefile` provides a convenient interface for all operations.

* **Run the game in a Docker container (easiest method):**
    ```sh
    make docker-run
    ```

* **Build the application binary locally:**
    ```sh
    make build
    ```

* **Run the application locally (after building):**
    ```sh
    make run
    ```

* **Clean up Docker resources (container, image, etc.):**
    ```sh
    make docker-clean
    ```

* **See all available commands:**
    ```sh
    make help
    ```

### Manual Commands (Without Make)

#### Running with Docker

1.  **Build and run the container using Docker Compose:**
    ```sh
    docker-compose up --build
    ```

2.  **Stop and remove the container:**
    ```sh
    docker-compose down
    ```

#### Running Locally

1.  **Run the application directly:**
    ```sh
    go run .
    ```

---

## How to Play

* **Move Cursor:** Use the **arrow keys** or **h, j, k, l** keys.
* **Place Marker:** Press **Enter** or **Spacebar**.
* **Reset Game:** Press **r**.
* **Quit:** Press **q** or **Ctrl+C**.

---

## Testing

To run the suite of unit tests for the game logic:

```sh
make test
```

To format the code using `go fmt`:

```sh
make fmt
```

Alternatively, you can run the tests manually with the Go tool:

```sh
go test -v ./...
```

---

## Built With

* [**Go**](https://golang.org/) - The core programming language.
* [**Bubble Tea**](https://github.com/charmbracelet/bubbletea) - A powerful TUI (Terminal User Interface) framework.
* [**Lipgloss**](https://github.com/charmbracelet/lipgloss) - A library for fancy, styled terminal output.
* [**Docker**](https://www.docker.com/) - For containerization and consistent environments.

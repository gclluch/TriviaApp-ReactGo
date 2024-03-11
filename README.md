# Trivia Application

## Overview

This Trivia Application is a dynamic, interactive platform designed to engage users with a wide range of trivia questions across various topics. Whether you're looking to test your knowledge in science, history, entertainment, or any other field, our application offers a rich, user-friendly experience that can be customized to suit any trivia enthusiast's preferences.

## Features

- **Dynamic Question Pool**: Access a vast database of trivia questions across multiple categories.
- **Single and Multiplayer Modes**: Enjoy solo challenges or compete against friends and other players in exciting multiplayer sessions.

## Technology Stack

- **Frontend**: React.js for a responsive and interactive UI.
- **Backend**: : Utilizes Go (Golang) with the Gin framework to handle API requests efficiently, providing a robust backend service.
- **Real-Time Communication**: WebSocket is implemented for real-time updates in multiplayer sessions, ensuring a seamless user experience.
- **Containerization**: Docker for easy application deployment and scaling.


## Getting Started

### Prerequisites

- Docker and Docker Compose
- Go environment set up for local development of the backend services.

### Installation

1. **Clone the Repository**

```bash
git clone [repository-url]
cd trivia-application
```

2. **Build and run with Docker**

```bash
docker-compose up --build
```
This command compiles the Go backend, bundles the React frontend, and starts the services defined in `docker-compose.yml`.

## Accessing the Application

After the Docker containers are up:

- The frontend can be accessed at `http://localhost:3000`.
- The backend API is available at `http://localhost:8080`.


## License ##
This project is licensed under the MIT License - see the LICENSE.md file for details.
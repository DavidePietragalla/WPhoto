# WasaPhoto 📸

WasaPhoto is a social media platform focused on photo sharing. It allows users to create accounts, upload photos, follow other users, like and comment on posts, and manage their network. The platform provides API endpoints for user authentication, post management, user interaction (following, liking, commenting), and moderation (banning). It uses a SQLite database to store user data, posts, and relationships. The API is documented using OpenAPI, enabling easy integration with frontend applications and other services.

## 🚀 Key Features

- **User Authentication:** Secure user login and authentication using Bearer tokens.
- **User Profiles:** Create and manage user profiles with nicknames and profile information.
- **Photo Posting:** Upload and share photos with other users.
- **Following/Followers:** Follow other users to see their posts in your home feed.
- **Liking:** Like posts to show appreciation.
- **Commenting:** Comment on posts to engage in discussions.
- **Banning:** Moderate content by banning users.
- **Search:** Search for users by nickname.
- **API Documentation:** Comprehensive API documentation using OpenAPI (doc/api.yaml).
- **Home Feed:** Get a personalized home feed of posts from users you follow.

## 🛠️ Tech Stack

- **Backend:** Go
- **Database:** SQLite3
- **API Router:** `github.com/julienschmidt/httprouter`
- **Configuration:** `github.com/ardanlabs/conf`
- **UUID Generation:** `github.com/gofrs/uuid`
- **HTTP Handlers:** `github.com/gorilla/handlers` (for logging/CORS)
- **Logging:** `github.com/sirupsen/logrus`
- **YAML Parsing:** `gopkg.in/yaml.v2`
- **API Specification:** OpenAPI 3.0.3

## 📦 Getting Started

### Prerequisites

- Go 1.17 or higher
- SQLite3

### Installation

1.  Clone the repository:

    ```bash
    git clone <repository_url>
    cd WasaPhoto
    ```

2.  Download dependencies:

    ```bash
    go mod download
    ```

### Running Locally

1.  Build the application:

    ```bash
    go build -o webapi cmd/webapi/main.go
    ```

2.  Run the application:

    ```bash
    ./webapi
    ```

    *Note: You might need to configure the database connection and other settings via environment variables or a configuration file.*

## 📂 Project Structure

```
WasaPhoto/
├── cmd/
│   └── webapi/             # Main application entry point
├── doc/
│   └── api.yaml            # OpenAPI specification
├── service/
│   ├── api/                # API handler logic
│   │   ├── api.go          # Core API router and configuration
│   │   ├── ban.go          # Ban user endpoints
│   │   ├── followers.go    # Follower endpoints
│   │   ├── search.go       # Search endpoint
│   │   ├── struct.go       # API data structures
│   │   ├── utils.go        # Utility functions (validation, auth)
│   │   └── reqcontext/     # Request context package (assumed)
│   ├── database/           # Database interaction layer
│   │   ├── ban-db.go       # Ban-related database operations
│   │   ├── comment-db.go   # Comment-related database operations
│   │   ├── database.go     # Database interface definition
│   │   ├── follow-db.go      # Follow-related database operations
│   │   ├── like-db.go      # Like-related database operations
│   │   ├── post-db.go      # Post-related database operations
│   │   └── user-db.go      # User-related database operations
├── go.mod                  # Go module definition
├── go.sum                  # Go module checksums
└── ...
```

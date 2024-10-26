# zkillanalytics

Provides basic analytics for zkillboard data - based on a list of corporations, characters, or alliances saved in the application config

## Usage

go run .

Access on localhost:8080, available routes are

- / - all available charts for the month
- /top/mtd  - kills by character in the month
- /top/ytd  - kills by character in the year
- /ourships/mtd - ships used by characters for kills in the month
- /ourships/ytd - ships used by characters for kills in the year
- /victims/mtd - victims by corporation in the month
- /victims/ytd - victims by corporation in the year

## Todo

- [ ] add weapon types
- [ ] add damage done
- [ ] add solo kills
- [ ] add top systems
- [ ] add CI/CD
- [ ] add tests


cmd/: Contains the application's entry point (main.go). Each subdirectory here can represent different executable commands if your project expands.

internal/: Houses the core application code. The internal directory restricts the visibility of these packages to your module, preventing external usage.

api/: Contains subpackages for each external API your application interacts with.

esi/: Manages all interactions with the ESI (EVE Swagger Interface) API.

zkill/: Manages all interactions with the ZKillboard API.

config/: Handles configuration loading and management (e.g., environment variables, config files).

model/: Defines the data models and structures used across the application.

persist/: Manages data persistence logic, such as database interactions.

repository/: Implements the Repository pattern, providing an abstraction layer over data sources (databases, APIs, etc.).

service/: Contains business logic that orchestrates between repositories and APIs.

utils/: Utility functions and helpers used across the application.

migrations/: Holds database migration files if you're using a relational database.

pkg/: For shared libraries or packages that could be used by external applications (optional based on project needs).
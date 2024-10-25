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
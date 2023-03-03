An app for fast domain whois checks (5-10k domains were checked on my desktop per hour)
It only displays if domain is available or not

## How to use:
1. Obtain an api key from apilayer.com
2. Run a DB (dockerized postgre for example)
3. Set envs:
    "DATABASE_DSN":"postgresql://log:pass@host:port/db",
    "API_KEY":"key"
4. Cron the checker OR fill the DB with domains and run once via go run /cmd/app/main.go

DB schema is being created automatically after the 1st run.

Links can be passed in any format: clean (google.com) or raw(https://www.google.com/search?client=safari&rls=en&q=never+gonna+give+you+up&ie=UTF-8&oe=UTF-8)

10 goroutines were relatively stable, but you can play with the exact amount to reach max speed

GLHF
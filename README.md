# Backend for Equitalytics.com
This repository houses all the backend components of [https://www.equitalytics.com](Equitalytics.com). 

The backend is entirely built in go, and hosted in an Azure Container Instance. The docker container's base image is ubuntu server. The data is stored in a PostgreSQL database server also hosted and maintained on Azure.

There are three entrypoints to the backend application, all held within the `cmd` directory. 

`cmd/main.go` is the entrypoint for the API server. The API serves a dynamic route `/info/{ticker}`, where company data for the specificied `{ticker}` is gathered and returned upon request. Requests are handled within their respective handler files, and then database retrieval operations take place at the `info_service` layer. The API server was built mostly using the Go standard library. Chi Router was used just for routing and a basic middleware stack. 

You may also notice `user_handler.go` and `user_service.go`. Currently, no user service nor user handler functions exist because there's no need to keep track of users for the endpoint that's available. However, if I can find a way for the data aggregation engine to pull the full amount of data needed for the companies that don't report their metrics uniformly, then I may start offering a commercial API service, in which case there'd be a need for user authentication.


`cmd/data/main.go` is the second entrypoint to the application. This houses a CLI tool to run the data retrieval engine. `go run cmd/data/main.go seed` seeds the database. Currently it takes about an hour and a half. `go run cmd/data/main.go update` updates the existing entries in the database with the most recent yearly filing data. `go run cmd/data/main.go market` updates the live market price, as well as metrics that need to be calculated using the live market price. 

The data retrieval engine works by first retrieve a list of all CIKs (CIK is unique identifier assigned to each reporting entity) that report to the SEC, Each CIK is ran against another endpoint that reveals more information about company name, industry, and type. Only companies determined to be of type "large-accelerated-filer" actually have their full metrics retrieved. Companies other than large-acclereatead filers do not report the same facts and metrics so I decided it's best to leave them out until I scale up. From there, each company's facts are validated, key ratios are calculated, and the company data is stored in the PostgreSQL database.

The third entrypoint is `cmd/migrate/main.go`, which is another CLI with 2 commands: `go run cmd/migrate/main.go up` and `go run cmd/migrate/main.go down`. Both these CLI commands migrate the database either up to its latest iteration or down to its earliest depending on the command used. 

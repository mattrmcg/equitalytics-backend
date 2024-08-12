# Backend for Equitalytics.com
This repository houses all the backend components of [https://www.equitalytics.com](Equitalytics.com). 

The backend is entirely built in go, and hosted in an Azure Container Instance. The docker container's base image is ubuntu server. The data is stored in a PostgreSQL database server also hosted and maintained on Azure.

There are three entrypoints to the backend application, all held within the `cmd` directory. 

`cmd/main.go` is the entrypoint for the API server. The API serves a dynamic route `/info/{ticker}`, where company data for the specificied `{ticker}` is gathered and returned upon request. Requests are handled within their respective handler files, and then database retrieval operations take place at the `info_service` layer. 

You may also notice `user_handler.go` and `user_service.go`. Currently, no user service nor user handler functions exist because there's no need to keep track of users for the endpoint that's available. However, if I can find a way for the data aggregation engine to pull the full amount of data needed for the companies that don't report their metrics uniformly, then I may start offering a commercial API service, in which case there'd have to be user authentication.

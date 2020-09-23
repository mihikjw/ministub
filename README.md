# ministub
A small API stubbing tool for microservice dependency simulation, allows the developer to define an API schema in YAML as well as follow-on actions allowing two-way communication between microservices, or to further contact a third API or similar. I developed this tool primarily to allow the mocking of microservices where a two-way communication stream was required, either as part of a CI pipeline or as a CLI tool on a local dev environment.

- Define a YAML file with an API and actions on request
- Define a series of input requests to recieve and return different status codes on different occurances (including random injection)
- Define a series of follow-on subsiquent external requests, also injecting different requests to allow for failure and resilience testing
- Extract metrics from an API

## Metrics
Metrics are gathered from succesful requests, i.e. there was no error in the original request. For example, if you configured a weighting for 40% of client requests to a specific endpoint should return HTTP status 400, but the client request was actually warranting a response of status 400, these requests would not be considered for the weighting of the next succesful request.
# ministub
A small API stubbing tool for microservice dependency simulation, allows the developer to define an API schema in YAML as well as follow-on actions allowing two-way communication between microservices, or to further contact a third API or similar. I developed this tool primarily to allow the mocking of microservices where a two-way communication stream was required, either as part of a CI pipeline or as a CLI tool on a local dev environment.

- Define a YAML file with an API and actions on request
- Define a series of input requests to recieve and return different status codes on different occurances, with a percentage weighting for each response
- Define a series of follow-on subsiquent actions upon an incoming request
- Extract metrics from an API from a `/stats` endpoint
- Shutdown the application from an `/exit` endpoint

## Install
1. `git clone https://github.com/MichaelWittgreffe/ministub`
2. `cd ministub`
3. `make`
4. Optional Docker Build: `make docker-build`
4. Optional Install On PATH: `sudo make install`

## Usage
ministub is a CLI tool, also provided is a docker version if you would prefer.

## File Format
ministub uses a YAML format to define an API to host. A fully-featured example can be found at `/examples/v1demopapi.yml`.
Put more stuff here ...

## TO DO
- Docs
- Unit Tests!!
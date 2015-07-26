jsonapi
=======
developer's note: jsonapi is not complete yet. this readme refers to features that are not complete yet, but will be before this project hits 1.0. Volenti non fit iniuria.

jsonapi is a Golang http framework for producing JSON-based APIs. jsonapi's mission is to provide rails-like simplicity without sacrificing the flexibility, performance, scalability, and concurrency that Go is known for. 

[Fully-functional Blog Example](example/main.go)

jsonapi is:
* Easy to understand. Code written for jsonapi reads like a graph of your software architecture model, and everything else is CRUD primitives.
* Powerful by default. Several Resources (SQL, Redis, etc) come built in, as well as Caching and Pagination layers that can be easily tacked on anywhere.
* Very fast. Most requests incur under 2ms of framework overhead, even on my low-grade laptop.
* Flexible for high-volume enviornments. jsonapi is designed to be the framework that you write your prototype with, and the framework that you scale to millions with, without substantial rewrites in between.
* Concurrent to the extent possible. jsonapi builds an internal dependency tree that allows Resources to compute at their own speed, independent of one another until they must be brought together.
* Fail-safe. Resource-agnostic mechanisms exist to roll back transactions after a critical failure
* Built with authentication in mind. jsonapi provides a generalized authentication scheme that can be easily extended to check for login, permission level, or outright refuse certain types of requests for certain resources
* [Compliant with the json-api spec](https://github.com/json-api/json-api/blob/gh-pages/format/index.md)

jsonapi is not:
* capable of being a drop in replacement for a pre-existing API, unless that API is already json-api spec compliant.

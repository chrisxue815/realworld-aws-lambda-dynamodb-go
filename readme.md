# ![RealWorld Example App](logo.png)

> ### AWS Lambda + DynamoDB + Go codebase containing real world examples (CRUD, auth, advanced patterns, etc) that adheres to the [RealWorld](https://github.com/gothinkster/realworld) spec and API.

### [Demo](https://chrisxue815.github.io/realworld/build/#/)

This codebase was created to demonstrate a fully fledged fullstack application built with **AWS Lambda + DynamoDB + Go** including CRUD operations, authentication, routing, pagination, and more.

We've gone to great lengths to adhere to the **AWS Lambda + DynamoDB + Go** community styleguides & best practices.

For more information on how to this works with other frontends/backends, head over to the [RealWorld](https://github.com/gothinkster/realworld) repo.

# Getting started

## Prerequisite

* Install Go, Node.js, Serverless CLI
* In `angularjs-realworld-example-app`, run `npm install`

## Build and deploy backend

In the root directory of this project:

* `make build`
* `sls deploy --stage dev`

## Build and serve frontend

In `angularjs-realworld-example-app`:

* `npx gulp`

# How it works

Routes and their handlers are defined in `serverless.yml`.

For example, the following section means `POST /users` is handled by `bin/users-post`, which is built from `route/users-post/main.go`.

```
  users-post:
    handler: bin/users-post
    events:
      - http:
          path: users
          method: post
          cors: true
```

# Design choices
* Scrypt-based password hasing
* Input validation
* Data consistency

These tradeoffs are made for simpler code:
* Hardcoded secret. Cons: can't invalidate tokens
* Shared states (like DB and RNG) are singletons, no dependency injection. Cons: can't control their lifecycles, potential memory leak, unit-test-unfriendly
* Immutable username
* Case-sensitive username
* Performance bottlenecks around global secondary indices with a single hash-key value, like ArticleTable.CreatedAt and TagTable.ArticleCount.

<p align="center" style="font-size: 24px; font-family: Arial, sans-serif;">
  <h1 style="font-family: 'Playfair Display', serif; font-size: 36px;">Onde Hoje ?</h1>
  <p style="font-family: 'Lato', sans-serif; font-size: 16px;">A Single place to share the underground</p>
</p>

Onde Hoje is an innovative application that allows users to share underground events in their city. This app is designed for people who are looking for unique experiences and want to explore the hidden gems of their city.
With Onde Hoje, users can find events that are not advertised on mainstream platforms such as Facebook or Instagram. This app is perfect for those who want to discover new artists, musicians, and performers before they become mainstream.


<p>
  <img src="./assets/under.png" alt="image_alt_text">
</p>





To run the API locally and start shipping 🚢 new features is easy:


```bash
# start running the development environment. For now, a PostgreSQL database.
make dev/start
```

After that, to run the API, apply the command:

```go
go run cmd/ondehoje/main.go
```

# API Requests

Just access your browser at http://localhost:8000/docs. That's it, all routes grouped in one place!


# Core Concepts
Some concepts that this simple API will build under:


* API
* Database
* Database Migration
* Structured Logs
* Dev Container environment
* Go Unit test
* Metrics
* Swagger/OpenAPI

# Heroku Database

Before connecting to the PostgreSQL database, make sure you have the [Heroku CLI installed](https://devcenter.heroku.com/articles/heroku-cli)

```bash
sudo apt-get install postgresql  #to install psql

heroku pg:psql -a ondehoje # to access database and execute admin commands
```

# Ship New code to production

Although there is a CI for that, it's possible to SHIP🚀 new code to production just using your machine.

Just run the following commands, in this order, and be happy😃:

```
make image
make publish
make heroku/release

```

# Useful Vscode extensions

[OpenAPI Swagger Editor](https://42crunch.com/tutorial-openapi-swagger-extension-vs-code/)
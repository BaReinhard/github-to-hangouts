# Github to Hangouts Alerting

Simple way to Alert Hangout Chatrooms with Webhooks from Github

### Setup

- Enable Webhook on Specific Repo
- Set /endpoint and ?key=KEY
- Ensure content-type is `application/json`
- Setup app.yaml

```
runtime: go
api_version: go1

service: github-integrations

env_variables:
 SECURE_ENDPOINT: "YOURENDPOINT"
 SECURE_KEY: "YOUR_PRESHARED_KEY"
handlers:


# All URLs are handled by the Go application script
- url: /.*
  script: _go_app
```

## Go Chat Router

A simple go chat router, running in GAE standard env. Hangouts chat sends a Message Payload to the App Engine, it checks the space.name and sends the payload to another App Engine instance, cloud function, pub/sub, etc.

**Why?**

Currently, GCP projects can only handle 1 interactive bot. This is not scalable. To be able to scale a single Bot, routing the Payloads based on Hangouts Chat Rooms is a simple work around. The one downside is that the Bot name must be generic.

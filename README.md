# Hora

> "*Hora*: In Greco-Roman mythology, any one of the personifications of the seasons and goddesses of natural order; in
> the Iliad they were the custodians of the gates of Olympus." - [Britannica][1] 

This is a design framework that can be used for creating a single point of entry into multiple validation services for
a given payload. A gatekeeper if you will...

**Use Case** - This framework can and should be used if you want to dynamically validate the structure of a data payload
against multiple validators. The dynamic portion is represented by the fact that a given validator is only run if the
given payload matches the validator's criteria to run.

The main goal here is not to demonstrate how to implement a specific payload validator, but instead is to show the
flexible framework that worked for us, regardless of what validation is being run.

The flow chart below depicts the described design at a high level

![flow][2]

## Example

We have provided an example json schema validator that can be easily tested out by doing the following.

1. `make compile`
2. `make run`
3. view Swagger documentation @ `localhost:8080`
4. Provide a list of validation "candidates" (individual data objects to be validated) to the `/validate` endpoint and
compare results to the schema defined in `./src/validators/example/schema.json`

<!-- Links -->
[1]: https://www.britannica.com/topic/Hora-Greek-mythology
[2]: documentation/img/flow.png

# Damo's thing

you'll need `make` and `docker` to do dev things...

## Run

Note that you can run this same binary as a local server by passing the `-l` flag and optionally specifying the port with `-p <port>`.

```
make run
```

### Environment

The lambda is configured with a few env vars see `serverless.yml`. To facilitate sharing between local execution and lambda these live as exports in `config/<ENVIRONMENT>-env.sh` (so long as they aren't secrets ;)

## Deploy

Deploy api to lambda.

```
make deploy
```

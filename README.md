# Optibus Go CLI tools

This might temporary, but for now looking a place to
throw some clis that help in devx

* aws-secrets - will fetch secrets and output them in a way a tool
like direnv or any tool that know how to work .env files can use
* check-js-deps - specific for aramda/mithra mono repo. This tool tries
to let the developer know when to `pnpm instsll`

Until we sign these clis with apple developer id,
developers might need to build them locally

## building

```sh
make build
```

```

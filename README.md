#  Service generator

## Stats
Lines of code
 - Total:          43355
    - Infra & domain: 27694(64%)
    - Application:    15661(36%)

Number of commits
 - Total:          551
   - Infra & domain: 394(70%)
   - Application:    168(30%)

```
$ git log --oneline -- ./internal/application | wc -l 
168

$ find ./internal/application  -name '*.go' | xargs  wc -l
15661 total

git log --oneline -- ./internal/domain ./internal/infra/{clients,http,kafka,marshaling,observability,rabbitmq} ./ahoy ./.github | wc -l
394

$ find ./internal/domain ./internal/infra/{clients,http,kafka,marshaling,observability,rabbitmq} -name '*.go' | xargs  wc -l
27694 total
```


## Scope of work
gradual conversion
versioning

domain types

api contract
api codegen

clients, discovery and transport layer

DTOs and serialization

auth, middlewares

config and env variable reads, updates, vault integration

kafka channels and schema registry

cron jobs

background jobs

command line tools

db connections and pools

caches

logs
metrics

ACL

escape hatches

dapr

## Roadmap
Domain types
API layer
Pluggable transport layer
DTO
Service configuration
Vault integration
Internal dependencies
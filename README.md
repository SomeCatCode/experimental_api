# Just a golang Playground

## Routes

### index
GET /organisation 

### get by id
GET /organisation/{id}

### create
POST /organisation

### update by id
PATCH /organisation/{id}

### delete by id
DELETE /organisation/{id}


## TODO
- write basic crud operations
- combine model, repo and handler
- make interface for repo, for diffrent databases
- field validation
- better error response
- sideloads
- filters
- jwt
- rolebased responses/tranformer
- rolebased setter
- add tests
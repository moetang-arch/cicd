# cicd
CI/CD tool for your tech stack

# 1. Architecture

events:

```
(git push event) -> event processors -> (push event)
    -> build environment setup -> (setup event)
    -> testing, building, and packaging -> (release event)
    -> (concurrent)  -> release environment
                     -> deployment
```

sync:

```
(a hour/day) -> sync with repository -> (push event)
    -> ...
```

# 2. Event Processor - 1st stage

contains:

```
github processor
```

# 3. Build Environment Setup / Release Environment - 2nd stage

contains:

```
setup controller
setup agent
```

### Requirement

* Use `go mod` for your project. Then only way to support dependencies.
* Keep your `go.mod` file up-to-date

### Planning

* [ ] support cancelling job
* [ ] support go.mod
* [ ] release environment

# 4. Testing, Building, Packaging - 3rd stage

```
tbp agent
tbp service
```

### Planning

* [ ] support cancelling job

# 5. Deployment - 4th stage

```
deployment service
deployment agent
```

# 6. Web UI - Controller

# 7. Planning

* [ ] event delivery failed(git push)
* [ ] need sync with repository once a hour/day if commits/branches/tags are out-of-date

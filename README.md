# cicd
CI/CD tool for your tech stack

# 1. Architecture

events:

```
(git push event) -> event processors -> (push event)
    -> build environment setup -> (setup event)
    -> testing, building, and packaging -> (release event)
    -> deployment
```

sync:

```
(a hour/day) -> sync with repository -> (push event)
    -> ...
```

# 2. Event Processor

contains:

```
github processor
```

# 3. Build Environment Setup

contains:

```
setup master
setup slave
```

### Planning

* support cancelling job

# 4. Testing, Building, Packaging

```
tbp agent
tbp service
```

### Planning

* support cancelling job

# 5. Deployment

```
deployment service
deployment agent
```

# 6. Web UI

# 7. Planning

* [ ] event delivery failed(git push)
* [ ] need sync with repository once a hour/day if commits/branches/tags are out-of-date

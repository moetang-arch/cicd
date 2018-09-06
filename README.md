# cicd
CI/CD tool for your tech stack

# 1. Architecture

(git push event) -> processors -> (push event)
    -> build environment setup -> (setup event)
    -> testing, building, and packaging -> (release event)
    -> deployment

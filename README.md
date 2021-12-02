Your task is to:
- prepare a validation system to check if a number is a prime number and asynchronously store the result
- prepare presentation and demonstrate execution

The system should have the following characteristic:
1. Results should be stored to postgres database
2. The prime validation algorithm should be brute-force and coded by yourself (please, suggest an alternative algorithm that you would use in real-life scenario)
3. There should be 2 additional micro-services:
	I) Input service running on HTTP protocol, accepting POST requests in the following JSON format: {n=123325435} and calling background service
	II) Background service running on gRPC protocol, that will be called from input service and will validate if a number received is a prime number and will store the result to postgres database.
	    The following fields should be stored: number_tested, is_prime_number, validation_time, auto_incremening_seq, time_needed_to_validate_microsec
4. All components should run in docker image, preferable (not necessary) configured in Kubernetes
5. System should be able to handle multiple concurrent request and should be writing to postgres database in an asynchronous manner

### /cmd
Main applications for this project.

### /internal
Private application and library code. This is the code you don't want others importing in their applications or libraries.

### /pkg
Library code that's ok to use by external applications.

### /migrations
DB migrations (using migrate CLI)

### /build
Packaging and Continuous Integration.

### /deployments
Kubernetes files

[You can read more about the structure here](https://github.com/golang-standards/project-layout)

## TODO
- Tests
- Documentation
- Code review
- Next.JS app
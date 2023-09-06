# crawler

## Running instructions
### Manually
In order to run the application in development mode the user must follow the next steps
```
cd ./cmd
go get
go run main.go
```
### Docker
The application also can be executed in docker 
```
docker build --rm -t crawler:latest .
docker run -d -p 8080:8080  --name crawler crawler:latest
docker logs -f crawler 
```

* Environment variabes can be passed as parameters in docker run command using the flag --env

## Configuration
If no environment variables the application will start using the default values:

| Variable | Default | Description |
| --- | --- |--- |
WEB_PAGE | https://parserdigital.com/ | root site to start crawling
WORKERS| 10 | max number of concurrent workers exploring for links

## Desing considerations
In order to speed up the crawling process, the application uses concurrent workers to navigate through the different links. By default, the application will **dynamically** spin up to the maximum number of workers specified in the WORKERS environment variable.

### Improvement oportunities
- Use an external queue system to queue pending links for exploration.
- Use an external database in the store package.
- Increase code coverage.
- Add a larger mock HTML page for the current benchmark test.
- Read Robots.txt from the root site and follow rules such as:
        User-agent
        Request delay
        Allowed and disallowed paths


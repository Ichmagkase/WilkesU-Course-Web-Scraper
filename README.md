# Wilkes University Course Web Scraper

This software parses the Wilkes University's Course Catalog.

Websites that can be parsed must be pertain to a year starting at 2020.

## Building the application:

The parsers uses ```docker-compose``` to build the application.

Note, you will need to have port 27017 and port 8080 open for MongoDB and the scraper respectively.

In the root directory of this project, run:

```
docker-compose up --build
```

## Arguments

The scraper allows for command line arguments

```
scraper semester year
```


The semester can either be:
 - ```F```
 - ```Sp```

The year can be any number , but note the parsers can only parses from years after 2019.

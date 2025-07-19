# Wilkes University Course Web Scraper

## The Idea

The current course website Wilkes University uses does not have 2 main things:
- Course Search Functionality
- Ease of use on mobile devices

With this project, we plan to fix these and make a better platform for all to use. Our goal is to be able to search courses not just from the current semester, but from prior ones, all in one, easy to use, website.

## Long Term Goals

1. Local running website that shows data found on rosters.wilkes.edu/scheds/*.
2. Local running website fits on most devices and looks visually pleasing.
3. Deployed API with proper endpoints.
4. Deployed website using the deployed API.
5. Allowing courses to be selected for users to create there own course schedule.
6. Mobile app development

## Getting Started

Running the website is done with `docker compose` as we have 2 services: the website and the website scraper.
You can execute:
```bash
docker compose up
```
Which will build and start up the services. The scraper API runs on `localhost:8080` by default and the website on `localhost:5173` by default.

## Just Scraping

To run just the course website scraper you can use `scraper.go` as a CLI tool

```bash
go run scraper.go [F | Sp] year
```

The year is formatted as XXXX and the semesters are abbreviated to there shortnames F and Sp.

**NOTE**: The years can only be parsed from 2020 and onward as the website format changed between 2019 and 2020.

## Maintainers

[Nathaniel Martes](https://github.com/NateMartes)
[Zackery Drake](https://github.com/Ichmagkase)

package main

import (
	"awesomeProject/config"
	"fmt"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"strings"
	"sync"
)

var driver neo4j.Driver
var err error
var once sync.Once

func getDriver() neo4j.Driver {
	config := config.GetConfig()
	useConsoleLogger := func(level neo4j.LogLevel) func(config *neo4j.Config) {
		return func(config *neo4j.Config) {
			config.Log = neo4j.ConsoleLogger(level)
		}
	}
	once.Do(func() {
		driver, err = neo4j.NewDriver(config.DataBaseUrl, neo4j.BasicAuth(config.DataBaseUserName, config.DataBasePassword, ""), useConsoleLogger(neo4j.ERROR))
		if err != nil {
			check(err)
		}
	})
	return driver
}

func insertPage(page FormattedPage) {
	var (
		dbDriver neo4j.Driver
		result   neo4j.Result
		session  neo4j.Session
		err      error
	)

	dbDriver = getDriver()

	if session, err = dbDriver.Session(neo4j.AccessModeWrite); err != nil {
		check(err)
	}

	for i := 0; i < len(page.Links); i++ {
		var singleRelation SingleRelation
		singleRelation.Link = formatStringToWikiLinkFormat(extractLink(page.Links[i]))
		singleRelation.Title = formatStringToWikiLinkFormat(page.Title)
		insertLink(result, err, session, singleRelation)
	}
}

func insertLink(result neo4j.Result, e error, session neo4j.Session, singleRelation SingleRelation) {
	const INSERT_CYPHER_QUERY = `MERGE (page:Page {name: $name}) MERGE (link:Page {name: $link}) MERGE (page)-[:LINKS_TO]->(link) RETURN page.name,link.name`
	result, err = session.Run(INSERT_CYPHER_QUERY, map[string]interface{}{
		"name": singleRelation.Title,
		"link": singleRelation.Link,
	})
	if err != nil {
		check(err)
	}
	for result.Next() {
		fmt.Printf("Created A Link Between '%s' and '%s'\n", result.Record().GetByIndex(0).(string), result.Record().GetByIndex(1).(string))
	}
	if err = result.Err(); err != nil {
		check(err)
	}
}

func extractLink(link string) string {
	link = strings.Trim(link, "[")
	link = strings.Trim(link, "]")
	link = strings.Split(link, "|")[0]
	return link
}

func formatStringToWikiLinkFormat(value string) string {
	trimmedValue := strings.TrimSpace(value)
	lowerCasedValue := strings.ToLower(trimmedValue)
	camelCasedValue := strings.Title(lowerCasedValue)
	spacesReplacedValue := strings.Replace(camelCasedValue, " ", "_", -1)
	return spacesReplacedValue
}

type SingleRelation struct {
	Title string
	Link  string
}

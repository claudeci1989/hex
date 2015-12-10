package connectors

import (
	"log"
	"reflect"
)

var List = make(map[string]reflect.Type)

func init() {
	List["alias"] = reflect.TypeOf(Alias{})
	List["cli"] = reflect.TypeOf(Cli{})
	List["bamboo"] = reflect.TypeOf(Bamboo{})
	List["email"] = reflect.TypeOf(Email{})
	List["exec"] = reflect.TypeOf(Exec{})
	List["imageme"] = reflect.TypeOf(ImageMe{})
	List["jira"] = reflect.TypeOf(Jira{})
	List["monitor"] = reflect.TypeOf(Monitor{})
	List["slack"] = reflect.TypeOf(Slack{})
	List["response"] = reflect.TypeOf(Response{})
	List["rss"] = reflect.TypeOf(Rss{})
	List["website"] = reflect.TypeOf(Website{})
	List["wolfram"] = reflect.TypeOf(Wolfram{})
}

func MakeConnector(connType string) interface{} {
	if ct, ok := List[connType]; ok {
		c := (reflect.New(ct).Elem().Interface())
		return c
	} else {
		log.Print("Error in configuration, connector type '" + connType + "' not supported")
		log.Fatal("Exiting due to configuration error")
		return nil
	}
}

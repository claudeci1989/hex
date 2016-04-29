package connectors

import (
	"log"
	"reflect"
)

var List = make(map[string]reflect.Type)

func init() {
	List["cli"] = reflect.TypeOf(Cli{})
	List["bamboo"] = reflect.TypeOf(Bamboo{})
	List["email"] = reflect.TypeOf(Email{})
	List["exec"] = reflect.TypeOf(Exec{})
	List["hipchat"] = reflect.TypeOf(Hipchat{})
	List["imageme"] = reflect.TypeOf(ImageMe{})
	List["jira"] = reflect.TypeOf(Jira{})
	List["logging"] = reflect.TypeOf(Logging{})
	List["monitor"] = reflect.TypeOf(Monitor{})
	List["slack"] = reflect.TypeOf(Slack{})
	List["response"] = reflect.TypeOf(Response{})
	List["rss"] = reflect.TypeOf(Rss{})
	List["ssh"] = reflect.TypeOf(Ssh{})
	List["website"] = reflect.TypeOf(Website{})
	List["wolfram"] = reflect.TypeOf(Wolfram{})
	List["redis"] = reflect.TypeOf(Redis{})
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

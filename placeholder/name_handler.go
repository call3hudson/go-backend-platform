package placeholder

import (
	"fmt"
	"platform/http/actionresults"
	"platform/logging"
)

var names = []string{"Alice", "Bob", "Charlie", "Dora"}

type NameHandler struct {
	logging.Logger
}

func (n NameHandler) GetName(i int) actionresults.ActionResult {
	n.Logger.Debugf("GetName method invoked with argument: %v", i)
	var response string
	if i < len(names) {
		response = fmt.Sprintf("Name #%v: %v", i, names[i])
	} else {
		response = fmt.Sprintf("Index out of bounds")
	}
	return actionresults.NewTemplateAction("simple_message.html", response)
}

func (n NameHandler) GetNames() actionresults.ActionResult {
	n.Logger.Debug("GetNames method invoked")
	return actionresults.NewTemplateAction("simple_message.html", names)
}

type NewName struct {
	Name          string
	InsertAtStart bool
}

func (n NameHandler) PostName(new NewName) actionresults.ActionResult {
	n.Logger.Debugf("PostName method invoked with argument %v", new)
	if new.InsertAtStart {
		names = append([]string{new.Name}, names...)
	} else {
		names = append(names, new.Name)
	}
	return actionresults.NewRedirectAction("/names")
}

func (n NameHandler) GetJsonData() actionresults.ActionResult {
	return actionresults.NewJsonAction(names)
}

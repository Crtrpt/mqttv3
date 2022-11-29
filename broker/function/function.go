package function

// topic a->function->topic b

// broker a topic a ->function-> broker b topic b

type Function interface {
	Handle(input interface{}, state *interface{}) (output interface{}, err error)
}

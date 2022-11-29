package function

// topic a->function->topic b

// broker a topic a ->function-> broker b topic b

type Function interface {
	Handle(input any, state *any) (output any, err error)
}

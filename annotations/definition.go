package annotations

type Annotation struct {
	Name string
	Data map[string]interface{}
}

type AnnotationDef interface {
	// Returns the name of the annotation (e.g. @<name>)
	GetName() string
}

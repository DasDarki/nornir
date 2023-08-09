package annotations

type Annotation struct {
	Name string
	Data map[string]interface{}
}

type AnnotationDef interface {
	// Returns the name of the annotation (e.g. @<name>)
	GetName() string
	// Returns the allowed usages for the annotation
	GetUsages() []UsageKind
}

// The allowed usages for an annotation
type UsageKind string

const (
	UsageKindPackage  UsageKind = "package"
	UsageKindFunction UsageKind = "func"
	UsageKindVar      UsageKind = "var"
)

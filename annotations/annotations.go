package annotations

type ControllerAnnotation struct {
	Path     string `annotation:"name=path,numeric=0,default=$empty"`
	Frontend bool   `annotation:"name=frontend,numeric=1,default=true"`
}

func (a *ControllerAnnotation) GetName() string {
	return "Controller"
}

func (a *ControllerAnnotation) GetUsages() []UsageKind {
	return []UsageKind{UsageKindPackage}
}

type GetAnnotation struct {
	Path     string `annotation:"name=path,numeric=0,default=$empty"`
	Frontend bool   `annotation:"name=frontend,numeric=1,default=true"`
}

func (a *GetAnnotation) GetName() string {
	return "Get"
}

func (a *GetAnnotation) GetUsages() []UsageKind {
	return []UsageKind{UsageKindFunction}
}

type PostAnnotation struct {
	Path     string `annotation:"name=path,numeric=0,default=$empty"`
	Frontend bool   `annotation:"name=frontend,numeric=1,default=true"`
}

func (a *PostAnnotation) GetName() string {
	return "Post"
}

func (a *PostAnnotation) GetUsages() []UsageKind {
	return []UsageKind{UsageKindFunction}
}

type PutAnnotation struct {
	Path     string `annotation:"name=path,numeric=0,default=$empty"`
	Frontend bool   `annotation:"name=frontend,numeric=1,default=true"`
}

func (a *PutAnnotation) GetName() string {
	return "Put"
}

func (a *PutAnnotation) GetUsages() []UsageKind {
	return []UsageKind{UsageKindFunction}
}

type DeleteAnnotation struct {
	Path     string `annotation:"name=path,numeric=0,default=$empty"`
	Frontend bool   `annotation:"name=frontend,numeric=1,default=true"`
}

func (a *DeleteAnnotation) GetName() string {
	return "Delete"
}

func (a *DeleteAnnotation) GetUsages() []UsageKind {
	return []UsageKind{UsageKindFunction}
}

type PatchAnnotation struct {
	Path     string `annotation:"name=path,numeric=0,default=$empty"`
	Frontend bool   `annotation:"name=frontend,numeric=1,default=true"`
}

func (a *PatchAnnotation) GetName() string {
	return "Patch"
}

func (a *PatchAnnotation) GetUsages() []UsageKind {
	return []UsageKind{UsageKindFunction}
}

type HeadAnnotation struct {
	Path     string `annotation:"name=path,numeric=0,default=$empty"`
	Frontend bool   `annotation:"name=frontend,numeric=1,default=true"`
}

func (a *HeadAnnotation) GetName() string {
	return "Head"
}

func (a *HeadAnnotation) GetUsages() []UsageKind {
	return []UsageKind{UsageKindFunction}
}

type OptionsAnnotation struct {
	Path     string `annotation:"name=path,numeric=0,default=$empty"`
	Frontend bool   `annotation:"name=frontend,numeric=1,default=true"`
}

func (a *OptionsAnnotation) GetName() string {
	return "Options"
}

func (a *OptionsAnnotation) GetUsages() []UsageKind {
	return []UsageKind{UsageKindFunction}
}

type RequestMappingAnnotation struct {
	Method   string `annotation:"name=method,numeric=0,default=GET"`
	Path     string `annotation:"name=path,numeric=1,default=$empty"`
	Frontend bool   `annotation:"name=frontend,numeric=2,default=true"`
}

func (a *RequestMappingAnnotation) GetName() string {
	return "RequestMapping"
}

func (a *RequestMappingAnnotation) GetUsages() []UsageKind {
	return []UsageKind{UsageKindFunction}
}

type InjectNornirAnnotation struct {
}

func (a *InjectNornirAnnotation) GetName() string {
	return "InjectNornir"
}

func (a *InjectNornirAnnotation) GetUsages() []UsageKind {
	return []UsageKind{UsageKindVar}
}

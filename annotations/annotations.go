package annotations

type ControllerAnnotation struct {
	Path string `annotation:"name=path,numeric=0,default=$empty"`
}

func (a *ControllerAnnotation) GetName() string {
	return "Controller"
}

type GetAnnotation struct {
	Path string `annotation:"name=path,numeric=0,default=$empty"`
}

func (a *GetAnnotation) GetName() string {
	return "Get"
}

type PostAnnotation struct {
	Path string `annotation:"name=path,numeric=0,default=$empty"`
}

func (a *PostAnnotation) GetName() string {
	return "Post"
}

type PutAnnotation struct {
	Path string `annotation:"name=path,numeric=0,default=$empty"`
}

func (a *PutAnnotation) GetName() string {
	return "Put"
}

type DeleteAnnotation struct {
	Path string `annotation:"name=path,numeric=0,default=$empty"`
}

func (a *DeleteAnnotation) GetName() string {
	return "Delete"
}

type PatchAnnotation struct {
	Path string `annotation:"name=path,numeric=0,default=$empty"`
}

func (a *PatchAnnotation) GetName() string {
	return "Patch"
}

type HeadAnnotation struct {
	Path string `annotation:"name=path,numeric=0,default=$empty"`
}

func (a *HeadAnnotation) GetName() string {
	return "Head"
}

type OptionsAnnotation struct {
	Path string `annotation:"name=path,numeric=0,default=$empty"`
}

func (a *OptionsAnnotation) GetName() string {
	return "Options"
}

type RequestMappingAnnotation struct {
	Method string `annotation:"name=method,numeric=0,default=Get"`
	Path   string `annotation:"name=path,numeric=1,default=$empty"`
}

func (a *RequestMappingAnnotation) GetName() string {
	return "RequestMapping"
}

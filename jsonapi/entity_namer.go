package jsonapi

// The EntityNamer interface can be optionally implemented to directly return the
// name of resource used for the "type" field.
//
// Note: By default the name is guessed from the struct name.
type EntityNamer interface {
	GetName() string
}

// The EntityPather interface can be optionally implemented to directly return the
// path for the resource to be used for routing to the resource.
// The resulting base URL will take the form `/{apiPrefix}/{path}/{id}`
// Note: By default the path is guessed from the struct name.
type EntityPather interface {
	GetPath() string
}

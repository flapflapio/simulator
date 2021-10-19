package machine

// Some methods for serializing types
type Marshalable interface {
	Json() string
	JsonMap() map[string]interface{}
	String() string
}

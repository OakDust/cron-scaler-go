package handleradapter

// adapter for the external handler
// makes grpc request to external service called scaler-handler
type Adapter struct {
}

func New() *Adapter {
	return &Adapter{}
}

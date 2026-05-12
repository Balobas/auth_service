package deliveryGrpcInterceptors

type Provider struct {
	ucAuth UcAuth
}

func NewProvider(ucAuth UcAuth) *Provider {
	return &Provider{
		ucAuth: ucAuth,
	}
}

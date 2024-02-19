package grpc

import (
	"google.golang.org/grpc/resolver"
)

func NewResolverBuilder(addresses ...string) *ResolverBuilder {
	return &ResolverBuilder{
		addresses: addresses,
	}
}

type ResolverBuilder struct {
	addresses []string
}

func (rb *ResolverBuilder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	r := &Resolver{
		addresses: rb.addresses,
		cc:        cc,
	}
	r.ResolveNow(resolver.ResolveNowOptions{})
	return r, nil
}

func (rb *ResolverBuilder) Scheme() string {
	return "registrar"
}

type Resolver struct {
	addresses []string
	cc        resolver.ClientConn
}

func (r *Resolver) ResolveNow(options resolver.ResolveNowOptions) {
	addresses := make([]resolver.Address, 0, len(r.addresses))
	for _, address := range r.addresses {
		addresses = append(addresses, resolver.Address{Addr: address})
	}
	state := resolver.State{
		Addresses: addresses,
	}
	err := r.cc.UpdateState(state)
	if err != nil {
		r.cc.ReportError(err)
	}
}

func (r *Resolver) Close() {
}

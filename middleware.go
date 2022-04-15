package pepper

type MiddlewareFunc func(p *Pepper, res Response, req *Request) bool
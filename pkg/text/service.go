package text

type Service struct {
	Resolver
}

func NewService(r Resolver) *Service {
	return &Service{r}
}

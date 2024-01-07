package text

import (
	"fmt"
	"math/rand"

	"github.com/zostay/today/pkg/ref"
)

type Service struct {
	Resolver
}

func NewService(r Resolver) *Service {
	return &Service{r}
}

type randomOpts struct {
	category string
	book     string
}

type RandomReferenceOption func(*randomOpts)

func FromBook(name string) RandomReferenceOption {
	return func(o *randomOpts) {
		o.book = name
	}
}

func FromCategory(name string) RandomReferenceOption {
	return func(o *randomOpts) {
		o.category = name
	}
}

// Random pulls a random reference from the Bible and returns it. You can use the
// options to help narrow down where the passages are selected from.
func (s *Service) Random(opt ...RandomReferenceOption) (string, error) {
	o := &randomOpts{}
	for _, f := range opt {
		f(o)
	}

	var (
		b  *ref.Book
		vs []ref.VerseRef
	)
	if o.category != "" {
		exs, err := ref.LookupCategory(o.category)
		if err != nil {
			return "", err
		}

		// lazy way to weight the books by the number of verses they have
		bag := make([]*ref.BookExtract, 0, len(exs))
		for i := range exs {
			for range exs[i].Verses() {
				bag = append(bag, &exs[i])
			}
		}

		be := bag[rand.Int()%len(bag)] //nolint:gosec // weak random is fine here
		b = be.Book
		vs = ref.RandomPassageFromExtract(be)
	} else {
		if o.book != "" {
			ex, err := ref.LookupBook(o.book)
			if err != nil {
				return "", err
			}

			b = ex.Book
		} else {
			b = ref.RandomCanonical()
		}

		vs = ref.RandomPassage(b)
	}

	v1, v2 := vs[0], vs[len(vs)-1]

	if len(vs) > 1 {
		return fmt.Sprintf("%s %s-%s", b.Name, v1.Ref(), v2.Ref()), nil
	}

	return fmt.Sprintf("%s %s", b.Name, v1.Ref()), nil
}

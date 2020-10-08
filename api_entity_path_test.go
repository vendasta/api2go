package api2go

import (
	"net/http"
	"net/http/httptest"
	"strings"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type SconeTaste struct {
	ID    string `json:"-"`
	Taste string `json:"taste"`
}

func (s SconeTaste) GetID() string {
	return s.ID
}

func (s *SconeTaste) SetID(ID string) error {
	s.ID = ID
	return nil
}

func (s SconeTaste) GetName() string {
	return "Scone-tastes"
}

func (s SconeTaste) GetPath() string {
	//Leading and trailing slashes are not required. Including in test to make sure we don't end up with double slashes
	return "/food/sconeTastes/"
}

type SconeResource struct{}

func (s SconeResource) FindOne(ID string, req Request) (Responder, error) {
	return &Response{Res: SconeTaste{ID: "blubb", Taste: "Very Bad"}}, nil
}

func (s SconeResource) FindAll(req Request) (Responder, error) {
	return &Response{Res: []SconeTaste{
		{
			ID:    "1",
			Taste: "Very Good",
		},
		{
			ID:    "2",
			Taste: "Very Bad",
		},
	}}, nil
}

func (s SconeResource) Create(obj interface{}, req Request) (Responder, error) {
	e := obj.(SconeTaste)
	e.ID = "newID"
	return &Response{
		Res:  e,
		Code: http.StatusCreated,
	}, nil
}

func (s SconeResource) Delete(ID string, req Request) (Responder, error) {
	return &Response{
		Res:  SconeTaste{ID: ID},
		Code: http.StatusNoContent,
	}, nil
}

func (s SconeResource) Update(obj interface{}, req Request) (Responder, error) {
	return &Response{
		Res:  obj,
		Code: http.StatusNoContent,
	}, nil
}

var _ = Describe("Test route renaming with EntityPather interface", func() {
	var (
		api  *API
		rec  *httptest.ResponseRecorder
		body *strings.Reader
	)
	BeforeEach(func() {
		api = NewAPIWithRouting(testPrefix, NewStaticResolver(""), newTestRouter())
		api.AddResource(SconeTaste{}, SconeResource{})
		rec = httptest.NewRecorder()
		body = strings.NewReader(`
		{
			"data": {
				"attributes": {
					"taste": "smells awful"
				},
				"id": "blubb",
				"type": "Scone-tastes"
			}
		}
		`)
	})

	// check that renaming works, we do not test every single route here, the name variable is used
	// for each route, we just check the 5 basic ones. Marshalling and Unmarshalling is tested with
	// this again too.
	It("FindAll returns 200", func() {
		req, err := http.NewRequest("GET", "/v1/food/sconeTastes", nil)
		Expect(err).ToNot(HaveOccurred())
		api.Handler().ServeHTTP(rec, req)
		Expect(rec.Code).To(Equal(http.StatusOK))
	})

	It("FindOne", func() {
		req, err := http.NewRequest("GET", "/v1/food/sconeTastes/12345", nil)
		Expect(err).ToNot(HaveOccurred())
		api.Handler().ServeHTTP(rec, req)
		Expect(rec.Code).To(Equal(http.StatusOK))
	})

	It("Delete", func() {
		req, err := http.NewRequest("DELETE", "/v1/food/sconeTastes/12345", nil)
		Expect(err).ToNot(HaveOccurred())
		api.Handler().ServeHTTP(rec, req)
		Expect(rec.Code).To(Equal(http.StatusNoContent))
	})

	It("Create", func() {
		req, err := http.NewRequest("POST", "/v1/food/sconeTastes", body)
		Expect(err).ToNot(HaveOccurred())
		api.Handler().ServeHTTP(rec, req)
		// the response is always the one record returned by FindOne, the implementation does not
		// check the ID here and returns something new ...
		Expect(rec.Body.String()).To(MatchJSON(`
		{
			"data": {
				"attributes": {
					"taste": "smells awful"
				},
				"id": "newID",
				"type": "Scone-tastes"
			}
		}
		`))
		Expect(rec.Code).To(Equal(http.StatusCreated))
	})

	It("Update", func() {
		req, err := http.NewRequest("PATCH", "/v1/food/sconeTastes/blubb", body)
		Expect(err).ToNot(HaveOccurred())
		api.Handler().ServeHTTP(rec, req)
		Expect(rec.Body.String()).To(Equal(""))
		Expect(rec.Code).To(Equal(http.StatusNoContent))
	})
})

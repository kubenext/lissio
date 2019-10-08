package errors

import (
	"sort"

	lru "github.com/hashicorp/golang-lru"
)

const maxErrors = 50

type ErrorStore interface {
	List() []InternalError
	Get(string) (InternalError, bool)
	Add(InternalError) (found bool)
	AddError(error) InternalError
}

type errorStore struct {
	recentErrors *lru.Cache
}

var _ ErrorStore = (*errorStore)(nil)

func NewErrorStore() (ErrorStore, error) {
	cache, err := lru.New(maxErrors)
	if err != nil {
		return nil, err
	}

	return &errorStore{
		recentErrors: cache,
	}, nil
}

// Get returns an InternalError and if it was found in the store.
func (e *errorStore) Get(id string) (err InternalError, found bool) {
	v, ok := e.recentErrors.Peek(id)
	if !ok {
		return nil, ok
	}
	return v.(InternalError), ok
}

// Add adds a new InternalError directly in to the store.
func (e *errorStore) Add(intErr InternalError) (found bool) {
	ok, _ := e.recentErrors.ContainsOrAdd(intErr.ID(), intErr)
	return ok
}

func (e *errorStore) AddError(err error) InternalError {
	intErr := convertError(err)
	e.Add(intErr)
	return intErr
}

// List returns a list of all the IntrenalError objects in the store from newest to oldest.
func (e *errorStore) List() []InternalError {
	var intErrList []InternalError
	for _, key := range e.recentErrors.Keys() {
		v, ok := e.recentErrors.Peek(key)
		if !ok {
			continue
		}
		intErrList = append(intErrList, v.(InternalError))
	}
	sort.Slice(intErrList, func(i, j int) bool {
		return intErrList[i].Timestamp().After(intErrList[j].Timestamp())
	})
	return intErrList
}

func convertError(err error) InternalError {
	return NewGenericError(err)
}

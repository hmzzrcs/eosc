package store

import (
	"context"
	"io/ioutil"
	"time"

	"github.com/eolinker/eosc"
	"github.com/ghodss/yaml"
)

type Store struct {
	dispatcher *eosc.StoreEventDispatcher

	data eosc.IUntyped
}

func (s *Store) AddListen(h eosc.IStoreEventHandler) error {
	if s.dispatcher.AddListen(h) {
		list := s.All()
		return h.OnInit(list)
	}

	return nil
}

func NewStore(file string) (eosc.IStore, error) {
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	c := new(Config)

	if err = yaml.Unmarshal(data, c); err != nil {
		return nil, err
	}

	s := &Store{
		data:       eosc.NewUntyped(),
		dispatcher: eosc.NewStoreDispatcher(),
	}
	now := time.Now().Format("2006-01-02 15:04:05")
	err = s.setData(c.Router, "router", now)
	if err != nil {
		return nil, err
	}

	err = s.setData(c.Service, "service", now)
	if err != nil {
		return nil, err
	}

	err = s.setData(c.Upstream, "upstream", now)
	if err != nil {
		return nil, err
	}

	err = s.setData(c.Discovery, "discovery", now)
	if err != nil {
		return nil, err
	}
	err = s.setData(c.Auth, "auth", now)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (s *Store) setData(items []Item, profession string, now string) error {
	for _, r := range items {
		v, err := r.newStoreValue(profession, now)
		if err != nil {
			return err
		}
		s.data.Set(v.Id, v)
	}
	return nil
}

func (s *Store) ReadLock(ctx context.Context) (bool, error) {
	return true, nil
}

func (s *Store) ReadUnLock() error {
	return nil
}

func (s *Store) TryLock(ctx context.Context, expire int) (bool, error) {
	return true, nil
}

func (s *Store) UnLock() error {
	return nil
}

func (s *Store) Initialization() error {
	return nil
}

func (s *Store) All() []eosc.StoreValue {
	list := s.data.List()
	res := make([]eosc.StoreValue, len(list))
	for i, v := range list {
		res[i] = *(v.(*eosc.StoreValue))
	}
	return res
}

func (s *Store) Get(id string) (eosc.StoreValue, bool) {
	if o, has := s.data.Get(id); has {
		return *o.(*eosc.StoreValue), true
	}
	return eosc.StoreValue{}, false
}

func (s *Store) Set(v eosc.StoreValue) error {
	return ErrorReadOnly
}

func (s *Store) Del(id string) error {
	return ErrorReadOnly
}

func (s *Store) ReadOnly() bool {
	return true
}

func (s *Store) GetListener() eosc.IStoreListener {
	return s
}

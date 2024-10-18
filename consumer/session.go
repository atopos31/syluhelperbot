package consumer

import "sync"

type chatSession struct {
	chatIds map[string]string
	lock    sync.Mutex
}

func NewChatSession() *chatSession {
	return &chatSession{chatIds: make(map[string]string), lock: sync.Mutex{}}
}

func (s *chatSession) Get(key string) (string, bool) {
	s.lock.Lock()
	defer s.lock.Unlock()
	val, ok := s.chatIds[key]
	return val, ok
}

func (s *chatSession) Set(key, value string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.chatIds[key] = value
}

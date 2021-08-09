package store

type Notification struct {
	action string
	key    string
	oldVal DataItem
	newVal DataItem
}

type NotifyFunc = func(<-chan Notification)

type Watcher struct {
	chMap map[string][]chan Notification
}

func NewWatcher() *Watcher {
	return &Watcher{
		chMap: make(map[string][]chan Notification, 0),
	}
}

func (w *Watcher) Add(key string, fn NotifyFunc) error {
	_, ok := w.chMap[key]
	if !ok {
		w.chMap[key] = make([]chan Notification, 0)
	}
	c := make(chan Notification)
	w.chMap[key] = append(w.chMap[key], c)

	go fn(c)
	return nil

}

func (w *Watcher) Notify(action string, key string, oldVal, newVal DataItem) error {
	chans, ok := w.chMap[key]
	if ok {
		n := Notification{
			action: action,
			key:    key,
			oldVal: oldVal,
			newVal: newVal,
		}
		for _, ch := range chans {
			ch <- n
		}
		delete(w.chMap, key)
	}
	return nil
}

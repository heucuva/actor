package actor

import (
	"context"
	"reflect"
	"sync"
	"time"

	"github.com/pkg/errors"
)

var (
	// ErrManagerStopped is for when the manager has been asked to stop
	ErrManagerStopped = errors.New("manager stopped")

	// ErrActorNotFound is for when an actor is not found in the manager's list(s)
	ErrActorNotFound = errors.New("actor not found")

	// ErrActorAlreadyAdded is for when an actor is already in the manager's list(s)
	ErrActorAlreadyAdded = errors.New("actor already added")

	// ErrTickIntervalCannotBeZero is for when someone tries to pass a zero value into the tick interval
	// ... that's a special case: see TickEveryFrame()
	ErrTickIntervalCannotBeZero = errors.New("tick interval cannot be zero")
)

// DefaultTickInterval is the default tick interval for actors
const DefaultTickInterval = time.Millisecond * 200

type actorSettings struct {
	tickInterval time.Duration
}

// Option is a function that sets up an option during the AddActor function
type Option func(*actorSettings) error

// TickInterval sets the tick interval for the actor
func TickInterval(interval time.Duration) Option {
	return func(s *actorSettings) error {
		if interval == time.Duration(0) {
			return ErrTickIntervalCannotBeZero
		}

		s.tickInterval = interval
		return nil
	}
}

// TickEveryFrame sets the tick interval for the actor to Every-Frame
// see: Manager.TickFrame()
func TickEveryFrame() Option {
	return func(s *actorSettings) error {
		s.tickInterval = time.Duration(0) // special Every-Frame interval
		return nil
	}
}

type actorList struct {
	list     map[Actor]struct{}
	lastTick time.Time
}

type actorMgrInfo struct {
	tickGroup *time.Ticker
}

// Manager manages actors - and isn't paid enough to deal with their crap
type Manager struct {
	mu                  sync.RWMutex
	actors              map[Actor]actorMgrInfo
	tickGroups          map[*time.Ticker]*actorList
	tickGroupTickers    map[time.Duration]*time.Ticker
	tickGroupsUpdatedCh chan struct{}
	tickStoppedCh       chan struct{}
	tickFrameCh         chan struct{}
	stopping            bool

	cancelFunc context.CancelFunc
}

// NewManager creates a new actor manager
func NewManager() *Manager {
	m := Manager{
		actors:              make(map[Actor]actorMgrInfo),
		tickGroups:          make(map[*time.Ticker]*actorList),
		tickGroupTickers:    make(map[time.Duration]*time.Ticker),
		tickGroupsUpdatedCh: make(chan struct{}, 1),
		tickStoppedCh:       make(chan struct{}, 1),
		tickFrameCh:         make(chan struct{}, 1),
	}

	return &m
}

// Stop stops all actors ticks and shuts down the the manager,
// effectively rendering it useless
func (m *Manager) Stop() {
	if m.stopping {
		return
	}

	m.stopping = true

	defer func() {
		close(m.tickGroupsUpdatedCh)
		close(m.tickStoppedCh)
		close(m.tickFrameCh)
	}()

	if m.cancelFunc != nil {
		m.cancelFunc()
		// wait for done signal from ticker
		<-m.tickStoppedCh
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	actors := m.actors
	m.actors = make(map[Actor]actorMgrInfo)
	m.tickGroupTickers = nil
	m.tickGroups = nil

	for _, a := range actors {
		m.stopActor(a, ErrManagerStopped)
	}
}

func (m *Manager) stopActor(a Actor, reason error) error {
	if err := EndPlay(a, reason); err != nil {
		return err
	}

	return nil
}

// RemoveActor removes the actor from any tick groups and from the managed list of actors
func (m *Manager) RemoveActor(a Actor, reason error) error {
	if err := m.removeActorFromLists(a); err != nil {
		return err
	}

	m.stopActor(a, reason)

	return nil
}

func (m *Manager) removeActorFromLists(a Actor) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	ami, found := m.actors[a]
	if !found {
		return ErrActorNotFound
	}

	delete(m.actors, a)

	ticker := ami.tickGroup

	tg, ok := m.tickGroups[ticker]
	if !ok {
		// not in a tick group
		return nil
	}

	delete(tg.list, a)

	if len(tg.list) == 0 {
		tickerIntv := time.Duration(0)
		for k, v := range m.tickGroupTickers {
			if v == ticker {
				tickerIntv = k
				break
			}
		}

		delete(m.tickGroups, ticker)
		delete(m.tickGroupTickers, tickerIntv)
	}

	return nil
}

// AddActor adds an actor to the various lists internally and sets up the tick interval
func (m *Manager) AddActor(a Actor, opts ...Option) error {
	if m.stopping {
		return ErrManagerStopped
	}

	m.mu.RLock()
	_, found := m.actors[a]
	m.mu.RUnlock()
	if found {
		return ErrActorAlreadyAdded
	}

	s := actorSettings{
		tickInterval: DefaultTickInterval,
	}

	for _, opt := range opts {
		opt(&s)
	}

	if err := BeginPlay(a); err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	updatedGroups := false

	var ticker *time.Ticker
	if s.tickInterval != 0 {
		if tgt, ok := m.tickGroupTickers[s.tickInterval]; ok {
			ticker = tgt
		} else {
			ticker = time.NewTicker(s.tickInterval)
			m.tickGroupTickers[s.tickInterval] = ticker
			updatedGroups = true
		}
	} else {
		// special case for Every-Frame ticking actors
		ticker = nil
	}

	m.actors[a] = actorMgrInfo{
		tickGroup: ticker,
	}

	tg, ok := m.tickGroups[ticker]
	if !ok {
		tg = &actorList{
			list:     make(map[Actor]struct{}),
			lastTick: time.Now(),
		}
		m.tickGroups[ticker] = tg
		updatedGroups = true // just in case
	}

	tg.list[a] = struct{}{}

	if updatedGroups {
		m.tickGroupsUpdatedCh <- struct{}{}
	}
	return nil
}

// TickFrame triggers a single (manually-fired) frame tick for actors attached to the Every-Frame (interval == 0) tick interval
func (m *Manager) TickFrame() error {
	if m.stopping {
		return ErrManagerStopped
	}

	m.tickFrameCh <- struct{}{}
	return nil
}

type waitList struct {
	cases []reflect.SelectCase
	tgs   []*actorList
}

func (m *Manager) generateWaitList(ctx context.Context) *waitList {
	mgr.mu.RLock()
	defer mgr.mu.RUnlock()
	wl := waitList{
		cases: make([]reflect.SelectCase, 0),
		tgs:   make([]*actorList, 0),
	}
	wl.cases = append(wl.cases, reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(ctx.Done()),
	})
	wl.cases = append(wl.cases, reflect.SelectCase{
		Dir:  reflect.SelectRecv,
		Chan: reflect.ValueOf(mgr.tickGroupsUpdatedCh),
	})
	if actors, ok := mgr.tickGroups[nil]; ok {
		// special case for ticker==nil, which is the Every-Frame group
		wl.cases = append(wl.cases, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(mgr.tickFrameCh),
		})
		wl.tgs = append(wl.tgs, actors)
	}
	for ticker, actors := range mgr.tickGroups {
		if ticker == nil {
			continue
		}
		wl.cases = append(wl.cases, reflect.SelectCase{
			Dir:  reflect.SelectRecv,
			Chan: reflect.ValueOf(ticker.C),
		})
		wl.tgs = append(wl.tgs, actors)
	}
	return &wl
}

func (m *Manager) processTickGroups(ctx context.Context) {
	wl := m.generateWaitList(ctx)

	go func() {
		defer m.Stop()

	mainTickLoop:
		for {
			chosen, _, _ := reflect.Select(wl.cases)
			if chosen == 0 {
				// done!
				break mainTickLoop
			} else if chosen == 1 {
				// updated!
				wl = m.generateWaitList(ctx)
				continue mainTickLoop
			}

			tg := wl.tgs[chosen-2]
			// copy the actor list so we can unlock it for other folks
			mgr.mu.RLock()
			actors := make([]Actor, len(tg.list))
			i := 0
			for a := range tg.list {
				actors[i] = a
				i++
			}
			mgr.mu.RUnlock()

			now := time.Now()
			deltaTime := now.Sub(tg.lastTick)
		actorTickLoop:
			for _, a := range actors {
				if canTick, err := WantTick(a); err != nil {
					panic(err)
				} else if !canTick {
					continue actorTickLoop
				}
				if err := Tick(a, deltaTime); err != nil {
					panic(err)
				}
			}
			tg.lastTick = now
		}
		// we're done, signal a stop
		m.tickStoppedCh <- struct{}{}
	}()
}

// StartTicking starts the manager ticking
func (m *Manager) StartTicking(ctx context.Context) {
	ctx, cancel := context.WithCancel(ctx)
	m.cancelFunc = cancel

	m.processTickGroups(ctx)
}

var mgr *Manager

// GetManager returns the global actor manager
func GetManager() *Manager {
	return mgr
}

func init() {
	mgr = NewManager()

	mgr.StartTicking(context.Background())
}

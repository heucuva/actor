package actor

import (
	"reflect"
	"time"

	"github.com/pkg/errors"
)

var (
	// ErrActorSpawn might happen if an actor spawn error has occurred
	ErrActorSpawn = errors.New("could not create Actor")
)

// Actor is an interface to actor-based shenanigans
type Actor interface{}

type spawnActorSettings struct {
	deferredSpawn bool
}

// SpawnActorOption is a function that sets up an option during the SpawnActor/FinishSpawningActor functions
type SpawnActorOption func(*spawnActorSettings) error

// DeferredSpawnActor enables Deferred Spawning for the actor, such that the creator must call FinishSpawningActor to complete the process
func DeferredSpawnActor() SpawnActorOption {
	return func(s *spawnActorSettings) error {
		s.deferredSpawn = true
		return nil
	}
}

// SpawnActor will spawn an actor of the type provided
func SpawnActor(typ reflect.Type, opts ...SpawnActorOption) (Actor, error) {
	s := spawnActorSettings{}
	for _, opt := range opts {
		if err := opt(&s); err != nil {
			return nil, err
		}
	}

	a, ok := reflect.New(typ).Interface().(Actor)
	if !ok {
		return nil, errors.Wrapf(ErrActorSpawn, "unexpected type %v", typ)
	}

	if err := PostSpawnInitialize(a); err != nil {
		return nil, err
	}

	if s.deferredSpawn {
		return a, nil
	}

	if err := FinishSpawningActor(a, opts...); err != nil {
		return nil, err
	}

	return a, nil
}

// FinishSpawningActor finishes the spawning process for actors created with DeferredSpawnActor enabled
func FinishSpawningActor(a Actor, opts ...SpawnActorOption) error {
	s := spawnActorSettings{}
	for _, opt := range opts {
		if err := opt(&s); err != nil {
			return err
		}
	}

	if err := ExecuteConstruction(a); err != nil {
		return err
	}

	if err := OnConstruction(a); err != nil {
		return err
	}

	if err := PostActorConstruction(a); err != nil {
		return err
	}

	if err := PreInitializeComponents(a); err != nil {
		return err
	}

	if err := InitializeComponents(a); err != nil {
		return err
	}

	if err := PostInitializeComponents(a); err != nil {
		return err
	}

	if err := OnActorSpawned(a); err != nil {
		return err
	}

	return nil
}

// PostSpawnInitializeIntf is for actors that want to have PostSpawnInitialize() called right before PostActorCreated() is called
type PostSpawnInitializeIntf interface {
	PostSpawnInitialize() error
}

// ExecuteConstructionIntf is for actors that want to have ExecuteConstruction() called right before OnConstruction() is called
type ExecuteConstructionIntf interface {
	ExecuteConstruction() error
}

// OnConstructionIntf is for actors that want to have OnConstruction() called right before PostActorConstruction() is called
type OnConstructionIntf interface {
	OnConstruction() error
}

// PostActorConstructionIntf is for actors that want to have PostActorConstruction() called right before PreInitializeComponents() is called
type PostActorConstructionIntf interface {
	PostActorConstruction() error
}

// PreInitializeComponentsIntf is for actors that want to have PreInitializeComponents() called right before InitializeComponents() is called
type PreInitializeComponentsIntf interface {
	PreInitializeComponents() error
}

// InitializeComponentsIntf is for actors that want to have InitializeComponents() called right before PostInitializeComponents() is called
type InitializeComponentsIntf interface {
	InitializeComponents() error
}

// PostInitializeComponentsIntf is for actors that want to have InitializeComponents() called right before PostInitializeComponents() is called
type PostInitializeComponentsIntf interface {
	PostInitializeComponents() error
}

// OnActorSpawnedIntf is for actors that want to have OnActorSpawned() called right before BeginPlay() is called
type OnActorSpawnedIntf interface {
	OnActorSpawned() error
}

// BeginPlayIntf is for actors that want to have BeginPlay() called right before the Tick() loop starts
type BeginPlayIntf interface {
	BeginPlay() error
}

// WantTickIntf is for actors that want to announce that they can't Tick() sometimes
type WantTickIntf interface {
	WantTick() (bool, error)
}

// TickIntf is for actors that want to have Tick() called during the Tick() loop (and before EndPlay() is called)
type TickIntf interface {
	Tick(deltaTime time.Duration) error
}

// EndPlayIntf is for actors that want to have EndPlay() called after the Tick() loop ends and before BeginDestroy() is called
type EndPlayIntf interface {
	EndPlay(endPlayReason error) error
}

// BeginDestroyIntf is for actors that want to have BeginDestroy() called right before FinishDestroy() is called
type BeginDestroyIntf interface {
	BeginDestroy() error
}

// FinishDestroyIntf is for actors that want to have FinishDestroy() called
type FinishDestroyIntf interface {
	FinishDestroy() error
}

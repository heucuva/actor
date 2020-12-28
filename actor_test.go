package actor_test

import (
	"math/rand"
	"reflect"
	"testing"

	"github.com/heucuva/actor"
)

type spawnActorTest struct {
	counter int

	hitPostSpawnInitialize     int
	hitExecuteConstruction     int
	hitOnConstruction          int
	hitPostActorConstruction   int
	hitPreInitializeComponents int
	hitInitializeComponents    int
	hitOnActorSpawned          int
}

func (a *spawnActorTest) PostSpawnInitialize() error {
	a.counter++
	a.hitPostSpawnInitialize = a.counter
	return nil
}

func (a *spawnActorTest) ExecuteConstruction() error {
	a.counter++
	a.hitExecuteConstruction = a.counter
	return nil
}

func (a *spawnActorTest) OnConstruction() error {
	a.counter++
	a.hitOnConstruction = a.counter
	return nil
}

func (a *spawnActorTest) PostActorConstruction() error {
	a.counter++
	a.hitPostActorConstruction = a.counter
	return nil
}

func (a *spawnActorTest) PreInitializeComponents() error {
	a.counter++
	a.hitPreInitializeComponents = a.counter
	return nil
}

func (a *spawnActorTest) InitializeComponents() error {
	a.counter++
	a.hitInitializeComponents = a.counter
	return nil
}

func (a *spawnActorTest) OnActorSpawned() error {
	a.counter++
	a.hitOnActorSpawned = a.counter
	return nil
}

var defaultSpawnActorTest = spawnActorTest{}

func TestSpawnActor(t *testing.T) {
	act, err := actor.SpawnActor(reflect.TypeOf(defaultSpawnActorTest))
	if err != nil {
		t.Fatal(err)
	}

	c := 0

	a, ok := act.(*spawnActorTest)
	if !ok {
		t.Fatalf("expected spawnActorTest, got %v", reflect.TypeOf(a))
	}

	c++
	if a.hitPostSpawnInitialize == 0 {
		t.Fatal("PostSpawnInitialize not triggered")
	} else if a.hitPostSpawnInitialize != c {
		t.Fatalf("PostSpawnInitialize triggered at the wrong time - expected %d, got %v", c, a.hitPostSpawnInitialize)
	}

	c++
	if a.hitExecuteConstruction == 0 {
		t.Fatal("ExecuteConstruction not triggered")
	} else if a.hitExecuteConstruction != c {
		t.Fatalf("ExecuteConstruction triggered at the wrong time - expected %d, got %v", c, a.hitExecuteConstruction)
	}

	c++
	if a.hitOnConstruction == 0 {
		t.Fatal("OnConstruction not triggered")
	} else if a.hitOnConstruction != c {
		t.Fatalf("OnConstruction triggered at the wrong time - expected %d, got %v", c, a.hitOnConstruction)
	}

	c++
	if a.hitPostActorConstruction == 0 {
		t.Fatal("PostActorConstruction not triggered")
	} else if a.hitPostActorConstruction != c {
		t.Fatalf("PostActorConstruction triggered at the wrong time - expected %d, got %v", c, a.hitPostActorConstruction)
	}

	c++
	if a.hitPreInitializeComponents == 0 {
		t.Fatal("PreInitializeComponents not triggered")
	} else if a.hitPreInitializeComponents != c {
		t.Fatalf("PreInitializeComponents triggered at the wrong time - expected %d, got %v", c, a.hitPreInitializeComponents)
	}

	c++
	if a.hitInitializeComponents == 0 {
		t.Fatal("InitializeComponents not triggered")
	} else if a.hitInitializeComponents != c {
		t.Fatalf("InitializeComponents triggered at the wrong time - expected %d, got %v", c, a.hitInitializeComponents)
	}

	c++
	if a.hitOnActorSpawned == 0 {
		t.Fatal("OnActorSpawned not triggered")
	} else if a.hitOnActorSpawned != c {
		t.Fatalf("OnActorSpawned triggered at the wrong time - expected %d, got %v", c, a.hitOnActorSpawned)
	}
}

func TestSpawnActorDeferred(t *testing.T) {
	act, err := actor.SpawnActor(reflect.TypeOf(defaultSpawnActorTest), actor.DeferredSpawnActor())
	if err != nil {
		t.Fatal(err)
	}

	a, ok := act.(*spawnActorTest)
	if !ok {
		t.Fatalf("expected spawnActorTest, got %v", reflect.TypeOf(a))
	}

	if a.hitPostSpawnInitialize == 0 {
		t.Fatal("PostSpawnInitialize not triggered")
	} else if a.hitPostSpawnInitialize != 1 {
		t.Fatalf("PostSpawnInitialize triggered at the wrong time - expected 1, got %v", a.hitPostSpawnInitialize)
	}

	if a.hitExecuteConstruction != 0 {
		t.Fatalf("ExecuteConstruction triggered at the wrong time - expected 0, got %v", a.hitExecuteConstruction)
	}

	if a.hitOnConstruction != 0 {
		t.Fatalf("OnConstruction triggered at the wrong time - expected 0, got %v", a.hitOnConstruction)
	}

	if a.hitPostActorConstruction != 0 {
		t.Fatalf("PostActorConstruction triggered at the wrong time - expected 0, got %v", a.hitPostActorConstruction)
	}

	if a.hitPreInitializeComponents != 0 {
		t.Fatalf("PreInitializeComponents triggered at the wrong time - expected 0, got %v", a.hitPreInitializeComponents)
	}

	if a.hitInitializeComponents != 0 {
		t.Fatalf("InitializeComponents triggered at the wrong time - expected 0, got %v", a.hitInitializeComponents)
	}

	if a.hitOnActorSpawned != 0 {
		t.Fatalf("OnActorSpawned triggered at the wrong time - expected 0, got %v", a.hitOnActorSpawned)
	}
}

func TestSpawnActorDeferredFinish(t *testing.T) {
	opts := []actor.SpawnActorOption{
		actor.DeferredSpawnActor(),
	}

	act, err := actor.SpawnActor(reflect.TypeOf(defaultSpawnActorTest), opts...)
	if err != nil {
		t.Fatal(err)
	}

	a, ok := act.(*spawnActorTest)
	if !ok {
		t.Fatalf("expected spawnActorTest, got %v", reflect.TypeOf(a))
	}

	if a.hitPostSpawnInitialize == 0 {
		t.Fatal("PostSpawnInitialize not triggered")
	} else if a.hitPostSpawnInitialize != 1 {
		t.Fatalf("PostSpawnInitialize triggered at the wrong time - expected 1, got %v", a.hitPostSpawnInitialize)
	}

	if a.hitExecuteConstruction != 0 {
		t.Fatalf("ExecuteConstruction triggered at the wrong time - expected 0, got %v", a.hitExecuteConstruction)
	}

	if a.hitOnConstruction != 0 {
		t.Fatalf("OnConstruction triggered at the wrong time - expected 0, got %v", a.hitOnConstruction)
	}

	if a.hitPostActorConstruction != 0 {
		t.Fatalf("PostActorConstruction triggered at the wrong time - expected 0, got %v", a.hitPostActorConstruction)
	}

	if a.hitPreInitializeComponents != 0 {
		t.Fatalf("PreInitializeComponents triggered at the wrong time - expected 0, got %v", a.hitPreInitializeComponents)
	}

	if a.hitInitializeComponents != 0 {
		t.Fatalf("InitializeComponents triggered at the wrong time - expected 0, got %v", a.hitInitializeComponents)
	}

	if a.hitOnActorSpawned != 0 {
		t.Fatalf("OnActorSpawned triggered at the wrong time - expected 0, got %v", a.hitOnActorSpawned)
	}

	// create a random base for count so that we can test expectations
	countBase := int(rand.Int31()>>2) + 10
	a.counter = countBase
	c := countBase

	if err := actor.FinishSpawningActor(a, opts...); err != nil {
		t.Fatal(err)
	}

	if a.hitPostSpawnInitialize != 1 {
		t.Fatalf("PostSpawnInitialize triggered at the wrong time - expected 1, got %v", a.hitPostSpawnInitialize)
	}

	c++
	if a.hitExecuteConstruction == 0 {
		t.Fatal("ExecuteConstruction not triggered")
	} else if a.hitExecuteConstruction != c {
		t.Fatalf("ExecuteConstruction triggered at the wrong time - expected %d, got %v", c, a.hitExecuteConstruction)
	}

	c++
	if a.hitOnConstruction == 0 {
		t.Fatal("OnConstruction not triggered")
	} else if a.hitOnConstruction != c {
		t.Fatalf("OnConstruction triggered at the wrong time - expected %d, got %v", c, a.hitOnConstruction)
	}

	c++
	if a.hitPostActorConstruction == 0 {
		t.Fatal("PostActorConstruction not triggered")
	} else if a.hitPostActorConstruction != c {
		t.Fatalf("PostActorConstruction triggered at the wrong time - expected %d, got %v", c, a.hitPostActorConstruction)
	}

	c++
	if a.hitPreInitializeComponents == 0 {
		t.Fatal("PreInitializeComponents not triggered")
	} else if a.hitPreInitializeComponents != c {
		t.Fatalf("PreInitializeComponents triggered at the wrong time - expected %d, got %v", c, a.hitPreInitializeComponents)
	}

	c++
	if a.hitInitializeComponents == 0 {
		t.Fatal("InitializeComponents not triggered")
	} else if a.hitInitializeComponents != c {
		t.Fatalf("InitializeComponents triggered at the wrong time - expected %d, got %v", c, a.hitInitializeComponents)
	}

	c++
	if a.hitOnActorSpawned == 0 {
		t.Fatal("OnActorSpawned not triggered")
	} else if a.hitOnActorSpawned != c {
		t.Fatalf("OnActorSpawned triggered at the wrong time - expected %d, got %v", c, a.hitOnActorSpawned)
	}
}

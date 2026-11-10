package bulkops

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"
	"time"
)

func newCancellationTestManager() *Manager {
	return NewManager(nil, nil, nil, "datacenter")
}

func (m *Manager) queuedPayloadExists(op *Operation) bool {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	res, exists := m.resources[op.ResourceType]
	if !exists {
		return false
	}
	switch op.OperationType {
	case "PUT":
		_, found := res.Put[op.StoreKey]
		return found
	case "PATCH":
		_, found := res.Patch[op.StoreKey]
		return found
	case "DELETE":
		for _, name := range res.Delete {
			if name == op.StoreKey {
				return true
			}
		}
	}
	return false
}

func (m *Manager) operationByID(id string) *Operation {
	m.operationMutex.Lock()
	defer m.operationMutex.Unlock()
	op := m.pendingOperations[id]
	if op == nil {
		return nil
	}
	copied := *op
	return &copied
}

func waitResult(t *testing.T, m *Manager, id string) error {
	t.Helper()
	done := make(chan error, 1)
	go func() { done <- m.WaitForOperation(context.Background(), id, time.Second) }()
	select {
	case err := <-done:
		return err
	case <-time.After(5 * time.Second):
		t.Fatalf("waiter for %s did not return", id)
		return nil
	}
}

func TestFailAllPendingOperationsRemovesQueuedPayloads(t *testing.T) {
	m := newCancellationTestManager()
	ctx := context.Background()

	put := m.AddPut(ctx, "tenant", "t1", map[string]interface{}{"enable": true})
	patch := m.AddPatch(ctx, "service", "s1", map[string]interface{}{"enable": false})
	del := m.AddDelete(ctx, "bundle", "b1")
	aclPut := m.AddPut(ctx, "acl", "a1", map[string]interface{}{}, map[string]string{"ip_version": "4"})

	m.FailAllPendingOperations(ctx, errors.New("tenant PUT failed"))

	for _, id := range []string{put, patch, del, aclPut} {
		if err := waitResult(t, m, id); err == nil {
			t.Errorf("operation %s: expected a cancellation error", id)
		}
		op := m.operationByID(id)
		if op.Status != OperationFailed {
			t.Errorf("operation %s: status %v, want failed", id, op.Status)
		}
		if m.queuedPayloadExists(op) {
			t.Errorf("operation %s: payload still queued after cancellation", id)
		}
	}
	if m.hasPendingOperations() {
		t.Error("queue still holds payloads after cancellation")
	}
	m.mutex.Lock()
	params, origNames := len(m.resourceHeaderParams), len(m.resourceOriginalNames)
	m.mutex.Unlock()
	if params != 0 || origNames != 0 {
		t.Errorf("header-split metadata left behind: %d params, %d original names", params, origNames)
	}
}

func TestFailAllPendingOperationsLeavesExecutingOperationsAlone(t *testing.T) {
	m := newCancellationTestManager()
	ctx := context.Background()

	executing := m.AddPut(ctx, "tenant", "t1", map[string]interface{}{})
	m.createExtractor("tenant", "PUT")()
	m.markOperationsAsExecuting("tenant", "PUT", []string{"t1"})

	m.FailAllPendingOperations(ctx, errors.New("other failure"))

	if op := m.operationByID(executing); op.Status != OperationExecuting {
		t.Fatalf("executing operation was changed by cancellation: status %v", op.Status)
	}
}

func TestAddAndCancelInterleavingIsAtomic(t *testing.T) {
	m := newCancellationTestManager()
	ctx := context.Background()

	reached := make(chan struct{})
	resume := make(chan struct{})
	afterStoreHook = func() {
		close(reached)
		<-resume
	}
	started := make(chan struct{})
	cancelStartedHook = func() { close(started) }
	t.Cleanup(func() { afterStoreHook, cancelStartedHook = nil, nil })

	added := make(chan string, 1)
	go func() { added <- m.AddPut(ctx, "tenant", "t1", map[string]interface{}{}) }()
	<-reached
	afterStoreHook = nil

	if m.mutex.TryLock() {
		m.mutex.Unlock()
		close(resume)
		t.Fatal("the queue lock is free between storing the payload and registering the operation; a cancellation could interleave")
	}

	cancelled := make(chan struct{})
	go func() {
		m.FailAllPendingOperations(ctx, errors.New("failure during add"))
		close(cancelled)
	}()
	<-started

	close(resume)
	id := <-added
	<-cancelled

	op := m.operationByID(id)
	if op.Status != OperationFailed {
		t.Fatalf("status %v, want failed: the cancellation must cover the completed add", op.Status)
	}
	if m.queuedPayloadExists(op) {
		t.Fatal("payload still queued for a cancelled operation")
	}
	if err := waitResult(t, m, id); err == nil {
		t.Fatal("waiter returned success for a cancelled operation")
	}
}

func TestAddAndCancelRaceStress(t *testing.T) {
	m := newCancellationTestManager()
	ctx := context.Background()
	types := []string{"tenant", "service", "bundle", "switchpoint"}

	var mu sync.Mutex
	var ids []string
	var wg sync.WaitGroup
	stop := make(chan struct{})

	cancelDone := make(chan struct{})
	go func() {
		defer close(cancelDone)
		for {
			select {
			case <-stop:
				return
			default:
				m.FailAllPendingOperations(ctx, errors.New("stress failure"))
			}
		}
	}()

	for worker := 0; worker < 8; worker++ {
		wg.Add(1)
		go func(worker int) {
			defer wg.Done()
			for i := 0; i < 200; i++ {
				rt := types[(worker+i)%len(types)]
				name := fmt.Sprintf("w%d-%d", worker, i)
				var id string
				switch i % 3 {
				case 0:
					id = m.AddPut(ctx, rt, name, map[string]interface{}{})
				case 1:
					id = m.AddPatch(ctx, rt, name, map[string]interface{}{})
				default:
					id = m.AddDelete(ctx, rt, name)
				}
				mu.Lock()
				ids = append(ids, id)
				mu.Unlock()
			}
		}(worker)
	}
	wg.Wait()
	close(stop)
	<-cancelDone

	for _, id := range ids {
		op := m.operationByID(id)
		if op == nil {
			t.Fatalf("operation %s not registered", id)
		}
		switch op.Status {
		case OperationPending:
			if !m.queuedPayloadExists(op) {
				t.Fatalf("operation %s is pending without a queued payload", id)
			}
		case OperationFailed:
			if m.queuedPayloadExists(op) {
				t.Fatalf("operation %s was cancelled but its payload is still queued", id)
			}
			if err := waitResult(t, m, id); err == nil {
				t.Fatalf("waiter for cancelled operation %s returned success", id)
			}
		default:
			t.Fatalf("operation %s has unexpected status %v", id, op.Status)
		}
	}
}

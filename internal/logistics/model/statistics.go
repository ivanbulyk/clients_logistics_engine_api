package model

import (
    "sync"
    "time"
)

// Statistics about world operations
type Statistics struct {
    Operation []*Operation
    ExecTime  time.Time
}

// Operation kind
type Operation struct {
    Name string
    A    uint64
    B    uint64

    sync.Mutex
}

// AddA safe incrementation
func (o *Operation) AddA() {
    o.Lock()
    defer o.Unlock()
    o.A++
}

// AddB safe incrementation
func (o *Operation) AddB() {
    o.Lock()
    defer o.Unlock()
    o.B++
}

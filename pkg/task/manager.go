package task

import (
	. "m7s.live/m7s/v5/pkg/util"
)

type ManagerItem[K comparable] interface {
	ITask
	GetKey() K
}

type Manager[K comparable, T ManagerItem[K]] struct {
	Work
	Collection[K, T]
}

func (m *Manager[K, T]) Add(ctx T, opt ...any) *Task {
	ctx.OnStart(func() {
		m.Collection.Add(ctx)
		if m.Logger != nil {
			m.Logger.Debug("add", "key", ctx.GetKey(), "count", m.Length)
		}
	})
	ctx.OnDispose(func() {
		m.Remove(ctx)
		if m.Logger != nil {
			m.Logger.Debug("remove", "key", ctx.GetKey(), "count", m.Length)
		}
	})
	return m.AddTask(ctx, opt...)
}

// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/suyuan32/simple-admin-job/ent/predicate"
	"github.com/suyuan32/simple-admin-job/ent/task"
	"github.com/suyuan32/simple-admin-job/ent/tasklog"
)

// TaskUpdate is the builder for updating Task entities.
type TaskUpdate struct {
	config
	hooks    []Hook
	mutation *TaskMutation
}

// Where appends a list predicates to the TaskUpdate builder.
func (tu *TaskUpdate) Where(ps ...predicate.Task) *TaskUpdate {
	tu.mutation.Where(ps...)
	return tu
}

// SetUpdatedAt sets the "updated_at" field.
func (tu *TaskUpdate) SetUpdatedAt(t time.Time) *TaskUpdate {
	tu.mutation.SetUpdatedAt(t)
	return tu
}

// SetStatus sets the "status" field.
func (tu *TaskUpdate) SetStatus(u uint8) *TaskUpdate {
	tu.mutation.ResetStatus()
	tu.mutation.SetStatus(u)
	return tu
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (tu *TaskUpdate) SetNillableStatus(u *uint8) *TaskUpdate {
	if u != nil {
		tu.SetStatus(*u)
	}
	return tu
}

// AddStatus adds u to the "status" field.
func (tu *TaskUpdate) AddStatus(u int8) *TaskUpdate {
	tu.mutation.AddStatus(u)
	return tu
}

// ClearStatus clears the value of the "status" field.
func (tu *TaskUpdate) ClearStatus() *TaskUpdate {
	tu.mutation.ClearStatus()
	return tu
}

// SetName sets the "name" field.
func (tu *TaskUpdate) SetName(s string) *TaskUpdate {
	tu.mutation.SetName(s)
	return tu
}

// SetNillableName sets the "name" field if the given value is not nil.
func (tu *TaskUpdate) SetNillableName(s *string) *TaskUpdate {
	if s != nil {
		tu.SetName(*s)
	}
	return tu
}

// SetTaskGroup sets the "task_group" field.
func (tu *TaskUpdate) SetTaskGroup(s string) *TaskUpdate {
	tu.mutation.SetTaskGroup(s)
	return tu
}

// SetNillableTaskGroup sets the "task_group" field if the given value is not nil.
func (tu *TaskUpdate) SetNillableTaskGroup(s *string) *TaskUpdate {
	if s != nil {
		tu.SetTaskGroup(*s)
	}
	return tu
}

// SetCronExpression sets the "cron_expression" field.
func (tu *TaskUpdate) SetCronExpression(s string) *TaskUpdate {
	tu.mutation.SetCronExpression(s)
	return tu
}

// SetNillableCronExpression sets the "cron_expression" field if the given value is not nil.
func (tu *TaskUpdate) SetNillableCronExpression(s *string) *TaskUpdate {
	if s != nil {
		tu.SetCronExpression(*s)
	}
	return tu
}

// SetPattern sets the "pattern" field.
func (tu *TaskUpdate) SetPattern(s string) *TaskUpdate {
	tu.mutation.SetPattern(s)
	return tu
}

// SetNillablePattern sets the "pattern" field if the given value is not nil.
func (tu *TaskUpdate) SetNillablePattern(s *string) *TaskUpdate {
	if s != nil {
		tu.SetPattern(*s)
	}
	return tu
}

// SetPayload sets the "payload" field.
func (tu *TaskUpdate) SetPayload(s string) *TaskUpdate {
	tu.mutation.SetPayload(s)
	return tu
}

// SetNillablePayload sets the "payload" field if the given value is not nil.
func (tu *TaskUpdate) SetNillablePayload(s *string) *TaskUpdate {
	if s != nil {
		tu.SetPayload(*s)
	}
	return tu
}

// AddTaskLogIDs adds the "task_logs" edge to the TaskLog entity by IDs.
func (tu *TaskUpdate) AddTaskLogIDs(ids ...uint64) *TaskUpdate {
	tu.mutation.AddTaskLogIDs(ids...)
	return tu
}

// AddTaskLogs adds the "task_logs" edges to the TaskLog entity.
func (tu *TaskUpdate) AddTaskLogs(t ...*TaskLog) *TaskUpdate {
	ids := make([]uint64, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return tu.AddTaskLogIDs(ids...)
}

// Mutation returns the TaskMutation object of the builder.
func (tu *TaskUpdate) Mutation() *TaskMutation {
	return tu.mutation
}

// ClearTaskLogs clears all "task_logs" edges to the TaskLog entity.
func (tu *TaskUpdate) ClearTaskLogs() *TaskUpdate {
	tu.mutation.ClearTaskLogs()
	return tu
}

// RemoveTaskLogIDs removes the "task_logs" edge to TaskLog entities by IDs.
func (tu *TaskUpdate) RemoveTaskLogIDs(ids ...uint64) *TaskUpdate {
	tu.mutation.RemoveTaskLogIDs(ids...)
	return tu
}

// RemoveTaskLogs removes "task_logs" edges to TaskLog entities.
func (tu *TaskUpdate) RemoveTaskLogs(t ...*TaskLog) *TaskUpdate {
	ids := make([]uint64, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return tu.RemoveTaskLogIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (tu *TaskUpdate) Save(ctx context.Context) (int, error) {
	tu.defaults()
	return withHooks(ctx, tu.sqlSave, tu.mutation, tu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (tu *TaskUpdate) SaveX(ctx context.Context) int {
	affected, err := tu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (tu *TaskUpdate) Exec(ctx context.Context) error {
	_, err := tu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tu *TaskUpdate) ExecX(ctx context.Context) {
	if err := tu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (tu *TaskUpdate) defaults() {
	if _, ok := tu.mutation.UpdatedAt(); !ok {
		v := task.UpdateDefaultUpdatedAt()
		tu.mutation.SetUpdatedAt(v)
	}
}

func (tu *TaskUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(task.Table, task.Columns, sqlgraph.NewFieldSpec(task.FieldID, field.TypeUint64))
	if ps := tu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := tu.mutation.UpdatedAt(); ok {
		_spec.SetField(task.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := tu.mutation.Status(); ok {
		_spec.SetField(task.FieldStatus, field.TypeUint8, value)
	}
	if value, ok := tu.mutation.AddedStatus(); ok {
		_spec.AddField(task.FieldStatus, field.TypeUint8, value)
	}
	if tu.mutation.StatusCleared() {
		_spec.ClearField(task.FieldStatus, field.TypeUint8)
	}
	if value, ok := tu.mutation.Name(); ok {
		_spec.SetField(task.FieldName, field.TypeString, value)
	}
	if value, ok := tu.mutation.TaskGroup(); ok {
		_spec.SetField(task.FieldTaskGroup, field.TypeString, value)
	}
	if value, ok := tu.mutation.CronExpression(); ok {
		_spec.SetField(task.FieldCronExpression, field.TypeString, value)
	}
	if value, ok := tu.mutation.Pattern(); ok {
		_spec.SetField(task.FieldPattern, field.TypeString, value)
	}
	if value, ok := tu.mutation.Payload(); ok {
		_spec.SetField(task.FieldPayload, field.TypeString, value)
	}
	if tu.mutation.TaskLogsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   task.TaskLogsTable,
			Columns: []string{task.TaskLogsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(tasklog.FieldID, field.TypeUint64),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tu.mutation.RemovedTaskLogsIDs(); len(nodes) > 0 && !tu.mutation.TaskLogsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   task.TaskLogsTable,
			Columns: []string{task.TaskLogsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(tasklog.FieldID, field.TypeUint64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tu.mutation.TaskLogsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   task.TaskLogsTable,
			Columns: []string{task.TaskLogsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(tasklog.FieldID, field.TypeUint64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, tu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{task.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	tu.mutation.done = true
	return n, nil
}

// TaskUpdateOne is the builder for updating a single Task entity.
type TaskUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *TaskMutation
}

// SetUpdatedAt sets the "updated_at" field.
func (tuo *TaskUpdateOne) SetUpdatedAt(t time.Time) *TaskUpdateOne {
	tuo.mutation.SetUpdatedAt(t)
	return tuo
}

// SetStatus sets the "status" field.
func (tuo *TaskUpdateOne) SetStatus(u uint8) *TaskUpdateOne {
	tuo.mutation.ResetStatus()
	tuo.mutation.SetStatus(u)
	return tuo
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (tuo *TaskUpdateOne) SetNillableStatus(u *uint8) *TaskUpdateOne {
	if u != nil {
		tuo.SetStatus(*u)
	}
	return tuo
}

// AddStatus adds u to the "status" field.
func (tuo *TaskUpdateOne) AddStatus(u int8) *TaskUpdateOne {
	tuo.mutation.AddStatus(u)
	return tuo
}

// ClearStatus clears the value of the "status" field.
func (tuo *TaskUpdateOne) ClearStatus() *TaskUpdateOne {
	tuo.mutation.ClearStatus()
	return tuo
}

// SetName sets the "name" field.
func (tuo *TaskUpdateOne) SetName(s string) *TaskUpdateOne {
	tuo.mutation.SetName(s)
	return tuo
}

// SetNillableName sets the "name" field if the given value is not nil.
func (tuo *TaskUpdateOne) SetNillableName(s *string) *TaskUpdateOne {
	if s != nil {
		tuo.SetName(*s)
	}
	return tuo
}

// SetTaskGroup sets the "task_group" field.
func (tuo *TaskUpdateOne) SetTaskGroup(s string) *TaskUpdateOne {
	tuo.mutation.SetTaskGroup(s)
	return tuo
}

// SetNillableTaskGroup sets the "task_group" field if the given value is not nil.
func (tuo *TaskUpdateOne) SetNillableTaskGroup(s *string) *TaskUpdateOne {
	if s != nil {
		tuo.SetTaskGroup(*s)
	}
	return tuo
}

// SetCronExpression sets the "cron_expression" field.
func (tuo *TaskUpdateOne) SetCronExpression(s string) *TaskUpdateOne {
	tuo.mutation.SetCronExpression(s)
	return tuo
}

// SetNillableCronExpression sets the "cron_expression" field if the given value is not nil.
func (tuo *TaskUpdateOne) SetNillableCronExpression(s *string) *TaskUpdateOne {
	if s != nil {
		tuo.SetCronExpression(*s)
	}
	return tuo
}

// SetPattern sets the "pattern" field.
func (tuo *TaskUpdateOne) SetPattern(s string) *TaskUpdateOne {
	tuo.mutation.SetPattern(s)
	return tuo
}

// SetNillablePattern sets the "pattern" field if the given value is not nil.
func (tuo *TaskUpdateOne) SetNillablePattern(s *string) *TaskUpdateOne {
	if s != nil {
		tuo.SetPattern(*s)
	}
	return tuo
}

// SetPayload sets the "payload" field.
func (tuo *TaskUpdateOne) SetPayload(s string) *TaskUpdateOne {
	tuo.mutation.SetPayload(s)
	return tuo
}

// SetNillablePayload sets the "payload" field if the given value is not nil.
func (tuo *TaskUpdateOne) SetNillablePayload(s *string) *TaskUpdateOne {
	if s != nil {
		tuo.SetPayload(*s)
	}
	return tuo
}

// AddTaskLogIDs adds the "task_logs" edge to the TaskLog entity by IDs.
func (tuo *TaskUpdateOne) AddTaskLogIDs(ids ...uint64) *TaskUpdateOne {
	tuo.mutation.AddTaskLogIDs(ids...)
	return tuo
}

// AddTaskLogs adds the "task_logs" edges to the TaskLog entity.
func (tuo *TaskUpdateOne) AddTaskLogs(t ...*TaskLog) *TaskUpdateOne {
	ids := make([]uint64, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return tuo.AddTaskLogIDs(ids...)
}

// Mutation returns the TaskMutation object of the builder.
func (tuo *TaskUpdateOne) Mutation() *TaskMutation {
	return tuo.mutation
}

// ClearTaskLogs clears all "task_logs" edges to the TaskLog entity.
func (tuo *TaskUpdateOne) ClearTaskLogs() *TaskUpdateOne {
	tuo.mutation.ClearTaskLogs()
	return tuo
}

// RemoveTaskLogIDs removes the "task_logs" edge to TaskLog entities by IDs.
func (tuo *TaskUpdateOne) RemoveTaskLogIDs(ids ...uint64) *TaskUpdateOne {
	tuo.mutation.RemoveTaskLogIDs(ids...)
	return tuo
}

// RemoveTaskLogs removes "task_logs" edges to TaskLog entities.
func (tuo *TaskUpdateOne) RemoveTaskLogs(t ...*TaskLog) *TaskUpdateOne {
	ids := make([]uint64, len(t))
	for i := range t {
		ids[i] = t[i].ID
	}
	return tuo.RemoveTaskLogIDs(ids...)
}

// Where appends a list predicates to the TaskUpdate builder.
func (tuo *TaskUpdateOne) Where(ps ...predicate.Task) *TaskUpdateOne {
	tuo.mutation.Where(ps...)
	return tuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (tuo *TaskUpdateOne) Select(field string, fields ...string) *TaskUpdateOne {
	tuo.fields = append([]string{field}, fields...)
	return tuo
}

// Save executes the query and returns the updated Task entity.
func (tuo *TaskUpdateOne) Save(ctx context.Context) (*Task, error) {
	tuo.defaults()
	return withHooks(ctx, tuo.sqlSave, tuo.mutation, tuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (tuo *TaskUpdateOne) SaveX(ctx context.Context) *Task {
	node, err := tuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (tuo *TaskUpdateOne) Exec(ctx context.Context) error {
	_, err := tuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tuo *TaskUpdateOne) ExecX(ctx context.Context) {
	if err := tuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (tuo *TaskUpdateOne) defaults() {
	if _, ok := tuo.mutation.UpdatedAt(); !ok {
		v := task.UpdateDefaultUpdatedAt()
		tuo.mutation.SetUpdatedAt(v)
	}
}

func (tuo *TaskUpdateOne) sqlSave(ctx context.Context) (_node *Task, err error) {
	_spec := sqlgraph.NewUpdateSpec(task.Table, task.Columns, sqlgraph.NewFieldSpec(task.FieldID, field.TypeUint64))
	id, ok := tuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "Task.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := tuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, task.FieldID)
		for _, f := range fields {
			if !task.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != task.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := tuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := tuo.mutation.UpdatedAt(); ok {
		_spec.SetField(task.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := tuo.mutation.Status(); ok {
		_spec.SetField(task.FieldStatus, field.TypeUint8, value)
	}
	if value, ok := tuo.mutation.AddedStatus(); ok {
		_spec.AddField(task.FieldStatus, field.TypeUint8, value)
	}
	if tuo.mutation.StatusCleared() {
		_spec.ClearField(task.FieldStatus, field.TypeUint8)
	}
	if value, ok := tuo.mutation.Name(); ok {
		_spec.SetField(task.FieldName, field.TypeString, value)
	}
	if value, ok := tuo.mutation.TaskGroup(); ok {
		_spec.SetField(task.FieldTaskGroup, field.TypeString, value)
	}
	if value, ok := tuo.mutation.CronExpression(); ok {
		_spec.SetField(task.FieldCronExpression, field.TypeString, value)
	}
	if value, ok := tuo.mutation.Pattern(); ok {
		_spec.SetField(task.FieldPattern, field.TypeString, value)
	}
	if value, ok := tuo.mutation.Payload(); ok {
		_spec.SetField(task.FieldPayload, field.TypeString, value)
	}
	if tuo.mutation.TaskLogsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   task.TaskLogsTable,
			Columns: []string{task.TaskLogsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(tasklog.FieldID, field.TypeUint64),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tuo.mutation.RemovedTaskLogsIDs(); len(nodes) > 0 && !tuo.mutation.TaskLogsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   task.TaskLogsTable,
			Columns: []string{task.TaskLogsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(tasklog.FieldID, field.TypeUint64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := tuo.mutation.TaskLogsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   task.TaskLogsTable,
			Columns: []string{task.TaskLogsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(tasklog.FieldID, field.TypeUint64),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &Task{config: tuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, tuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{task.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	tuo.mutation.done = true
	return _node, nil
}

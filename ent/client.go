// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"log"
	"reflect"

	"github.com/suyuan32/simple-admin-job/ent/migrate"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/suyuan32/simple-admin-job/ent/task"
	"github.com/suyuan32/simple-admin-job/ent/tasklog"
)

// Client is the client that holds all ent builders.
type Client struct {
	config
	// Schema is the client for creating, migrating and dropping schema.
	Schema *migrate.Schema
	// Task is the client for interacting with the Task builders.
	Task *TaskClient
	// TaskLog is the client for interacting with the TaskLog builders.
	TaskLog *TaskLogClient
}

// NewClient creates a new client configured with the given options.
func NewClient(opts ...Option) *Client {
	cfg := config{log: log.Println, hooks: &hooks{}, inters: &inters{}}
	cfg.options(opts...)
	client := &Client{config: cfg}
	client.init()
	return client
}

func (c *Client) init() {
	c.Schema = migrate.NewSchema(c.driver)
	c.Task = NewTaskClient(c.config)
	c.TaskLog = NewTaskLogClient(c.config)
}

type (
	// config is the configuration for the client and its builder.
	config struct {
		// driver used for executing database requests.
		driver dialect.Driver
		// debug enable a debug logging.
		debug bool
		// log used for logging on debug mode.
		log func(...any)
		// hooks to execute on mutations.
		hooks *hooks
		// interceptors to execute on queries.
		inters *inters
	}
	// Option function to configure the client.
	Option func(*config)
)

// options applies the options on the config object.
func (c *config) options(opts ...Option) {
	for _, opt := range opts {
		opt(c)
	}
	if c.debug {
		c.driver = dialect.Debug(c.driver, c.log)
	}
}

// Debug enables debug logging on the ent.Driver.
func Debug() Option {
	return func(c *config) {
		c.debug = true
	}
}

// Log sets the logging function for debug mode.
func Log(fn func(...any)) Option {
	return func(c *config) {
		c.log = fn
	}
}

// Driver configures the client driver.
func Driver(driver dialect.Driver) Option {
	return func(c *config) {
		c.driver = driver
	}
}

// Open opens a database/sql.DB specified by the driver name and
// the data source name, and returns a new client attached to it.
// Optional parameters can be added for configuring the client.
func Open(driverName, dataSourceName string, options ...Option) (*Client, error) {
	switch driverName {
	case dialect.MySQL, dialect.Postgres, dialect.SQLite:
		drv, err := sql.Open(driverName, dataSourceName)
		if err != nil {
			return nil, err
		}
		return NewClient(append(options, Driver(drv))...), nil
	default:
		return nil, fmt.Errorf("unsupported driver: %q", driverName)
	}
}

// ErrTxStarted is returned when trying to start a new transaction from a transactional client.
var ErrTxStarted = errors.New("ent: cannot start a transaction within a transaction")

// Tx returns a new transactional client. The provided context
// is used until the transaction is committed or rolled back.
func (c *Client) Tx(ctx context.Context) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, ErrTxStarted
	}
	tx, err := newTx(ctx, c.driver)
	if err != nil {
		return nil, fmt.Errorf("ent: starting a transaction: %w", err)
	}
	cfg := c.config
	cfg.driver = tx
	return &Tx{
		ctx:     ctx,
		config:  cfg,
		Task:    NewTaskClient(cfg),
		TaskLog: NewTaskLogClient(cfg),
	}, nil
}

// BeginTx returns a transactional client with specified options.
func (c *Client) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	if _, ok := c.driver.(*txDriver); ok {
		return nil, errors.New("ent: cannot start a transaction within a transaction")
	}
	tx, err := c.driver.(interface {
		BeginTx(context.Context, *sql.TxOptions) (dialect.Tx, error)
	}).BeginTx(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("ent: starting a transaction: %w", err)
	}
	cfg := c.config
	cfg.driver = &txDriver{tx: tx, drv: c.driver}
	return &Tx{
		ctx:     ctx,
		config:  cfg,
		Task:    NewTaskClient(cfg),
		TaskLog: NewTaskLogClient(cfg),
	}, nil
}

// Debug returns a new debug-client. It's used to get verbose logging on specific operations.
//
//	client.Debug().
//		Task.
//		Query().
//		Count(ctx)
func (c *Client) Debug() *Client {
	if c.debug {
		return c
	}
	cfg := c.config
	cfg.driver = dialect.Debug(c.driver, c.log)
	client := &Client{config: cfg}
	client.init()
	return client
}

// Close closes the database connection and prevents new queries from starting.
func (c *Client) Close() error {
	return c.driver.Close()
}

// Use adds the mutation hooks to all the entity clients.
// In order to add hooks to a specific client, call: `client.Node.Use(...)`.
func (c *Client) Use(hooks ...Hook) {
	c.Task.Use(hooks...)
	c.TaskLog.Use(hooks...)
}

// Intercept adds the query interceptors to all the entity clients.
// In order to add interceptors to a specific client, call: `client.Node.Intercept(...)`.
func (c *Client) Intercept(interceptors ...Interceptor) {
	c.Task.Intercept(interceptors...)
	c.TaskLog.Intercept(interceptors...)
}

// Mutate implements the ent.Mutator interface.
func (c *Client) Mutate(ctx context.Context, m Mutation) (Value, error) {
	switch m := m.(type) {
	case *TaskMutation:
		return c.Task.mutate(ctx, m)
	case *TaskLogMutation:
		return c.TaskLog.mutate(ctx, m)
	default:
		return nil, fmt.Errorf("ent: unknown mutation type %T", m)
	}
}

// TaskClient is a client for the Task schema.
type TaskClient struct {
	config
}

// NewTaskClient returns a client for the Task from the given config.
func NewTaskClient(c config) *TaskClient {
	return &TaskClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `task.Hooks(f(g(h())))`.
func (c *TaskClient) Use(hooks ...Hook) {
	c.hooks.Task = append(c.hooks.Task, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `task.Intercept(f(g(h())))`.
func (c *TaskClient) Intercept(interceptors ...Interceptor) {
	c.inters.Task = append(c.inters.Task, interceptors...)
}

// Create returns a builder for creating a Task entity.
func (c *TaskClient) Create() *TaskCreate {
	mutation := newTaskMutation(c.config, OpCreate)
	return &TaskCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of Task entities.
func (c *TaskClient) CreateBulk(builders ...*TaskCreate) *TaskCreateBulk {
	return &TaskCreateBulk{config: c.config, builders: builders}
}

// MapCreateBulk creates a bulk creation builder from the given slice. For each item in the slice, the function creates
// a builder and applies setFunc on it.
func (c *TaskClient) MapCreateBulk(slice any, setFunc func(*TaskCreate, int)) *TaskCreateBulk {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return &TaskCreateBulk{err: fmt.Errorf("calling to TaskClient.MapCreateBulk with wrong type %T, need slice", slice)}
	}
	builders := make([]*TaskCreate, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		builders[i] = c.Create()
		setFunc(builders[i], i)
	}
	return &TaskCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for Task.
func (c *TaskClient) Update() *TaskUpdate {
	mutation := newTaskMutation(c.config, OpUpdate)
	return &TaskUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *TaskClient) UpdateOne(t *Task) *TaskUpdateOne {
	mutation := newTaskMutation(c.config, OpUpdateOne, withTask(t))
	return &TaskUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *TaskClient) UpdateOneID(id uint64) *TaskUpdateOne {
	mutation := newTaskMutation(c.config, OpUpdateOne, withTaskID(id))
	return &TaskUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for Task.
func (c *TaskClient) Delete() *TaskDelete {
	mutation := newTaskMutation(c.config, OpDelete)
	return &TaskDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *TaskClient) DeleteOne(t *Task) *TaskDeleteOne {
	return c.DeleteOneID(t.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *TaskClient) DeleteOneID(id uint64) *TaskDeleteOne {
	builder := c.Delete().Where(task.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &TaskDeleteOne{builder}
}

// Query returns a query builder for Task.
func (c *TaskClient) Query() *TaskQuery {
	return &TaskQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeTask},
		inters: c.Interceptors(),
	}
}

// Get returns a Task entity by its id.
func (c *TaskClient) Get(ctx context.Context, id uint64) (*Task, error) {
	return c.Query().Where(task.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *TaskClient) GetX(ctx context.Context, id uint64) *Task {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryTaskLogs queries the task_logs edge of a Task.
func (c *TaskClient) QueryTaskLogs(t *Task) *TaskLogQuery {
	query := (&TaskLogClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := t.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(task.Table, task.FieldID, id),
			sqlgraph.To(tasklog.Table, tasklog.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, task.TaskLogsTable, task.TaskLogsColumn),
		)
		fromV = sqlgraph.Neighbors(t.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *TaskClient) Hooks() []Hook {
	return c.hooks.Task
}

// Interceptors returns the client interceptors.
func (c *TaskClient) Interceptors() []Interceptor {
	return c.inters.Task
}

func (c *TaskClient) mutate(ctx context.Context, m *TaskMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&TaskCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&TaskUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&TaskUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&TaskDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown Task mutation op: %q", m.Op())
	}
}

// TaskLogClient is a client for the TaskLog schema.
type TaskLogClient struct {
	config
}

// NewTaskLogClient returns a client for the TaskLog from the given config.
func NewTaskLogClient(c config) *TaskLogClient {
	return &TaskLogClient{config: c}
}

// Use adds a list of mutation hooks to the hooks stack.
// A call to `Use(f, g, h)` equals to `tasklog.Hooks(f(g(h())))`.
func (c *TaskLogClient) Use(hooks ...Hook) {
	c.hooks.TaskLog = append(c.hooks.TaskLog, hooks...)
}

// Intercept adds a list of query interceptors to the interceptors stack.
// A call to `Intercept(f, g, h)` equals to `tasklog.Intercept(f(g(h())))`.
func (c *TaskLogClient) Intercept(interceptors ...Interceptor) {
	c.inters.TaskLog = append(c.inters.TaskLog, interceptors...)
}

// Create returns a builder for creating a TaskLog entity.
func (c *TaskLogClient) Create() *TaskLogCreate {
	mutation := newTaskLogMutation(c.config, OpCreate)
	return &TaskLogCreate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// CreateBulk returns a builder for creating a bulk of TaskLog entities.
func (c *TaskLogClient) CreateBulk(builders ...*TaskLogCreate) *TaskLogCreateBulk {
	return &TaskLogCreateBulk{config: c.config, builders: builders}
}

// MapCreateBulk creates a bulk creation builder from the given slice. For each item in the slice, the function creates
// a builder and applies setFunc on it.
func (c *TaskLogClient) MapCreateBulk(slice any, setFunc func(*TaskLogCreate, int)) *TaskLogCreateBulk {
	rv := reflect.ValueOf(slice)
	if rv.Kind() != reflect.Slice {
		return &TaskLogCreateBulk{err: fmt.Errorf("calling to TaskLogClient.MapCreateBulk with wrong type %T, need slice", slice)}
	}
	builders := make([]*TaskLogCreate, rv.Len())
	for i := 0; i < rv.Len(); i++ {
		builders[i] = c.Create()
		setFunc(builders[i], i)
	}
	return &TaskLogCreateBulk{config: c.config, builders: builders}
}

// Update returns an update builder for TaskLog.
func (c *TaskLogClient) Update() *TaskLogUpdate {
	mutation := newTaskLogMutation(c.config, OpUpdate)
	return &TaskLogUpdate{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOne returns an update builder for the given entity.
func (c *TaskLogClient) UpdateOne(tl *TaskLog) *TaskLogUpdateOne {
	mutation := newTaskLogMutation(c.config, OpUpdateOne, withTaskLog(tl))
	return &TaskLogUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// UpdateOneID returns an update builder for the given id.
func (c *TaskLogClient) UpdateOneID(id uint64) *TaskLogUpdateOne {
	mutation := newTaskLogMutation(c.config, OpUpdateOne, withTaskLogID(id))
	return &TaskLogUpdateOne{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// Delete returns a delete builder for TaskLog.
func (c *TaskLogClient) Delete() *TaskLogDelete {
	mutation := newTaskLogMutation(c.config, OpDelete)
	return &TaskLogDelete{config: c.config, hooks: c.Hooks(), mutation: mutation}
}

// DeleteOne returns a builder for deleting the given entity.
func (c *TaskLogClient) DeleteOne(tl *TaskLog) *TaskLogDeleteOne {
	return c.DeleteOneID(tl.ID)
}

// DeleteOneID returns a builder for deleting the given entity by its id.
func (c *TaskLogClient) DeleteOneID(id uint64) *TaskLogDeleteOne {
	builder := c.Delete().Where(tasklog.ID(id))
	builder.mutation.id = &id
	builder.mutation.op = OpDeleteOne
	return &TaskLogDeleteOne{builder}
}

// Query returns a query builder for TaskLog.
func (c *TaskLogClient) Query() *TaskLogQuery {
	return &TaskLogQuery{
		config: c.config,
		ctx:    &QueryContext{Type: TypeTaskLog},
		inters: c.Interceptors(),
	}
}

// Get returns a TaskLog entity by its id.
func (c *TaskLogClient) Get(ctx context.Context, id uint64) (*TaskLog, error) {
	return c.Query().Where(tasklog.ID(id)).Only(ctx)
}

// GetX is like Get, but panics if an error occurs.
func (c *TaskLogClient) GetX(ctx context.Context, id uint64) *TaskLog {
	obj, err := c.Get(ctx, id)
	if err != nil {
		panic(err)
	}
	return obj
}

// QueryTasks queries the tasks edge of a TaskLog.
func (c *TaskLogClient) QueryTasks(tl *TaskLog) *TaskQuery {
	query := (&TaskClient{config: c.config}).Query()
	query.path = func(context.Context) (fromV *sql.Selector, _ error) {
		id := tl.ID
		step := sqlgraph.NewStep(
			sqlgraph.From(tasklog.Table, tasklog.FieldID, id),
			sqlgraph.To(task.Table, task.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, tasklog.TasksTable, tasklog.TasksColumn),
		)
		fromV = sqlgraph.Neighbors(tl.driver.Dialect(), step)
		return fromV, nil
	}
	return query
}

// Hooks returns the client hooks.
func (c *TaskLogClient) Hooks() []Hook {
	return c.hooks.TaskLog
}

// Interceptors returns the client interceptors.
func (c *TaskLogClient) Interceptors() []Interceptor {
	return c.inters.TaskLog
}

func (c *TaskLogClient) mutate(ctx context.Context, m *TaskLogMutation) (Value, error) {
	switch m.Op() {
	case OpCreate:
		return (&TaskLogCreate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdate:
		return (&TaskLogUpdate{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpUpdateOne:
		return (&TaskLogUpdateOne{config: c.config, hooks: c.Hooks(), mutation: m}).Save(ctx)
	case OpDelete, OpDeleteOne:
		return (&TaskLogDelete{config: c.config, hooks: c.Hooks(), mutation: m}).Exec(ctx)
	default:
		return nil, fmt.Errorf("ent: unknown TaskLog mutation op: %q", m.Op())
	}
}

// hooks and interceptors per client, for fast access.
type (
	hooks struct {
		Task, TaskLog []ent.Hook
	}
	inters struct {
		Task, TaskLog []ent.Interceptor
	}
)

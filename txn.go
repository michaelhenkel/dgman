package dgman

import (
	"context"

	"github.com/dgraph-io/dgo/v2"
)

// TxnContext is dgo transaction coupled with context
type TxnContext struct {
	txn *dgo.Txn
	ctx context.Context
}

// Commit calls Commit on the dgo transaction.
func (t *TxnContext) Commit() error {
	return t.txn.Commit(t.ctx)
}

// Discard calls Discard on the dgo transaction.
func (t *TxnContext) Discard() error {
	return t.txn.Discard(t.ctx)
}

// Txn returns the dgo transaction
func (t *TxnContext) Txn() *dgo.Txn {
	return t.txn
}

// WithContext replaces the current transaction context
func (t *TxnContext) WithContext(ctx context.Context) {
	t.ctx = ctx
}

// Context returns the transaction context
func (t *TxnContext) Context() context.Context {
	return t.ctx
}

// Mutate is a shortcut to create mutations from data to be marshalled into JSON,
// it will inject the node type from the Struct name
func (t *TxnContext) Mutate(data interface{}, options ...*MutateOptions) error {
	opt := &MutateOptions{}
	if len(options) > 0 {
		opt = options[0]
	}

	mType, err := newMutateType(data)
	if err != nil {
		return err
	}

	assigned, err := mutate(t.ctx, t.txn, data, opt)
	if err != nil {
		return err
	}

	return mType.saveUID(assigned.Uids)
}

// Create create node(s) with field unique checking, similar to Mutate,
// will inject node type from the Struct name
func (t *TxnContext) Create(data interface{}, options ...*MutateOptions) error {
	return mutateWithConstraints(t.ctx, t.txn, data, false, options...)
}

// Update updates a node by their UID with field unique checking, similar to Mutate,
// will inject node type from the Struct name
func (t *TxnContext) Update(data interface{}, options ...*MutateOptions) error {
	return mutateWithConstraints(t.ctx, t.txn, data, true, options...)
}

// Delete prepares a delete mutation using a query
func (t *TxnContext) Delete(model interface{}, opt ...*MutateOptions) *Deleter {
	mutateOpt := &MutateOptions{}
	if len(opt) > 0 {
		mutateOpt = opt[0]
	}

	q := &Query{ctx: t.ctx, tx: t.txn, model: model}
	return &Deleter{q: q, ctx: t.ctx, tx: t.txn, mutateOpt: mutateOpt}
}

// Get prepares a query for a model
func (t *TxnContext) Get(model interface{}) *Query {
	return &Query{ctx: t.ctx, tx: t.txn, model: model}
}

// NewTxnContext creates a new transaction coupled with a context
func NewTxnContext(ctx context.Context, c *dgo.Dgraph) *TxnContext {
	return &TxnContext{
		txn: c.NewTxn(),
		ctx: ctx,
	}
}

// NewTxn creates a new transaction
func NewTxn(c *dgo.Dgraph) *TxnContext {
	return NewTxnContext(context.Background(), c)
}

// NewReadOnlyTxnContext creates a new read only transaction coupled with a context
func NewReadOnlyTxnContext(ctx context.Context, c *dgo.Dgraph) *TxnContext {
	return &TxnContext{
		txn: c.NewReadOnlyTxn(),
		ctx: ctx,
	}
}

// NewReadOnlyTxn creates a new read only transaction
func NewReadOnlyTxn(c *dgo.Dgraph) *TxnContext {
	return NewReadOnlyTxnContext(context.Background(), c)
}
package builtinfilters

import (
	"sync"

	"github.com/iotaledger/goshimmer/packages/binary/transaction/payload/valuetransfer"

	"github.com/iotaledger/goshimmer/packages/binary/transaction"
	"github.com/iotaledger/hive.go/async"
)

type ValueTransactionSignatureFilter struct {
	onAcceptCallback func(tx *transaction.Transaction)
	onRejectCallback func(tx *transaction.Transaction)
	workerPool       async.WorkerPool

	onAcceptCallbackMutex sync.RWMutex
	onRejectCallbackMutex sync.RWMutex
}

func NewValueTransactionSignatureFilter() (result *ValueTransactionSignatureFilter) {
	result = &ValueTransactionSignatureFilter{}

	return
}

func (filter *ValueTransactionSignatureFilter) Filter(tx *transaction.Transaction) {
	filter.workerPool.Submit(func() {
		if payload := tx.GetPayload(); payload.GetType() == valuetransfer.Type {
			if valueTransfer, ok := payload.(*valuetransfer.ValueTransfer); ok && valueTransfer.VerifySignatures() {
				filter.getAcceptCallback()(tx)
			} else {
				filter.getRejectCallback()(tx)
			}
		} else {
			filter.getAcceptCallback()(tx)
		}
	})
}

func (filter *ValueTransactionSignatureFilter) OnAccept(callback func(tx *transaction.Transaction)) {
	filter.onAcceptCallbackMutex.Lock()
	filter.onAcceptCallback = callback
	filter.onAcceptCallbackMutex.Unlock()
}

func (filter *ValueTransactionSignatureFilter) OnReject(callback func(tx *transaction.Transaction)) {
	filter.onRejectCallbackMutex.Lock()
	filter.onRejectCallback = callback
	filter.onRejectCallbackMutex.Unlock()
}

func (filter *ValueTransactionSignatureFilter) Shutdown() {
	filter.workerPool.ShutdownGracefully()
}

func (filter *ValueTransactionSignatureFilter) getAcceptCallback() (result func(tx *transaction.Transaction)) {
	filter.onAcceptCallbackMutex.RLock()
	result = filter.onAcceptCallback
	filter.onAcceptCallbackMutex.RUnlock()

	return
}

func (filter *ValueTransactionSignatureFilter) getRejectCallback() (result func(tx *transaction.Transaction)) {
	filter.onRejectCallbackMutex.RLock()
	result = filter.onRejectCallback
	filter.onRejectCallbackMutex.RUnlock()

	return
}
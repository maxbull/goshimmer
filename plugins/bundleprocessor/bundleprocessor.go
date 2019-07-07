package bundleprocessor

import (
	"github.com/iotaledger/goshimmer/packages/errors"
	"github.com/iotaledger/goshimmer/packages/model/bundle"
	"github.com/iotaledger/goshimmer/packages/model/transactionmetadata"
	"github.com/iotaledger/goshimmer/packages/model/value_transaction"
	"github.com/iotaledger/goshimmer/packages/ternary"
	"github.com/iotaledger/goshimmer/plugins/tangle"
)

func ProcessSolidBundleHead(headTransaction *value_transaction.ValueTransaction) (*bundle.Bundle, errors.IdentifiableError) {
	// only process the bundle if we didn't process it, yet
	return tangle.GetBundle(headTransaction.GetHash(), func(headTransactionHash ternary.Trytes) (*bundle.Bundle, errors.IdentifiableError) {
		// abort if bundle syntax is wrong
		if !headTransaction.IsHead() {
			return nil, ErrProcessBundleFailed.Derive(errors.New("invalid parameter"), "transaction needs to be head of bundle")
		}

		// initialize event variables
		newBundle := bundle.New(headTransactionHash)
		bundleTransactions := make([]*value_transaction.ValueTransaction, 0)

		// iterate through trunk transactions until we reach the tail
		currentTransaction := headTransaction
		for {
			// abort if we reached a previous head
			if currentTransaction.IsHead() && currentTransaction != headTransaction {
				newBundle.SetTransactionHashes(mapTransactionsToTransactionHashes(bundleTransactions))

				Events.InvalidBundle.Trigger(newBundle, bundleTransactions)

				return nil, ErrProcessBundleFailed.Derive(errors.New("invalid bundle found"), "missing bundle tail")
			}

			// update bundle transactions
			bundleTransactions = append(bundleTransactions, currentTransaction)

			// retrieve & update metadata
			currentTransactionMetadata, dbErr := tangle.GetTransactionMetadata(currentTransaction.GetHash(), transactionmetadata.New)
			if dbErr != nil {
				return nil, ErrProcessBundleFailed.Derive(dbErr, "failed to retrieve transaction metadata")
			}
			currentTransactionMetadata.SetBundleHeadHash(headTransactionHash)

			// update value bundle flag
			if !newBundle.IsValueBundle() && currentTransaction.GetValue() != 0 {
				newBundle.SetValueBundle(true)
			}

			// if we are done -> trigger events
			if currentTransaction.IsTail() {
				newBundle.SetTransactionHashes(mapTransactionsToTransactionHashes(bundleTransactions))

				Events.BundleSolid.Trigger(newBundle, bundleTransactions)

				return newBundle, nil
			}

			// try to iterate to next turn
			if nextTransaction, err := tangle.GetTransaction(currentTransaction.GetTrunkTransactionHash()); err != nil {
				return nil, ErrProcessBundleFailed.Derive(err, "failed to retrieve trunk while processing bundle")
			} else {
				currentTransaction = nextTransaction
			}
		}
	})
}

func mapTransactionsToTransactionHashes(transactions []*value_transaction.ValueTransaction) (result []ternary.Trytes) {
	result = make([]ternary.Trytes, len(transactions))
	for k, v := range transactions {
		result[k] = v.GetHash()
	}

	return
}
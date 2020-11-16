/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package msgprocessor

import (
	"github.com/hyperledger/fabric-protos-go/common"
	"github.com/hyperledger/fabric/core/ledger/kvledger/txmgmt/rwsetutil"
	"github.com/hyperledger/fabric/gdpr"
	"github.com/pkg/errors"
)

type GDPRFilter struct {
}

func (G GDPRFilter) Apply(message *common.Envelope) error {
	preImages := gdpr.HashedPreImages(message.PreImages)
	var finalError error
	gdpr.ProcessEnvelope(message, nil, 0, func(block *common.Block, i int, rws *rwsetutil.TxRwSet) {
		for _, nsRWS := range rws.NsRwSets {
			for _, kvWrite := range nsRWS.KvRwSet.Writes {
				if !preImages.Exists(string(kvWrite.ValueHash)) {
					finalError = errors.Errorf("key %s is missing a pre-image", kvWrite.Key)
					return
				}
			}
		}
	}, func(err error) {
		finalError = err
	})
	return finalError
}
package datastore

import (
    "sync"

	"github.com/plan-systems/go-plan/pdi"
	"github.com/plan-systems/go-plan/plan"
	"github.com/plan-systems/go-plan/ski"
)

const (
	txnEncodingDesc1 = "/plan/pdi/encoding/datastore/1"
)

// dsDecoder implements pdi.TxnDecoder
type dsDecoder struct {
	pdi.TxnDecoder

    theadsafe    bool
    mutex        sync.Mutex
	encodingDesc string
	hashKits     map[ski.HashKitID]ski.HashKit
}

// NewTxnDecoder creates a TxnDecoder for use with pdi-datastore
func NewTxnDecoder(
    inMakeThreadsafe bool,
) pdi.TxnDecoder {

	dec := &dsDecoder{
        theadsafe:    inMakeThreadsafe,
		hashKits:     map[ski.HashKitID]ski.HashKit{},
		encodingDesc: txnEncodingDesc1,
	}

	return dec
}

// EncodingDesc -- see TxnDecoder
func (dec *dsDecoder) EncodingDesc() string {
	return dec.encodingDesc
}

// DecodeRawTxn -- See TxnDecoder
func (dec *dsDecoder) DecodeRawTxn(
	rawTxn []byte,
	outInfo *pdi.TxnInfo,
) ([]byte, error) {

    const (
        sigLenLen = 2
    )

	txnLen := uint32(len(rawTxn))
	if txnLen < 50 {
		return nil, plan.Errorf(nil, plan.FailedToUnmarshal, "txn is too small (txnLen=%v)", txnLen)
	}

	// 1) Unmarshal the txn info
	var txnInfo pdi.TxnInfo
	pos := 2 + uint32(rawTxn[0]) | ( uint32(rawTxn[1]) << 8 )
    if pos > txnLen {
		return nil, plan.Error(nil, plan.FailedToUnmarshal, "txnInfo len exceeds txn buf size")
    }
	if err := txnInfo.Unmarshal(rawTxn[2:pos]); err != nil {
		return nil, plan.Error(err, plan.FailedToUnmarshal, "failed to unmarshal txnInfo")
	}

	// 2) Extract the payload buf
	end := pos + txnInfo.SegSz
	if end > txnLen {
		return nil, plan.Errorf(nil, plan.FailedToUnmarshal, "payload buffer EOS (txnLen=%v, pos=%v, end=%v)", txnLen, pos, end)
	}
	payloadBuf := rawTxn[pos:end]

	// 3) Extract the sig -- the last 2 bytes
    sigLen := uint32(rawTxn[txnLen-2])| (uint32(rawTxn[txnLen-1]) << 8)
	txnLen -= sigLenLen + sigLen
	if txnLen < 10 {
		return nil, plan.Errorf(nil, plan.FailedToUnmarshal, "txn sig len is wrong (txnLen=%v, sigLen=%v)", txnLen, sigLen)
	}
	sig := rawTxn[txnLen : txnLen+sigLen]

    if dec.theadsafe {
        dec.mutex.Lock()
    }

	// 4) Prep the hasher so we can generate a digest
	hashKit, ok := dec.hashKits[txnInfo.HashKitId]
	if ! ok {
		var err error
		hashKit, err = ski.NewHashKit(txnInfo.HashKitId)
		if err != nil {
			return nil, err
		}
		dec.hashKits[txnInfo.HashKitId] = hashKit
	}

    miscBuf := make([]byte, 0, pdi.UTIDBinarySz + hashKit.HashSz)

	// 5) Calculate the hash digest and thus UTID of the raw txn
	hashKit.Hasher.Reset()
	hashKit.Hasher.Write(rawTxn[:txnLen])
	txnInfo.TxnHashname = hashKit.Hasher.Sum(miscBuf)
    txnInfo.UTID = pdi.UTIDFromInfo(miscBuf[hashKit.HashSz:hashKit.HashSz], txnInfo.TimeSealed, txnInfo.TxnHashname)

    if dec.theadsafe {
        dec.mutex.Unlock()
    }

	// 6) Verify the sig
	pubKey := &ski.PubKey{
		KeyDomain: ski.KeyDomain_PERSONAL,
		Bytes:     txnInfo.From,
	}
	if err := ski.VerifySignatureFrom(sig, txnInfo.TxnHashname, pubKey); err != nil {
		return nil, err
	}

	if outInfo != nil {
		*outInfo = txnInfo
	}

	return payloadBuf, nil
}

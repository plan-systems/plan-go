package datastore


import (
    //"bytes"
	"sync"

	"github.com/plan-systems/go-plan/pdi"
	"github.com/plan-systems/go-plan/plan"
	"github.com/plan-systems/go-plan/ski"
)






const (

    // ProtocolDesc appears in StorageEpochs originating from this impl
    ProtocolDesc = "plan/storage/pdi-datastore/1"

)


// NewStorageEpoch generates a new StorageEpoch, needed when creating a new community.
func NewStorageEpoch(
    skiSession  ski.Session,
    inCommunity *pdi.CommunityEpoch,
) (*pdi.StorageEpoch, error) {

    epoch := &pdi.StorageEpoch{
        StorageProtocol: ProtocolDesc,
        CommunityID: inCommunity.CommunityID,
        TxnHashKit: ski.HashKitID_LegacyKeccak_256,
        Name: inCommunity.CommunityName,
        TxnMaxSize: 32000,
        CommunityChIDs: make([]byte, pdi.CommunityChID_NumChannels * plan.ChIDSz),
    }

    var err error
    epoch.OriginKey, err = ski.GenerateNewKey(
        skiSession,
        inCommunity.FormGenesisKeyringName(),
        ski.KeyInfo{
            KeyType: ski.KeyType_SigningKey,
            CryptoKit: ski.CryptoKitID_NaCl,
        },
    )
    if err != nil {
        return nil, err
    }

    return epoch, err
}





// dsEncoder implements pdi.TxnEncoder
type dsEncoder struct {
	pdi.TxnEncoder

    StorageEpoch  pdi.StorageEpoch
    threadsafe    bool
    mutex         sync.Mutex

    packer        ski.PayloadPacker
    from          ski.KeyInfo

    // Used to marshal TxnInfo
    scrap         []byte
}


// NewTxnEncoder creates a new StorageProviderAgent for use with a pdi-datastore StorageProvider.
// If inSegmentMaxSz == 0, then a default size is chosen
func NewTxnEncoder(
    inMakeThreadsafe bool,
	inStorageEpoch pdi.StorageEpoch,
) pdi.TxnEncoder {

	enc := &dsEncoder{
        StorageEpoch: inStorageEpoch,
        threadsafe:   inMakeThreadsafe,
        packer:       ski.NewPacker(false),
	}

	return enc
}


// ResetSigner --see TxnEncoder
func (enc *dsEncoder) ResetSigner(
	inSession ski.Session,
	inFrom []byte,
) error {

	return enc.packer.ResetSession(
        inSession, 
        ski.KeyRef{
            KeyringName: enc.StorageEpoch.StorageKeyringName(),
            PubKey: inFrom,
        }, 
        enc.StorageEpoch.TxnHashKit,
        &enc.from)

}


// EncodeToTxns -- See StorageProviderAgent.EncodeToTxns()
func (enc *dsEncoder) EncodeToTxns(
	inPayload         []byte,
    inPayloadID       []byte,
	inPayloadEnc      plan.Encoding,
	inTransfers       []*pdi.Transfer,
    timeSealed        int64,
    outTxns           *pdi.PayloadTxns,
) error {

	segs, err := pdi.SegmentIntoTxns(
		inPayload,
        inPayloadID,
		inPayloadEnc,
		enc.StorageEpoch.TxnMaxSize,
	)
	if err != nil {
		return err
	}


    outTxns.PayloadEncoding = inPayloadEnc
    outTxns.PayloadID       = nil
    outTxns.NewlyAuthored   = true

    payloadIDLen := len(inPayloadID)

    if cap(outTxns.Txns) < len(segs) {
	    outTxns.Txns = make([]pdi.RawTxn, len(segs))
    } else {
	    outTxns.Txns = outTxns.Txns[:len(segs)]
    }

    // Put the transfers in the last segment
    segs[len(segs)-1].Transfers = inTransfers

    // Use the same time stamp for the entire batch
    if timeSealed == 0 {
        timeSealed = plan.Now()
    }

    {
        pos := uint32(0)

        // We have a mutex b/c of the shared scrap
        if enc.threadsafe {
            enc.mutex.Lock()
        }

        scrap := enc.scrap

		for i, seg := range segs {
            
            // Set the rest of the txn fields
            seg.TimeSealed = timeSealed
            if i > 0 {
                seg.PrevURID = outTxns.Txns[i-1].URID
            }

            headerSz := seg.Size()
            if headerSz > len(scrap) {
                enc.scrap = make([]byte, headerSz + 5000)
                scrap = enc.scrap
            }
            headerSz, err = seg.MarshalTo(scrap)
            if err != nil {
                return plan.Error(nil, plan.MarshalFailed, "failed to marshal txn info")
            }

            packingInfo := ski.PackingInfo{}
            err = enc.packer.PackAndSign(
                plan.Encoding_TxnPayloadSegment,
                scrap[:headerSz],
                inPayload[pos:pos+seg.SegSz],
                pdi.URIDSz + payloadIDLen,
                &packingInfo,
            )
            if err != nil {
                return err
            }

            outTxns.Txns[i].Bytes = packingInfo.SignedBuf
            outTxns.Txns[i].URID = pdi.URIDFromInfo(packingInfo.Extra, seg.TimeSealed, packingInfo.Hash)

            // Allocate bytes for the the payload ID so we don't have to do an additional alloc -- but do it only once.
            if payloadIDLen > 0 {
                outTxns.PayloadID = append(packingInfo.Extra[pdi.URIDSz:pdi.URIDSz], inPayloadID...)
                payloadIDLen = 0
            }

            pos += seg.SegSz
        }

        if enc.threadsafe {
            enc.mutex.Unlock()
        }

        if int(pos) != len(inPayload) {
			return plan.Error(nil, plan.AssertFailed, "payload size failed check")
		}
	}

	return nil
}



// dsDecoder implements pdi.TxnDecoder
type dsDecoder struct {
	pdi.TxnDecoder

    unpacker        ski.PayloadUnpacker
}

// NewTxnDecoder creates a TxnDecoder for use with pdi-datastore
func NewTxnDecoder(
    inMakeThreadsafe bool,
) pdi.TxnDecoder {

	return &dsDecoder{
        unpacker: ski.NewUnpacker(inMakeThreadsafe),
	}
}

// DecodeRawTxn -- See TxnDecoder
func (dec *dsDecoder) DecodeRawTxn(
	inRawTxn []byte, 
	outInfo  *pdi.TxnInfo,
) ([]byte, error) {

    out := ski.SignedPayload{
        Hash: make([]byte, 128),
    }

    err := dec.unpacker.UnpackAndVerify(inRawTxn, &out)
    if err != nil {
        return nil, err
    }

	var txnInfo pdi.TxnInfo
	if err := txnInfo.Unmarshal(out.Header); err != nil {
		return nil, plan.Error(err, plan.UnmarshalFailed, "failed to unmarshal TxnInfo")
	}

    txnInfo.From = out.Signer.PubKey
    txnInfo.TxnHashname = out.Hash
    txnInfo.URID = pdi.URIDFromInfo(out.Hash[len(out.Hash):], txnInfo.TimeSealed, out.Hash)

	// 2) Isolate the payload buf
	if txnInfo.SegSz != uint32(len(out.Body)) {
		return nil, plan.Errorf(nil, plan.UnmarshalFailed, "txn payload buf length doesn't match")
	}

	if outInfo != nil {
		*outInfo = txnInfo
	}

	return out.Body, nil
}

package repo


import (
    "bytes"
    "os"
    "path"
    //"io/ioutil"
    //"strings"
    "sync"
    //"context"
    "fmt"
    //"sort"
    //"encoding/hex"
    //"encoding/json"
    //"encoding/base64"
    //"encoding/binary"
    "sync/atomic"

    "github.com/plan-systems/plan-core/tools"
    "github.com/plan-systems/plan-core/tools/ctx"
    "github.com/plan-systems/plan-core/pdi"
    "github.com/plan-systems/plan-core/plan"
    "github.com/plan-systems/plan-core/ski"


    "github.com/dgraph-io/badger"

)

// ChSessID identifies a sub session within/under MemberSession, the obj associated with a client that is connected.
type ChSessID uint32

// ChMgr is the top level interface for a community's channels.
type ChMgr struct {
    ctx.Context

    HomePath                string
    prevSessID              uint32
    StorageEpoch            pdi.StorageEpoch               
    CR                      *CommunityRepo
    chMutex                 sync.RWMutex
    chLoaded                map[plan.ChIDBlob]*ChStore


    //db                      *badger.DB
    //State                   ChMgrState
}


// NewChMgr creates and inits a new ChMgr instance.
func NewChMgr(
    inHomeDir string,
    CR *CommunityRepo,
) *ChMgr {

    chMgr := &ChMgr{
        chLoaded: make(map[plan.ChIDBlob]*ChStore),
        HomePath: path.Join(inHomeDir, "ChMgr"),
        prevSessID: 0,
        StorageEpoch: *CR.GenesisSeed.StorageEpoch,
        CR: CR,
    }

    chMgr.SetLogLabel(CR.GetLogLabel())

    return chMgr

}

// Startup -- see plan.Flow.Startup()
func (chMgr *ChMgr) Startup() error {

    err := chMgr.CtxStart(
        chMgr.ctxStartup,
        nil,
        nil,
        chMgr.ctxStopping,
    )

    return err
}



func (chMgr *ChMgr) ctxStartup() error {

    /*
    opts := badger.DefaultOptions
    opts.Dir = path.Join(chMgr.HomePath, "chMgr.badger")
    opts.ValueDir = opts.Dir

    var err error
    if chMgr.db, err = badger.Open(opts); err != nil {
        return plan.Error(err, plan.StorageNotReady, "chMgr db failure")
    }

    // Load latest chMgr info
    if err == nil {
        dbTxn := chMgr.db.NewTransaction(false)

        var item *badger.Item

        item, err = dbTxn.Get(chMgrStateKey)
        if err == nil {
            err = item.Value(func(val []byte) error {
                return chMgr.State.Unmarshal(val)
            })
        } else if err == badger.ErrKeyNotFound {
            chMgr.State = ChMgrState{}
            err = nil
        }

        dbTxn.Discard()
    }*/

    return nil
}
  

func (chMgr *ChMgr) ctxStopping() {

    // The assumption at this point is that no new entries are being sent to channels to be merged.
    // This means that we can close each channel's entriesToMerge to proceed with shutdown
    var waitingOn []*ChStore
    for {

        channels := &sync.WaitGroup{} 

        // First, cause all the ChStores to exit their main entry processing loop.
        chMgr.chMutex.RLock()
        {
            // Make a list to hold the channels we're waiting on
            N := len(chMgr.chLoaded)
            if cap(waitingOn) < N {
                waitingOn = make([]*ChStore, 0, N)
            } else {
                waitingOn = waitingOn[:0]
            }
            channels.Add(N)
            for _, i := range chMgr.chLoaded {
                waitingOn = append(waitingOn, i)
                go func(chSt *ChStore) {
                    chSt.ShutdownEntryProcessing(true)
                    channels.Done()
                }(i)
            }
        }
        chMgr.chMutex.RUnlock()

        if len(waitingOn) == 0 {
            break
        }

        chMgr.CR.Infof(1, "stopping %d open channels", len(waitingOn))

        // Wait until every loaded channel has exited it main validation loop
        channels.Wait()

        // At this point, nothing is running (so no futher calls to FetchChannel() should/will be made), so now we can fully shutdown each channel.
        channels.Add(len(waitingOn))
        for _, i := range waitingOn {
            chSt := i
            chMgr.detachChannel(chSt.ChID())
            go chSt.Shutdown(channels)
        }

        // Wait until all channels have fully shutdown (dbs closed, etc)
        channels.Wait()

        chMgr.CR.Infof(1, "all channels shutdown")

    }
}




type chEntryOrigin int32
const (
    entryUseAnonymous chEntryOrigin = iota
	entryFromStorageProvider
    entryRevalidating
    entryWasAuthored 
)



var chEntryPool = sync.Pool{
	New: func() interface{} {
		return new(chEntry)
	},
}

// newChEntry is equivalent to instantiating a fresh chEntry (but saves on allocations)
func newChEntry(inOrigin chEntryOrigin) *chEntry {
    entry := chEntryPool.Get().(*chEntry)
    entry.Reset()

    entry.origin = inOrigin

    return entry
}

// RecycleChEntry should be used whenever a chEntry is no longer referenced (and whose parts can be reused)
func RecycleChEntry(entry *chEntry) {

    entry.onMergeComplete = nil
    if entry.PayloadTxnSet != nil {
        pdi.RecycleTxnSet(entry.PayloadTxnSet)
        entry.PayloadTxnSet = nil
    }

    chEntryPool.Put(entry)
}


type chEntryPartStatus int32
const (
    partNotLoaded chEntryPartStatus = iota
	partLoaded
    partTouched
)

// chEntry is the working struct for an entry being merged and validated in a repo channel.
type chEntry struct {
    Body                []byte
    
    Info                pdi.EntryInfo
    State               EntryState
    StatePrev           EntryState

    PayloadTxnSet       *pdi.PayloadTxnSet

    origin              chEntryOrigin

    // TODO: make bytes.Buffer
    scrap               []byte
    
    //Asset               pcore.Marshaller
    //newChEpoch          *pdi.ChannelEpoch
    stateStatus          chEntryPartStatus
    bodyStatus           chEntryPartStatus

    onMergeComplete     func(entry *chEntry, chSt *ChStore, inErr error)

}



func (entry *chEntry) Reset() {

    entry.Info.EntrySubOp = 0
    entry.Info.SupersedesEntryID = entry.Info.SupersedesEntryID[:0]
    entry.Info.Extensions = nil
    entry.Body = entry.Body[:0]
    entry.State.LiveIDs = entry.State.LiveIDs[:0]
    entry.onMergeComplete = nil

    entry.PayloadTxnSet = nil
}




func (entry *chEntry) LivenessChanged() bool {
    return (entry.State.Status == EntryStatus_LIVE) != (entry.StatePrev.Status == EntryStatus_LIVE)
}

func (entry *chEntry) IsLive() bool {
    return entry.State.Status == EntryStatus_LIVE
}


func (entry *chEntry) AssignFromDecrytedEntry(
    inPayload *ski.SignedPayload,
) error {

    entry.State.Status              = EntryStatus_AWAITING_MERGE
    entry.State.Flags               = 1 << byte(EntryFlag_WELL_FORMED)
    entry.State.LiveIDs             = entry.State.LiveIDs[:0]

    entry.StatePrev.Status          = EntryStatus_AWAITING_MERGE
    entry.StatePrev.Flags           = 0
    entry.StatePrev.LiveIDs         = entry.StatePrev.LiveIDs[:0]

    entry.stateStatus               = partTouched
    entry.bodyStatus                = partTouched

    entry.Info.Recycle()

    err := entry.Info.Unmarshal(inPayload.Header)
    if err != nil {
        err = plan.Error(err, plan.ChEntryIsMalformed, "EntryInfo.Unmarshal() error")
    }

    if err == nil {
    
        entry.Info.AuthorSig = append(entry.Info.AuthorSig[:0], inPayload.Signer.PubKey...)
        entry.Body           = append(entry.Body[:0],           inPayload.Body...)

        if len(entry.Info.TIDs) < int(pdi.EntryTID_NormalNumTIDs) * plan.TIDSz {
            err = plan.Error(nil, plan.ChEntryIsMalformed, "entry missing required TIDs")
        }
    }

    if err != nil {
        entry.ThrowMalformed(err)
        return err
    }

    entryID := entry.Info.EntryID()

    // The entry ID's latter bytes are rightmost bytes of the entry hashname.  
    // This has to always be set since the hashname can never be known when the entry is being packed
    entryID.SetHash(inPayload.Hash)
    
    copy(entry.PayloadTxnSet.PayloadTID[:], entryID)

    return nil
}



// MissingFlag returns true if the given flag is NOT set
func (entry *chEntry) MissingFlag(inFlag EntryFlag) bool {
    flag := uint32(1) << byte(inFlag)
    return (entry.State.Flags & flag) == 0
}

// ClearFlag clears the given flag
func (entry *chEntry) ClearFlag(inFlag EntryFlag) {
    flag := uint32(1) << byte(inFlag)
    if entry.State.Flags & flag != 0 {
        entry.State.Flags ^= flag
        entry.TouchState()
    }
}

// HasFlag returns true if the given flag is set
func (entry *chEntry) HasFlag(inFlag EntryFlag) bool {
    flag := uint32(1) << byte(inFlag)
    return (entry.State.Flags & flag) == flag
}

// HasFlags returns true if ALL of the passed flags are present
func (entry *chEntry) HasFlags(inFlags...EntryFlag) bool {

    chFlags := entry.State.Flags

    for _, flag := range inFlags {
        bit := uint32(1) << uint8(flag)
        if (chFlags & bit) == 0 {
            return false
        }
    }

    return true
}

// AddFlags ORs the given flags to this
func (entry *chEntry) AddFlags(inFlags...EntryFlag) {

    chFlags := entry.State.Flags

    for _, flag := range inFlags {
        bit := uint32(1) << uint8(flag)
        chFlags |= bit
    }

    if entry.State.Flags != chFlags {
        entry.State.Flags = chFlags
        entry.TouchState()
    }
}

// ThrowMalformed marks this entry as disbarred due to not being self-consistent for the given reason.
func (entry *chEntry) ThrowMalformed(inErr error) {
    entry.ClearFlag(EntryFlag_WELL_FORMED)
    entry.SetStatus(EntryStatus_DISBARRED)
}

// SetStatus sets this ch entry status to the given value
func (entry *chEntry) SetStatus(inStatus EntryStatus) {
    if entry.State.Status != inStatus {
        entry.State.Status = inStatus
        entry.TouchState()
    }
}

/*
// Defer marks this entry as deferred for the given reason/error but never overwrites a EntryStatus_DISBARRED status.
func (entry *chEntry) Defer(inErr error) {
    if ntry.State.Status != EntryStatus_DISBARRED
    entry.SetStatus(EntryStatus_DEFERRED)
}*/

// TouchState marks for entry.State to be flushed to the channel db
func (entry *chEntry) TouchState() {
    entry.stateStatus = partTouched
}

// IsWellFormed returns true if EntryFlag_WELL_FORMED is set for this entry.
func (entry *chEntry) IsWellFormed() bool {
    return (entry.State.Flags & (1 << byte(EntryFlag_WELL_FORMED))) != 0
}

// QueueEntryForMerge is called when an entry arrives to a repo and must be merged into the repo.
func (chMgr *ChMgr) QueueEntryForMerge(
    entryCommEpoch *pdi.CommunityEpoch,
    entry *chEntry,
) {

    var (
        chGenesisEpoch *pdi.ChannelEpoch
        err error
    )

    switch entry.Info.EntryOp {
        
        case pdi.EntryOp_NEW_CHANNEL_EPOCH: {
            chEpoch := &pdi.ChannelEpoch{}
            if err = chEpoch.Unmarshal(entry.Body); err != nil {
                err = plan.Error(err, plan.ChEntryIsMalformed, "error unmarshalling NEW_CHANNEL_EPOCH")
                entry.ThrowMalformed(err)
            } else {
                entryID := entry.Info.EntryID()

                // Warning: when we use slices from chEntry, watch out for escapes since chEntry is recycled (and will be overwritten)
                chEpoch.EpochTID = entryID.Clone()
                chEpoch.CommunityEpochID = entryCommEpoch.EpochTID

                // Having set the above fields, remarshal the body
                entry.Body = tools.SmartMarshal(chEpoch, entry.Body)
    
                // Is the the channel genesis epoch?  If so, the channel ID derives from the channel ID.
                if chEpoch.IsChannelGenesis() {
                    entry.Info.ChannelID = append(entry.Info.ChannelID[:0], entryID.ExtractChID()...)
                    chGenesisEpoch = chEpoch
                    entry.AddFlags(EntryFlag_IS_CHANNEL_GENESIS) 
                }
            }
        }
    }

    chID := entry.Info.ChID()

    if len(chID) != plan.ChIDSz {
        err = plan.Error(err, plan.ChEntryIsMalformed, "malformed channel ID")
        entry.ThrowMalformed(err)
    }

    var chSt *ChStore

    if entry.IsWellFormed() {
        chSt, err = chMgr.fetchChannel(entry.Info.ChannelID, true, chGenesisEpoch)

        if err != nil {
            chMgr.Error("failed to fetch channel: ", plan.BinEncode(entry.Info.ChannelID))
        }
    }


    if chSt == nil {
        //chEntry.Status = EntryStatus_DEFERRED
        panic("TODO")
    } else {

        // From here on out, an entry is always accompanied by the ChStore, so we drop the entry chID since it's now redundant.
        // We zero it here to ensure entry.Info.ChannelID isn't used accidentally. 
        entry.Info.ChannelID = entry.Info.ChannelID[:0]

        chSt.entriesToMerge <- entry
    }


    //chMgr.Log(entryErr, "error merging entry")
    //chMgr.Log(processErr, "error processing entry") */
}


// FetchACC returns the given channel known/presumed to be an ACC.
func (chMgr *ChMgr) FetchACC(
    inChID plan.ChID,
) (*ChACC, error) {

    ch, err := chMgr.ChannelExists(inChID)
    if err != nil {
        return nil, err
    }

    var acc *ChACC

    if ch != nil && ch.agent != nil {
        acc, _ = ch.agent.(*ChACC)
    }

    if acc == nil {
        return nil, plan.Errorf(nil, plan.NotAnACC, "cited channel %v is not an access control channel", inChID.SuffixStr())
    }

    return acc, nil
}


// FetchMemberRegistry returns the community's member registry channel
func (chMgr *ChMgr) FetchMemberRegistry() (*ChMemberRegistry, error) {

    ch, err := chMgr.FetchChannel(chMgr.StorageEpoch.MemberRegistry())
    if err != nil {
        return nil, err
    }

    reg, ok := ch.agent.(*ChMemberRegistry)
    if ! ok {
        return nil, plan.Error(nil, plan.AssertFailed, "member registry not available")
    }

    return reg, nil
}


// FetchCommunityEpochsChannel returns the community's epoch channel
func (chMgr *ChMgr) FetchCommunityEpochsChannel() *ChCommunityEpochs {

    ch, err := chMgr.FetchChannel(chMgr.StorageEpoch.CommunityEpochs())
    if err != nil {
        return nil
    }

    chCE, _ := ch.agent.(*ChCommunityEpochs)

    return chCE
}



// ChannelExists returns the given channel, loading it if necessary.  If the channel is not found, then it isn not autho
func (chMgr *ChMgr) ChannelExists(
    inChID plan.ChID,
) (*ChStore, error) {

    return chMgr.fetchChannel(inChID, false, nil)
}


// FetchChannel returns the owning ChAgent for the given channel ID.
// If inGenesisEpoch is set, channel genesis is performed. 
func (chMgr *ChMgr) FetchChannel(
    inChID plan.ChID,
) (*ChStore, error) {

    return chMgr.fetchChannel(inChID, true, nil)
}


// FetchChannel returns the owning ChAgent for the given channel ID.
// If inGenesisEpoch is set, channel genesis is performed. 
func (chMgr *ChMgr) fetchChannel(
    inChID plan.ChID,
    inAutoCreate bool,
    inChGenesis *pdi.ChannelEpoch,
) (*ChStore, error) {

    var (
        err error
    )

    if len(inChID) != plan.ChIDSz {
        return nil, plan.Error(nil, plan.ParamErr, "malformed ChID")
    }

    chID := inChID.Blob()

    chMgr.chMutex.RLock()
    chSt := chMgr.chLoaded[chID]
    chMgr.chMutex.RUnlock()

    loading := false

    if chSt == nil {  
        chMgr.chMutex.Lock()
        loading = true
        defer chMgr.chMutex.Unlock()
    
        // If the ch store isn't already loaded, try to load it (recheck for race condition)
        if chSt = chMgr.chLoaded[chID]; chSt == nil { 
            chSt, err = chMgr.loadChannel(inChID, inAutoCreate)

            if err == nil {
                err = chSt.startEntryProcessor()
            }

            // We've already mutexed for this write condition (above)
            if err == nil {
                chMgr.chLoaded[chID] = chSt
            }
        }
    }

    if err == nil && chSt.agent == nil {
        msg := "loading channel (%s)"

        // Are we processing a channel genesis entry but the channel's agent is nil?  
        // MOSTLY, this occurs when the first entry to be processed is a channel genesis entry.
        // But this also occurs when entries are merged into a channel BEFORE the channel's genesis entry arrives.
        if inChGenesis != nil {
            if len(chSt.State.ChProtocol) == 0 {
                chSt.State.ChProtocol = inChGenesis.ChProtocol
                if loading {
                    msg = "initializing channel with protocol '%s'"
                } else {
                    msg = "retroactively assigning channel protocol to '%s'"
                }
            } else if chSt.State.ChProtocol != inChGenesis.ChProtocol {
                err = plan.Errorf(nil, plan.FailedToLoadChannel, "ch protocol mismatch")
            }
        }

        if err == nil {
            chSt.Infof(2, msg, chSt.State.ChProtocol)

            // TODO: set chSt state so it doesn't keep failing to load agent every time it get fetched
            err = chSt.setupChAgent()
            if err != nil {
                chSt.Error("failed to setup channel agent: ", err)
            }
        }
    }

    // TODO: rebuild channel if db load fails?
    if err != nil {
        chMgr.Error("failed to load channel ", inChID.Str(), ": ", err)

        chSt = nil
    }

    return chSt, err
}



func (chMgr *ChMgr) detachChannel(
    inChID plan.ChID,
) {

    chID := inChID.Blob()

    chMgr.chMutex.Lock()
    ch := chMgr.chLoaded[chID]
    if ch != nil {
        delete(chMgr.chLoaded, chID)
    }
    chMgr.chMutex.Unlock()

}


func (chMgr *ChMgr) loadChannel(
    inChID plan.ChID,
    inAutoCreate bool,
) (*ChStore, error) {

    if ! chMgr.CtxRunning() {
        return nil, plan.Error(nil, plan.ServiceShutdown, "chMgr is not available")
    }

    opts := badger.DefaultOptions(path.Join(chMgr.HomePath, inChID.Str()))

    if ! inAutoCreate {
        if _, err := os.Stat(opts.Dir); os.IsNotExist(err) {
            return nil, plan.Errorf(nil, plan.ChannelNotFound, "channel %v not found", inChID.Str())
        }
    }

    var err error
    if err = os.MkdirAll(opts.Dir, plan.DefaultFileMode); err != nil {
        return nil, err
    }

    chSt := &ChStore{
        chMgr: chMgr,
        entriesToMerge: make(chan *chEntry, 2),
        revalAfter: pdi.URIDTimestampMax,
    }

    chSt.db, err = badger.Open(opts)
    if err != nil {
        return nil, err
    }

    // First load the latest channel state object from the chDB
    if err == nil {
        chTxn := chSt.db.NewTransaction(false)

        var item *badger.Item

        item, dbErr := chTxn.Get(chStateKey)
        if dbErr == nil {
            err = item.Value(func(val []byte) error {
                return chSt.State.Unmarshal(val)
            })
            if err == nil {
                if ! bytes.Equal(chSt.State.ChID, inChID) {
                    err = plan.Errorf(nil, plan.FailedToLoadChannel, "chID in state does not match %v, got %v", inChID, chSt.State.ChID)
                }
            }
        } else if dbErr == badger.ErrKeyNotFound {
            chSt.State.ChID =  append(chSt.State.ChID[:0], inChID...)
        }

        chTxn.Discard()
    }
    
    if err != nil {
        // TODO: close db and mark as failed
        chSt = nil
    }

    return chSt, err
}








        /*
        chSt = &ChStore{}

        dbTxn := chMgr.db.NewTransaction(true)

        item, err = dbTxn.Get(inChannelID)
        if err == nil {
            err = item.Value(func(val []byte) error {
                return chSt.State.Unmarshal(val)
            })

        // If we didn't find the channel ID, designate (and/or allocate) a ChStore for this channel ID.
        } else if err == badger.ErrKeyNotFound {
            chMgr.State.NumChStores++
            chSt.State.ChStoreID = chMgr.State.NumChStores

            var (
                buf [200]byte
                bufSz int
            )
            bufSz, err = chSt.State.MarshalTo(buf[:])

            // Store the new ChStore ID assignment
		    if err == nil {
                err = dbTxn.Set(inChannelID, buf[:bufSz])
                
                // Store the new ch mgr state
                if err == nil {
                    bufSz, err = chMgr.State.MarshalTo(buf[:])
                    err = dbTxn.Set(chMgrStateKey, buf[:bufSz])
                }
            }
        }
        }*/





// StartChannelSession instantiates a nre channel session for the given channel ID (and accompanying params)
func (chMgr *ChMgr) StartChannelSession(
    inMemberSession *MemberSession,
    inChID plan.ChID,
) (*ChSession, error) {

    chSt, err := chMgr.FetchChannel(inChID)
    if err != nil {
        return nil, err
    }

    cs := &ChSession{
        MemberSession: inMemberSession,
        ChSessID: ChSessID(atomic.AddUint32(&chMgr.prevSessID, 1)),
        chSt: chSt,
        msgInbox: make(chan *Msg, 1),
        entryUpdates: make(chan plan.TIDBlob, 1),
        readerCmdQueue: make(chan uint32, 1),
    }

    cs.SetLogLabel(fmt.Sprint("ChSess_", cs.ChSessID))

    err = cs.CtxStart(
        cs.ctxStartup,
        cs.ctxAboutToStop,
        nil,
        cs.ctxStopping,
    )

    //inMemberSession.CtxAddChild(cs)
 
    return cs, err
}


// ChSession instances correspond to client rpc calls to StartChannelSession().
//
// A ChSession is the bridge between a ChAgent and a MemberSession
type ChSession struct {
    ctx.Context
    //plan.ContextWorker

    ChSessID            ChSessID

    // Parent member session
    MemberSession       *MemberSession

    chSt                *ChStore

    msgInbox            chan *Msg
    readerCmdQueue      chan uint32

    // EntryTIDs that have changed liveness.
    entryUpdates        chan plan.TIDBlob  
}

func (cs *ChSession) ctxStartup() error {

    cs.chSt.attachChSession(cs)

    cs.CtxGo(func() {

        for msg := range cs.msgInbox {
            switch msg.Op {
                case MsgOp_POST_CH_ENTRY:
                    cs.AuthorNewEntry(msg)
                case MsgOp_LATEST_CH_EPOCH:
                    msg.BUF0 = cs.chSt.ExportLatestChEpoch(msg.BUF0)
                    cs.MemberSession.msgOutbox <- msg
                case MsgOp_LATEST_CH_INFO:
                    msg.BUF0 = cs.chSt.ExportLatestChInfo(msg.BUF0)
                    cs.MemberSession.msgOutbox <- msg
                case MsgOp_CLOSE_CH_SESSION:
                    cs.CtxStop("closed by client", nil)
                case MsgOp_RESET_ENTRY_READER:
                    cs.readerCmdQueue <- msg.FLAGS
            }
        }

        // Causes chSessionEntryReader() to exit
        close(cs.entryUpdates)
    })

    cs.CtxGo(func() {
        chSessionEntryReader(cs)
    })

    return nil
}

func (cs *ChSession) ctxAboutToStop() {
    cs.chSt.detachChSession(cs)
}

func (cs *ChSession) ctxStopping() {

    if cs.MemberSession.CtxRunning() {
        cs.MemberSession.msgOutbox <- &Msg{
            Op: MsgOp_CH_SESSION_CLOSED,
            ChSessID: uint32(cs.ChSessID),
        }
    }
    close(cs.msgInbox)
}

// AuthorNewEntry recycles the given repo.Msg, seting its EntryInfo so that it's compliant with the latest ChannelEpoch from the channel associated with this ChSession.
func (cs *ChSession) AuthorNewEntry(
    msg *Msg,
) {

    var (
        err error
    )

    info := msg.EntryInfo

    info.SetTimeAuthored(0)
    info.AuthorSig = nil
    copy(info.AuthorEntryID(), cs.MemberSession.MemberEpoch.EpochTID)

    info.ChannelID = append(info.ChannelID[:0], cs.chSt.ChID()...)

    err = cs.chSt.FetchAuthoringTIDs(
        info.ChannelEpochID(),
        info.ACCEntryID(),
    )

    msg.Op = MsgOp_CH_NEW_ENTRY_READY

    // TODO: on err, return msg w/ err, etc.
    if err != nil {
        msg.Error = err.Error()
    }
    //msg.BUF0 = plan.SmartMarshal(&info, msg.BUF0)
    cs.MemberSession.msgOutbox <- msg

}


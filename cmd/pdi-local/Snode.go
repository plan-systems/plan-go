

package main

import (

    "os"
    "path"
    "io/ioutil"
    //"fmt"
    "sync"
    "time"
    "net"
    crand "crypto/rand"
    //"encoding/hex"
    "encoding/json"

    "github.com/plan-systems/go-plan/pcore"
    "github.com/plan-systems/go-plan/plan"
    "github.com/plan-systems/go-plan/pdi"

    ds "github.com/plan-systems/go-plan/pdi/StorageProviders/datastore"
    _ "github.com/plan-systems/go-plan/ski/CryptoKits/nacl"    

    "github.com/ethereum/go-ethereum/common/hexutil"

    "golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

)


// GenesisParams is entered by humans
type GenesisParams struct {
    CommunityName           string                  `json:"community_name"`
    CommunityID             hexutil.Bytes           `json:"community_id"`
    //GenesisAddr             hexutil.Bytes           `json:"genesis_addr"`
}




const (

    // DefaultGrpcNetworkName is the default net.Listen() network layer name
    //DefaultGrpcNetworkName      = "tcp"

    // DefaultGrpcNetworkAddr is the default net.Listen() local network address
    //DefaultGrpcNetworkAddr      = ":50053"

    // CurrentSnodeVersion specifies the Snode version 
    //CurrentSnodeVersion         = "0.1"

    // ConfigFilename is the file name the root stage config resides in
    ConfigFilename              = "SnodeConfig.json"

)



// Snode represents an instance of a running Snode daemon.  Multiple disk-independent instances
//     can be instantiated and offer service in parallel, this is not typical operation. Rather,
//     one instance runs and hosts service for one or more communities.
type Snode struct {
    plan.Logger

    flow                        plan.Flow

    storesMutex                 sync.RWMutex
    stores                      map[plan.CommunityID]*ds.Store
    
    activeSessions              pcore.SessionGroup

    BasePath                    string
    Config                      Config

    grpcServer                  *grpc.Server
    listener                    net.Listener
    grpcDone                    chan struct{}
}



// Config specifies all operating parameters if a Snode (PLAN's p2p/server node)
type Config struct {

    Name                        string                          `json:"node_name"`
    NodeID                      hexutil.Bytes                   `json:"node_id"`

    StorageConfigs              []ds.StorageConfig              `json:"storage_configs"`

    DefaultFileMode             os.FileMode                     `json:"default_file_mode"`

    GrpcNetworkName             string                          `json:"grpc_network"`
    GrpcNetworkAddr             string                          `json:"grpc_addr"`

    Version                     int32                           `json:"version"`

}




// ApplyDefaults sets std fields and values
func (config *Config) ApplyDefaults() {

    config.DefaultFileMode = plan.DefaultFileMode
    config.GrpcNetworkName = "tcp"
    config.GrpcNetworkAddr = ":50053"
    config.Version = 1

}


// NewSnode creates and initializes a new Snode instance
func NewSnode(
    inBasePath *string,
    inDoInit bool,
) (*Snode, error) {


    sn := &Snode{
        stores: make(map[plan.CommunityID]*ds.Store),
        activeSessions: pcore.NewSessionGroup(),
    }
        
    var err error

    if inBasePath == nil || len(*inBasePath) == 0 {
        sn.BasePath, err = plan.UseLocalDir("pdi-local")
    } else {
        sn.BasePath = *inBasePath
    }
    if err != nil { return nil, err }

    if err = os.MkdirAll(sn.BasePath, plan.DefaultFileMode); err != nil {
        return nil, err
    }

    if err = sn.readConfig(inDoInit); err != nil {
        return nil, err
    }

    sn.SetLogLabel("pdi-local")

    return sn, nil
}

// Startup -- see plan.Flow.Startup()
func (sn *Snode) Startup(inCtx context.Context) (context.Context, error) {

    err := sn.flow.Startup(
        inCtx,
        sn.GetLogLabel(),
        sn.onInternalStartup,
        sn.onInternalShutdown,
    )

    return sn.flow.Ctx, err
}

// Shutdown -- see plan.Flow.Shutdown()
func (sn *Snode) Shutdown(
    inReason string,
) {

    sn.flow.Shutdown(inReason)
}


func (sn *Snode) onInternalStartup() error {

    var err error


    for i := range sn.Config.StorageConfigs {
        info := &sn.Config.StorageConfigs[i]
        
        St := ds.NewStore(info, sn.BasePath)

        sn.registerStore(St)
    }

    for _, St := range sn.stores {
        St.Startup(sn.flow.Ctx, false)
    }

    if err == nil {
        err = sn.startServer()
    }

    return err
}



func (sn *Snode) onInternalShutdown() {

    // Shutdown the Stores FIRST so that all we have to do is wait on the server to stop.
    storesRunning := &sync.WaitGroup{}

    storesRunning.Add(len(sn.stores))
    
    for _, v := range sn.stores {
        St := v
        go func() {
            St.Shutdown(sn.flow.ShutdownReason)
            storesRunning.Done()
        }()
    }

    if sn.grpcServer != nil {
        sn.Info(0, "stopping StorageProvider grpc service")
        sn.grpcServer.GracefulStop()

        _, _ = <- sn.grpcDone
        sn.Info(1, "StorageProvider service stopped")
    }

    storesRunning.Wait()

}




// readConfig uses BasePath to read in the node's config file
func (sn *Snode) readConfig(inFirstTime bool) error {

    pathname := path.Join(sn.BasePath, ConfigFilename)
 
    buf, err := ioutil.ReadFile(pathname)
    if err == nil { 
        err = json.Unmarshal(buf, &sn.Config)
    }
    if inFirstTime {
        if os.IsNotExist(err) {
            sn.Config.ApplyDefaults()
            sn.Config.NodeID = make([]byte, plan.CommunityIDSz)
            crand.Read(sn.Config.NodeID)

            err = sn.writeConfig()
        } else if err == nil {
            err = plan.Errorf(nil, plan.ConfigFailure, "Init failed: node config %v already exists", pathname)
        }
    } else if err != nil {
        err = plan.Errorf(err, plan.ConfigFailure, "Failed to load node config")
    }

    return err
}

// writeConfig writes out the node config file based on BasePath
func (sn *Snode) writeConfig() error {

    buf, err := json.MarshalIndent(&sn.Config, "", "\t")
    if err == nil {
        pathname := path.Join(sn.BasePath, ConfigFilename)

        err = ioutil.WriteFile(pathname, buf, sn.Config.DefaultFileMode)
    }

    if err != nil {
        return plan.Errorf(err, plan.FailedToAccessPath, "Failed to write node config")
    }

    return nil
}


// CreateNewStore creates a new data store and adds it to this nodes list of stores (and updates the config on disk)
func (sn *Snode) CreateNewStore(
    inImplName string,
    inDeposits []*pdi.Transfer,
    inGenesisTxns []pdi.RawTxn,
    inEpoch pdi.StorageEpoch,
) error {
    
    if sn.flow.IsRunning() {
        return plan.Error(nil, plan.AssertFailed, "can't create store while running")
    }

    stConfig := &ds.StorageConfig{
        HomePath: path.Join("datastore", inEpoch.FormSuggestedDirName()),
        ImplName: inImplName,
        StorageEpoch: inEpoch,
    }

    _, err := plan.CreateNewDir(sn.BasePath, stConfig.HomePath)
    if err != nil { return err }

    St := ds.NewStore(
        stConfig, 
        sn.BasePath,
    )
    if err = St.Startup(context.Background(), true); err != nil {
        return err
    }

    if err = St.DepositTransfers(inDeposits); err != nil {
        return err
    }

    // Commit the given txns
    for _, txn := range inGenesisTxns {
        if err != nil {
            break
        }
        err = St.DoCommitJob(ds.CommitJob{
            Txn: pdi.DecodedTxn{
                RawTxn: txn.Bytes,
            },
        })
    }

    if err == nil {
        sn.Config.StorageConfigs = append(sn.Config.StorageConfigs, *stConfig)

        if err = sn.writeConfig(); err != nil {
            return err
        }
    }

    // Sleep a little so the log messages show up in a nice order for such an important occasion!
    time.Sleep(100 * time.Millisecond)

    St.Shutdown("creation complete")

    return nil
}

func (sn *Snode) registerStore(St *ds.Store) {
   
    communityID := plan.GetCommunityID(St.Config.StorageEpoch.CommunityID)

    sn.storesMutex.Lock()
    sn.stores[communityID] = St
    sn.storesMutex.Unlock()
}

func (sn *Snode) fetchStore(inCommunityID []byte) *ds.Store {

    communityID := plan.GetCommunityID(inCommunityID)

    sn.storesMutex.RLock()
    St := sn.stores[communityID]
    sn.storesMutex.RUnlock()

    return St

}


func (sn *Snode) startServer() error {

    sn.grpcDone = make(chan struct{})

    sn.Infof(0, "starting StorageProvider service on %v %v", sn.Config.GrpcNetworkName, sn.Config.GrpcNetworkAddr)
    listener, err := net.Listen(sn.Config.GrpcNetworkName, sn.Config.GrpcNetworkAddr)
    if err != nil {
        return err
    }

    // TODO: turn off compression since we're dealing w/ encrypted data
    sn.grpcServer = grpc.NewServer()
    pdi.RegisterStorageProviderServer(sn.grpcServer, sn)
    
    // Register reflection service on gRPC server.
    reflection.Register(sn.grpcServer)
    go func() {

        if err := sn.grpcServer.Serve(listener); err != nil {
            sn.flow.LogErr(err, "grpc Serve() failed")
        }
        
        listener.Close()

        close(sn.grpcDone)
    }()

    return nil
}

// StartSession -- see service StorageProvider in pdi.proto
func (sn *Snode) StartSession(ctx context.Context, in *pdi.SessionReq) (*pdi.StorageInfo, error) {
    if in.StorageEpoch == nil {
        return nil, plan.Errorf(nil, plan.ParamMissing, "missing StorageEpoch")
    }

    St := sn.fetchStore(in.StorageEpoch.CommunityID)
    if St == nil {
        return nil, plan.Errorf(nil, plan.CommunityNotFound, "community not found: %v", in.StorageEpoch.CommunityID)
    }

    // TODO security checks to prevent DoS
    session :=  sn.activeSessions.NewSession(ctx, nil)
    session.Cookie = St

    info := &pdi.StorageInfo{
    }

	return info, nil
}

// FetchSessionStore uses the metadata in the given session to recover the session info and associated Store
func (sn *Snode) FetchSessionStore(ctx context.Context) (*ds.Store, error) {
    session, err := sn.activeSessions.FetchSession(ctx)
    if err != nil {
        return nil, err
    }

    St, _ := session.Cookie.(*ds.Store)
    if St == nil {
        return nil, plan.Errorf(nil, plan.AssertFailed, "internal type assertion err")
    }

    err = St.CheckStatus()
    if err != nil {
        return nil, err
    }

    return St, nil
}

// Scan -- see service StorageProvider in pdi.proto
func (sn *Snode) Scan(inScanPB *pdi.TxnScan, inOutlet pdi.StorageProvider_ScanServer) error {
    St, err := sn.FetchSessionStore(inOutlet.Context())
    if err != nil {
        return err
    }

    job := ds.ScanJob{
        TxnScan:   inScanPB,
        Outlet:    inOutlet,
        OnComplete: make(chan error),
    }
    
    St.DoScanJob(job)

    err = <-job.OnComplete
    return err
}

// FetchTxns -- see service StorageProvider in pdi.proto
func (sn *Snode) FetchTxns(inTxnList *pdi.TxnList, inOutlet pdi.StorageProvider_FetchTxnsServer) error {
    St, err := sn.FetchSessionStore(inOutlet.Context())
    if err != nil {
        return err
    }
        
    job := ds.SendJob{
        URIDs:     inTxnList.URIDs,
        Outlet:    inOutlet,
        OnComplete: make(chan error),
    }

    St.DoSendJob(job)

    err = <-job.OnComplete
    return err
}

// CommitTxn -- see service StorageProvider in pdi.proto
func (sn *Snode) CommitTxn(ctx context.Context, inRawTxn *pdi.RawTxn) (*plan.Status, error) {
    St, err := sn.FetchSessionStore(ctx)
    if err != nil {
        return nil, err
    }

    err = St.DoCommitJob(ds.CommitJob{
        Txn: pdi.DecodedTxn{
            RawTxn: inRawTxn.Bytes,
        },
    })

    if err != nil {
        return nil, err
    }

    return &plan.Status{}, nil
}


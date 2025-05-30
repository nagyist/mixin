package storage

import (
	"github.com/MixinNetwork/mixin/common"
	"github.com/MixinNetwork/mixin/crypto"
)

type Store interface {
	Close() error

	CheckGenesisLoad(snapshots []*common.SnapshotWithTopologicalOrder) (bool, error)
	LoadGenesis(rounds []*common.Round, snapshots []*common.SnapshotWithTopologicalOrder, transactions []*common.VersionedTransaction) error
	ReadAssetWithBalance(id crypto.Hash) (*common.Asset, common.Integer, error)
	ReadAllNodes(threshold uint64, withState bool) []*common.Node
	AddNodeOperation(tx *common.VersionedTransaction, timestamp, threshold uint64, finalized bool) error
	ReadTransaction(hash crypto.Hash) (*common.VersionedTransaction, string, error)
	WriteTransaction(tx *common.VersionedTransaction) error
	StartNewRound(node crypto.Hash, number uint64, references *common.RoundLink, finalStart uint64) error
	UpdateEmptyHeadRound(node crypto.Hash, number uint64, references *common.RoundLink) error
	LastSnapshot() (*common.SnapshotWithTopologicalOrder, *common.VersionedTransaction)
	WriteConsensusSnapshot(snap *common.Snapshot, tx *common.VersionedTransaction, hack *common.Snapshot) error
	ReadLastConsensusSnapshot() (*common.Snapshot, error)

	ReadUTXOKeys(hash crypto.Hash, index uint) (*common.UTXOKeys, error)
	ReadUTXOLock(hash crypto.Hash, index uint) (*common.UTXOWithLock, error)
	LockUTXOs(inputs []*common.Input, tx crypto.Hash, fork bool) error
	ReadDepositLock(deposit *common.DepositData) (crypto.Hash, error)
	LockDepositInput(deposit *common.DepositData, tx crypto.Hash, fork bool) error
	ReadWithdrawalClaim(hash crypto.Hash) (*common.VersionedTransaction, string, error)
	ReadGhostKeyLock(key crypto.Key) (*crypto.Hash, error)
	LockGhostKeys(keys []*crypto.Key, tx crypto.Hash, fork bool) error
	ReadSnapshot(hash crypto.Hash) (*common.SnapshotWithTopologicalOrder, error)
	ReadSnapshotsSinceTopology(offset, count uint64) ([]*common.SnapshotWithTopologicalOrder, error)
	ReadSnapshotWithTransactionsSinceTopology(topologyOffset, count uint64) ([]*common.SnapshotWithTopologicalOrder, []*common.VersionedTransaction, error)
	ReadSnapshotsForNodeRound(nodeIdWithNetwork crypto.Hash, round uint64) ([]*common.SnapshotWithTopologicalOrder, error)
	ReadRound(hash crypto.Hash) (*common.Round, error)
	ReadLink(from, to crypto.Hash) (uint64, error)
	WriteSnapshot(*common.SnapshotWithTopologicalOrder, []crypto.Hash) error
	ReadCustodian(ts uint64) (*common.CustodianUpdateRequest, error)
	ListCustodianUpdates() ([]*common.CustodianUpdateRequest, error)

	CachePutTransaction(tx *common.VersionedTransaction) error
	CacheGetTransaction(hash crypto.Hash) (*common.VersionedTransaction, error)
	CacheRetrieveTransactions(limit int) ([]*common.VersionedTransaction, error)
	CacheRemoveTransactions([]crypto.Hash) error

	ReadLastMintDistribution(batch uint64) (*common.MintDistribution, error)
	LockMintInput(mint *common.MintData, tx crypto.Hash, fork bool) error
	ReadMintDistributions(offset, count uint64) ([]*common.MintDistribution, []*common.VersionedTransaction, error)
	ReadSnapshotWorksForNodeRound(nodeId crypto.Hash, round uint64) ([]*common.SnapshotWork, error)
	ListWorkOffsets(cids []crypto.Hash) (map[crypto.Hash]uint64, error)
	ListNodeWorks(cids []crypto.Hash, day uint32) (map[crypto.Hash][2]uint64, error)
	ReadWorkOffset(nodeId crypto.Hash) (uint64, error)
	WriteRoundWork(nodeId crypto.Hash, round uint64, snapshots []*common.SnapshotWork, credit bool) error

	ReadRoundSpaceCheckpoint(nodeId crypto.Hash) (uint64, uint64, error)
	WriteRoundSpaceAndState(space *common.RoundSpace) error
	ListAggregatedRoundSpaceCheckpoints(cids []crypto.Hash) (map[crypto.Hash]*common.RoundSpace, error)
	ReadNodeRoundSpacesForBatch(nodeId crypto.Hash, batch uint64) ([]*common.RoundSpace, error)

	RemoveGraphEntries(prefix string) (int, error)
	ValidateGraphEntries(networkId crypto.Hash, depth uint64) (int, int, error)
}

package tracker

import (
	"log"

	"github.com/trichner/eveapi"
)

const (
	// Limit of WalletJournal entries to fetch at once
	fetchChunkSize = 50

	// Default account key, this is the corp main wallet
	defaultAccountKey = 1000
)

// Replicator represents a service that will fetch the WalletJournal periodically and
// store the new entries into the repo
type Replicator struct {
	api        eveapi.API
	repo       Repo
	accountKey int64
}

type ReplicatorConfig func(r *Replicator) (error)

func ApiCredentials(key string, vcode string) ReplicatorConfig {
	return func(r *Replicator) error {
		r.api.APIKey = eveapi.Key{ID: key, VCode: vcode}
		return nil
	}
}

func AccountKey(accountKey int64) ReplicatorConfig {
	return func(r *Replicator) error {
		r.accountKey = accountKey
		return nil
	}
}

func SqlRepo(r *Replicator) error {
	r.repo = NewSqliteRepo()
	return nil
}

// NewReplicator sets up a new Replicator that may be configured according to ones needs.
// At least the ApiCredentials should be specified. In case there is no repo defined this will use a sqlite3 backed
// database. Furthermore the corp main wallet will be tracked.
func NewReplicator(conf ...ReplicatorConfig) (Replicator, error) {
	r := Replicator{}

	// setup defaults
	r.repo = NewSqliteRepo()
	r.accountKey = defaultAccountKey
	r.api = eveapi.API{
		Server:    eveapi.Tranquility,
		APIKey:    eveapi.Key{},
		UserAgent: "evedt",
		Debug:     false,
	}

	var err error
	for _, c := range conf {
		err = c(&r)
		if err != nil {
			return r, err
		}
	}

	if err = r.repo.Open(); err != nil {
		return r, err
	}

	return r, nil
}

// Run executes the Replicator for once, this fetches all new WalletJournal entries
// and stores them locally
func (r *Replicator) Run() error {

	last := r.repo.LastDonation()

	fromRefID := int64(0)
again:
	for {

		log.Printf("Fetching wallet from refID: %d\n", fromRefID)
		walletJournal, err := r.api.CorpWalletJournal(r.accountKey, fromRefID, fetchChunkSize)
		if err != nil {
			return err
		}

		txs := walletJournal.Transactions

		// No more transactions?
		if len(txs) == 0 {
			break again
		}

		for _, t := range txs {

			// Do we already have this one?
			if last != nil && t.RefID == last.RefID {
				break again
			}

			if t.RefTypeID == 10 { // Player donation
				d := Donation{
					CharacterID:   t.OwnerID1,
					CharacterName: t.OwnerName1,
					RefID:         t.RefID,
					Amount:        t.Amount,
					Date:          t.TransactionDateTime.Time,
				}
				log.Printf("Storing transaction from %s, txId:%s\n", d.CharacterName, d.RefID)
				r.repo.StoreDonation(&d)
			}
		}

		// the lowest 'RefID' will be the next 'fromRefID', which is used as an upper bound
		fromRefID = txs[len(txs)-1].RefID
	}

	return nil
}

func (r *Replicator) Close() error {
	return r.repo.Close()
}

func (r *Replicator) LastDonation() *Donation {
	return r.repo.LastDonation()
}

func (r *Replicator) FindDonations(limit int, days int) Donations {
	return r.repo.FindDonations(limit, days)
}

func (r *Replicator) FindTopDonations(limit int, days int) Donations {
	return r.repo.FindTopDonations(limit, days)
}

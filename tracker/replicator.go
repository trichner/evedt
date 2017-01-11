package tracker

import (
	"log"

	"github.com/trichner/eveapi"
)

const (
	// Limit of WalletJournal entries to fetch at once
	fetchChunkSize = 50
)

// Replicator represents a service that will fetch the WalletJournal periodically and
// store the new entries into the repo
type Replicator struct {
	api        *eveapi.API
	repo       Repo
	accountKey int64
}

// Init initializes a replicator with the appropriate configuration
func (r *Replicator) Init(repo Repo, apiKey string, apiVCode string, accountKey int64) {

	r.accountKey = accountKey
	r.repo = repo

	key := eveapi.Key{ID: apiKey, VCode: apiVCode}
	r.api = &eveapi.API{
		Server:    eveapi.Tranquility,
		APIKey:    key,
		UserAgent: "evedt",
		Debug:     false,
	}

}

// Run executes the replicator for once, this fetches all new WalletJournal entries
// and stores them locally
func (r *Replicator) Run() error {

	last := r.repo.LastDonation()

	fromRefID := int64(0)
Loop:
	for {

		log.Printf("Fetching wallet from refID: %d\n", fromRefID)
		walletJournal, err := r.api.CorpWalletJournal(r.accountKey, fromRefID, fetchChunkSize)
		if err != nil {
			return err
		}

		txs := walletJournal.Transactions

		// No more transactions?
		if len(txs) == 0 {
			break Loop
		}

		for _, t := range txs {

			// Do we already have this one?
			if last != nil && t.RefID == last.RefID {
				break Loop
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

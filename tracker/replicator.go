package tracker

import (
	"log"

	"github.com/trichner/eveapi"
	"strings"
	"os"
	"encoding/json"
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

	f, err := os.OpenFile("corp-wallet-journal.json", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Print(err)
	}

	defer f.Close()

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
		log.Printf("Received %d transactions.\n", len(txs))

		// No more transactions?
		if len(txs) == 0 {
			break again
		}

		for _, t := range txs {

			// Do we already have this one?
			if last != nil && t.RefID == last.RefID {
				break again
			}

			reason := t.Reason
			if t.RefTypeID == 10 || t.RefTypeID == 37 { // Player donation
				reason = strings.TrimPrefix(reason, "DESC: ")
			}
			d := Donation{
				CharacterID:   t.OwnerID1,
				CharacterName: t.OwnerName1,
				RefID:         t.RefID,
				RefTypeID:     t.RefTypeID,
				Amount:        t.Amount,
				Date:          t.TransactionDateTime.Time,
				Reason:        reason,
			}
			log.Printf("Nex TX (%s) refTypeId:%v, txId:%v: %s\n", d.CharacterName, d.RefTypeID, d.RefID, d.Reason)
			r.repo.StoreDonation(&d)

			// also log the entire transaction to a file
			bytes, err := json.Marshal(&t)
			if err != nil {
				continue;
			}

			line := string(bytes) + "\n"
			_, err = f.WriteString(line)
			if err != nil {
				log.Print(err)
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

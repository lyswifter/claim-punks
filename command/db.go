package command

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"

	"github.com/ipfs/go-datastore"
	"github.com/ipfs/go-datastore/query"
	levelds "github.com/ipfs/go-ds-leveldb"
	"github.com/mitchellh/go-homedir"
	ldbopts "github.com/syndtr/goleveldb/leveldb/opt"
	"go.uber.org/multierr"
	"golang.org/x/xerrors"
)

var DB datastore.Batching

func setupLevelDs(path string, readonly bool) (datastore.Batching, error) {
	if _, err := os.ReadDir(path); err != nil {
		if os.IsNotExist(err) {
			//mkdir
			err = os.MkdirAll(path, 0777)
			if err != nil {
				return nil, err
			}
		}
	}

	db, err := levelds.NewDatastore(path, &levelds.Options{
		Compression: ldbopts.NoCompression,
		NoSync:      false,
		Strict:      ldbopts.StrictAll,
		ReadOnly:    readonly,
	})
	if err != nil {
		fmt.Printf("NewDatastore: %s\n", err)
		return nil, err
	}

	return db, err
}

func DataStores() {
	repodir, err := homedir.Expand(repoPath)
	if err != nil {
		return
	}

	idb, err := setupLevelDs(path.Join(repodir, "datastore"), false)
	if err != nil {
		log.Printf("setup infodb: err %s", err)
		return
	}
	DB = idb

	log.Printf("DB: %+v", DB)
}

func saveinfo(info ResultPunk) error {
	key := datastore.NewKey(info.Wal)

	ishas, err := DB.Has(key)
	if err != nil {
		return err
	}

	if ishas {
		return xerrors.New("is exist")
	}

	in, err := json.Marshal(info)
	if err != nil {
		return err
	}

	err = DB.Put(key, in)
	if err != nil {
		return err
	}

	return nil
}

func listinfos() ([]ResultPunk, error) {
	res, err := DB.Query(query.Query{})
	if err != nil {
		return nil, err
	}

	defer res.Close()

	infos := []ResultPunk{}

	var errs error

	for {
		res, ok := res.NextSync()
		if !ok {
			break
		}

		if res.Error != nil {
			return nil, res.Error
		}

		info := &ResultPunk{}
		err := json.Unmarshal(res.Value, info)
		if err != nil {
			errs = multierr.Append(errs, xerrors.Errorf("decoding state for key '%s': %w", res.Key, err))
			continue
		}

		infos = append(infos, *info)
	}

	log.Printf("read infos ok, len %d", len(infos))

	return infos, nil
}

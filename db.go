package main

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/ipfs/go-datastore"
	levelds "github.com/ipfs/go-ds-leveldb"
	"github.com/mitchellh/go-homedir"
	ldbopts "github.com/syndtr/goleveldb/leveldb/opt"
)

const repoPath = "~/.claimpunks"

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

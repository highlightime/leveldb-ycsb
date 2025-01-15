// Copyright 2019 The LevelDB-Go and Pebble Authors. All rights reserved. Use
// of this source code is governed by a BSD-style license that can be found in
// the LICENSE file.

package main

import (
	"log"

	"github.com/cockroachdb/pebble"

	// "github.com/ethereum/go-ethereum/ethdb"
	// "github.com/ethereum/go-ethereum/ethdb/leveldb"
	"github.com/syndtr/goleveldb/leveldb"
)

// DB specifies the minimal interfaces that need to be implemented to support
// the pebble command.
type DB interface {
	NewIter(*pebble.IterOptions) iterator
	NewBatch() batch
	Stat() (string, error)
	Scan(iter iterator, key []byte, count int64, reverse bool) error
	Metrics() *pebble.Metrics
	Flush() error
}

func (d levelDB) NewIter(opts *pebble.IterOptions) iterator {
	return nil
}
func (d levelDB) Stat() (string, error) {
	return "", nil
}
func (d levelDB) Metrics() *pebble.Metrics {
	return nil
}
func (d levelDB) Flush() error {
	return nil
}

type pebbleDB struct {
	d       *pebble.DB
	ballast []byte
}

func (p pebbleDB) Flush() error {
	return nil
}

func (p pebbleDB) NewIter(opts *pebble.IterOptions) iterator {
	return nil
}

func (p pebbleDB) NewBatch() batch {
	return nil
}

func (p pebbleDB) Scan(iter iterator, key []byte, count int64, reverse bool) error {
	return nil
}

func (p pebbleDB) Metrics() *pebble.Metrics {
	return nil
}

func (p pebbleDB) Stat() (string, error) {
	return "", nil
}

type iterator interface {
	SeekLT(key []byte) bool
	SeekGE(key []byte) bool
	Valid() bool
	Key() []byte
	Value() []byte
	First() bool
	Next() bool
	Last() bool
	Prev() bool
	Close() error
}

type levelDBIter struct {
}
type batch interface {
	Close() error
	Commit(opts *pebble.WriteOptions) error
	Set(key, value []byte, opts *pebble.WriteOptions) error
	Delete(key []byte, opts *pebble.WriteOptions) error
	LogData(data []byte, opts *pebble.WriteOptions) error
}

type levelDBBatch struct {
	b *leveldb.Batch
	d levelDB
}
type levelDB struct {
	d *leveldb.DB
}

func (d levelDB) NewBatch() batch {
	return levelDBBatch{b: new(leveldb.Batch), d: d}
}
func (b levelDBBatch) LogData(data []byte, opts *pebble.WriteOptions) error {
	return nil
}
func (b levelDBBatch) Delete(key []byte, opts *pebble.WriteOptions) error {
	return nil
}

func (b levelDBBatch) Commit(opts *pebble.WriteOptions) error {
	return b.d.d.Write(b.b, nil)
}

func (b levelDBBatch) Set(key, value []byte, opts *pebble.WriteOptions) error {
	b.b.Put(key, value)
	return nil
}

func (b levelDBBatch) Close() error {
	return nil
}

func newLevelDB(dir string) DB {
	d, err := leveldb.OpenFile(dir, nil)
	if err != nil {
		log.Fatal(err)
	}
	return levelDB{d: d}
}

func (d levelDB) Scan(iter iterator, key []byte, count int64, reverse bool) error {
	return nil
}

// func newPebbleDB(dir string) DB {
// 	cache := pebble.NewCache(cacheSize)
// 	defer cache.Unref()
// 	opts := &pebble.Options{
// 		Cache:                       cache,
// 		Comparer:                    &cockroachkvs.Comparer,
// 		DisableWAL:                  disableWAL,
// 		FormatMajorVersion:          pebble.FormatNewest,
// 		KeySchema:                   cockroachkvs.KeySchema.Name,
// 		KeySchemas:                  sstable.MakeKeySchemas(&cockroachkvs.KeySchema),
// 		L0CompactionThreshold:       2,
// 		L0StopWritesThreshold:       1000,
// 		LBaseMaxBytes:               64 << 20, // 64 MB
// 		Levels:                      make([]pebble.LevelOptions, 7),
// 		MaxOpenFiles:                16384,
// 		MemTableSize:                64 << 20,
// 		MemTableStopWritesThreshold: 4,
// 		Merger: &pebble.Merger{
// 			Name: "cockroach_merge_operator",
// 		},
// 		MaxConcurrentCompactions: func() int {
// 			return 3
// 		},
// 	}
// 	// In FormatColumnarBlocks (the value of FormatNewest at the time of
// 	// writing), columnar blocks are only written if explicitly opted into.
// 	opts.Experimental.EnableColumnarBlocks = func() bool { return true }

// 	for i := 0; i < len(opts.Levels); i++ {
// 		l := &opts.Levels[i]
// 		l.BlockSize = 32 << 10       // 32 KB
// 		l.IndexBlockSize = 256 << 10 // 256 KB
// 		l.FilterPolicy = bloom.FilterPolicy(10)
// 		l.FilterType = pebble.TableFilter
// 		if i > 0 {
// 			l.TargetFileSize = opts.Levels[i-1].TargetFileSize * 2
// 		}
// 		l.EnsureDefaults()
// 	}
// 	opts.Levels[6].FilterPolicy = nil
// 	opts.FlushSplitBytes = opts.Levels[0].TargetFileSize

// 	opts.EnsureDefaults()

// 	if verbose {
// 		lel := pebble.MakeLoggingEventListener(nil)
// 		opts.EventListener = &lel
// 		opts.EventListener.TableDeleted = nil
// 		opts.EventListener.TableIngested = nil
// 		opts.EventListener.WALCreated = nil
// 		opts.EventListener.WALDeleted = nil
// 	}

// 	if pathToLocalSharedStorage != "" {
// 		opts.Experimental.RemoteStorage = remote.MakeSimpleFactory(map[remote.Locator]remote.Storage{
// 			// Store all shared objects on local disk, for convenience.
// 			"": remote.NewLocalFS(pathToLocalSharedStorage, vfs.Default),
// 		})
// 		opts.Experimental.CreateOnShared = remote.CreateOnSharedAll
// 		if secondaryCacheSize != 0 {
// 			opts.Experimental.SecondaryCacheSizeBytes = secondaryCacheSize
// 		}
// 	}

// 	p, err := pebble.Open(dir, opts)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	if pathToLocalSharedStorage != "" {
// 		if err := p.SetCreatorID(1); err != nil {
// 			log.Fatal(err)
// 		}
// 	}
// 	return pebbleDB{
// 		d:       p,
// 		ballast: make([]byte, 1<<30),
// 	}
// }

// func (p pebbleDB) Flush() error {
// 	return p.d.Flush()
// }

// func (p pebbleDB) NewIter(opts *pebble.IterOptions) iterator {
// 	iter, _ := p.d.NewIter(opts)
// 	return iter
// }

// func (p pebbleDB) NewBatch() batch {
// 	return p.d.NewBatch()
// }

// func (p pebbleDB) Scan(iter iterator, key []byte, count int64, reverse bool) error {
// 	var data bytealloc.A
// 	if reverse {
// 		for i, valid := 0, iter.SeekLT(key); valid; valid = iter.Prev() {
// 			data, _ = data.Copy(iter.Key())
// 			data, _ = data.Copy(iter.Value())
// 			i++
// 			if i >= int(count) {
// 				break
// 			}
// 		}
// 	} else {
// 		for i, valid := 0, iter.SeekGE(key); valid; valid = iter.Next() {
// 			data, _ = data.Copy(iter.Key())
// 			data, _ = data.Copy(iter.Value())
// 			i++
// 			if i >= int(count) {
// 				break
// 			}
// 		}
// 	}
// 	return nil
// }

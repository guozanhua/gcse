package main

import (
	"log"
	"runtime"

	"github.com/daviddengcn/gcse"
	"github.com/daviddengcn/sophie"
	"github.com/daviddengcn/sophie/kv"
)

func clearOutdatedIndex() error {
	segm, err := gcse.IndexSegments.FindMaxDone()
	if err != nil {
		return err
	}
	all, err := gcse.IndexSegments.ListAll()
	if err != nil {
		return err
	}

	for _, s := range all {
		if s == segm {
			continue
		}

		err := s.Remove()
		if err != nil {
			return err
		}
		log.Printf("Outdated segment %v removed!", s)
	}

	return nil
}

func doIndex() bool {
	idxSegm, err := gcse.IndexSegments.GenMaxSegment()
	if err != nil {
		log.Printf("GenMaxSegment failed: %v", err)
		return false
	}

	runtime.GC()
	gcse.DumpMemStats()

	log.Printf("Indexing to %v ...", idxSegm)

	fpDocDB := sophie.LocalFsPath(gcse.DocsDBPath.S())

	ts, err := gcse.Index(kv.DirInput(fpDocDB))
	if err != nil {
		log.Printf("Indexing failed: %v", err)
		return false
	}

	f, err := idxSegm.Join(gcse.IndexFn).Create()
	if err != nil {
		log.Printf("Create index file failed: %v", err)
		return false
	}
	//defer f.Close()
	log.Printf("Saving index to %v ...", idxSegm)
	if err := ts.Save(f); err != nil {
		log.Printf("ts.Save failed: %v", err)
		return false
	}
	f.Close()
	f = nil
	runtime.GC()
	gcse.DumpMemStats()

	if err := idxSegm.Done(); err != nil {
		log.Printf("segm.Done failed: %v", err)
		return false
	}

	log.Printf("Indexing success: %s (%d)", idxSegm, ts.DocCount())

	ts = nil
	gcse.DumpMemStats()
	runtime.GC()
	gcse.DumpMemStats()

	return true
}

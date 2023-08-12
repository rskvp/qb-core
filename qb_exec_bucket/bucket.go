package qb_exec_bucket

import (
	"github.com/rskvp/qb-core/qb_utils"
)

var Bucket *BucketHelper

const (
	wpName = "bucket"
)

func init() {
	Bucket = NewBucketHelper()
}

// BucketHelper
// main executable container
type BucketHelper struct {
	root string // all buckets are created under this path
}

func NewBucketHelper() (instance *BucketHelper) {
	instance = new(BucketHelper)
	instance.SetRoot(".")

	return
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *BucketHelper) SetRoot(root string) {
	qb_utils.Paths.GetWorkspace(wpName).SetPath(root)
	instance.root = qb_utils.Paths.GetWorkspace(wpName).GetPath()
}

//----------------------------------------------------------------------------------------------------------------------
//	b u i l d e r
//----------------------------------------------------------------------------------------------------------------------

func (instance *BucketHelper) NewBucket(execPath string, global bool) *BucketExec {
	bucket := NewBucketExec(instance.root, execPath, global)
	return bucket
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

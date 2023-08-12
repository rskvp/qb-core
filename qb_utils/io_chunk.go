package qb_utils

import (
	"math"
	"os"
)

const (
	ChunkSizeSmall  = int64(256 * 1024)      // 256k
	ChunkSizeMedium = int64(512 * 1024)      // 512k
	ChunkSizeLarge  = int64(1 * 1024 * 1024) // 1Mb
)

//----------------------------------------------------------------------------------------------------------------------
//	FileChunker
//----------------------------------------------------------------------------------------------------------------------

type FileChunker struct {
	chunkSize int64
}

func NewFileChunkerSmall() (instance *FileChunker) {
	return NewFileChunker(ChunkSizeSmall)
}

func NewFileChunkerMedium() (instance *FileChunker) {
	return NewFileChunker(ChunkSizeMedium)
}

func NewFileChunkerLarge() (instance *FileChunker) {
	return NewFileChunker(ChunkSizeLarge)
}

func NewFileChunker(chunkSizeBytes int64) (instance *FileChunker) {
	instance = new(FileChunker)
	instance.chunkSize = chunkSizeBytes

	return
}

//----------------------------------------------------------------------------------------------------------------------
//	p u b l i c
//----------------------------------------------------------------------------------------------------------------------

func (instance *FileChunker) CalculateChunks(filename string) (fileSize int64, totalPartsNum uint64, err error) {
	file, e := os.Open(filename)
	if e != nil {
		err = e
		return
	}
	defer file.Close()

	fileSize, totalPartsNum = instance.calculateChunks(file)
	return
}

func (instance *FileChunker) SplitAll(filename string) (response [][]byte, err error) {
	if nil != instance {
		err = instance.SplitWalk(filename, func(data []byte) {
			response = append(response, data)
		})
	}
	return
}

func (instance *FileChunker) SplitWalk(filename string, callback func(data []byte)) (err error) {
	if nil != instance && nil != callback {
		file, e := os.Open(filename)
		if e != nil {
			err = e
			return
		}
		defer file.Close()

		fileSize, totalPartsNum := instance.calculateChunks(file)

		for i := uint64(0); i < totalPartsNum; i++ {
			partSize := int(math.Min(float64(instance.chunkSize), float64(fileSize-int64(i*uint64(instance.chunkSize)))))
			partBuffer := make([]byte, partSize)
			_, err = file.Read(partBuffer)
			if nil != err {
				return
			}
			callback(partBuffer)
		}
	}
	return
}

//----------------------------------------------------------------------------------------------------------------------
//	p r i v a t e
//----------------------------------------------------------------------------------------------------------------------

func (instance *FileChunker) calculateChunks(file *os.File) (fileSize int64, count uint64) {
	fileInfo, _ := file.Stat()
	fileSize = fileInfo.Size()

	count = uint64(math.Ceil(float64(fileSize) / float64(instance.chunkSize)))

	return
}

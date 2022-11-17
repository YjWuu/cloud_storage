package meta

import (
	"cloud_storage/fileStore_mysql/db"
	"sort"
)

// FileMeta 文件元信息结构
type FileMeta struct {
	FileSha1 string //唯一标识
	FileName string
	FileSize int64
	Location string //文件路径
	UploadAt string //时间戳
}

var fileMetas map[string]FileMeta

func init() {
	fileMetas = make(map[string]FileMeta)
}

// UpdateFileMeta 新增/更新文件元信息
func UpdateFileMeta(fmeta FileMeta) {
	fileMetas[fmeta.FileSha1] = fmeta
}

// UpdateFileMetaDB  新增/更新文件元信息到mysql中
func UpdateFileMetaDB(fmeta FileMeta) bool {
	return db.OnFileUploadFinished(fmeta.FileSha1, fmeta.FileName, fmeta.FileSize, fmeta.Location)
}

// GetFileMeta 通过sha1值获取文件的元信息对象
func GetFileMeta(fileSha1 string) FileMeta {
	return fileMetas[fileSha1]
}

// GetFileMetaDB 通过sha1值从mysql获取文件的元信息对象
func GetFileMetaDB(fileSha1 string) (*FileMeta, error) {
	tfile, err := db.GetFileMeta(fileSha1)
	if err != nil {
		return nil, err
	}
	fmeta := FileMeta{
		FileSha1: tfile.FileHash,
		FileName: tfile.FileName.String,
		FileSize: tfile.FileSize.Int64,
		Location: tfile.FileAddr.String,
	}
	return &fmeta, nil

}

// GetLastFileMeta 获取批量的文件元信息列表
func GetLastFileMeta(count int) []FileMeta {
	fMetaArray := make([]FileMeta, len(fileMetas))
	for _, v := range fileMetas {
		fMetaArray = append(fMetaArray, v)
	}
	sort.Sort(ByUploadTime{})
	return fMetaArray[0:count]
}

// RemoveFileMeta 删除元信息
func RemoveFileMeta(filesha1 string) {
	delete(fileMetas, filesha1)
}

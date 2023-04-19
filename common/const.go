package common

const (
	MetaDateDirName   = ".fe_meta_data"
	ExtraInfoFilename = "extra_info.yml"

	FileTypeDir     byte = 1
	FileTypeRegular byte = 2
	FileTypeDriver  byte = 3

	TabTypeFile       = "File"
	TabTypeExtraInfo  = "Extra Info"
	TabTypeFileCreate = "Create"
	TabTypeOpList     = "Operate List"

	FavoritesNotExistNotify byte = 0
	FavoritesNotExistDelete byte = 1
	FavoritesNotExistIgnore byte = 2
)

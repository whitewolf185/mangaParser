package config

var SliceOfFileExtensions = []string{".png", ".jpg"}


const (
	// ParentPathToDownloadPattern паттерн для определения пути скачивания манги. Первым параметром должен приниматься personID, 2-ой -- manga name
	ParentPathToDownloadPattern = "./data/%s/%s"
	// ParentPathPdfResultPattern паттерн для определения пути для папки формирования итогового pdf файла.
	// 1-ым аргументом приходит personID
	ParentPathPdfFolderPattern = "./data/%s/result"
)


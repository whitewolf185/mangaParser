package mangalib

// структура управления получения ссылок на картинки из mangalib
type mangaLibController struct {
	urlGetter UrlGetter
}

// NewMangaLibController конструктор mangaLibController
func NewMangaLibController(urlGetter UrlGetter) mangaLibController {
	return mangaLibController{
		urlGetter: urlGetter,
	}
}

// структура, которая используется для парсинга данных об имеющихся картинках в главе для получения ссылок на них
type pageList struct {
	Page int `json:"p"`
	Url string `json:"u"`
}
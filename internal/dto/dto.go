package dto

// Book нужна для передачи из репозиториев в юзкейс данных по книге. Нужна,
// так как может иметь место ситуация, когда для получения одной entity.Book
// нужно обратиться в разные репозитории. Соответственно, из нескольких DTO
// затем будет собрана одна entity.Book.
type Book struct {
	ID            int
	Author        string
	Title         string
	PublisherYear int
}

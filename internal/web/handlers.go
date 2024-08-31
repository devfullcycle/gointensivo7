package web

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"books/internal/service"
)

// BookHandlers lida com as requisições HTTP relacionadas a livros.
type BookHandlers struct {
	service *service.BookService
}

// NewBookHandlers cria uma nova instância de BookHandlers.
func NewBookHandlers(service *service.BookService) *BookHandlers {
	return &BookHandlers{service: service}
}

// GetBooks lida com a requisição GET /books.
func (h *BookHandlers) GetBooks(w http.ResponseWriter, r *http.Request) {
	books, err := h.service.GetBooks()
	if err != nil {
		log.Println("Error fetching books from database:", err)
		http.Error(w, "failed to get books", http.StatusInternalServerError)
		return
	}

	if len(books) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

// CreateBook lida com a requisição POST /books.
func (h *BookHandlers) CreateBook(w http.ResponseWriter, r *http.Request) {
	var book service.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		log.Println("Error decoding request payload:", err)
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	if err := h.service.CreateBook(&book); err != nil {
		log.Println("Error creating book in database:", err)
		http.Error(w, "failed to create book", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

// GetBookByID lida com a requisição GET /books/{id}.
func (h *BookHandlers) GetBookByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("Invalid book ID:", idStr, "Error:", err)
		http.Error(w, "invalid book ID", http.StatusBadRequest)
		return
	}

	book, err := h.service.GetBookByID(id)
	if err != nil {
		log.Println("Error fetching book with ID", id, "from database:", err)
		http.Error(w, "failed to get book", http.StatusInternalServerError)
		return
	}
	if book == nil {
		log.Println("Book with ID", id, "not found")
		http.Error(w, "book not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(book)
}

// UpdateBook lida com a requisição PUT /books/{id}.
func (h *BookHandlers) UpdateBook(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("Invalid book ID:", idStr, "Error:", err)
		http.Error(w, "invalid book ID", http.StatusBadRequest)
		return
	}

	var book service.Book
	if err := json.NewDecoder(r.Body).Decode(&book); err != nil {
		log.Println("Error decoding request payload:", err)
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}
	book.ID = id

	if err := h.service.UpdateBook(&book); err != nil {
		log.Println("Error updating book with ID", id, "in database:", err)
		http.Error(w, "failed to update book", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(book)
}

// DeleteBook lida com a requisição DELETE /books/{id}.
func (h *BookHandlers) DeleteBook(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Println("Invalid book ID:", idStr, "Error:", err)
		http.Error(w, "invalid book ID", http.StatusBadRequest)
		return
	}

	if err := h.service.DeleteBook(id); err != nil {
		log.Println("Error deleting book with ID", id, "from database:", err)
		http.Error(w, "failed to delete book", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *BookHandlers) SimulateReading(w http.ResponseWriter, r *http.Request) {
	var request struct {
		BookIDs []int `json:"book_ids"`
	}

	// Decodifica o JSON recebido no corpo da requisição
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if len(request.BookIDs) == 0 {
		http.Error(w, "No book IDs provided", http.StatusBadRequest)
		return
	}

	// Chama o serviço para simular a leitura de múltiplos livros
	response := h.service.SimulateMultipleReadings(request.BookIDs, 2*time.Second)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

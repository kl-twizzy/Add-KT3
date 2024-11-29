package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gorilla/mux"
)

// Структура для животных
type Animal struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

var animals = map[int]Animal{
	1: {ID: 1, Name: "Лев"},
	2: {ID: 2, Name: "Тигр"},
}

// Обработчик для главной страницы
func homeHandler(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/index.html"))
	tmpl.Execute(w, nil)
}

// Обработчик для получения списка животных
func getAnimals(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(animals)
}

// Обработчик для создания нового животного
func createAnimal(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Парсим форму
		err := r.ParseForm()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Получаем имя животного из формы
		name := r.FormValue("name")
		if name == "" {
			http.Error(w, "Имя животного не может быть пустым", http.StatusBadRequest)
			return
		}

		// Создаем новое животное
		newAnimal := Animal{
			ID:   len(animals) + 1,
			Name: name,
		}
		animals[newAnimal.ID] = newAnimal

		// Перенаправляем пользователя на главную страницу
		http.Redirect(w, r, "/", http.StatusSeeOther)
	} else {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}

// Обработчик для получения информации о животном по ID
func getAnimalByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	animalID := 0
	fmt.Sscanf(id, "%d", &animalID)

	animal, exists := animals[animalID]
	if !exists {
		http.Error(w, "Животное не найдено", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(animal)
}

// Обработчик для обновления информации о животном
func updateAnimal(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	animalID := 0
	fmt.Sscanf(id, "%d", &animalID)

	var updatedAnimal Animal
	err := json.NewDecoder(r.Body).Decode(&updatedAnimal)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	updatedAnimal.ID = animalID
	animals[animalID] = updatedAnimal

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(updatedAnimal)
}

// Обработчик для удаления животного
func deleteAnimal(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	animalID := 0
	fmt.Sscanf(id, "%d", &animalID)

	_, exists := animals[animalID]
	if !exists {
		http.Error(w, "Животное не найдено", http.StatusNotFound)
		return
	}

	delete(animals, animalID)
	w.WriteHeader(http.StatusNoContent)
}

func main() {
	r := mux.NewRouter()

	// Маршрутизация для статических файлов
	fs := http.FileServer(http.Dir("static/"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	// Маршрутизация для HTML-шаблонов
	r.HandleFunc("/", homeHandler)

	// Маршрутизация для REST API
	r.HandleFunc("/animals", getAnimals).Methods("GET")
	r.HandleFunc("/animals", createAnimal).Methods("POST")
	r.HandleFunc("/animals/{id}", getAnimalByID).Methods("GET")
	r.HandleFunc("/animals/{id}", updateAnimal).Methods("PUT")
	r.HandleFunc("/animals/{id}", deleteAnimal).Methods("DELETE")

	fmt.Println("Сервер запущен на порту :8080")
	http.ListenAndServe(":8080", r)
}

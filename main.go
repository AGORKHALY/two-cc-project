package main

import (
	"log"
	"net/http"
	"os"

	"github.com/AGORKHALY/two-cc-project/models"
	"github.com/AGORKHALY/two-cc-project/storage"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

// MARK: Repository
type Repository struct {
	DB *gorm.DB
}

// MARK: Routes
func (r *Repository) SetupRoutes(app *fiber.App) {
	api := app.Group("/api")
	api.Post("/create_cars", r.CreateCars)
	api.Delete("/delete_cars/:id", r.DeleteCars)
	api.Get("/get_cars/:id", r.GetCarById)
	api.Get("/cars", r.GetCars)
}

// MARK: Car struct
type Car struct {
	Company string `json:"company"`
	Model   string `json:"model"`
	Color   string `json:"color"`
}

// MARK: CreateCar
func (r *Repository) CreateCars(c *fiber.Ctx) error {
	car := Car{}

	err := c.BodyParser(&car)

	if err != nil {
		c.Status(http.StatusUnprocessableEntity).JSON(
			&fiber.Map{"message": "request failed"})
		return err
	}

	err = r.DB.Create(&car).Error

	if err != nil {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "an error occurred while creating the car"})
		return err
	}

	c.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "car created successfully", "data": car})

	return nil

}

// MARK: DeleteBook
func (r *Repository) DeleteCars(c *fiber.Ctx) error {
	carModel := models.Cars{}

	id := c.Params("id")

	if id == "" {
		c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": "id cannot be empty"})
		return nil
	}

	err := r.DB.Delete(carModel, id)

	if err.Error != nil {
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "could not delete car"})
		return err.Error
	}

	c.Status(http.StatusOK).JSON(&fiber.Map{"message": "car deleted successfully"})

	return nil
}

// MARK: GetBookById
func (r *Repository) GetCarById(c *fiber.Ctx) error {
	carModel := models.Cars{}
	id := c.Params("id")

	if id == "" {
		c.Status(http.StatusInternalServerError).JSON(
			&fiber.Map{"message": "id cannot be empty"})
		return nil
	}

	err := r.DB.First(&carModel, id)

	if err.Error != nil {
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{"message": "could not find car"})
		return err.Error
	}

	c.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "car found successfully", "data": carModel})

	return nil
}

// MARK: Get All Cars
func (r *Repository) GetCars(c *fiber.Ctx) error {
	carModels := &[]models.Cars{}

	err := r.DB.Find(carModels).Error

	if err != nil {
		c.Status(http.StatusBadRequest).JSON(
			&fiber.Map{"message": "an error occurred while fetching the cars"})
		return err
	}

	c.Status(http.StatusOK).JSON(
		&fiber.Map{"message": "cars fetched successfully", "data": carModels})
	return nil
}

// MARK: Main
func main() {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//MARK: Database connection
	config := &storage.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASSWORD"),
		User:     os.Getenv("DB_USER"),
		SSLMode:  os.Getenv("DB_SSL_MODE"),
		DBName:   os.Getenv("DB_NAME"),
	}

	db, err := storage.NewConnection(config)

	if err != nil {
		log.Fatal("Error connecting to database")
	}

	err = models.MigrateCars(db)

	if err != nil {
		log.Fatal("Error migrating database")
	}

	r := Repository{
		DB: db,
	}

	app := fiber.New()
	r.SetupRoutes(app)
	app.Listen(":8080")
}

package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Appointment represents a scheduled appointment
type Appointment struct {
	ID          uint `gorm:"primaryKey"`
	ClientName  string
	DateTime    time.Time
	Description string
	Confirmed   bool
}

// AppointmentService handles operations related to appointments
type AppointmentService struct {
	DB *gorm.DB
}

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Get database credentials from environment variables
	username := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	database := os.Getenv("DB_DATABASE")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	// Construct DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		username, password, dbHost, dbPort, database)

	// Connect to MySQL database
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}

	// Drop the appointments table if it exists
	err = db.Migrator().DropTable(&Appointment{})
	if err != nil {
		log.Fatalf("Error dropping table: %v", err)
	}

	// Auto-migrate schema
	if err := db.Migrator().AutoMigrate(&Appointment{}); err != nil {
		log.Fatalf("Error migrating database schema: %v", err)
	}

	// Initialize AppointmentService
	appointmentService := AppointmentService{DB: db}

	// Schedule multiple appointments
	appointmentsToSchedule := []struct {
		ClientName  string
		DateTime    time.Time
		Description string
		Confirmed   bool // Whether the appointment is confirmed
	}{
		{"John Doe", time.Now().AddDate(0, 0, 7), "Meeting with client", true},
		{"Jane Smith", time.Now().AddDate(0, 0, 8), "Lunch meeting", false},
		{"James", time.Now().AddDate(0, 0, 9), "Strategy session", true},
		{"Victor", time.Now().AddDate(0, 0, 10), "Product demo", false},
		{"Bolu", time.Now().AddDate(0, 0, 11), "Team meeting", true},
		{"Ben Carson", time.Now().AddDate(0, 0, 12), "Client consultation", false},
		{"Ven Peril", time.Now().AddDate(0, 0, 13), "Marketing presentation", true},
		{"Tinubu", time.Now().AddDate(0, 0, 14), "Project kickoff", false},
		{"Ambode", time.Now().AddDate(0, 0, 15), "Budget review", true},
		{"Stephen", time.Now().AddDate(0, 0, 16), "Training workshop", false},
		{"Ruth", time.Now().AddDate(0, 0, 17), "Business development meeting", true},
		{"Mark", time.Now().AddDate(0, 0, 18), "Client pitch", false},
		{"Sarah", time.Now().AddDate(0, 0, 19), "Team brainstorming", true},
		{"Alon", time.Now().AddDate(0, 0, 20), "Conference call", false},
		{"Karma", time.Now().AddDate(0, 0, 21), "Client onboarding", true},
		{"Venis", time.Now().AddDate(0, 0, 22), "Sales strategy meeting", false},
		{"Micheal", time.Now().AddDate(0, 0, 23), "Client feedback session", true},
		{"Saka", time.Now().AddDate(0, 0, 24), "Project status update", false},
		{"Thomas", time.Now().AddDate(0, 0, 25), "Marketing campaign review", true},
		{"Etoo", time.Now().AddDate(0, 0, 26), "Business lunch", false},
		{"Fashola", time.Now().AddDate(0, 0, 27), "Product demo", true},
		{"Djspark", time.Now().AddDate(0, 0, 28), "Client meeting", false},
		{"Sam", time.Now().AddDate(0, 0, 29), "Team meeting", true},
		{"Jennifa", time.Now().AddDate(0, 0, 30), "Project discussion", false},
	}

	// Schedule appointments
	for _, appt := range appointmentsToSchedule {
		err = appointmentService.ScheduleAppointment(appt.ClientName, appt.DateTime, appt.Description, appt.Confirmed)
		if err != nil {
			log.Fatalf("Error scheduling appointment: %v", err)
		}
	}

	// Get appointments for the next week
	appointments, err := appointmentService.GetAppointments(time.Now(), time.Now().AddDate(0, 0, 7))
	if err != nil {
		log.Fatalf("Error getting appointments: %v", err)
	}

	// Display appointments
	fmt.Println("Scheduled Appointments:")
	for _, appt := range appointments {
		fmt.Printf("ID: %d | Client: %s | Date & Time: %s | Description: %s | Confirmed: %t\n",
			appt.ID, appt.ClientName, appt.DateTime.Format(time.RFC3339), appt.Description, appt.Confirmed)
	}
}

// ScheduleAppointment schedules a new appointment
func (s *AppointmentService) ScheduleAppointment(clientName string, dateTime time.Time, description string, confirmed bool) error {
	appointment := &Appointment{
		ClientName:  clientName,
		DateTime:    dateTime,
		Description: description,
		Confirmed:   confirmed,
	}
	err := s.DB.Create(appointment).Error
	if err != nil {
		return fmt.Errorf("failed to schedule appointment: %v", err)
	}
	return nil
}

// GetAppointments retrieves appointments between startDateTime and endDateTime
func (s *AppointmentService) GetAppointments(startDateTime, endDateTime time.Time) ([]Appointment, error) {
	var appointments []Appointment
	err := s.DB.Where("date_time BETWEEN ? AND ?", startDateTime, endDateTime).Find(&appointments).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get appointments: %v", err)
	}
	return appointments, nil
}

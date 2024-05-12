package main

import (
	"fmt"
	"log"
	"time"

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
	// Connect to MySQL database
	dsn := "joshua468:Temitope2080@tcp(127.0.0.1:3306)/appointments?charset=utf8mb4&parseTime=True&loc=Local"
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

	// Schedule an appointment
	err = appointmentService.ScheduleAppointment("John Doe", time.Now().AddDate(0, 0, 7), "Meeting with client")
	if err != nil {
		log.Fatalf("Error scheduling appointment: %v", err)
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
func (s *AppointmentService) ScheduleAppointment(clientName string, dateTime time.Time, description string) error {
	appointment := &Appointment{
		ClientName:  clientName,
		DateTime:    dateTime,
		Description: description,
		Confirmed:   false,
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

package entity

import (
	"github.com/google/uuid"
	"time"
)

// JobSeeker represents the job-seekers table
type JobSeeker struct {
	ID                uuid.UUID `db:"id" json:"id"`
	UserID            uuid.UUID `db:"user_id" json:"user_id"`
	FirstName         string    `db:"first_name" json:"first_name"`
	LastName          string    `db:"last_name" json:"last_name"`
	DateOfBirth       time.Time `db:"date_of_birth" json:"date_of_birth"`
	Bio               string    `db:"bio" json:"bio"`
	ProfilePictureURL *string   `db:"profile_picture_url" json:"profile_picture_url"`
	CreatedAt         time.Time `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time `db:"updated_at" json:"updated_at"`
}

// JobSeekerSkill represents the job_seeker_skills table
type JobSeekerSkill struct {
	ID               uuid.UUID `db:"id" json:"id"`
	JobSeekerID      uuid.UUID `db:"job_seeker_id" json:"job_seeker_id"`
	SkillName        string    `db:"skill_name" json:"skill_name"`
	ProficiencyLevel string    `db:"proficiency_level" json:"proficiency_level"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
	UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`
}

// JobSeekerProject represents the job_seeker_projects table
type JobSeekerProject struct {
	ID          uuid.UUID  `db:"id" json:"id"`
	JobSeekerID uuid.UUID  `db:"job_seeker_id" json:"job_seeker_id"`
	Title       string     `db:"title" json:"title"`
	Description string     `db:"description" json:"description"`
	StartDate   time.Time  `db:"start_date" json:"start_date"`
	EndDate     *time.Time `db:"end_date" json:"end_date"` // Use a pointer to allow null
	Link        string     `db:"link" json:"link"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
}

// JobSeekerPreference represents the job_seeker_preferences table
type JobSeekerPreference struct {
	ID                uuid.UUID `db:"id" json:"id"`
	JobSeekerID       uuid.UUID `db:"job_seeker_id" json:"job_seeker_id"`
	PreferredJobType  string    `db:"preferred_job_type" json:"preferred_job_type"`
	PreferredIndustry string    `db:"preferred_industry" json:"preferred_industry"`
	PreferredLocation string    `db:"preferred_location" json:"preferred_location"`
	SalaryExpectation *int      `db:"salary_expectation" json:"salary_expectation"` // Use a pointer to allow null
	RemotePreference  bool      `db:"remote_preference" json:"remote_preference"`
	CreatedAt         time.Time `db:"created_at" json:"created_at"`
	UpdatedAt         time.Time `db:"updated_at" json:"updated_at"`
}

// JobSeekerLink represents the job_seeker_links table
type JobSeekerLink struct {
	ID          uuid.UUID `db:"id" json:"id"`
	JobSeekerID uuid.UUID `db:"job_seeker_id" json:"job_seeker_id"`
	LinkType    string    `db:"link_type" json:"link_type"`
	URL         string    `db:"url" json:"url"`
	CreatedAt   time.Time `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time `db:"updated_at" json:"updated_at"`
}

// JobSeekerExperience represents the job_seeker_experiences table
type JobSeekerExperience struct {
	ID          uuid.UUID  `db:"id" json:"id"`
	JobSeekerID uuid.UUID  `db:"job_seeker_id" json:"job_seeker_id"`
	JobTitle    string     `db:"job_title" json:"job_title"`
	CompanyName string     `db:"company_name" json:"company_name"`
	StartDate   time.Time  `db:"start_date" json:"start_date"`
	EndDate     *time.Time `db:"end_date" json:"end_date"` // Use a pointer to allow null
	Description string     `db:"description" json:"description"`
	CreatedAt   time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updated_at"`
}

// JobSeekerEducation represents the job_seeker_educations table
type JobSeekerEducation struct {
	ID              uuid.UUID  `db:"id" json:"id"`
	JobSeekerID     uuid.UUID  `db:"job_seeker_id" json:"job_seeker_id"`
	InstitutionName string     `db:"institution_name" json:"institution_name"`
	Degree          string     `db:"degree" json:"degree"`
	FieldOfStudy    string     `db:"field_of_study" json:"field_of_study"`
	StartDate       time.Time  `db:"start_date" json:"start_date"`
	EndDate         *time.Time `db:"end_date" json:"end_date"` // Use a pointer to allow null
	Grade           string     `db:"grade" json:"grade"`
	Description     string     `db:"description" json:"description"`
	CreatedAt       time.Time  `db:"created_at" json:"created_at"`
	UpdatedAt       time.Time  `db:"updated_at" json:"updated_at"`
}

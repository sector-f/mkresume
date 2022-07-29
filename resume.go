package main

type Resume struct {
	Header     Header
	Experience []Company
	Skills     []Skill
	Projects   []Project
	Education  []Edu
}

func (r *Resume) redact() {
	r.Header.LastName = "[Redacted]"
	r.Header.Location = "[Location]"
	r.Header.Email.Display = "[email address]"
	r.Header.Email.Link = "mailto:me@example.com"

	for i := range r.Experience {
		r.Experience[i].Name = "[Company Name]"
		r.Experience[i].Location = "[Company Location]"
	}

	for i := range r.Education {
		r.Education[i].Institution = "[Institution]"
	}
}

type Header struct {
	FirstName string `json:"first_name" toml:"first_name"`
	LastName  string `json:"last_name" toml:"last_name"`
	JobTitle  string `json:"job_title" toml:"job_title"`
	Location  string
	Email     Link
	Github    Link
	Website   Link
}

type Link struct {
	Display string
	Link    string
}

type Company struct {
	Name     string `json:"company" toml:"company"`
	JobTitle string `json:"job_title" toml:"job_title"`
	Location string
	Time     string
	Bullets  []string
}

type Skill struct {
	Heading     string
	Description string
}

type Project struct {
	Name        string
	Description string
	Url         string
	Hide        bool
}

type Edu struct {
	Name        string
	Institution string
	Date        string
}

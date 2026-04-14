package main

type Job struct {
	Role        string
	Company     string
	Period      string
	Description string
}

type Project struct {
	Title       string
	Description string
	Tags        []string
	Year        string
	Status      string
}

type ContactLink struct {
	Label string
	Value string
	URL   string
}

var bio = `I'm Arkadiusz, a tech lead and product developer based in Poland, working remotely with teams and clients across the US and Europe.

I've spent the last 10 years building software — and more importantly, building the bridges between the people who make it and the people who need it.

I started in 2016 writing backend PHP, moved into fullstack, and spent the last five years in frontend-focused product development. Each step brought me closer to where I thrive — leading projects at the intersection of engineering, product, and communication.

I translate founder vision into shipped products. My day-to-day is making sure what we build actually solves the right problem — sitting with stakeholders, aligning teams, and turning business goals into clear technical direction.

Currently going deeper into engineering leadership, building my own SaaS product, and exploring how AI is reshaping the way we develop software.`

var experience = []Job{
	{
		Role:    "Tech Lead",
		Company: "Tonik",
		Period:  "2024 — Present",
		Description: `Working at the intersection of engineering and leadership for a studio that builds products with VC-backed US startups. My work spans the full spectrum — from early-stage MVPs in Next.js and POCs in React Native to team augmentation on established web platforms, across fintech, entertainment, and AI-driven products.

I still write production code daily, but increasingly my role is about raising the bar for the whole team: mentoring junior and mid-level developers, owning client communication, and establishing standards around code quality and delivery process.`,
	},
	{
		Role:    "Senior Software Developer",
		Company: "iTeamly",
		Period:  "2022 — 2025",
		Description: `A software house role that gave me rare breadth — I was placed across three distinct client engagements, each one pushing me further. I worked on a document management platform for a Swiss company, then on core feature development for a Finnish event management product.

My most impactful chapter was with a US-based translation management company, where I became the lead developer and primary contributor to their internal design system. The defining project: a sweeping, company-wide design system overhaul.`,
	},
	{
		Role:    "Frontend Developer",
		Company: "Sundose",
		Period:  "2020 — 2022",
		Description: `My first startup and my first time owning an entire frontend layer end-to-end. I was responsible for building and evolving the customer-facing platform for a personalized nutrition product.

By the end of my time at Sundose, I was the sole person responsible for all frontend development, having grown into full ownership of that domain.`,
	},
	{
		Role:        "Full Stack Developer",
		Company:     "LIBRUS",
		Period:      "2018 — 2020",
		Description: `Joined as a PHP developer working on LIBRUS's information and marketing platform — a gateway product serving millions of students, teachers, and parents across Poland. Midway through my tenure, I spearheaded the company's transition to modern frontend development by introducing React.js and Material UI.`,
	},
	{
		Role:        "Full Stack Developer",
		Company:     "Life in Mobile",
		Period:      "2017 — 2018",
		Description: `My introduction to agency life and working with US-based clients. I led the adoption of Laravel and Vue.js and co-created a project boilerplate that became the new standard for every client engagement that followed.`,
	},
	{
		Role:        "Junior Backend Developer",
		Company:     "ENCJA.COM",
		Period:      "2016",
		Description: `Where it all began. My first commercial role gave me the foundation of professional software development — building an internal product on a custom PHP framework and learning what it truly means to write code that runs in production.`,
	},
}

var stackFrontend = []string{"React", "Next.js", "React Native", "TypeScript", "Tailwind CSS", "Three.js", "Storybook"}
var stackBackend = []string{"Node.js", "PostgreSQL", "PHP", "Redis", "Supabase"}
var stackTools = []string{"Vercel", "Git", "Agile/Scrum", "AI-assisted dev", "Design Systems"}

var projects = []Project{
	{
		Title:       "Terminal Portfolio",
		Description: "You're looking at it right now. An SSH-accessible TUI version of my portfolio, built in Go with the Charm ecosystem.",
		Tags:        []string{"Go", "Bubble Tea", "Wish", "SSH"},
		Year:        "2026",
		Status:      "Live",
	},
	{
		Title:       "SaaS for Photographers",
		Description: "A SaaS product for the photography industry. Currently in early development — building in public.",
		Tags:        []string{"Next.js", "TypeScript", "AI-assisted development"},
		Year:        "2025",
		Status:      "Coming Soon",
	},
}

var contacts = []ContactLink{
	{Label: "Email", Value: "hello@arkadiuszjuszczyk.com", URL: "mailto:hello@arkadiuszjuszczyk.com"},
	{Label: "GitHub", Value: "github.com/4rek", URL: "https://github.com/4rek"},
	{Label: "LinkedIn", Value: "linkedin.com/in/arkadiuszjuszczyk", URL: "https://www.linkedin.com/in/arkadiuszjuszczyk/"},
	{Label: "Twitter", Value: "@4juszczyk", URL: "https://x.com/4juszczyk"},
}

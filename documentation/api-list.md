API Specification
│
├── 1. Overview
├── 2. API Conventions
├── 3. Authentication
│
├── 4. Authentication Endpoints
│   ├── POST /auth/register
│   ├── POST /auth/login
│   └── POST /auth/logout
│
├── 5. User Profile
│   ├── GET /me
│   └── PATCH /me
│
├── 6. Anime Discovery
│   ├── GET /animes/search
│   ├── GET /animes/search/advanced
│   └── GET /animes/:external_id
│
├── 7. Anime Identification
│   └── POST /animes/identify
│
├── 8. Personal Library
│   ├── POST /me/library
│   ├── GET /me/library
│   ├── GET /me/library/:id
│   ├── PATCH /me/library/:id
│   └── DELETE /me/library/:id
│
├── 9. Groups
│   ├── POST /me/groups
│   ├── GET /me/groups
│   ├── GET /me/groups/:id
│   ├── PATCH /me/groups/:id
│   ├── DELETE /me/groups/:id
│   ├── POST /me/groups/:id/animes
│   └── DELETE /me/groups/:id/animes/:library_anime_id
│
├── 10. Error Handling
│
└── 11. External API Integration
    ├── Third-Party Anime API
    └── trace.moe
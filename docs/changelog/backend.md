# Changelog

- 04.05.2026
  - Updated the `Term` model
  - Added `StartDate` (JSON `start_date`) and `EndDate` (JSON `end_date`) fields to the `Term` model.
  - Requires UI changes:
    - Status: `Open`
- 04.05.2026
  - Add `school_id` to the `Class` model
  - Changes:
    - Added `SchoolId` (JSON `school_id`) field to the `Class` model to associate classes with specific schools.
    - Updated the data migration that creates initial classes to include the `school_id` field.
    - Added a new route to fetch all classes for a specific school - `GET /api/school/{school_id}/classes`
  - Requires UI changes:
    - Status: `Open`
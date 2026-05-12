-- +goose Up
-- +goose StatementBegin
-- Mock Data Migration for Testing
-- Password for all users: 12345678
-- Hash: $2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy

-- Insert Admin User
INSERT INTO users (first_name, last_name, email, password) VALUES ('Admin', 'User', 'admin@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');

-- Assign Admin Role
INSERT INTO user_roles (user_id, role_id) SELECT u.user_id, r.role_id FROM users u, roles r WHERE u.email = 'admin@app.com' AND r.role_name = 'ADMIN';

-- Insert Directors
INSERT INTO users (first_name, last_name, email, password) VALUES ('Director', '1', 'director1@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Director', '2', 'director2@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Director', '3', 'director3@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');

-- Assign Director Roles
INSERT INTO user_roles (user_id, role_id) SELECT u.user_id, r.role_id FROM users u, roles r WHERE u.email IN ('director1@app.com', 'director2@app.com', 'director3@app.com') AND r.role_name = 'DIRECTOR';

-- Assign second role to Director 1
INSERT INTO user_roles (user_id, role_id) SELECT u.user_id, r.role_id FROM users u, roles r WHERE u.email IN ('director1@app.com') AND r.role_name = 'PARENT';

-- Assign third role to Director 1
INSERT INTO user_roles (user_id, role_id) SELECT u.user_id, r.role_id FROM users u, roles r WHERE u.email IN ('director1@app.com') AND r.role_name = 'ADMIN';

-- Create Director Entries
INSERT INTO directors (user_id, school_id) SELECT u.user_id, 1 FROM users u WHERE u.email = 'director1@app.com';
INSERT INTO directors (user_id, school_id) SELECT u.user_id, 2 FROM users u WHERE u.email = 'director2@app.com';
INSERT INTO directors (user_id, school_id) SELECT u.user_id, 3 FROM users u WHERE u.email = 'director3@app.com';

-- Insert Teachers
INSERT INTO users (first_name, last_name, email, password) VALUES ('Teacher', '1', 'teacher1@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Teacher', '2', 'teacher2@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Teacher', '3', 'teacher3@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Teacher', '4', 'teacher4@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Teacher', '5', 'teacher5@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Teacher', '6', 'teacher6@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Teacher', '7', 'teacher7@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Teacher', '8', 'teacher8@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Teacher', '9', 'teacher9@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Teacher', '10', 'teacher10@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Teacher', '11', 'teacher11@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Teacher', '12', 'teacher12@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Teacher', '13', 'teacher13@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Teacher', '14', 'teacher14@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Teacher', '15', 'teacher15@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');

-- Assign Teacher Roles
INSERT INTO user_roles (user_id, role_id) SELECT u.user_id, r.role_id FROM users u, roles r WHERE u.email LIKE 'teacher%@app.com' AND r.role_name = 'TEACHER';

-- Assign second role to Teacher 5
INSERT INTO user_roles (user_id, role_id) SELECT u.user_id, r.role_id FROM users u, roles r WHERE u.email IN ('teacher5@app.com') AND r.role_name = 'PARENT';

-- Create Teacher Entries (5 teachers per school)
INSERT INTO teachers (user_id, school_id) SELECT u.user_id, 1 FROM users u WHERE u.email = 'teacher1@app.com';
INSERT INTO teachers (user_id, school_id) SELECT u.user_id, 1 FROM users u WHERE u.email = 'teacher2@app.com';
INSERT INTO teachers (user_id, school_id) SELECT u.user_id, 1 FROM users u WHERE u.email = 'teacher3@app.com';
INSERT INTO teachers (user_id, school_id) SELECT u.user_id, 1 FROM users u WHERE u.email = 'teacher4@app.com';
INSERT INTO teachers (user_id, school_id) SELECT u.user_id, 1 FROM users u WHERE u.email = 'teacher5@app.com';
INSERT INTO teachers (user_id, school_id) SELECT u.user_id, 2 FROM users u WHERE u.email = 'teacher6@app.com';
INSERT INTO teachers (user_id, school_id) SELECT u.user_id, 2 FROM users u WHERE u.email = 'teacher7@app.com';
INSERT INTO teachers (user_id, school_id) SELECT u.user_id, 2 FROM users u WHERE u.email = 'teacher8@app.com';
INSERT INTO teachers (user_id, school_id) SELECT u.user_id, 2 FROM users u WHERE u.email = 'teacher9@app.com';
INSERT INTO teachers (user_id, school_id) SELECT u.user_id, 2 FROM users u WHERE u.email = 'teacher10@app.com';
INSERT INTO teachers (user_id, school_id) SELECT u.user_id, 3 FROM users u WHERE u.email = 'teacher11@app.com';
INSERT INTO teachers (user_id, school_id) SELECT u.user_id, 3 FROM users u WHERE u.email = 'teacher12@app.com';
INSERT INTO teachers (user_id, school_id) SELECT u.user_id, 3 FROM users u WHERE u.email = 'teacher13@app.com';
INSERT INTO teachers (user_id, school_id) SELECT u.user_id, 3 FROM users u WHERE u.email = 'teacher14@app.com';
INSERT INTO teachers (user_id, school_id) SELECT u.user_id, 3 FROM users u WHERE u.email = 'teacher15@app.com';

-- Assign Subjects to Teachers (each teacher teaches 2-3 subjects)
INSERT INTO teacher_subjects (teacher_id, subject_id) SELECT t.teacher_id, s.subject_id FROM teachers t, subjects s WHERE t.teacher_id = (SELECT teacher_id FROM teachers ORDER BY teacher_id LIMIT 1 OFFSET 0) AND s.subject_name = 'Bulgarian';
INSERT INTO teacher_subjects (teacher_id, subject_id) SELECT t.teacher_id, s.subject_id FROM teachers t, subjects s WHERE t.teacher_id = (SELECT teacher_id FROM teachers ORDER BY teacher_id LIMIT 1 OFFSET 0) AND s.subject_name = 'Literature';
INSERT INTO teacher_subjects (teacher_id, subject_id) SELECT t.teacher_id, s.subject_id FROM teachers t, subjects s WHERE t.teacher_id = (SELECT teacher_id FROM teachers ORDER BY teacher_id LIMIT 1 OFFSET 1) AND s.subject_name = 'Maths';
INSERT INTO teacher_subjects (teacher_id, subject_id) SELECT t.teacher_id, s.subject_id FROM teachers t, subjects s WHERE t.teacher_id = (SELECT teacher_id FROM teachers ORDER BY teacher_id LIMIT 1 OFFSET 1) AND s.subject_name = 'Physics';
INSERT INTO teacher_subjects (teacher_id, subject_id) SELECT t.teacher_id, s.subject_id FROM teachers t, subjects s WHERE t.teacher_id = (SELECT teacher_id FROM teachers ORDER BY teacher_id LIMIT 1 OFFSET 2) AND s.subject_name = 'English';
INSERT INTO teacher_subjects (teacher_id, subject_id) SELECT t.teacher_id, s.subject_id FROM teachers t, subjects s WHERE t.teacher_id = (SELECT teacher_id FROM teachers ORDER BY teacher_id LIMIT 1 OFFSET 2) AND s.subject_name = 'Geography';
INSERT INTO teacher_subjects (teacher_id, subject_id) SELECT t.teacher_id, s.subject_id FROM teachers t, subjects s WHERE t.teacher_id = (SELECT teacher_id FROM teachers ORDER BY teacher_id LIMIT 1 OFFSET 3) AND s.subject_name = 'Chemistry';
INSERT INTO teacher_subjects (teacher_id, subject_id) SELECT t.teacher_id, s.subject_id FROM teachers t, subjects s WHERE t.teacher_id = (SELECT teacher_id FROM teachers ORDER BY teacher_id LIMIT 1 OFFSET 3) AND s.subject_name = 'Biology';
INSERT INTO teacher_subjects (teacher_id, subject_id) SELECT t.teacher_id, s.subject_id FROM teachers t, subjects s WHERE t.teacher_id = (SELECT teacher_id FROM teachers ORDER BY teacher_id LIMIT 1 OFFSET 4) AND s.subject_name = 'History';
INSERT INTO teacher_subjects (teacher_id, subject_id) SELECT t.teacher_id, s.subject_id FROM teachers t, subjects s WHERE t.teacher_id = (SELECT teacher_id FROM teachers ORDER BY teacher_id LIMIT 1 OFFSET 4) AND s.subject_name = 'Physical Education';
INSERT INTO teacher_subjects (teacher_id, subject_id) SELECT t.teacher_id, s.subject_id FROM teachers t, subjects s WHERE t.teacher_id = (SELECT teacher_id FROM teachers ORDER BY teacher_id LIMIT 1 OFFSET 5) AND s.subject_name = 'Bulgarian';
INSERT INTO teacher_subjects (teacher_id, subject_id) SELECT t.teacher_id, s.subject_id FROM teachers t, subjects s WHERE t.teacher_id = (SELECT teacher_id FROM teachers ORDER BY teacher_id LIMIT 1 OFFSET 5) AND s.subject_name = 'Music';
INSERT INTO teacher_subjects (teacher_id, subject_id) SELECT t.teacher_id, s.subject_id FROM teachers t, subjects s WHERE t.teacher_id = (SELECT teacher_id FROM teachers ORDER BY teacher_id LIMIT 1 OFFSET 6) AND s.subject_name = 'Maths';
INSERT INTO teacher_subjects (teacher_id, subject_id) SELECT t.teacher_id, s.subject_id FROM teachers t, subjects s WHERE t.teacher_id = (SELECT teacher_id FROM teachers ORDER BY teacher_id LIMIT 1 OFFSET 6) AND s.subject_name = 'Information Technology';
INSERT INTO teacher_subjects (teacher_id, subject_id) SELECT t.teacher_id, s.subject_id FROM teachers t, subjects s WHERE t.teacher_id = (SELECT teacher_id FROM teachers ORDER BY teacher_id LIMIT 1 OFFSET 7) AND s.subject_name = 'English';
INSERT INTO teacher_subjects (teacher_id, subject_id) SELECT t.teacher_id, s.subject_id FROM teachers t, subjects s WHERE t.teacher_id = (SELECT teacher_id FROM teachers ORDER BY teacher_id LIMIT 1 OFFSET 7) AND s.subject_name = 'Literature';
INSERT INTO teacher_subjects (teacher_id, subject_id) SELECT t.teacher_id, s.subject_id FROM teachers t, subjects s WHERE t.teacher_id = (SELECT teacher_id FROM teachers ORDER BY teacher_id LIMIT 1 OFFSET 8) AND s.subject_name = 'Chemistry';
INSERT INTO teacher_subjects (teacher_id, subject_id) SELECT t.teacher_id, s.subject_id FROM teachers t, subjects s WHERE t.teacher_id = (SELECT teacher_id FROM teachers ORDER BY teacher_id LIMIT 1 OFFSET 8) AND s.subject_name = 'Physics';
INSERT INTO teacher_subjects (teacher_id, subject_id) SELECT t.teacher_id, s.subject_id FROM teachers t, subjects s WHERE t.teacher_id = (SELECT teacher_id FROM teachers ORDER BY teacher_id LIMIT 1 OFFSET 9) AND s.subject_name = 'History';
INSERT INTO teacher_subjects (teacher_id, subject_id) SELECT t.teacher_id, s.subject_id FROM teachers t, subjects s WHERE t.teacher_id = (SELECT teacher_id FROM teachers ORDER BY teacher_id LIMIT 1 OFFSET 9) AND s.subject_name = 'Geography';
INSERT INTO teacher_subjects (teacher_id, subject_id) SELECT t.teacher_id, s.subject_id FROM teachers t, subjects s WHERE t.teacher_id = (SELECT teacher_id FROM teachers ORDER BY teacher_id LIMIT 1 OFFSET 10) AND s.subject_name = 'Bulgarian';
INSERT INTO teacher_subjects (teacher_id, subject_id) SELECT t.teacher_id, s.subject_id FROM teachers t, subjects s WHERE t.teacher_id = (SELECT teacher_id FROM teachers ORDER BY teacher_id LIMIT 1 OFFSET 10) AND s.subject_name = 'English';
INSERT INTO teacher_subjects (teacher_id, subject_id) SELECT t.teacher_id, s.subject_id FROM teachers t, subjects s WHERE t.teacher_id = (SELECT teacher_id FROM teachers ORDER BY teacher_id LIMIT 1 OFFSET 11) AND s.subject_name = 'Maths';
INSERT INTO teacher_subjects (teacher_id, subject_id) SELECT t.teacher_id, s.subject_id FROM teachers t, subjects s WHERE t.teacher_id = (SELECT teacher_id FROM teachers ORDER BY teacher_id LIMIT 1 OFFSET 11) AND s.subject_name = 'Biology';
INSERT INTO teacher_subjects (teacher_id, subject_id) SELECT t.teacher_id, s.subject_id FROM teachers t, subjects s WHERE t.teacher_id = (SELECT teacher_id FROM teachers ORDER BY teacher_id LIMIT 1 OFFSET 12) AND s.subject_name = 'Physics';
INSERT INTO teacher_subjects (teacher_id, subject_id) SELECT t.teacher_id, s.subject_id FROM teachers t, subjects s WHERE t.teacher_id = (SELECT teacher_id FROM teachers ORDER BY teacher_id LIMIT 1 OFFSET 12) AND s.subject_name = 'Chemistry';
INSERT INTO teacher_subjects (teacher_id, subject_id) SELECT t.teacher_id, s.subject_id FROM teachers t, subjects s WHERE t.teacher_id = (SELECT teacher_id FROM teachers ORDER BY teacher_id LIMIT 1 OFFSET 13) AND s.subject_name = 'Geography';
INSERT INTO teacher_subjects (teacher_id, subject_id) SELECT t.teacher_id, s.subject_id FROM teachers t, subjects s WHERE t.teacher_id = (SELECT teacher_id FROM teachers ORDER BY teacher_id LIMIT 1 OFFSET 13) AND s.subject_name = 'Physical Education';
INSERT INTO teacher_subjects (teacher_id, subject_id) SELECT t.teacher_id, s.subject_id FROM teachers t, subjects s WHERE t.teacher_id = (SELECT teacher_id FROM teachers ORDER BY teacher_id LIMIT 1 OFFSET 14) AND s.subject_name = 'History';
INSERT INTO teacher_subjects (teacher_id, subject_id) SELECT t.teacher_id, s.subject_id FROM teachers t, subjects s WHERE t.teacher_id = (SELECT teacher_id FROM teachers ORDER BY teacher_id LIMIT 1 OFFSET 14) AND s.subject_name = 'Information Technology';

-- Insert Parents
INSERT INTO users (first_name, last_name, email, password) VALUES ('Parent', '1', 'parent1@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Parent', '2', 'parent2@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Parent', '3', 'parent3@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Parent', '4', 'parent4@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Parent', '5', 'parent5@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Parent', '6', 'parent6@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Parent', '7', 'parent7@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Parent', '8', 'parent8@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Parent', '9', 'parent9@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Parent', '10', 'parent10@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Parent', '11', 'parent11@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Parent', '12', 'parent12@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Parent', '13', 'parent13@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Parent', '14', 'parent14@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Parent', '15', 'parent15@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Parent', '16', 'parent16@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Parent', '17', 'parent17@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Parent', '18', 'parent18@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Parent', '19', 'parent19@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Parent', '20', 'parent20@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');

-- Assign Parent Roles
INSERT INTO user_roles (user_id, role_id) SELECT u.user_id, r.role_id FROM users u, roles r WHERE u.email LIKE 'parent%@app.com' AND r.role_name = 'PARENT';

-- Create Parent Entries
INSERT INTO parents (user_id) SELECT user_id FROM users WHERE email LIKE 'parent%@app.com';
-- Create Parent Entry for Director 1 (who also has PARENT role)
INSERT INTO parents (user_id) SELECT user_id FROM users WHERE email = 'director1@app.com';
-- Create Parent Entry for Teacher 5 (who also has PARENT role)
INSERT INTO parents (user_id) SELECT user_id FROM users WHERE email = 'teacher5@app.com';

-- Insert Students
INSERT INTO users (first_name, last_name, email, password) VALUES ('Student', '1', 'student1@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Student', '2', 'student2@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Student', '3', 'student3@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Student', '4', 'student4@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Student', '5', 'student5@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Student', '6', 'student6@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Student', '7', 'student7@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Student', '8', 'student8@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Student', '9', 'student9@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Student', '10', 'student10@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Student', '11', 'student11@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Student', '12', 'student12@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Student', '13', 'student13@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Student', '14', 'student14@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Student', '15', 'student15@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Student', '16', 'student16@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Student', '17', 'student17@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Student', '18', 'student18@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Student', '19', 'student19@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Student', '20', 'student20@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Student', '21', 'student21@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Student', '22', 'student22@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Student', '23', 'student23@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Student', '24', 'student24@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Student', '25', 'student25@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Student', '26', 'student26@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Student', '27', 'student27@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Student', '28', 'student28@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Student', '29', 'student29@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Student', '30', 'student30@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Student', '31', 'student31@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Student', '32', 'student32@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Student', '33', 'student33@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');
INSERT INTO users (first_name, last_name, email, password) VALUES ('Student', '34', 'student34@app.com', '$2a$06$m.UbLZoVTw7H.tDQcOtf3.eamAhlKaFfhjzMgddFlbWiDy6bx8vBy');

-- Assign Student Roles
INSERT INTO user_roles (user_id, role_id) SELECT u.user_id, r.role_id FROM users u, roles r WHERE u.email LIKE 'student%@app.com' AND r.role_name = 'STUDENT';

-- Create Student Entries and Assign to Classes (School 1 - GPAE)
INSERT INTO students (user_id, school_id, class_id) SELECT u.user_id, 1, c.class_id FROM users u, classes c WHERE u.email = 'student1@app.com' AND c.school_id = 1 AND c.grade_level = 1 AND c.class_name = 'A';
INSERT INTO students (user_id, school_id, class_id) SELECT u.user_id, 1, c.class_id FROM users u, classes c WHERE u.email = 'student2@app.com' AND c.school_id = 1 AND c.grade_level = 1 AND c.class_name = 'A';
INSERT INTO students (user_id, school_id, class_id) SELECT u.user_id, 1, c.class_id FROM users u, classes c WHERE u.email = 'student3@app.com' AND c.school_id = 1 AND c.grade_level = 1 AND c.class_name = 'B';
INSERT INTO students (user_id, school_id, class_id) SELECT u.user_id, 1, c.class_id FROM users u, classes c WHERE u.email = 'student4@app.com' AND c.school_id = 1 AND c.grade_level = 2 AND c.class_name = 'A';
INSERT INTO students (user_id, school_id, class_id) SELECT u.user_id, 1, c.class_id FROM users u, classes c WHERE u.email = 'student5@app.com' AND c.school_id = 1 AND c.grade_level = 2 AND c.class_name = 'B';
INSERT INTO students (user_id, school_id, class_id) SELECT u.user_id, 1, c.class_id FROM users u, classes c WHERE u.email = 'student6@app.com' AND c.school_id = 1 AND c.grade_level = 3 AND c.class_name = 'A';
INSERT INTO students (user_id, school_id, class_id) SELECT u.user_id, 1, c.class_id FROM users u, classes c WHERE u.email = 'student7@app.com' AND c.school_id = 1 AND c.grade_level = 3 AND c.class_name = 'B';
INSERT INTO students (user_id, school_id, class_id) SELECT u.user_id, 1, c.class_id FROM users u, classes c WHERE u.email = 'student8@app.com' AND c.school_id = 1 AND c.grade_level = 4 AND c.class_name = 'A';
INSERT INTO students (user_id, school_id, class_id) SELECT u.user_id, 1, c.class_id FROM users u, classes c WHERE u.email = 'student9@app.com' AND c.school_id = 1 AND c.grade_level = 4 AND c.class_name = 'B';
INSERT INTO students (user_id, school_id, class_id) SELECT u.user_id, 1, c.class_id FROM users u, classes c WHERE u.email = 'student10@app.com' AND c.school_id = 1 AND c.grade_level = 4 AND c.class_name = 'C';
-- Students for Director 1 (who is also a parent)
INSERT INTO students (user_id, school_id, class_id) SELECT u.user_id, 1, c.class_id FROM users u, classes c WHERE u.email = 'student31@app.com' AND c.school_id = 1 AND c.grade_level = 1 AND c.class_name = 'A';
INSERT INTO students (user_id, school_id, class_id) SELECT u.user_id, 1, c.class_id FROM users u, classes c WHERE u.email = 'student32@app.com' AND c.school_id = 1 AND c.grade_level = 3 AND c.class_name = 'A';
-- Students for Teacher 5 (who is also a parent)
INSERT INTO students (user_id, school_id, class_id) SELECT u.user_id, 1, c.class_id FROM users u, classes c WHERE u.email = 'student33@app.com' AND c.school_id = 1 AND c.grade_level = 2 AND c.class_name = 'A';
INSERT INTO students (user_id, school_id, class_id) SELECT u.user_id, 1, c.class_id FROM users u, classes c WHERE u.email = 'student34@app.com' AND c.school_id = 1 AND c.grade_level = 4 AND c.class_name = 'A';

-- Create Student Entries and Assign to Classes (School 2 - SMG)
INSERT INTO students (user_id, school_id, class_id) SELECT u.user_id, 2, c.class_id FROM users u, classes c WHERE u.email = 'student11@app.com' AND c.school_id = 2 AND c.grade_level = 5 AND c.class_name = 'A';
INSERT INTO students (user_id, school_id, class_id) SELECT u.user_id, 2, c.class_id FROM users u, classes c WHERE u.email = 'student12@app.com' AND c.school_id = 2 AND c.grade_level = 5 AND c.class_name = 'B';
INSERT INTO students (user_id, school_id, class_id) SELECT u.user_id, 2, c.class_id FROM users u, classes c WHERE u.email = 'student13@app.com' AND c.school_id = 2 AND c.grade_level = 5 AND c.class_name = 'C';
INSERT INTO students (user_id, school_id, class_id) SELECT u.user_id, 2, c.class_id FROM users u, classes c WHERE u.email = 'student14@app.com' AND c.school_id = 2 AND c.grade_level = 6 AND c.class_name = 'A';
INSERT INTO students (user_id, school_id, class_id) SELECT u.user_id, 2, c.class_id FROM users u, classes c WHERE u.email = 'student15@app.com' AND c.school_id = 2 AND c.grade_level = 6 AND c.class_name = 'B';
INSERT INTO students (user_id, school_id, class_id) SELECT u.user_id, 2, c.class_id FROM users u, classes c WHERE u.email = 'student16@app.com' AND c.school_id = 2 AND c.grade_level = 6 AND c.class_name = 'C';
INSERT INTO students (user_id, school_id, class_id) SELECT u.user_id, 2, c.class_id FROM users u, classes c WHERE u.email = 'student17@app.com' AND c.school_id = 2 AND c.grade_level = 7 AND c.class_name = 'A';
INSERT INTO students (user_id, school_id, class_id) SELECT u.user_id, 2, c.class_id FROM users u, classes c WHERE u.email = 'student18@app.com' AND c.school_id = 2 AND c.grade_level = 7 AND c.class_name = 'B';
INSERT INTO students (user_id, school_id, class_id) SELECT u.user_id, 2, c.class_id FROM users u, classes c WHERE u.email = 'student19@app.com' AND c.school_id = 2 AND c.grade_level = 7 AND c.class_name = 'C';
INSERT INTO students (user_id, school_id, class_id) SELECT u.user_id, 2, c.class_id FROM users u, classes c WHERE u.email = 'student20@app.com' AND c.school_id = 2 AND c.grade_level = 7 AND c.class_name = 'D';

-- Create Student Entries and Assign to Classes (School 3 - GPNE)
INSERT INTO students (user_id, school_id, class_id) SELECT u.user_id, 3, c.class_id FROM users u, classes c WHERE u.email = 'student21@app.com' AND c.school_id = 3 AND c.grade_level = 8 AND c.class_name = 'A';
INSERT INTO students (user_id, school_id, class_id) SELECT u.user_id, 3, c.class_id FROM users u, classes c WHERE u.email = 'student22@app.com' AND c.school_id = 3 AND c.grade_level = 9 AND c.class_name = 'A';
INSERT INTO students (user_id, school_id, class_id) SELECT u.user_id, 3, c.class_id FROM users u, classes c WHERE u.email = 'student23@app.com' AND c.school_id = 3 AND c.grade_level = 9 AND c.class_name = 'B';
INSERT INTO students (user_id, school_id, class_id) SELECT u.user_id, 3, c.class_id FROM users u, classes c WHERE u.email = 'student24@app.com' AND c.school_id = 3 AND c.grade_level = 10 AND c.class_name = 'A';
INSERT INTO students (user_id, school_id, class_id) SELECT u.user_id, 3, c.class_id FROM users u, classes c WHERE u.email = 'student25@app.com' AND c.school_id = 3 AND c.grade_level = 10 AND c.class_name = 'B';
INSERT INTO students (user_id, school_id, class_id) SELECT u.user_id, 3, c.class_id FROM users u, classes c WHERE u.email = 'student26@app.com' AND c.school_id = 3 AND c.grade_level = 11 AND c.class_name = 'A';
INSERT INTO students (user_id, school_id, class_id) SELECT u.user_id, 3, c.class_id FROM users u, classes c WHERE u.email = 'student27@app.com' AND c.school_id = 3 AND c.grade_level = 11 AND c.class_name = 'B';
INSERT INTO students (user_id, school_id, class_id) SELECT u.user_id, 3, c.class_id FROM users u, classes c WHERE u.email = 'student28@app.com' AND c.school_id = 3 AND c.grade_level = 12 AND c.class_name = 'A';
INSERT INTO students (user_id, school_id, class_id) SELECT u.user_id, 3, c.class_id FROM users u, classes c WHERE u.email = 'student29@app.com' AND c.school_id = 3 AND c.grade_level = 12 AND c.class_name = 'B';
INSERT INTO students (user_id, school_id, class_id) SELECT u.user_id, 3, c.class_id FROM users u, classes c WHERE u.email = 'student30@app.com' AND c.school_id = 3 AND c.grade_level = 12 AND c.class_name = 'C';

-- Link Students to Parents
INSERT INTO student_parents (parent_id, student_id) SELECT p.parent_id, s.student_id FROM parents p, students s, users u1, users u2 WHERE p.user_id = u1.user_id AND s.user_id = u2.user_id AND u1.email = 'parent1@app.com' AND u2.email = 'student1@app.com';
INSERT INTO student_parents (parent_id, student_id) SELECT p.parent_id, s.student_id FROM parents p, students s, users u1, users u2 WHERE p.user_id = u1.user_id AND s.user_id = u2.user_id AND u1.email = 'parent1@app.com' AND u2.email = 'student2@app.com';
INSERT INTO student_parents (parent_id, student_id) SELECT p.parent_id, s.student_id FROM parents p, students s, users u1, users u2 WHERE p.user_id = u1.user_id AND s.user_id = u2.user_id AND u1.email = 'parent2@app.com' AND u2.email = 'student1@app.com';
INSERT INTO student_parents (parent_id, student_id) SELECT p.parent_id, s.student_id FROM parents p, students s, users u1, users u2 WHERE p.user_id = u1.user_id AND s.user_id = u2.user_id AND u1.email = 'parent2@app.com' AND u2.email = 'student2@app.com';
INSERT INTO student_parents (parent_id, student_id) SELECT p.parent_id, s.student_id FROM parents p, students s, users u1, users u2 WHERE p.user_id = u1.user_id AND s.user_id = u2.user_id AND u1.email = 'parent3@app.com' AND u2.email = 'student3@app.com';
INSERT INTO student_parents (parent_id, student_id) SELECT p.parent_id, s.student_id FROM parents p, students s, users u1, users u2 WHERE p.user_id = u1.user_id AND s.user_id = u2.user_id AND u1.email = 'parent4@app.com' AND u2.email = 'student4@app.com';
INSERT INTO student_parents (parent_id, student_id) SELECT p.parent_id, s.student_id FROM parents p, students s, users u1, users u2 WHERE p.user_id = u1.user_id AND s.user_id = u2.user_id AND u1.email = 'parent4@app.com' AND u2.email = 'student5@app.com';
INSERT INTO student_parents (parent_id, student_id) SELECT p.parent_id, s.student_id FROM parents p, students s, users u1, users u2 WHERE p.user_id = u1.user_id AND s.user_id = u2.user_id AND u1.email = 'parent5@app.com' AND u2.email = 'student4@app.com';
INSERT INTO student_parents (parent_id, student_id) SELECT p.parent_id, s.student_id FROM parents p, students s, users u1, users u2 WHERE p.user_id = u1.user_id AND s.user_id = u2.user_id AND u1.email = 'parent5@app.com' AND u2.email = 'student5@app.com';
INSERT INTO student_parents (parent_id, student_id) SELECT p.parent_id, s.student_id FROM parents p, students s, users u1, users u2 WHERE p.user_id = u1.user_id AND s.user_id = u2.user_id AND u1.email = 'parent6@app.com' AND u2.email = 'student6@app.com';
INSERT INTO student_parents (parent_id, student_id) SELECT p.parent_id, s.student_id FROM parents p, students s, users u1, users u2 WHERE p.user_id = u1.user_id AND s.user_id = u2.user_id AND u1.email = 'parent6@app.com' AND u2.email = 'student7@app.com';
INSERT INTO student_parents (parent_id, student_id) SELECT p.parent_id, s.student_id FROM parents p, students s, users u1, users u2 WHERE p.user_id = u1.user_id AND s.user_id = u2.user_id AND u1.email = 'parent7@app.com' AND u2.email = 'student8@app.com';
INSERT INTO student_parents (parent_id, student_id) SELECT p.parent_id, s.student_id FROM parents p, students s, users u1, users u2 WHERE p.user_id = u1.user_id AND s.user_id = u2.user_id AND u1.email = 'parent8@app.com' AND u2.email = 'student9@app.com';
INSERT INTO student_parents (parent_id, student_id) SELECT p.parent_id, s.student_id FROM parents p, students s, users u1, users u2 WHERE p.user_id = u1.user_id AND s.user_id = u2.user_id AND u1.email = 'parent9@app.com' AND u2.email = 'student10@app.com';
INSERT INTO student_parents (parent_id, student_id) SELECT p.parent_id, s.student_id FROM parents p, students s, users u1, users u2 WHERE p.user_id = u1.user_id AND s.user_id = u2.user_id AND u1.email = 'parent10@app.com' AND u2.email = 'student11@app.com';
INSERT INTO student_parents (parent_id, student_id) SELECT p.parent_id, s.student_id FROM parents p, students s, users u1, users u2 WHERE p.user_id = u1.user_id AND s.user_id = u2.user_id AND u1.email = 'parent10@app.com' AND u2.email = 'student12@app.com';
INSERT INTO student_parents (parent_id, student_id) SELECT p.parent_id, s.student_id FROM parents p, students s, users u1, users u2 WHERE p.user_id = u1.user_id AND s.user_id = u2.user_id AND u1.email = 'parent11@app.com' AND u2.email = 'student13@app.com';
INSERT INTO student_parents (parent_id, student_id) SELECT p.parent_id, s.student_id FROM parents p, students s, users u1, users u2 WHERE p.user_id = u1.user_id AND s.user_id = u2.user_id AND u1.email = 'parent12@app.com' AND u2.email = 'student14@app.com';
INSERT INTO student_parents (parent_id, student_id) SELECT p.parent_id, s.student_id FROM parents p, students s, users u1, users u2 WHERE p.user_id = u1.user_id AND s.user_id = u2.user_id AND u1.email = 'parent13@app.com' AND u2.email = 'student15@app.com';
INSERT INTO student_parents (parent_id, student_id) SELECT p.parent_id, s.student_id FROM parents p, students s, users u1, users u2 WHERE p.user_id = u1.user_id AND s.user_id = u2.user_id AND u1.email = 'parent13@app.com' AND u2.email = 'student16@app.com';
INSERT INTO student_parents (parent_id, student_id) SELECT p.parent_id, s.student_id FROM parents p, students s, users u1, users u2 WHERE p.user_id = u1.user_id AND s.user_id = u2.user_id AND u1.email = 'parent14@app.com' AND u2.email = 'student17@app.com';
INSERT INTO student_parents (parent_id, student_id) SELECT p.parent_id, s.student_id FROM parents p, students s, users u1, users u2 WHERE p.user_id = u1.user_id AND s.user_id = u2.user_id AND u1.email = 'parent15@app.com' AND u2.email = 'student18@app.com';
INSERT INTO student_parents (parent_id, student_id) SELECT p.parent_id, s.student_id FROM parents p, students s, users u1, users u2 WHERE p.user_id = u1.user_id AND s.user_id = u2.user_id AND u1.email = 'parent15@app.com' AND u2.email = 'student19@app.com';
INSERT INTO student_parents (parent_id, student_id) SELECT p.parent_id, s.student_id FROM parents p, students s, users u1, users u2 WHERE p.user_id = u1.user_id AND s.user_id = u2.user_id AND u1.email = 'parent16@app.com' AND u2.email = 'student20@app.com';
INSERT INTO student_parents (parent_id, student_id) SELECT p.parent_id, s.student_id FROM parents p, students s, users u1, users u2 WHERE p.user_id = u1.user_id AND s.user_id = u2.user_id AND u1.email = 'parent17@app.com' AND u2.email = 'student21@app.com';
INSERT INTO student_parents (parent_id, student_id) SELECT p.parent_id, s.student_id FROM parents p, students s, users u1, users u2 WHERE p.user_id = u1.user_id AND s.user_id = u2.user_id AND u1.email = 'parent17@app.com' AND u2.email = 'student22@app.com';
INSERT INTO student_parents (parent_id, student_id) SELECT p.parent_id, s.student_id FROM parents p, students s, users u1, users u2 WHERE p.user_id = u1.user_id AND s.user_id = u2.user_id AND u1.email = 'parent18@app.com' AND u2.email = 'student23@app.com';
INSERT INTO student_parents (parent_id, student_id) SELECT p.parent_id, s.student_id FROM parents p, students s, users u1, users u2 WHERE p.user_id = u1.user_id AND s.user_id = u2.user_id AND u1.email = 'parent18@app.com' AND u2.email = 'student24@app.com';
INSERT INTO student_parents (parent_id, student_id) SELECT p.parent_id, s.student_id FROM parents p, students s, users u1, users u2 WHERE p.user_id = u1.user_id AND s.user_id = u2.user_id AND u1.email = 'parent19@app.com' AND u2.email = 'student25@app.com';
INSERT INTO student_parents (parent_id, student_id) SELECT p.parent_id, s.student_id FROM parents p, students s, users u1, users u2 WHERE p.user_id = u1.user_id AND s.user_id = u2.user_id AND u1.email = 'parent19@app.com' AND u2.email = 'student26@app.com';
INSERT INTO student_parents (parent_id, student_id) SELECT p.parent_id, s.student_id FROM parents p, students s, users u1, users u2 WHERE p.user_id = u1.user_id AND s.user_id = u2.user_id AND u1.email = 'parent20@app.com' AND u2.email = 'student27@app.com';
INSERT INTO student_parents (parent_id, student_id) SELECT p.parent_id, s.student_id FROM parents p, students s, users u1, users u2 WHERE p.user_id = u1.user_id AND s.user_id = u2.user_id AND u1.email = 'parent20@app.com' AND u2.email = 'student28@app.com';
INSERT INTO student_parents (parent_id, student_id) SELECT p.parent_id, s.student_id FROM parents p, students s, users u1, users u2 WHERE p.user_id = u1.user_id AND s.user_id = u2.user_id AND u1.email = 'parent20@app.com' AND u2.email = 'student29@app.com';
-- Link students to Director 1 (who is also a parent)
INSERT INTO student_parents (parent_id, student_id) SELECT p.parent_id, s.student_id FROM parents p, students s, users u1, users u2 WHERE p.user_id = u1.user_id AND s.user_id = u2.user_id AND u1.email = 'director1@app.com' AND u2.email = 'student31@app.com';
INSERT INTO student_parents (parent_id, student_id) SELECT p.parent_id, s.student_id FROM parents p, students s, users u1, users u2 WHERE p.user_id = u1.user_id AND s.user_id = u2.user_id AND u1.email = 'director1@app.com' AND u2.email = 'student32@app.com';
-- Link students to Teacher 5 (who is also a parent)
INSERT INTO student_parents (parent_id, student_id) SELECT p.parent_id, s.student_id FROM parents p, students s, users u1, users u2 WHERE p.user_id = u1.user_id AND s.user_id = u2.user_id AND u1.email = 'teacher5@app.com' AND u2.email = 'student33@app.com';
INSERT INTO student_parents (parent_id, student_id) SELECT p.parent_id, s.student_id FROM parents p, students s, users u1, users u2 WHERE p.user_id = u1.user_id AND s.user_id = u2.user_id AND u1.email = 'teacher5@app.com' AND u2.email = 'student34@app.com';
INSERT INTO student_parents (parent_id, student_id) SELECT p.parent_id, s.student_id FROM parents p, students s, users u1, users u2 WHERE p.user_id = u1.user_id AND s.user_id = u2.user_id AND u1.email = 'parent20@app.com' AND u2.email = 'student30@app.com';

-- Create Curricula for Spring 2026 term (ongoing term) - School 1 (GPAE) Grade 1A
INSERT INTO curricula (class_id, subject_id, teacher_id, term_id) SELECT c.class_id, s.subject_id, t.teacher_id, tm.term_id FROM classes c, subjects s, teachers t, terms tm WHERE c.school_id = 1 AND c.grade_level = 1 AND c.class_name = 'A' AND s.subject_name = 'Bulgarian' AND t.school_id = 1 AND tm.name = 'Spring 2026' LIMIT 1;
INSERT INTO curricula (class_id, subject_id, teacher_id, term_id) SELECT c.class_id, s.subject_id, t.teacher_id, tm.term_id FROM classes c, subjects s, teachers t, terms tm WHERE c.school_id = 1 AND c.grade_level = 1 AND c.class_name = 'A' AND s.subject_name = 'Maths' AND t.school_id = 1 AND tm.name = 'Spring 2026' LIMIT 1;
INSERT INTO curricula (class_id, subject_id, teacher_id, term_id) SELECT c.class_id, s.subject_id, t.teacher_id, tm.term_id FROM classes c, subjects s, teachers t, terms tm WHERE c.school_id = 1 AND c.grade_level = 1 AND c.class_name = 'A' AND s.subject_name = 'English' AND t.school_id = 1 AND tm.name = 'Spring 2026' LIMIT 1;

-- Create Curricula - School 1 (GPAE) Grade 3A
INSERT INTO curricula (class_id, subject_id, teacher_id, term_id) SELECT c.class_id, s.subject_id, t.teacher_id, tm.term_id FROM classes c, subjects s, teachers t, terms tm WHERE c.school_id = 1 AND c.grade_level = 3 AND c.class_name = 'A' AND s.subject_name = 'Bulgarian' AND t.school_id = 1 AND tm.name = 'Spring 2026' LIMIT 1;
INSERT INTO curricula (class_id, subject_id, teacher_id, term_id) SELECT c.class_id, s.subject_id, t.teacher_id, tm.term_id FROM classes c, subjects s, teachers t, terms tm WHERE c.school_id = 1 AND c.grade_level = 3 AND c.class_name = 'A' AND s.subject_name = 'Maths' AND t.school_id = 1 AND tm.name = 'Spring 2026' LIMIT 1;
INSERT INTO curricula (class_id, subject_id, teacher_id, term_id) SELECT c.class_id, s.subject_id, t.teacher_id, tm.term_id FROM classes c, subjects s, teachers t, terms tm WHERE c.school_id = 1 AND c.grade_level = 3 AND c.class_name = 'A' AND s.subject_name = 'English' AND t.school_id = 1 AND tm.name = 'Spring 2026' LIMIT 1;
INSERT INTO curricula (class_id, subject_id, teacher_id, term_id) SELECT c.class_id, s.subject_id, t.teacher_id, tm.term_id FROM classes c, subjects s, teachers t, terms tm WHERE c.school_id = 1 AND c.grade_level = 3 AND c.class_name = 'A' AND s.subject_name = 'History' AND t.school_id = 1 AND tm.name = 'Spring 2026' LIMIT 1;

-- Create Curricula - School 1 (GPAE) Grade 2A
INSERT INTO curricula (class_id, subject_id, teacher_id, term_id) SELECT c.class_id, s.subject_id, t.teacher_id, tm.term_id FROM classes c, subjects s, teachers t, terms tm WHERE c.school_id = 1 AND c.grade_level = 2 AND c.class_name = 'A' AND s.subject_name = 'Bulgarian' AND t.school_id = 1 AND tm.name = 'Spring 2026' LIMIT 1;
INSERT INTO curricula (class_id, subject_id, teacher_id, term_id) SELECT c.class_id, s.subject_id, t.teacher_id, tm.term_id FROM classes c, subjects s, teachers t, terms tm WHERE c.school_id = 1 AND c.grade_level = 2 AND c.class_name = 'A' AND s.subject_name = 'Maths' AND t.school_id = 1 AND tm.name = 'Spring 2026' LIMIT 1;
INSERT INTO curricula (class_id, subject_id, teacher_id, term_id) SELECT c.class_id, s.subject_id, t.teacher_id, tm.term_id FROM classes c, subjects s, teachers t, terms tm WHERE c.school_id = 1 AND c.grade_level = 2 AND c.class_name = 'A' AND s.subject_name = 'English' AND t.school_id = 1 AND tm.name = 'Spring 2026' LIMIT 1;

-- Create Curricula - School 1 (GPAE) Grade 4A
INSERT INTO curricula (class_id, subject_id, teacher_id, term_id) SELECT c.class_id, s.subject_id, t.teacher_id, tm.term_id FROM classes c, subjects s, teachers t, terms tm WHERE c.school_id = 1 AND c.grade_level = 4 AND c.class_name = 'A' AND s.subject_name = 'Bulgarian' AND t.school_id = 1 AND tm.name = 'Spring 2026' LIMIT 1;
INSERT INTO curricula (class_id, subject_id, teacher_id, term_id) SELECT c.class_id, s.subject_id, t.teacher_id, tm.term_id FROM classes c, subjects s, teachers t, terms tm WHERE c.school_id = 1 AND c.grade_level = 4 AND c.class_name = 'A' AND s.subject_name = 'Maths' AND t.school_id = 1 AND tm.name = 'Spring 2026' LIMIT 1;
INSERT INTO curricula (class_id, subject_id, teacher_id, term_id) SELECT c.class_id, s.subject_id, t.teacher_id, tm.term_id FROM classes c, subjects s, teachers t, terms tm WHERE c.school_id = 1 AND c.grade_level = 4 AND c.class_name = 'A' AND s.subject_name = 'English' AND t.school_id = 1 AND tm.name = 'Spring 2026' LIMIT 1;
INSERT INTO curricula (class_id, subject_id, teacher_id, term_id) SELECT c.class_id, s.subject_id, t.teacher_id, tm.term_id FROM classes c, subjects s, teachers t, terms tm WHERE c.school_id = 1 AND c.grade_level = 4 AND c.class_name = 'A' AND s.subject_name = 'Geography' AND t.school_id = 1 AND tm.name = 'Spring 2026' LIMIT 1;

-- Create Curricula - School 2 (SMG) Grade 5A
INSERT INTO curricula (class_id, subject_id, teacher_id, term_id) SELECT c.class_id, s.subject_id, t.teacher_id, tm.term_id FROM classes c, subjects s, teachers t, terms tm WHERE c.school_id = 2 AND c.grade_level = 5 AND c.class_name = 'A' AND s.subject_name = 'Bulgarian' AND t.school_id = 2 AND tm.name = 'Spring 2026' LIMIT 1;
INSERT INTO curricula (class_id, subject_id, teacher_id, term_id) SELECT c.class_id, s.subject_id, t.teacher_id, tm.term_id FROM classes c, subjects s, teachers t, terms tm WHERE c.school_id = 2 AND c.grade_level = 5 AND c.class_name = 'A' AND s.subject_name = 'Maths' AND t.school_id = 2 AND tm.name = 'Spring 2026' LIMIT 1;
INSERT INTO curricula (class_id, subject_id, teacher_id, term_id) SELECT c.class_id, s.subject_id, t.teacher_id, tm.term_id FROM classes c, subjects s, teachers t, terms tm WHERE c.school_id = 2 AND c.grade_level = 5 AND c.class_name = 'A' AND s.subject_name = 'English' AND t.school_id = 2 AND tm.name = 'Spring 2026' LIMIT 1;
INSERT INTO curricula (class_id, subject_id, teacher_id, term_id) SELECT c.class_id, s.subject_id, t.teacher_id, tm.term_id FROM classes c, subjects s, teachers t, terms tm WHERE c.school_id = 2 AND c.grade_level = 5 AND c.class_name = 'A' AND s.subject_name = 'Physics' AND t.school_id = 2 AND tm.name = 'Spring 2026' LIMIT 1;

-- Create Curricula - School 3 (GPNE) Grade 12A
INSERT INTO curricula (class_id, subject_id, teacher_id, term_id) SELECT c.class_id, s.subject_id, t.teacher_id, tm.term_id FROM classes c, subjects s, teachers t, terms tm WHERE c.school_id = 3 AND c.grade_level = 12 AND c.class_name = 'A' AND s.subject_name = 'Bulgarian' AND t.school_id = 3 AND tm.name = 'Spring 2026' LIMIT 1;
INSERT INTO curricula (class_id, subject_id, teacher_id, term_id) SELECT c.class_id, s.subject_id, t.teacher_id, tm.term_id FROM classes c, subjects s, teachers t, terms tm WHERE c.school_id = 3 AND c.grade_level = 12 AND c.class_name = 'A' AND s.subject_name = 'Maths' AND t.school_id = 3 AND tm.name = 'Spring 2026' LIMIT 1;
INSERT INTO curricula (class_id, subject_id, teacher_id, term_id) SELECT c.class_id, s.subject_id, t.teacher_id, tm.term_id FROM classes c, subjects s, teachers t, terms tm WHERE c.school_id = 3 AND c.grade_level = 12 AND c.class_name = 'A' AND s.subject_name = 'Physics' AND t.school_id = 3 AND tm.name = 'Spring 2026' LIMIT 1;
INSERT INTO curricula (class_id, subject_id, teacher_id, term_id) SELECT c.class_id, s.subject_id, t.teacher_id, tm.term_id FROM classes c, subjects s, teachers t, terms tm WHERE c.school_id = 3 AND c.grade_level = 12 AND c.class_name = 'A' AND s.subject_name = 'Chemistry' AND t.school_id = 3 AND tm.name = 'Spring 2026' LIMIT 1;

-- Add grades for Student 1
INSERT INTO grades (student_id, curriculum_id, grade_value, grade_date) SELECT s.student_id, cur.curriculum_id, 5.50, '2026-03-15' FROM students s, curricula cur, classes c, subjects sub, users u WHERE s.user_id = u.user_id AND u.email = 'student1@app.com' AND cur.class_id = s.class_id AND cur.class_id = c.class_id AND cur.subject_id = sub.subject_id AND sub.subject_name = 'Bulgarian';
INSERT INTO grades (student_id, curriculum_id, grade_value, grade_date) SELECT s.student_id, cur.curriculum_id, 6.00, '2026-04-10' FROM students s, curricula cur, classes c, subjects sub, users u WHERE s.user_id = u.user_id AND u.email = 'student1@app.com' AND cur.class_id = s.class_id AND cur.class_id = c.class_id AND cur.subject_id = sub.subject_id AND sub.subject_name = 'Maths';
INSERT INTO grades (student_id, curriculum_id, grade_value, grade_date) SELECT s.student_id, cur.curriculum_id, 5.00, '2026-04-20' FROM students s, curricula cur, classes c, subjects sub, users u WHERE s.user_id = u.user_id AND u.email = 'student1@app.com' AND cur.class_id = s.class_id AND cur.class_id = c.class_id AND cur.subject_id = sub.subject_id AND sub.subject_name = 'English';

-- Add grades for Student 2
INSERT INTO grades (student_id, curriculum_id, grade_value, grade_date) SELECT s.student_id, cur.curriculum_id, 4.50, '2026-03-15' FROM students s, curricula cur, classes c, subjects sub, users u WHERE s.user_id = u.user_id AND u.email = 'student2@app.com' AND cur.class_id = s.class_id AND cur.class_id = c.class_id AND cur.subject_id = sub.subject_id AND sub.subject_name = 'Bulgarian';
INSERT INTO grades (student_id, curriculum_id, grade_value, grade_date) SELECT s.student_id, cur.curriculum_id, 5.50, '2026-04-10' FROM students s, curricula cur, classes c, subjects sub, users u WHERE s.user_id = u.user_id AND u.email = 'student2@app.com' AND cur.class_id = s.class_id AND cur.class_id = c.class_id AND cur.subject_id = sub.subject_id AND sub.subject_name = 'Maths';

-- Add grades for Student 11
INSERT INTO grades (student_id, curriculum_id, grade_value, grade_date) SELECT s.student_id, cur.curriculum_id, 5.75, '2026-03-20' FROM students s, curricula cur, classes c, subjects sub, users u WHERE s.user_id = u.user_id AND u.email = 'student11@app.com' AND cur.class_id = s.class_id AND cur.class_id = c.class_id AND cur.subject_id = sub.subject_id AND sub.subject_name = 'Bulgarian';
INSERT INTO grades (student_id, curriculum_id, grade_value, grade_date) SELECT s.student_id, cur.curriculum_id, 6.00, '2026-04-05' FROM students s, curricula cur, classes c, subjects sub, users u WHERE s.user_id = u.user_id AND u.email = 'student11@app.com' AND cur.class_id = s.class_id AND cur.class_id = c.class_id AND cur.subject_id = sub.subject_id AND sub.subject_name = 'Maths';
INSERT INTO grades (student_id, curriculum_id, grade_value, grade_date) SELECT s.student_id, cur.curriculum_id, 5.25, '2026-04-12' FROM students s, curricula cur, classes c, subjects sub, users u WHERE s.user_id = u.user_id AND u.email = 'student11@app.com' AND cur.class_id = s.class_id AND cur.class_id = c.class_id AND cur.subject_id = sub.subject_id AND sub.subject_name = 'Physics';

-- Add grades for Student 28
INSERT INTO grades (student_id, curriculum_id, grade_value, grade_date) SELECT s.student_id, cur.curriculum_id, 5.50, '2026-03-10' FROM students s, curricula cur, classes c, subjects sub, users u WHERE s.user_id = u.user_id AND u.email = 'student28@app.com' AND cur.class_id = s.class_id AND cur.class_id = c.class_id AND cur.subject_id = sub.subject_id AND sub.subject_name = 'Bulgarian';
INSERT INTO grades (student_id, curriculum_id, grade_value, grade_date) SELECT s.student_id, cur.curriculum_id, 6.00, '2026-03-25' FROM students s, curricula cur, classes c, subjects sub, users u WHERE s.user_id = u.user_id AND u.email = 'student28@app.com' AND cur.class_id = s.class_id AND cur.class_id = c.class_id AND cur.subject_id = sub.subject_id AND sub.subject_name = 'Maths';
INSERT INTO grades (student_id, curriculum_id, grade_value, grade_date) SELECT s.student_id, cur.curriculum_id, 5.75, '2026-04-08' FROM students s, curricula cur, classes c, subjects sub, users u WHERE s.user_id = u.user_id AND u.email = 'student28@app.com' AND cur.class_id = s.class_id AND cur.class_id = c.class_id AND cur.subject_id = sub.subject_id AND sub.subject_name = 'Physics';
INSERT INTO grades (student_id, curriculum_id, grade_value, grade_date) SELECT s.student_id, cur.curriculum_id, 5.00, '2026-04-18' FROM students s, curricula cur, classes c, subjects sub, users u WHERE s.user_id = u.user_id AND u.email = 'student28@app.com' AND cur.class_id = s.class_id AND cur.class_id = c.class_id AND cur.subject_id = sub.subject_id AND sub.subject_name = 'Chemistry';

-- Add grades for Student 31 (Director 1's child - Grade 1A)
INSERT INTO grades (student_id, curriculum_id, grade_value, grade_date) SELECT s.student_id, cur.curriculum_id, 5.25, '2026-03-12' FROM students s, curricula cur, classes c, subjects sub, users u WHERE s.user_id = u.user_id AND u.email = 'student31@app.com' AND cur.class_id = s.class_id AND cur.class_id = c.class_id AND cur.subject_id = sub.subject_id AND sub.subject_name = 'Bulgarian';
INSERT INTO grades (student_id, curriculum_id, grade_value, grade_date) SELECT s.student_id, cur.curriculum_id, 5.75, '2026-04-08' FROM students s, curricula cur, classes c, subjects sub, users u WHERE s.user_id = u.user_id AND u.email = 'student31@app.com' AND cur.class_id = s.class_id AND cur.class_id = c.class_id AND cur.subject_id = sub.subject_id AND sub.subject_name = 'Maths';
INSERT INTO grades (student_id, curriculum_id, grade_value, grade_date) SELECT s.student_id, cur.curriculum_id, 5.50, '2026-04-15' FROM students s, curricula cur, classes c, subjects sub, users u WHERE s.user_id = u.user_id AND u.email = 'student31@app.com' AND cur.class_id = s.class_id AND cur.class_id = c.class_id AND cur.subject_id = sub.subject_id AND sub.subject_name = 'English';

-- Add grades for Student 32 (Director 1's child - Grade 3A)
INSERT INTO grades (student_id, curriculum_id, grade_value, grade_date) SELECT s.student_id, cur.curriculum_id, 4.75, '2026-03-10' FROM students s, curricula cur, classes c, subjects sub, users u WHERE s.user_id = u.user_id AND u.email = 'student32@app.com' AND cur.class_id = s.class_id AND cur.class_id = c.class_id AND cur.subject_id = sub.subject_id AND sub.subject_name = 'Bulgarian';
INSERT INTO grades (student_id, curriculum_id, grade_value, grade_date) SELECT s.student_id, cur.curriculum_id, 5.00, '2026-04-05' FROM students s, curricula cur, classes c, subjects sub, users u WHERE s.user_id = u.user_id AND u.email = 'student32@app.com' AND cur.class_id = s.class_id AND cur.class_id = c.class_id AND cur.subject_id = sub.subject_id AND sub.subject_name = 'Maths';
INSERT INTO grades (student_id, curriculum_id, grade_value, grade_date) SELECT s.student_id, cur.curriculum_id, 5.25, '2026-04-14' FROM students s, curricula cur, classes c, subjects sub, users u WHERE s.user_id = u.user_id AND u.email = 'student32@app.com' AND cur.class_id = s.class_id AND cur.class_id = c.class_id AND cur.subject_id = sub.subject_id AND sub.subject_name = 'English';
INSERT INTO grades (student_id, curriculum_id, grade_value, grade_date) SELECT s.student_id, cur.curriculum_id, 4.50, '2026-04-22' FROM students s, curricula cur, classes c, subjects sub, users u WHERE s.user_id = u.user_id AND u.email = 'student32@app.com' AND cur.class_id = s.class_id AND cur.class_id = c.class_id AND cur.subject_id = sub.subject_id AND sub.subject_name = 'History';

-- Add grades for Student 33 (Teacher 5's child - Grade 2A)
INSERT INTO grades (student_id, curriculum_id, grade_value, grade_date) SELECT s.student_id, cur.curriculum_id, 5.50, '2026-03-14' FROM students s, curricula cur, classes c, subjects sub, users u WHERE s.user_id = u.user_id AND u.email = 'student33@app.com' AND cur.class_id = s.class_id AND cur.class_id = c.class_id AND cur.subject_id = sub.subject_id AND sub.subject_name = 'Bulgarian';
INSERT INTO grades (student_id, curriculum_id, grade_value, grade_date) SELECT s.student_id, cur.curriculum_id, 5.00, '2026-04-06' FROM students s, curricula cur, classes c, subjects sub, users u WHERE s.user_id = u.user_id AND u.email = 'student33@app.com' AND cur.class_id = s.class_id AND cur.class_id = c.class_id AND cur.subject_id = sub.subject_id AND sub.subject_name = 'Maths';
INSERT INTO grades (student_id, curriculum_id, grade_value, grade_date) SELECT s.student_id, cur.curriculum_id, 5.75, '2026-04-16' FROM students s, curricula cur, classes c, subjects sub, users u WHERE s.user_id = u.user_id AND u.email = 'student33@app.com' AND cur.class_id = s.class_id AND cur.class_id = c.class_id AND cur.subject_id = sub.subject_id AND sub.subject_name = 'English';

-- Add grades for Student 34 (Teacher 5's child - Grade 4A)
INSERT INTO grades (student_id, curriculum_id, grade_value, grade_date) SELECT s.student_id, cur.curriculum_id, 4.25, '2026-03-11' FROM students s, curricula cur, classes c, subjects sub, users u WHERE s.user_id = u.user_id AND u.email = 'student34@app.com' AND cur.class_id = s.class_id AND cur.class_id = c.class_id AND cur.subject_id = sub.subject_id AND sub.subject_name = 'Bulgarian';
INSERT INTO grades (student_id, curriculum_id, grade_value, grade_date) SELECT s.student_id, cur.curriculum_id, 4.50, '2026-04-03' FROM students s, curricula cur, classes c, subjects sub, users u WHERE s.user_id = u.user_id AND u.email = 'student34@app.com' AND cur.class_id = s.class_id AND cur.class_id = c.class_id AND cur.subject_id = sub.subject_id AND sub.subject_name = 'Maths';
INSERT INTO grades (student_id, curriculum_id, grade_value, grade_date) SELECT s.student_id, cur.curriculum_id, 5.00, '2026-04-13' FROM students s, curricula cur, classes c, subjects sub, users u WHERE s.user_id = u.user_id AND u.email = 'student34@app.com' AND cur.class_id = s.class_id AND cur.class_id = c.class_id AND cur.subject_id = sub.subject_id AND sub.subject_name = 'English';
INSERT INTO grades (student_id, curriculum_id, grade_value, grade_date) SELECT s.student_id, cur.curriculum_id, 4.75, '2026-04-24' FROM students s, curricula cur, classes c, subjects sub, users u WHERE s.user_id = u.user_id AND u.email = 'student34@app.com' AND cur.class_id = s.class_id AND cur.class_id = c.class_id AND cur.subject_id = sub.subject_id AND sub.subject_name = 'Geography';

-- Add absences for Student 1
INSERT INTO absences (student_id, curriculum_id, absence_date, is_excused) SELECT s.student_id, cur.curriculum_id, '2026-03-05', true FROM students s, curricula cur, classes c, subjects sub, users u WHERE s.user_id = u.user_id AND u.email = 'student1@app.com' AND cur.class_id = s.class_id AND cur.class_id = c.class_id AND cur.subject_id = sub.subject_id AND sub.subject_name = 'Bulgarian';
INSERT INTO absences (student_id, curriculum_id, absence_date, is_excused) SELECT s.student_id, cur.curriculum_id, '2026-04-02', false FROM students s, curricula cur, classes c, subjects sub, users u WHERE s.user_id = u.user_id AND u.email = 'student1@app.com' AND cur.class_id = s.class_id AND cur.class_id = c.class_id AND cur.subject_id = sub.subject_id AND sub.subject_name = 'Maths';

-- Add absences for Student 2
INSERT INTO absences (student_id, curriculum_id, absence_date, is_excused) SELECT s.student_id, cur.curriculum_id, '2026-03-12', true FROM students s, curricula cur, classes c, subjects sub, users u WHERE s.user_id = u.user_id AND u.email = 'student2@app.com' AND cur.class_id = s.class_id AND cur.class_id = c.class_id AND cur.subject_id = sub.subject_id AND sub.subject_name = 'English';

-- Add absences for Student 11
INSERT INTO absences (student_id, curriculum_id, absence_date, is_excused) SELECT s.student_id, cur.curriculum_id, '2026-03-18', false FROM students s, curricula cur, classes c, subjects sub, users u WHERE s.user_id = u.user_id AND u.email = 'student11@app.com' AND cur.class_id = s.class_id AND cur.class_id = c.class_id AND cur.subject_id = sub.subject_id AND sub.subject_name = 'Physics';
INSERT INTO absences (student_id, curriculum_id, absence_date, is_excused) SELECT s.student_id, cur.curriculum_id, '2026-04-22', true FROM students s, curricula cur, classes c, subjects sub, users u WHERE s.user_id = u.user_id AND u.email = 'student11@app.com' AND cur.class_id = s.class_id AND cur.class_id = c.class_id AND cur.subject_id = sub.subject_id AND sub.subject_name = 'Bulgarian';

-- Add absences for Student 28
INSERT INTO absences (student_id, curriculum_id, absence_date, is_excused) SELECT s.student_id, cur.curriculum_id, '2026-03-22', false FROM students s, curricula cur, classes c, subjects sub, users u WHERE s.user_id = u.user_id AND u.email = 'student28@app.com' AND cur.class_id = s.class_id AND cur.class_id = c.class_id AND cur.subject_id = sub.subject_id AND sub.subject_name = 'Chemistry';

-- Add absences for Student 31 (Director 1's child - Grade 1A)
INSERT INTO absences (student_id, curriculum_id, absence_date, is_excused) SELECT s.student_id, cur.curriculum_id, '2026-03-08', true FROM students s, curricula cur, classes c, subjects sub, users u WHERE s.user_id = u.user_id AND u.email = 'student31@app.com' AND cur.class_id = s.class_id AND cur.class_id = c.class_id AND cur.subject_id = sub.subject_id AND sub.subject_name = 'Bulgarian';
INSERT INTO absences (student_id, curriculum_id, absence_date, is_excused) SELECT s.student_id, cur.curriculum_id, '2026-04-18', false FROM students s, curricula cur, classes c, subjects sub, users u WHERE s.user_id = u.user_id AND u.email = 'student31@app.com' AND cur.class_id = s.class_id AND cur.class_id = c.class_id AND cur.subject_id = sub.subject_id AND sub.subject_name = 'Maths';

-- Add absences for Student 32 (Director 1's child - Grade 3A)
INSERT INTO absences (student_id, curriculum_id, absence_date, is_excused) SELECT s.student_id, cur.curriculum_id, '2026-03-16', true FROM students s, curricula cur, classes c, subjects sub, users u WHERE s.user_id = u.user_id AND u.email = 'student32@app.com' AND cur.class_id = s.class_id AND cur.class_id = c.class_id AND cur.subject_id = sub.subject_id AND sub.subject_name = 'English';
INSERT INTO absences (student_id, curriculum_id, absence_date, is_excused) SELECT s.student_id, cur.curriculum_id, '2026-04-10', false FROM students s, curricula cur, classes c, subjects sub, users u WHERE s.user_id = u.user_id AND u.email = 'student32@app.com' AND cur.class_id = s.class_id AND cur.class_id = c.class_id AND cur.subject_id = sub.subject_id AND sub.subject_name = 'History';
INSERT INTO absences (student_id, curriculum_id, absence_date, is_excused) SELECT s.student_id, cur.curriculum_id, '2026-04-25', true FROM students s, curricula cur, classes c, subjects sub, users u WHERE s.user_id = u.user_id AND u.email = 'student32@app.com' AND cur.class_id = s.class_id AND cur.class_id = c.class_id AND cur.subject_id = sub.subject_id AND sub.subject_name = 'Bulgarian';

-- Add absences for Student 33 (Teacher 5's child - Grade 2A)
INSERT INTO absences (student_id, curriculum_id, absence_date, is_excused) SELECT s.student_id, cur.curriculum_id, '2026-03-09', false FROM students s, curricula cur, classes c, subjects sub, users u WHERE s.user_id = u.user_id AND u.email = 'student33@app.com' AND cur.class_id = s.class_id AND cur.class_id = c.class_id AND cur.subject_id = sub.subject_id AND sub.subject_name = 'Maths';
INSERT INTO absences (student_id, curriculum_id, absence_date, is_excused) SELECT s.student_id, cur.curriculum_id, '2026-04-19', true FROM students s, curricula cur, classes c, subjects sub, users u WHERE s.user_id = u.user_id AND u.email = 'student33@app.com' AND cur.class_id = s.class_id AND cur.class_id = c.class_id AND cur.subject_id = sub.subject_id AND sub.subject_name = 'Bulgarian';

-- Add absences for Student 34 (Teacher 5's child - Grade 4A)
INSERT INTO absences (student_id, curriculum_id, absence_date, is_excused) SELECT s.student_id, cur.curriculum_id, '2026-03-17', true FROM students s, curricula cur, classes c, subjects sub, users u WHERE s.user_id = u.user_id AND u.email = 'student34@app.com' AND cur.class_id = s.class_id AND cur.class_id = c.class_id AND cur.subject_id = sub.subject_id AND sub.subject_name = 'English';
INSERT INTO absences (student_id, curriculum_id, absence_date, is_excused) SELECT s.student_id, cur.curriculum_id, '2026-04-11', false FROM students s, curricula cur, classes c, subjects sub, users u WHERE s.user_id = u.user_id AND u.email = 'student34@app.com' AND cur.class_id = s.class_id AND cur.class_id = c.class_id AND cur.subject_id = sub.subject_id AND sub.subject_name = 'Geography';
INSERT INTO absences (student_id, curriculum_id, absence_date, is_excused) SELECT s.student_id, cur.curriculum_id, '2026-04-26', false FROM students s, curricula cur, classes c, subjects sub, users u WHERE s.user_id = u.user_id AND u.email = 'student34@app.com' AND cur.class_id = s.class_id AND cur.class_id = c.class_id AND cur.subject_id = sub.subject_id AND sub.subject_name = 'Bulgarian';

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM absences WHERE student_id IN (SELECT student_id FROM students WHERE user_id IN (SELECT user_id FROM users WHERE email LIKE '%@app.com'));
DELETE FROM grades WHERE student_id IN (SELECT student_id FROM students WHERE user_id IN (SELECT user_id FROM users WHERE email LIKE '%@app.com'));
DELETE FROM curricula WHERE curriculum_id IN (SELECT curriculum_id FROM curricula WHERE teacher_id IN (SELECT teacher_id FROM teachers WHERE user_id IN (SELECT user_id FROM users WHERE email LIKE '%@app.com')));
DELETE FROM student_parents WHERE student_id IN (SELECT student_id FROM students WHERE user_id IN (SELECT user_id FROM users WHERE email LIKE '%@app.com'));
DELETE FROM students WHERE user_id IN (SELECT user_id FROM users WHERE email LIKE 'student%@app.com');
DELETE FROM parents WHERE user_id IN (SELECT user_id FROM users WHERE email LIKE 'parent%@app.com');
DELETE FROM teacher_subjects WHERE teacher_id IN (SELECT teacher_id FROM teachers WHERE user_id IN (SELECT user_id FROM users WHERE email LIKE 'teacher%@app.com'));
DELETE FROM teachers WHERE user_id IN (SELECT user_id FROM users WHERE email LIKE 'teacher%@app.com');
DELETE FROM directors WHERE user_id IN (SELECT user_id FROM users WHERE email LIKE 'director%@app.com');
DELETE FROM user_roles WHERE user_id IN (SELECT user_id FROM users WHERE email LIKE '%@app.com');
DELETE FROM users WHERE email LIKE '%@app.com';
-- +goose StatementEnd


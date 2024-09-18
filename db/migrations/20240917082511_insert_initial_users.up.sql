INSERT INTO users (name, occupation, email, password_hash, avatar_file_name, role, token, created_at, updated_at)
VALUES 
('John Doe', 'Software Developer', 'johndoe@example.com', 'hashedpassword123', 'avatar.png', 'user', 'token123', NOW(), NOW()),
('Jane Smith', 'Product Manager', 'janesmith@example.com', 'hashedpassword456', 'avatar2.png', 'admin', 'token456', NOW(), NOW());
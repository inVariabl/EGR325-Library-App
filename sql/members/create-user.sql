CREATE USER IF NOT EXISTS 'libraryuser'@'localhost' IDENTIFIED BY 'password';
GRANT ALL PRIVILEGES ON library_db.* TO 'libraryuser'@'localhost';
FLUSH PRIVILEGES;
exit;

-- Generated via SQLite Manager (free Firefox plugin)
CREATE TABLE "uploads" ("id" INTEGER PRIMARY KEY  AUTOINCREMENT  NOT NULL  UNIQUE , "uuid" TEXT NOT NULL  UNIQUE , "original_filename" TEXT NOT NULL , "storage_path" TEXT NOT NULL  UNIQUE , "upload_date" DATETIME DEFAULT CURRENT_TIMESTAMP)
CREATE TABLE sqlite_sequence(name,seq)

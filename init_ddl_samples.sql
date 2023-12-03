--This should be connecting to the remote db and executing the following,
-- so our server will essentially host the API, services, and operate on DB connection
DROP TABLE IF EXISTS languages;
DROP TABLE IF EXISTS flashcards;

CREATE TABLE languages (
  id         INT AUTO_INCREMENT NOT NULL,
  host_lang      VARCHAR(128) NOT NULL,
  target_lang     VARCHAR(255) NOT NULL,
  priority      INT AUTO_INCREMENT NOT NULL,
  PRIMARY KEY (`id`)
);

CREATE TABLE flashcards (
  id         INT AUTO_INCREMENT NOT NULL,
  front      VARCHAR(128) NOT NULL,
  back     VARCHAR(255) NOT NULL,
  hint    VARCHAR(255),
  difficulty INT,
  lang      REFERENCES languages(ID),
  PRIMARY KEY(`id`) 

);

INSERT INTO languages
  (host_lang, target_lang, priority)
VALUES
  ('English', 'French', 1),
  ('English', 'Spanish', 2);

--For now, just set the init values. Maybe build on persistent store.
SET @engFren = (SELECT id from languages where host_lang = 'English' and target_lang = 'French')
SET @engSpan = (SELECT id from languages where host_lang = 'English' and target_lang = 'Spanish')


INSERT INTO flashcards
 (front, back, difficulty, hint, lang)
VALUES 
  ("Blue", "Bleu", 1, "Bleh", @engFren);

INSERT INTO flashcards
 (front, back, difficulty, hint, lang)
VALUES 
  ("Blue", "Azul", 1, "Spanglish", @engSpan);
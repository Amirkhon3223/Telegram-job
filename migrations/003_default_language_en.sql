-- Set default interface language to English
ALTER TABLE users ALTER COLUMN interface_language SET DEFAULT 'en';

-- Update existing users with no language set to 'en'
UPDATE users SET interface_language = 'en' WHERE interface_language IS NULL OR interface_language = '';

-- Set NOT NULL constraint
ALTER TABLE users ALTER COLUMN interface_language SET NOT NULL;

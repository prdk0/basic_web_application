CREATE OR REPLACE FUNCTION update_modified_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = current_timestamp;
    RETURN NEW;
END;

$$ language 'plpgsql';
CREATE TRIGGER update_room_restriction_modtime
BEFORE UPDATE ON room_restrictions
FOR EACH ROW
EXECUTE FUNCTION update_modified_column();
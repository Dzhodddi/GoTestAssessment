ALTER TABLE targets
    ADD CONSTRAINT unique_mission_target_name UNIQUE (mission_id, name);

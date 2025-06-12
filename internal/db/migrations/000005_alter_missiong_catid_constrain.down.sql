ALTER TABLE missions
DROP CONSTRAINT missions_cat_id_fkey;

ALTER TABLE missions
    ADD CONSTRAINT missions_cat_id_fkey
        FOREIGN KEY (cat_id) REFERENCES spycat(id) ON DELETE RESTRICT;

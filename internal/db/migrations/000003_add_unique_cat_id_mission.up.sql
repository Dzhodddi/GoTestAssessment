CREATE UNIQUE INDEX unique_cat_mission ON missions(cat_id) WHERE cat_id IS NOT NULL;

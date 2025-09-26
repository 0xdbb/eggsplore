-- +goose Up
-- +goose StatementBegin
ALTER TABLE accounts
ADD CONSTRAINT department_check 
CHECK (department IN (
  'Minerals Commission',
  'Lands Commission',
  'Forestry Commission',
  'Office of Administrator of Stool lands',
  'Environmental commission',
  'Goldbod',
  'Ghana Police',
  'Ghana Army',
  'National Security',
  'Geological Survey Authority',
  'Wildlife Division',
  'Ghana Space Science and Technology Institute',
  'Water Resources Commission',
  'Land Use and Spatial Planning',
  'National Anti-Illegal Mining Operations Secretariat',
  'GADE Team'
));
-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
ALTER TABLE accounts
DROP CONSTRAINT IF EXISTS department_check;
-- +goose StatementEnd

CREATE INDEX idx_capsules_user_id ON capsules(user_id);
CREATE INDEX idx_capsules_due_status ON capsules(due_date, status);
